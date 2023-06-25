package main

import (
	"context"
	"controller/api"
	"controller/pkg/logger"
	"controller/pkg/montior"
	"controller/task_tracker/config"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"controller/task_tracker/param"
	"controller/task_tracker/services"
	"controller/task_tracker/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"testing"
	"time"
)

var taskTracker *services.Service
var orderState *database.OrderState
var uploadRequest *database.UploadRequest
var downloadRequest *database.DownloadRequest
var fidReplicate *database.FidReplication
var taskInfo *database.Task

const (
	NEW_FILE = iota + 1
	REPEATE_FILE
	PART_REPEATE
)

func Init() error {
	c, err := config.LoadServerConfig(LocalServiceId)
	if err != nil {
		logger.Error(err)
		return err
	}
	config.ServerCfg = c

	logger.Infof("config: %v,  schedulerUrl: %v", *config.ServerCfg, config.ServerCfg.Scheduler.Url)

	ding := new(montior.DingTalk)
	ding.Init("task_tracker", config.ServerCfg.Montior.AccessToken, config.ServerCfg.Montior.Secret)
	//init logger
	logger.InitLoggerWithDingTalk(c.Logger.Level, ding)

	if err := database.InitDB(c.DB.Url, c.DB.DbUser, c.DB.DbPassword, c.DB.DbName, c.DB.Timeout); err != nil {
		logger.Errorf("init database: %v fail: %v", c.DB, err.Error())
		return nil
	}

	uploadRequest = new(database.UploadRequest)
	uploadRequest.Init(database.Db)

	orderState = new(database.OrderState)
	orderState.Init(database.Db)

	fidReplicate = new(database.FidReplication)
	fidReplicate.Init(database.Db)

	downloadRequest = new(database.DownloadRequest)
	downloadRequest.Init(database.Db)

	taskTracker = new(services.Service)
	taskTracker.Init(database.Db)

	taskInfo = new(database.Task)
	taskInfo.Init(database.Db)

	return nil
}
func GetDeleteRequest(order *dict.OrderStateInfo) *param.DeleteFidRequest {
	//获取上传成功的订单号：

	deleteFidRequest := &param.DeleteFidRequest{
		RequestId: uuid.NewV4().String(),
		UserId:    "46",
		OrderId:   order.OrderId,
		Fids:      make(map[string]bool),
	}

	for _, task := range order.Tasks {
		deleteFidRequest.Fids[task.Fid] = true
	}

	return deleteFidRequest
}

func TestDeleteFid(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	orders, err := orderState.LoadByStatus(dict.TASK_CHARGE_SUC)
	if err != nil {
		fmt.Printf("don't have order\n")
		return
	}

	num := 2
	if num > len(*orders) {
		num = len(*orders)
	}

	//to do gomonkey

	patch := gomonkey.ApplyFunc(utils.DeleteFid, func(request *param.DeleteOrderFidRequest) (*param.DeleteOrderFidResponse, error) {
		rsp := &param.DeleteOrderFidResponse{
			Status:  param.SUCCESS,
			OrderId: "orderId",
			Tasks:   make(map[string]*[]*param.UploadTask),
		}

		for fid, tasks := range request.Tasks {
			for _, task := range tasks {
				task.Status = 1
			}

			rsp.Tasks[fid] = &tasks
		}

		return rsp, nil
	})

	defer patch.Reset()

	for i := 0; i < num; i++ {
		ret, err := taskTracker.DeleteFid(GetDeleteRequest(&(*orders)[i]))
		fmt.Printf("orderId: %v, ret: %v, err: %v\n", (*orders)[i], ret, err)
		assert.Equal(t, nil, err)
		deleteRsp := ret.(*param.DeleteFidResponse)
		assert.Equal(t, param.SUCCESS, deleteRsp.Status)
		//check delete success.

		//再次查询订单，判断是否删除成功。
		states, err := orderState.GetOrderStateByOrderId((*orders)[i].OrderId)
		if err != nil || len(*states) < 1 {
			continue
		}

		for fid, task := range (*states)[0].Tasks {
			deleteStatus, ok := deleteRsp.Fids[fid]
			if !ok {
				continue
			}

			fidStatus := dict.TASK_CHARGE_SUC
			if deleteStatus == param.SUCCESS { //不重复文件删除
				fidStatus = dict.TASK_DEL_SUC
			}
			// 没有删除成功 或者是重复文件，保持状态不变。
			for _, rep := range task.Reps {
				assert.Equal(t, fidStatus, rep.Status)
			}
		}
	}
}

