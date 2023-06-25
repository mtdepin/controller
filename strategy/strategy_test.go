package main

import (
	"controller/pkg/logger"
	utilruntime "controller/pkg/runtime"
	"controller/strategy/algorithm"
	"controller/strategy/config"
	"controller/strategy/database"
	"controller/strategy/dict"
	"controller/strategy/param"
	"controller/strategy/services"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var strategy *services.Service
var taskDb *database.Task
var strategyDb *database.Strategy
var estimate *algorithm.Estimate

func Init() error {
	defer utilruntime.HandleCrash()
	c, err := config.LoadServerConfig(LocalServiceId)
	if err != nil {
		logger.Error(err)
		return err
	}
	config.ServerCfg = c

	if err := database.InitDB(c.DB.Url, c.DB.DbUser, c.DB.DbPassword, c.DB.DbName, c.DB.Timeout); err != nil {
		logger.Errorf("init database: %v fail: %v", c.DB, err.Error())
		return nil
	}

	//init logger
	logger.InitLogger(c.Logger.Level)
	logger.Info("server config:", *c)

	//init api server
	strategy = new(services.Service)
	strategy.Init(database.Db)

	taskDb = new(database.Task)
	taskDb.Init(database.Db)

	strategyDb = new(database.Strategy)
	strategyDb.Init(database.Db)

	estimate = new(algorithm.Estimate)

	return nil
}

//1.构建测试任务信息.
func getOrders() ([]string, error) {
	tasks, err := taskDb.Query(nil, "order_id", 100)
	if err != nil {
		return nil, err
	}
	if 0 == len(*tasks) {
		return nil, errors.New("task is empty")
	}

	orderIds := make([]string, 0, 100)
	orderId := (*tasks)[0].OrderId
	newOrderId := uuid.NewV4().String()
	orderIds = append(orderIds, newOrderId)

	for i, _ := range *tasks {
		if (*tasks)[i].OrderId != orderId {
			newOrderId = uuid.NewV4().String()
			orderId = (*tasks)[i].OrderId
			orderIds = append(orderIds, newOrderId)
		}

		(*tasks)[i].OrderId = newOrderId

		if err := taskDb.Add(&(*tasks)[i]); err != nil {
			return nil, err
		}
	}

	return orderIds, nil
}

func TestCreateStrategy(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	orderIds, err := getOrders()
	if err != nil {
		fmt.Printf("get orders fail, %v", err.Error())
		return
	}

	for i, _ := range orderIds {
		ret, err := strategy.CreateStrategy(&param.CreateStrategyRequest{
			//RequestId: fmt.Sprintf("requestId_%v", i),
			OrderId: orderIds[i],
		})
		assert.Equal(t, nil, err)
		rsp := ret.(*param.CreateStrategyResponse)
		assert.Equal(t, param.SUCCESS, rsp.Status)
		fmt.Printf("create order: %v strategy succ\n", orderIds[i])
	}
}

func TestGetReplicateStrategy(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	strategyInfos, err := strategyDb.Query(nil, "order_id", 10)
	if err != nil {
		fmt.Printf("get order strategy fail, %v\n", err.Error())
		return
	}
	if len(*strategyInfos) == 0 {
		fmt.Printf("get order strategy fail, strategy num is 0 \n")
		return
	}

	for i, _ := range *strategyInfos {
		ret, err := strategy.GetReplicateStrategy(&param.GetStrategyRequset{
			OrderId: (*strategyInfos)[i].OrderId,
		})
		assert.Equal(t, nil, err)
		rsp := ret.(*param.GetStrategyResponse)
		assert.Equal(t, param.SUCCESS, rsp.Status)
		assert.Equal(t, len((*strategyInfos)[i].Tasks), len(rsp.Strategy.Tasks))
		//fmt.Printf("create order: %v strategy succ\n", orderIds[i])
	}
}

func TestGetOrderDeleteStrategy(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	strategyInfos, err := strategyDb.Query(nil, "order_id", 10)
	if err != nil {
		fmt.Printf("get order strategy fail, %v\n", err.Error())
		return
	}
	if len(*strategyInfos) == 0 {
		fmt.Printf("get order strategy fail, strategy num is 0 \n")
		return
	}

	for i, _ := range *strategyInfos {
		ret, err := strategy.GetOrderDeleteStrategy(&param.GetStrategyRequset{
			OrderId: (*strategyInfos)[i].OrderId,
		})
		assert.Equal(t, nil, err)
		rsp := ret.(*param.GetStrategyResponse)
		assert.Equal(t, param.SUCCESS, rsp.Status)
		assert.Equal(t, len((*strategyInfos)[i].Tasks), len(rsp.Strategy.Tasks))
		//fmt.Printf("create order: %v strategy succ\n", orderIds[i])
	}
}

func GetFidDeleteStrategyRequest(strategyInfo *dict.StrategyInfo) *param.GetFidDeleteStrategyRequest {
	nLen := rand.Intn(len((*strategyInfo).Tasks))
	fids := make(map[string]bool)
	for i := 0; i < nLen; i++ {
		fids[(*strategyInfo).Tasks[i].Fid] = true
	}

	return &param.GetFidDeleteStrategyRequest{
		OrderId: strategyInfo.OrderId,
		Fids:    fids,
	}
}

func TestGetFidDeleteStrategy(t *testing.T) {
	if err := Init(); err != nil {
		return
	}

	strategyInfos, err := strategyDb.Query(nil, "order_id", 10)
	if err != nil {
		fmt.Printf("get order strategy fail, %v\n", err.Error())
		return
	}
	if len(*strategyInfos) == 0 {
		fmt.Printf("get order strategy fail, strategy num is 0 \n")
		return
	}

	for i, _ := range *strategyInfos {
		request := GetFidDeleteStrategyRequest(&(*strategyInfos)[i])
		ret, err := strategy.GetFidDeleteStrategy(request)

		assert.Equal(t, nil, err)
		rsp := ret.(*param.GetStrategyResponse)
		assert.Equal(t, param.SUCCESS, rsp.Status)
		assert.Equal(t, len(request.Fids), len(rsp.Strategy.Tasks))
		//fmt.Printf("create order: %v strategy succ\n", orderIds[i])
	}
}

//1.测试计算备份数.

func GetFidReps(num int) map[string]*dict.Rep {
	fidReps := make(map[string]*dict.Rep, num)
	for i := 0; i < num; i++ {
		fidReps[uuid.NewV4().String()] = &dict.Rep{
			MinRep: rand.Intn(5),
			MaxRep: 5 + rand.Intn(10),
		}
	}

	return fidReps
}

func setFidOrders(region, orderId string, mOrders map[string]*dict.Rep, rep *dict.RepInfo, minRep, maxRep, realMinRep, realMaxRep int) {
	if _, ok := mOrders[orderId]; !ok { //不存在，则添加，存在了，不做处理.
		mOrders[orderId] = &dict.Rep{
			Region:     region,
			MinRep:     minRep,
			MaxRep:     maxRep,
			RealRep:    rep.RealRep,
			RealMinRep: realMinRep,
			RealMaxRep: realMaxRep, //注意: 都用minRep,maxRep
			Expire:     rep.Expire,
			Status:     rep.Status,
			CreateTime: time.Now().UnixMilli(),
			UpdateTime: time.Now().UnixMilli(),
		}
	}
}

/*fid   orderid  region  minRep maxRep  realMinRep realMaxRep  Weight  Expire CreateTime
1111    order1    cd      1      2       1           2           1
1111    order2    cd      3      4       4           6           1.5
1111    order3    cd      5      6       7           9           1.5
1111    order4    cd      2      3       2           3           1.5
1111    order5    cd      6      7       6           7           1.5
*/
func TestCalculateRep(t *testing.T) {
	rep := &dict.RepInfo{
		MinRep: 1,
		MaxRep: 2,
	}

	fidReps := make(map[string]*dict.Rep, 10)

	fidReps["orderId_1"] = &dict.Rep{
		MinRep:     1,
		MaxRep:     2,
		RealMinRep: 1,
		RealMaxRep: 2,
	}

	initMinRep, initMaxRep, realMinRep, realMaxRep := estimate.CalculateRep(uuid.NewV4().String(), fidReps, rep)
	assert.Equal(t, rep.MinRep, initMinRep)
	assert.Equal(t, rep.MaxRep, initMaxRep)
	assert.Equal(t, rep.MinRep, realMinRep)
	assert.Equal(t, rep.MaxRep, realMaxRep)

	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_1", fidReps, rep)
	assert.Equal(t, rep.MinRep, initMinRep)
	assert.Equal(t, rep.MaxRep, initMaxRep)
	assert.Equal(t, rep.MinRep, realMinRep)
	assert.Equal(t, rep.MaxRep, realMaxRep)

	rep = &dict.RepInfo{
		MinRep: 3,
		MaxRep: 4,
	}

	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_2", fidReps, rep)
	fmt.Printf("initMinRep: %v, initMaxRep:%v, realMinRep: %v, realMaxRep: %v \n", initMinRep, initMaxRep, realMinRep, realMaxRep)
	assert.Equal(t, 4, rep.MinRep)
	assert.Equal(t, 6, rep.MaxRep)

	assert.Equal(t, 3, initMinRep)
	assert.Equal(t, 4, initMaxRep)
	assert.Equal(t, 4, realMinRep)
	assert.Equal(t, 6, realMaxRep)

	fidReps["orderId_2"] = &dict.Rep{
		MinRep:     3,
		MaxRep:     4,
		RealMinRep: 4,
		RealMaxRep: 6,
	}

	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_2", fidReps, rep)
	fmt.Printf("initMinRep: %v, initMaxRep:%v, realMinRep: %v, realMaxRep: %v \n", initMinRep, initMaxRep, realMinRep, realMaxRep)
	assert.Equal(t, 4, rep.MinRep)
	assert.Equal(t, 6, rep.MaxRep)
	assert.Equal(t, 3, initMinRep)
	assert.Equal(t, 4, initMaxRep)
	assert.Equal(t, 4, realMinRep)
	assert.Equal(t, 6, realMaxRep)

	rep = &dict.RepInfo{
		MinRep: 5,
		MaxRep: 6,
	}

	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_3", fidReps, rep)
	fmt.Printf("initMinRep: %v, initMaxRep:%v, realMinRep: %v, realMaxRep: %v \n", initMinRep, initMaxRep, realMinRep, realMaxRep)
	assert.Equal(t, 7, rep.MinRep)
	assert.Equal(t, 9, rep.MaxRep)

	assert.Equal(t, 5, initMinRep)
	assert.Equal(t, 6, initMaxRep)
	assert.Equal(t, 7, realMinRep)
	assert.Equal(t, 9, realMaxRep)

	fidReps["orderId_3"] = &dict.Rep{
		MinRep:     5,
		MaxRep:     6,
		RealMinRep: 7,
		RealMaxRep: 9,
	}
	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_3", fidReps, rep)
	fmt.Printf("initMinRep: %v, initMaxRep:%v, realMinRep: %v, realMaxRep: %v \n", initMinRep, initMaxRep, realMinRep, realMaxRep)
	assert.Equal(t, 7, rep.MinRep)
	assert.Equal(t, 9, rep.MaxRep)

	assert.Equal(t, 5, initMinRep)
	assert.Equal(t, 6, initMaxRep)
	assert.Equal(t, 7, realMinRep)
	assert.Equal(t, 9, realMaxRep)

	rep = &dict.RepInfo{
		MinRep: 2,
		MaxRep: 3,
	}

	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_4", fidReps, rep)
	fmt.Printf("initMinRep: %v, initMaxRep:%v, realMinRep: %v, realMaxRep: %v \n", initMinRep, initMaxRep, realMinRep, realMaxRep)
	assert.Equal(t, 7, rep.MinRep)
	assert.Equal(t, 9, rep.MaxRep)

	assert.Equal(t, 2, initMinRep)
	assert.Equal(t, 3, initMaxRep)
	assert.Equal(t, 2, realMinRep)
	assert.Equal(t, 3, realMaxRep)

	fidReps["orderId_4"] = &dict.Rep{
		MinRep:     2,
		MaxRep:     3,
		RealMinRep: 2,
		RealMaxRep: 3,
	}

	rep = &dict.RepInfo{
		MinRep: 6,
		MaxRep: 7,
	}

	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_5", fidReps, rep)
	fmt.Printf("initMinRep: %v, initMaxRep:%v, realMinRep: %v, realMaxRep: %v \n", initMinRep, initMaxRep, realMinRep, realMaxRep)
	assert.Equal(t, 7, rep.MinRep)
	assert.Equal(t, 9, rep.MaxRep)

	assert.Equal(t, 6, initMinRep)
	assert.Equal(t, 7, initMaxRep)
	assert.Equal(t, 6, realMinRep)
	assert.Equal(t, 7, realMaxRep)

	fidReps["orderId_5"] = &dict.Rep{
		MinRep:     6,
		MaxRep:     7,
		RealMinRep: 6,
		RealMaxRep: 7,
	}
	initMinRep, initMaxRep, realMinRep, realMaxRep = estimate.CalculateRep("orderId_5", fidReps, rep)
	fmt.Printf("initMinRep: %v, initMaxRep:%v, realMinRep: %v, realMaxRep: %v \n", initMinRep, initMaxRep, realMinRep, realMaxRep)
	assert.Equal(t, 7, rep.MinRep)
	assert.Equal(t, 9, rep.MaxRep)

	assert.Equal(t, 6, initMinRep)
	assert.Equal(t, 7, initMaxRep)
	assert.Equal(t, 6, realMinRep)
	assert.Equal(t, 7, realMaxRep)
}

/*fid   orderid  region  minRep maxRep  realMinRep realMaxRep  Weight  Expire CreateTime
1111    order1    cd      1      2       1           2           1
1111    order2    cd      3      4       4           6           1.5
1111    order3    cd      5      6       7           9           1.5
1111    order4    cd      2      3       2           3           1.5
1111    order5    cd      6      7       9           9           1.5
*/
func TestCalculateDeleteRep(t *testing.T) {
	rep := &dict.RepInfo{
		MinRep: 1,
		MaxRep: 2,
		Status: 0,
	}

	var fidReps = map[string]*dict.Rep{
		"orderId_1": &dict.Rep{MinRep: 1, MaxRep: 2, RealMinRep: 1, RealMaxRep: 2},
		"orderId_2": &dict.Rep{MinRep: 3, MaxRep: 4, RealMinRep: 4, RealMaxRep: 6},
		"orderId_3": &dict.Rep{MinRep: 5, MaxRep: 6, RealMinRep: 7, RealMaxRep: 9},
		"orderId_4": &dict.Rep{MinRep: 2, MaxRep: 3, RealMinRep: 2, RealMaxRep: 3},
		"orderId_5": &dict.Rep{MinRep: 6, MaxRep: 7, RealMinRep: 9, RealMaxRep: 9},
	}

	err := estimate.CalculateDeleteRep("orderId_1", fidReps, rep)
	assert.Equal(t, nil, err)
	assert.Equal(t, dict.TASK_DEL_SUC, rep.Status)
	_, ok := fidReps["orderId_1"]
	assert.Equal(t, false, ok)

	rep = &dict.RepInfo{
		MinRep: 9,
		MaxRep: 9,
		Status: 0,
	}
	err = estimate.CalculateDeleteRep("orderId_5", fidReps, rep)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, rep.Status)
	assert.Equal(t, 7, rep.MinRep)
	assert.Equal(t, 9, rep.MaxRep)

	_, ok = fidReps["orderId_5"]
	assert.Equal(t, false, ok)

	rep = &dict.RepInfo{
		MinRep: 7,
		MaxRep: 9,
		Status: 0,
	}
	err = estimate.CalculateDeleteRep("orderId_3", fidReps, rep)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, rep.Status)
	assert.Equal(t, 4, rep.MinRep)
	assert.Equal(t, 6, rep.MaxRep)

	_, ok = fidReps["orderId_3"]
	assert.Equal(t, false, ok)

	rep = &dict.RepInfo{
		MinRep: 4,
		MaxRep: 6,
		Status: 0,
	}
	err = estimate.CalculateDeleteRep("orderId_2", fidReps, rep)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, rep.Status)
	assert.Equal(t, 2, rep.MinRep)
	assert.Equal(t, 3, rep.MaxRep)

	_, ok = fidReps["orderId_2"]
	assert.Equal(t, false, ok)

	rep = &dict.RepInfo{
		MinRep: 2,
		MaxRep: 3,
		Status: 0,
	}
	err = estimate.CalculateDeleteRep("orderId_4", fidReps, rep)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, rep.Status)
	assert.Equal(t, 0, rep.MinRep)
	assert.Equal(t, 0, rep.MaxRep)

	_, ok = fidReps["orderId_4"]
	assert.Equal(t, false, ok)

	rep = &dict.RepInfo{
		MinRep: 1,
		MaxRep: 2,
		Status: 0,
	}
	err = estimate.CalculateDeleteRep("orderId_1", fidReps, rep)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, rep.Status)
	assert.Equal(t, 0, rep.MinRep)
	assert.Equal(t, 0, rep.MaxRep)

	_, ok = fidReps["orderId_1"]
	assert.Equal(t, false, ok)
}
