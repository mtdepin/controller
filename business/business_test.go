package main

import (
	"context"
	"controller/api"
	"controller/business/config"
	"controller/business/database"
	"controller/business/dict"
	"controller/business/param"
	"controller/business/processor"
	"controller/business/services"
	"controller/pkg/logger"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"testing"
)

var order *processor.Order
var business *services.Service

func Init() error {
	c, err := config.LoadServerConfig(LocalServiceId)
	if err != nil {
		logger.Error(err)
		return err
	}
	config.ServerCfg = c

	logger.Info("config Info: %v", *c)

	if err := database.InitDB(c.DB.Url, c.DB.DbUser, c.DB.DbPassword, c.DB.DbName, c.DB.Timeout); err != nil {
		logger.Errorf("init database: %v fail: %v", c.DB, err.Error())
		return nil
	}

	//init logger
	logger.InitLogger(c.Logger.Level)

	return nil
}

func GetUploadTask(nTaskNum int) *api.UploadTaskRequest {
	uploadTask := &api.UploadTaskRequest{
		RequestId:  uuid.NewV4().String(),
		UserId:     "46",
		UploadType: param.NoDevice,
		Group:      "chengdu",
		PieceNum:   nTaskNum,
		Ext:        &api.Extend{Ctx: context.Background()},
	}

	return uploadTask
}

func getUploadPieceFidRequest(orderId string, nTaskNum int) *api.UploadPieceFidRequest {
	uploadTask := &api.UploadPieceFidRequest{
		RequestId: uuid.NewV4().String(),
		Group:     "chengdu",
		OrderId:   orderId,
		Pieces:    make([]*api.PieceFid, 0, nTaskNum),
		Ext:       &api.Extend{Ctx: context.Background()},
	}

	for i := 0; i < nTaskNum; i++ {
		task := &api.PieceFid{
			Fid:        fmt.Sprintf("fid_%v_%v", rand.Intn(100), i),
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

		uploadTask.Pieces = append(uploadTask.Pieces, task)
	}

	return uploadTask
}

func TestCreateUploadTask(t *testing.T) {
	if err := Init(); err != nil {
		fmt.Printf("init business fail")
		return
	}

	order = new(processor.Order)
	order.Init(database.Db)

	patch := gomonkey.ApplyPrivateMethod(reflect.TypeOf(order), "getNodeList", func(_ *processor.Order) ([]*dict.Node, error) {
		nodes := make([]*dict.Node, 0, 1)
		nodes = append(nodes, &dict.Node{})

		return nodes, nil
	})

	patch = gomonkey.ApplyPrivateMethod(reflect.TypeOf(order), "CheckAccountBalance", func(_ *processor.Order) (bool, error) {
		return true, nil
	})

	patch = gomonkey.ApplyPrivateMethod(reflect.TypeOf(order), "selectUploadRegion", func(_ *processor.Order) (string, error) {
		return "home", nil
	})

	//selectUploadRegion

	defer patch.Reset()

	//init api server
	business = new(services.Service)
	business.InitService(order, nil)

	num := 50
	for i := 0; i < num; i++ {
		ret, err := business.UploadTask(GetUploadTask(i))

		if i == 0 {
			assert.NotEqual(t, nil, err)
			continue
		}

		assert.Equal(t, nil, err)
		uploadTaskResponse := ret.(*api.UploadTaskResponse)
		assert.Equal(t, param.SUCCESS, uploadTaskResponse.Status)

		rsp, err := business.UploadPieceFid(getUploadPieceFidRequest(uploadTaskResponse.OrderId, i))

		assert.Equal(t, nil, err)
		uploadPieceFidResponse := rsp.(*api.UploadPieceFidResponse)
		assert.Equal(t, param.SUCCESS, uploadPieceFidResponse.Status)
		assert.Equal(t, uploadTaskResponse.OrderId, uploadPieceFidResponse.OrderId)
	}
}

func GetDownloadTask(ntask int) *param.DownloadTaskRequest {
	downloadTask := &param.DownloadTaskRequest{
		RequestId:    uuid.NewV4().String(),
		UserId:       "46",
		Group:        "chengdu",
		DownloadType: param.NoDevice,
		Tasks:        make([]*dict.DownloadTask, 0, ntask),
		Ext:          &param.Extend{Ctx: context.Background()},
	}

	for i := 0; i < ntask; i++ {
		task := &dict.DownloadTask{
			Cid:  uuid.NewV4().String(),
			Name: fmt.Sprintf("file_%v", i),
			Size: rand.Intn(100000),
		}
		downloadTask.Tasks = append(downloadTask.Tasks, task)
	}

	return downloadTask
}

func TestCreateDownloadTask(t *testing.T) {
	if err := Init(); err != nil {
		fmt.Printf("init business fail")
		return
	}

	order = new(processor.Order)
	order.Init(database.Db)

	patch := gomonkey.ApplyPrivateMethod(reflect.TypeOf(order), "getDownloadNodesFromRM", func(_ *processor.Order) (*param.GetDownloadNodeResponse, error) {
		return &param.GetDownloadNodeResponse{
			Data: make(map[string][]*dict.Node),
		}, nil
	})
	defer patch.Reset()

	//init api server
	business = new(services.Service)
	business.InitService(order, nil)

	for i := 0; i < 100; i++ {
		ret, err := business.DownloadTask(GetDownloadTask(i))

		if i == 0 {
			assert.NotEqual(t, nil, err)
			continue
		}
		assert.Equal(t, nil, err)
		uploadTaskResponse := ret.(*param.DownloadTaskResponse)
		assert.Equal(t, param.SUCCESS, uploadTaskResponse.Status)
	}
}

func GetDeleteFid(orderId string, fids map[string]bool) *param.DeleteFidRequest {
	return &param.DeleteFidRequest{
		RequestId: uuid.NewV4().String(),
		UserId:    "46",
		OrderId:   orderId,
		Fids:      fids,
	}
}

/*func TestDeleteFid(t *testing.T) {
	if err := Init(); err != nil {
		fmt.Printf("init business fail")
		return
	}

	order = new(processor.Order)
	order.Init(database.Db)

	//init api server
	business = new(services.Service)
	business.InitService(order, nil)

	//fmt.Sprintf("orderId", %v)

	//for i := 0; i < 100; i++ {
	fids := make(map[string]bool)
	fids["QmPiqbwtNq37cCaVxT4mw7a3AnMJadk8ZWHJsimUDZkwCV"] = true
	ret, err := business.DeleteFid(GetDeleteFid("1a82e0f5-b9f3-4e7c-bcba-7659d8d16771", fids))

	assert.Equal(t, nil, err)
	uploadTaskResponse := ret.(*param.DeleteFidResponse)
	assert.Equal(t, param.SUCCESS, uploadTaskResponse.Status)
	//}
}*/