//nType : 0;

func GetUploadTask(nTaskNum, nType int) *dict.UploadRequestInfo {
	uploadTask := &dict.UploadRequestInfo{
		RequestId:  uuid.NewV4().String(),
		UserId:     "46",
		UploadType: param.NoDevice,
		Group:      "chengdu",
		PieceNum:   nTaskNum,
		//Tasks:      make([]*dict.UploadTask, 0, nTaskNum),
	}

	return uploadTask
}

func getUploadPieceFidRequest(orderId string, nTaskNum int, nType int) *api.UploadPieceFidRequest {
	uploadPieceFidRequest := &api.UploadPieceFidRequest{
		RequestId: uuid.NewV4().String(),
		OrderId:   orderId,
		Group:     "chengdu",
		Pieces:    make([]*api.PieceFid, 0, nTaskNum),
		Ext:       &api.Extend{Ctx: context.Background()},
	}

	repeateFidInfo, err := fidReplicate.Load(nTaskNum)
	if err != nil {
		fmt.Printf("load fidReplicate fail , err: %v\n", err.Error())
		return nil
	}

	nLen := len(*repeateFidInfo)

	if nType == REPEATE_FILE {
		nTaskNum = nLen
	} else if nType == PART_REPEATE { //一半重复fid, 一半新fid
		nLen = nLen / 2
	}
	//分类。

	index := 0
	for i := 0; i < nTaskNum; i++ {
		fid := fmt.Sprintf("fid_%v_%v_%v", time.Now().UnixNano(), rand.Intn(100000), i)
		if (nType == REPEATE_FILE || nType == PART_REPEATE) && index < nLen {
			fid = (*repeateFidInfo)[index].Fid
		}

		pieceFid := &api.PieceFid{
			Fid:        fid,
			Rep:        0,
			RepMin:     rand.Intn(3),
			RepMax:     3 + rand.Intn(5),
			Encryption: 0,
			Expire:     10000000,
			Level:      0,
			Size:       rand.Intn(10000),
			Name:       fmt.Sprintf("testfile_%v.txt", i),
			Ajust:      0,
		}

		uploadPieceFidRequest.Pieces = append(uploadPieceFidRequest.Pieces, pieceFid)
		index++
	}

	return uploadPieceFidRequest
}

//orderId;
func TestCreateUploadTask(t *testing.T) {
	//create replicate strategy, init task, init request.
	if err := Init(); err != nil {
		return
	}

	num := 100
	for i := 0; i < num; i++ {
		createUploadTask(t, i, NEW_FILE)
	}

	for i := 0; i < num; i++ {
		createUploadTask(t, i, PART_REPEATE)
	}

	for i := 0; i < num; i++ {
		createUploadTask(t, i, REPEATE_FILE)
	}
}

func createUploadTask(t *testing.T, nTaskNum, nType int) {
	taskInfo := GetUploadTask(nTaskNum, nType)

	if err := uploadRequest.Add(taskInfo); err != nil {
		fmt.Printf("upload request add fail")
		return
	}

	ret, err := taskTracker.CreateTask(&api.CreateTaskRequest{RequestId: taskInfo.RequestId, Type: param.UPLOAD})

	if nTaskNum == 0 {
		assert.NotEqual(t, nil, err)
		return
	}

	//check rsp
	assert.Equal(t, nil, err)
	uploadTaskResponse := ret.(*api.CreateTaskResponse)
	assert.Equal(t, param.SUCCESS, uploadTaskResponse.Status)

	//fmt.Printf("uploadTaskResponse.Fids: %v\n", uploadTaskResponse.Fids)
	assert.NotEmpty(t, uploadTaskResponse.OrderId)

	//

	uploadPieceFidRequest := getUploadPieceFidRequest(uploadTaskResponse.OrderId, nTaskNum, nType)
	if len(uploadPieceFidRequest.Pieces) != nTaskNum {
		return
	}

	saveTaskInfo(uploadTaskResponse.OrderId, uploadPieceFidRequest)

	rsp, err := taskTracker.UploadPieceFid(uploadPieceFidRequest)
	assert.Equal(t, nil, err)
	uploadPieceFidResponse := rsp.(*api.UploadPieceFidResponse)
	assert.Equal(t, param.SUCCESS, uploadPieceFidResponse.Status)

	//check order state.
	states, err := orderState.GetOrderStateByOrderId(uploadTaskResponse.OrderId)
	if err != nil {
		fmt.Printf("get orderstate from db fail: %v", err.Error())
		return
	}

	assert.Equal(t, 1, len(*states))
	assert.Equal(t, taskInfo.PieceNum, len((*states)[0].Tasks))

	bRepeate := false
	for _, task := range (*states)[0].Tasks { //订单初始化。
		if task.Repeate == dict.REPEATE {
			assert.Equal(t, dict.TASK_UPLOAD_SUC, task.Status)
			assert.NotEqual(t, "", task.Cid)
			_, ok := uploadPieceFidResponse.RepFids[task.Fid]
			assert.Equal(t, true, ok)

			bRepeate = true
		} else {
			assert.Equal(t, dict.TASK_INIT, task.Status)
			assert.Equal(t, "", task.Cid)
		}
	}
	if bRepeate {
		assert.NotEqual(t, 0, len(uploadPieceFidResponse.RepFids))
	} else {
		assert.Equal(t, 0, len(uploadPieceFidResponse.RepFids))
	}
}

func GetDownloadTask(nTaskNum, nType int) *dict.DownloadRequestInfo {
	downloadTask := &dict.DownloadRequestInfo{
		RequestId:   uuid.NewV4().String(),
		UserId:      "46",
		DownlodType: param.NoDevice,
		Tasks:       make([]*dict.DownloadTask, 0, nTaskNum),
	}

	if nTaskNum == 0 {
		return downloadTask
	}

	repeateFidInfo, err := fidReplicate.Load(nTaskNum)
	if err != nil {
		fmt.Printf("load fidReplicate fail , err: %v\n", err.Error())
		return nil
	}

	nLen := len(*repeateFidInfo)

	if nType == REPEATE_FILE {
		nTaskNum = nLen
	} else if nType == PART_REPEATE { //一半重复fid, 一半新fid
		nLen = nLen / 2
	}
	//分类。

	index := 0
	for i := 0; i < nTaskNum; i++ {
		cid := fmt.Sprintf("fid_%v_%v_%v", time.Now().UnixNano(), rand.Intn(100000), i)
		if (nType == REPEATE_FILE || nType == PART_REPEATE) && index < nLen {
			cid = (*repeateFidInfo)[index].Cid
		}

		task := &dict.DownloadTask{
			Cid:  cid,
			Size: rand.Intn(10000),
			Name: fmt.Sprintf("testfile_%v.txt", i),
		}

		downloadTask.Tasks = append(downloadTask.Tasks, task)
		index++
	}

	return downloadTask
}

func TestCreateDownloadTask(t *testing.T) {
	//create replicate strategy, init task, init request.
	if err := Init(); err != nil {
		return
	}

	num := 100
	for i := 0; i < num; i++ {
		createDownloadTask(t, i, NEW_FILE)
	}

	for i := 0; i < num; i++ {
		createDownloadTask(t, i, PART_REPEATE)
	}

	for i := 0; i < num; i++ {
		createDownloadTask(t, i, REPEATE_FILE)
	}
}

func createDownloadTask(t *testing.T, nTaskNum, nType int) {
	if err := Init(); err != nil {
		return
	}

	taskInfo := GetDownloadTask(nTaskNum, nType)

	if err := downloadRequest.Add(taskInfo); err != nil {
		fmt.Printf("upload request add fail")
		return
	}

	ret, err := taskTracker.CreateTask(&api.CreateTaskRequest{RequestId: taskInfo.RequestId, Type: param.DOWNLOAD})

	if nTaskNum == 0 {
		assert.NotEqual(t, nil, err)
		return
	}

	//check rsp
	assert.Equal(t, nil, err)
	downloadTaskResponse := ret.(param.DownloadTaskResponse)
	assert.Equal(t, param.SUCCESS, downloadTaskResponse.Status)

	//fmt.Printf("uploadTaskResponse.Fids: %v\n", uploadTaskResponse.Fids)
	assert.NotEmpty(t, downloadTaskResponse.OrderId)

	//check order
	states, err := orderState.GetOrderStateByOrderId(downloadTaskResponse.OrderId)
	if err != nil {
		fmt.Printf("get orderstate from db fail: %v", err.Error())
		return
	}

	assert.Equal(t, 1, len(*states))
	assert.Equal(t, len(taskInfo.Tasks), len((*states)[0].Tasks))

	for _, task := range (*states)[0].Tasks { //订单初始化。
		assert.Equal(t, dict.TASK_INIT, task.Status)
	}
	assert.Equal(t, dict.TASK_INIT, (*states)[0].Status)
}

//
func TestCallbackUpload(t *testing.T) {
	if err := Init(); err != nil {
		return
	}
	patch := gomonkey.ApplyFunc(utils.Replicate, func(request *param.ReplicationRequest) (*param.ReplicationResponse, error) {
		rsp := &param.ReplicationResponse{
			Status:  param.SUCCESS,
			OrderId: request.OrderId,
			Tasks:   make([]*param.TaskResponse, 0, len(request.Tasks)),
		}

		for fid, task := range request.Tasks { //1.所有集群都备份成功， 2. 部分集群备份成功
			rspTask := new(param.TaskResponse)
			rspTask.RegionStatus = make(map[string]int)
			for region, _ := range task.Reps {
				rspTask.RegionStatus[region] = param.SUCCESS
				rspTask.Fid = fid
				rspTask.Cid = task.Cid
			}

			rsp.Tasks = append(rsp.Tasks, rspTask)
		}

		return rsp, nil
	})

	defer patch.Reset()

	orders, err := orderState.Query(bson.M{"status": dict.TASK_INIT, "order_type": param.UPLOAD}, 10) //测试10个订单
	if err != nil {
		fmt.Printf("load init order fail: %v \n", err.Error())
	}

	//单测代码，有些订单，可能没有创建策略。

	for _, order := range *orders {

		rsp, err := utils.CreateStrategy(&param.CreateStrategyRequest{RequestId: uuid.NewV4().String(), OrderId: order.OrderId, Region: "chengdu"})
		if err != nil {
			fmt.Printf("create stratey fail : %v, orderId: %v\n", err.Error(), order.OrderId)
			return
		}
		if rsp.Status != param.SUCCESS {
			fmt.Printf("create stratey fail, status: %v, orderId: %v\n", rsp.Status, order.OrderId)
			return
		}

		bNewFile := false
		for _, task := range order.Tasks {
			if task.Repeate != dict.REPEATE {
				bNewFile = true
			}

			ret, err := taskTracker.CallbackUpload(&param.CallbackUploadRequest{
				OrderId: order.OrderId,
				Fid:     task.Fid,
				Cid:     fmt.Sprintf("test_cid_%v", uuid.NewV4().String()),
				Region:  "chengdu",
				Origins: "test_origins:127.0.0.1",
				Status:  param.SUCCESS,
			})

			assert.Equal(t, nil, err)

			rsp := ret.(param.CallbackUploadResponse)
			assert.Equal(t, param.SUCCESS, rsp.Status)
		}

		fmt.Printf("orderId: %v\n", order.OrderId)

		state, err := orderState.Query(bson.M{"order_id": order.OrderId}, 1)
		if err != nil {
			fmt.Printf("query order state fail: %v \n", err.Error())
		}

		if bNewFile { //包含新文件，则会备份， 全部是重复文件会等待上传完成通知。
			assert.Equal(t, dict.TASK_BEGIN_REP, (*state)[0].Status)
		}

	}
}

func TestCheckReplicate(t *testing.T) {
	patch := gomonkey.ApplyFunc(utils.SearchRep, func(request *dict.UploadFinishOrder) (*param.GetOrderRepResponse, error) {
		rsp := &param.GetOrderRepResponse{
			Status:  param.SUCCESS,
			OrderId: request.OrderId,
			Tasks:   make(map[string]*dict.TaskRepInfo),
		}

		for _, task := range request.Tasks { //1.所有集群都备份成功， 2. 部分集群备份成功
			rspTask := new(dict.TaskRepInfo)

			rspTask.Regions = make(map[string]*dict.RegionRep)

			for _, region := range task.Regions {
				rspTask.Regions[region] = &dict.RegionRep{
					Region:      region,
					CurRep:      3,
					Status:      param.SUCCESS,
					CheckStatus: param.SUCCESS,
				}
				rspTask.Fid = task.Fid
				rspTask.Cid = task.Cid
			}

			rsp.Tasks[task.Fid] = rspTask
		}

		return rsp, nil
	})

	defer patch.Reset()

	if err := Init(); err != nil {
		return
	}

	time.Sleep(100 * time.Second) //系统自动处理开始备份订单,下次程序启动订单全部备份成功。
	//to do proc, init success.

	/*time.Sleep(100)
	states, err := orderState.Query(bson.M{"status": dict.TASK_BEGIN_REP}, 10)
	if err != nil {
		fmt.Printf("query dict.TASK_BEGIN_REP order fail \n")
		return
	}*/

}

func TestCallbackDownload(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	limit := 10
	orders, err := orderState.Query(bson.M{"status": dict.TASK_INIT, "order_type": param.DOWNLOAD}, limit)
	if err != nil {
		fmt.Printf("query dict.TASK_INIT  download order order fail: %v \n", err.Error())
		return
	}

	for _, order := range *orders {
		for _, task := range order.Tasks {
			ret, err := taskTracker.CallbackDownload(&param.CallbackDownloadRequest{Cid: task.Cid, Region: task.Region, Status: param.SUCCESS, Origins: "test_origins/127.0.0.1", Ext: order.OrderId})
			assert.Equal(t, nil, err)

			rsp := ret.(param.CallbackDownloadResponse)
			assert.Equal(t, param.SUCCESS, rsp.Status)
		}
		time.Sleep(1 * time.Second) //等待计费成功。

		orders, err := orderState.Query(bson.M{"order_id": order.OrderId}, 1)
		if err != nil {
			fmt.Printf("query order_id: %v fail: %v \n", order.OrderId, err.Error())
			return
		}

		assert.Equal(t, dict.TASK_CHARGE_SUC, (*orders)[0].Status)
		//fmt.Printf("order: %v status: %v\n", (*orders)[0].OrderId, (*orders)[0].Status)
	}
}

func TestGetPageFids(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	nLen := 10
	lastValue := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		ret, err := taskTracker.GetPageFids(&param.FidPageQueryRequest{"create_time", lastValue, -1, 10, nil})
		assert.Equal(t, nil, err)
		rsp := ret.(*param.FidPageQueryResponse)

		first := rsp.FidInfos[0].CreateTime
		assert.Equal(t, true, first <= lastValue)
		assert.Equal(t, true, len(rsp.FidInfos) <= nLen)
		lastValue = rsp.FidInfos[nLen-1].CreateTime
		//rsp := &param.FidPageQueryResponse{}
		bt, err := json.Marshal(ret)
		fmt.Printf("%v\n", string(bt))
	}
}

func TestGetPageOrders(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	nLen := 10
	lastValue := time.Now().UnixMilli()
	for i := 0; i < 10; i++ {
		ret, err := taskTracker.GetPageOrders(&param.OrderPageQueryRequest{1, "create_time", lastValue, -1, 10, nil})
		assert.Equal(t, nil, err)
		rsp := ret.(*param.OrderPageQueryResponse)

		first := rsp.Orders[0].CreateTime
		assert.Equal(t, true, first <= lastValue)
		assert.Equal(t, true, len(rsp.Orders) <= nLen)
		//assert.Equal(t, true,rsp.Orders[0].or)
		lastValue = rsp.Orders[nLen-1].CreateTime
		//rsp := &param.FidPageQueryResponse{}
		bt, err := json.Marshal(ret)
		fmt.Printf("%v\n", string(bt))
	}
}

func TestJsonParse(t *testing.T) {

	req := &api.CreateTaskRequest{
		RequestId: "helo1",
		Ext:       &api.Extend{Ctx: context.Background()},
	}
	bt, _ := json.Marshal(req)

	//rsp := &api.CreateTaskRequest{Ext: &api.Extend{Ctx: context.Background()}}
	rsp := &api.CreateTaskRequest{}

	err := json.Unmarshal(bt, rsp)

	bt2, _ := json.Marshal(rsp)

	fmt.Printf("err :%v, btRep: %v, btrsp: %v \n", err, string(bt), string(bt2))
}

func saveTaskInfo(orderId string, request *api.UploadPieceFidRequest) error {
	if len(request.Pieces) == 0 {
		return errors.New("saveTaskInfo fail, len(request.Pieces) == 0 : order_id:" + orderId)
	}

	count, err := taskInfo.GetRequestTaskCount(request.RequestId)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(fmt.Sprintf("uploadtask saveTaskInfo request id: %v, task have exist", request.RequestId))
	}

	tasks := generateTaskInfo(orderId, request)

	for _, task := range tasks {
		if err := taskInfo.Add(task); err != nil {
			bt, _ := json.Marshal(task)
			logger.Warnf("saveTaskInfo to db fail, orderId: %v, err: %v, task:%v", orderId, err.Error(), string(bt))
			return err
		}
	}

	return nil
}

func generateTaskInfo(orderId string, req *api.UploadPieceFidRequest) []*dict.TaskInfo {
	if req == nil {
		return nil
	}

	tasks := make([]*dict.TaskInfo, 0, len(req.Pieces))
	for _, task := range req.Pieces {
		tasks = append(tasks, &dict.TaskInfo{
			Fid:        task.Fid,
			Cid:        "",
			OrderId:    orderId,
			RequestId:  req.RequestId,
			Rep:        task.Rep,
			VirtualRep: 0,
			RepMin:     task.RepMin,
			RepMax:     task.RepMax,
			Encryption: task.Encryption,
			Expire:     task.Expire,
			Level:      task.Level,
			Size:       task.Size,
			Name:       task.Name,
			Ajust:      task.Ajust,
			Desc:       "",
			Status:     0,
			CreateTime: time.Now().UnixMilli(),
			UpdateTime: time.Now().UnixMilli(),
		})
	}
	return tasks
}
