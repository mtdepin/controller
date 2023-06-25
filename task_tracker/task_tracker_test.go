package main

import (
	"context"
	"controller/api"
	"controller/pkg/cache"
	ctl "controller/pkg/http"
	"controller/pkg/logger"
	"controller/pkg/montior"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
	//reflect.AppendSlice()

	strJson := `{
    "order_id": "1670469465662018700", 
    "tasks": [
        {
            "fid": "fid2", 
            "cid": "QmQswPkNZwcw2fZjJncbjVjbUJQsUDpyCxFz72NyT69YLu", 
            "region": "chengdu", 
            "origins": "http://192.168.2.35:6001/api/v0/add", 
            "status": 1
        }, 
        {
            "fid": "fid3", 
            "cid": "QmZDiLATmRogoYeXUiMbfRfuRNP4cGvRe1gwDPqLkvYH1D", 
            "region": "chengdu", 
            "origins": "http://192.192.168.2.35:6001/api/v0/add", 
            "status": 1
        }, 
        {
            "fid": "fid3", 
            "cid": "QmZDiLATmRogoYeXUiMbfRfuRNP4cGvRe1gwDPqLkvYH1D", 
            "region": "chengdu", 
            "origins": "http://192.168.2.35:6001/api/v0/add", 
            "status": 1
        }, 
        {
            "fid": "fid2", 
            "cid": "QmQswPkNZwcw2fZjJncbjVjbUJQsUDpyCxFz72NyT69YLu", 
            "region": "chengdu", 
            "origins": "http://192.168.2.35:6001/api/v0/add", 
            "status": 1
        }
    ], 
    "status": 1
	}`

	request := &param.UploadFinishRequest{}

	err := json.Unmarshal([]byte(strJson), request)
	if err != nil {
		fmt.Printf(" json unmarshal fail, err: %v \n", err.Error())
	} else {
		fmt.Printf("request: %v\n", request)
	}

	return
	uploadTask := &param.CallbackUploadRequest{
		OrderId: "1669785814365258200",
		Fid:     "fid234567",
		Cid:     "QmdNY16tvjbTHwCteZ4VbX2rFrRgMKxZVZuXmifWq3tiLC",
		Region:  "chengdu",
		Origins: "chengdu",
		Status:  1,
	}

	if bt, err := json.Marshal(uploadTask); err == nil {
		fmt.Printf("%v\n", string(bt))
	} else {
		fmt.Printf("json marshal err: %v\n", err.Error())
	}
}

func TestGenergateFile(t *testing.T) {
	/*if err := cbg.WriteTupleEncodersToFile("cbor_gen.go", "statemachine",
		OrderState{},
	); err != nil {
		panic(err)
	}

	if err := cbg.WriteMapEncodersToFile("/cbor_map_gen.go", "database",
		types.Event{},
	); err != nil {
		panic(err)
	}*/
}

type Data struct {
	name string
}

var mIdx = make(map[int]*Data, 500)

func TestIndex(t *testing.T) {
	for i := 0; i < 100; i++ {
		go Set(i, &Data{"a"})
		//go Delete("1")
		//go Set("2", &Data{"b"})
		//go Set("3", &Data{"c"})
		//go Set("4", &Data{"b"})

	}
	time.Sleep(1 * time.Second)
	fmt.Printf("helo %v \n", mIdx)
}

func Set(key int, data *Data) {
	mIdx[key] = data
}

func Delete(key int) {
	delete(mIdx, key)
}

func TestUid(t *testing.T) {
	n := 100000000
	uidMap := make(map[string]bool, n)
	for i := 0; i < n; i++ {
		uidMap[uuid.NewV4().String()] = true
	}
	fmt.Printf(" size uid map, %d \n ", len(uidMap))
}

func TestMap(t *testing.T) {
	fmt.Printf("run begin \n")
	n := 200000
	time.Sleep(10 * time.Second)
	szUploadTasks := make([]*param.CallbackUploadRequest, n, n)
	mapUploadTasks := make(map[string]int, n)
	for i := 0; i < n; i++ {
		szUploadTasks[i] = &param.CallbackUploadRequest{
			OrderId: "1669785814365258200",
			Fid:     "fid234567",
			Cid:     "QmdNY16tvjbTHwCteZ4VbX2rFrRgMKxZVZuXmifWq3tiLC",
			Region:  "chengdu",
			Origins: "chengdu",
			Status:  1,
		}
		mapUploadTasks[fmt.Sprintf("1669785814365258200%d", i)] = i
	}
	time.Sleep(10 * time.Second)
	fmt.Printf("run finish \n")
}

func TestMap1(t *testing.T) {
	fmt.Printf("run begin \n")
	n := 200000
	time.Sleep(10 * time.Second)
	mapUploadTasks := make(map[string]*param.CallbackUploadRequest, n)
	for i := 0; i < n; i++ {
		mapUploadTasks[fmt.Sprintf("1669785814365258200%d", i)] = &param.CallbackUploadRequest{
			OrderId: "1669785814365258200",
			Fid:     "fid234567",
			Cid:     "QmdNY16tvjbTHwCteZ4VbX2rFrRgMKxZVZuXmifWq3tiLC",
			Region:  "chengdu",
			Origins: "chengdu",
			Status:  1,
		}
	}
	time.Sleep(10 * time.Second)
	fmt.Printf("run finish \n")
}

func TestSize(t *testing.T) {
	fmt.Printf("maxsize : %d \n", (8<<20)/2)
}

func Handler(sz []int, n int) {
	//time.Sleep(10 * time.Minute)
	for i := 0; i < n; i++ {
		sz[i] = i
	}
	//time.Sleep(10 * time.Minute)
}

func TestIOInter(t *testing.T) {
	n := 10000
	for i := 0; i < n; i++ {
		sz := make([]int, n, n)
		go Handler(sz, n)
	}
	fmt.Printf("helo \n")
	//time.Sleep(10 * time.Hour)
}

func TestCache(t *testing.T) {
	cache := new(cache.Cache)
	cache.InitCache(2)
	for i := 0; i < 10; i++ {
		cache.Add(i)
		cache.Print()
	}
	cache.Add(5)
	cache.Print()
	cache.Add(6)
	cache.Print()
	cache.Add(7)
	cache.Print()
	cache.Add(6)
	cache.Print()
	cache.Add(6)
	cache.Print()

	cache.Delete(8)
	cache.Print()
	cache.Delete(9)
	cache.Print()

	cache.Delete(7)
	cache.Print()

	cache.Delete(6)
	cache.Print()

	cache.Delete(5)
	cache.Print()
	for i := 0; i < 10; i++ {
		cache.Add(i)
		cache.Print()
	}
	//nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/searchRep"
	//cache.Add(7)
	//cache.Add()
}

func TestDoRequest(t *testing.T) {
	//str := `{"order_id":"ffe5b5aa-5678-4ea7-b7f3-8cdc27a6f6b8dfdsfdsdfdf","tasks":[{"fid":"177","cid":"Qmb3wsY43c8NmWLwiBY7a8zMP57U8uxG2EZscd75a74MjE","regions":["chengdu"]},{"fid":"178","cid":"Qmd7z6DMPGMwomU9MNbE3H7BdwCRyvDDpZjaMyUSfmcut6","regions":["chengdu"]},{"fid":"161","cid":"QmSDKo1bwWLXxvoXFCHr74FyMdwKsXRVZTkaCB2254WYtp","regions":["chengdu"]},{"fid":"162","cid":"QmUohPbYJHf4TXtbdLzqMZpmgaNjbqTQLnMSHNN6JV8Ae8","regions":["chengdu"]},{"fid":"167","cid":"QmcPfR6LfJAtxbtBVdNCZx2ppPY6dsBQSY78etqMEV6ngd","regions":["chengdu"]}]}`

	//http://192.168.2.98:8612/scheduler/v1/searchRep
	//ctl.DoRequestNew(http.MethodGet, "http://192.168.2.99:8612/scheduler/v1/searchRep1", nil, []byte(str))
	//ctl.DoRequestNew(http.MethodGet, "http://192.168.1.1:8080/v1/searchRep", nil, []byte("123456"))

}

type ChargeRequest1 struct {
	UserId    string   `json:"user_id"`
	OrderId   string   `json:"order_id"`
	OrderType int      `json:"order_type"`
	Tasks     []*Task1 `json:"tasks"`
}

type Task1 struct {
	Fid     string               `bson:"fid,omitempty" json:"fid"`
	Cid     string               `bson:"cid,omitempty" json:"cid"`
	Region  string               `bson:"region,omitempty" json:"region"`   //文件上传区域
	Origins string               `bson:"origins,omitempty" json:"origins"` //文件上传节点
	Reps    map[string]*RepInfo1 `bson:"reps,omitempty" json:"reps"`       //key region: val: 备份详情
	Status  int                  `bson:"status,omitempty" json:"status"`
	Size    int                  `bson:"size,omitempty" json:"size"` // 文件大小，单位字节
}

type RepInfo1 struct {
	Region     string `json:"region"`
	VirtualRep int    `json:"virtual_rep"`
	RealRep    int    `json:"real_rep"`
	MinRep     int    `json:"min_rep"`
	MaxRep     int    `json:"max_rep"`
	Expire     uint64 `json:"expire"`
	Encryption int    `json:"encryption"`
	Status     int    `json:"status"`
}

func TestCharge(t *testing.T) {
	str := `{"order_id":"203b0c42-5165-4c7f-b7d1-ecf32af688ef","order_type":1,"tasks":[{"fid":"109","cid":"QmQPigK1c5MG2P6RCQPEpe3V93uVZUKcr9aT7Pdud55XQs","region":"chengdu","origins":"/ip4/154.53.61.17/tcp/4001/p2p/12D3KooWL4SrmZJSN6myX84aHZ6DWj1ApDUEcjSgG53qVuVXnY3W/p2p-circuit/p2p/12D3KooWLvPtoiRf28KgbxFsjHbtFUSsfNhszwYr4mv65FsxJj6t","reps":{"chengdu":{"region":"chengdu","virtual_rep":0,"real_rep":4,"min_rep":2,"max_rep":4,"expire":100000,"encryption":0,"status":3}},"status":3},{"fid":"105","cid":"QmPGahgYdcbyKfQixQor7YWjn6CdYMmwRjBQXW3VDDszPA","region":"chengdu","origins":"/ip4/154.53.61.17/tcp/4001/p2p/12D3KooWL4SrmZJSN6myX84aHZ6DWj1ApDUEcjSgG53qVuVXnY3W/p2p-circuit/p2p/12D3KooWLvPtoiRf28KgbxFsjHbtFUSsfNhszwYr4mv65FsxJj6t","reps":{"chengdu":{"region":"chengdu","virtual_rep":0,"real_rep":4,"min_rep":2,"max_rep":4,"expire":100000,"encryption":0,"status":3}},"status":3},{"fid":"106","cid":"Qmdf6cA7Gv4FNtzwcFy6VprqG6EHj2joKkMdD4K51UtGVZ","region":"chengdu","origins":"/ip4/154.53.61.17/tcp/4001/p2p/12D3KooWL4SrmZJSN6myX84aHZ6DWj1ApDUEcjSgG53qVuVXnY3W/p2p-circuit/p2p/12D3KooWLvPtoiRf28KgbxFsjHbtFUSsfNhszwYr4mv65FsxJj6t","reps":{"chengdu":{"region":"chengdu","virtual_rep":0,"real_rep":4,"min_rep":2,"max_rep":4,"expire":100000,"encryption":0,"status":3}},"status":3},{"fid":"107","cid":"QmWVyYkYGKiDZLAxoaaeFSHehRd8U33XCEYQxgDJ6YxuvA","region":"chengdu","origins":"/ip4/154.53.61.17/tcp/4001/p2p/12D3KooWL4SrmZJSN6myX84aHZ6DWj1ApDUEcjSgG53qVuVXnY3W/p2p-circuit/p2p/12D3KooWLvPtoiRf28KgbxFsjHbtFUSsfNhszwYr4mv65FsxJj6t","reps":{"chengdu":{"region":"chengdu","virtual_rep":0,"real_rep":4,"min_rep":2,"max_rep":4,"expire":100000,"encryption":0,"status":3}},"status":3},{"fid":"108","cid":"QmSQyqsmymDptzsitk7jiuaWdkkc8C6dUUKtrTYC5NFbvb","region":"chengdu","origins":"/ip4/154.53.61.17/tcp/4001/p2p/12D3KooWL4SrmZJSN6myX84aHZ6DWj1ApDUEcjSgG53qVuVXnY3W/p2p-circuit/p2p/12D3KooWLvPtoiRf28KgbxFsjHbtFUSsfNhszwYr4mv65FsxJj6t","reps":{"chengdu":{"region":"chengdu","virtual_rep":0,"real_rep":4,"min_rep":2,"max_rep":4,"expire":100000,"encryption":0,"status":3}},"status":3}]}
`
	request := &ChargeRequest1{}
	err := json.Unmarshal([]byte(str), request)
	if err != nil {
		fmt.Printf(" err: %v \n", err.Error())
	}

	//http://192.168.2.98:8612/scheduler/v1/searchRep
	ctl.DoRequestNew(http.MethodGet, "http://192.168.2.99:8612/scheduler/v1/searchRep1", nil, []byte(str))
	//ctl.DoRequestNew(http.MethodGet, "http://192.168.1.1:8080/v1/searchRep", nil, []byte("123456"))
	/*for p := range s.providerFinder.FindProvidersAsync(ctx, c) {
		// When a provider indicates that it has a cid, it's equivalent to
		// the providing peer sending a HAVE
		s.sws.Update(p, nil, []cid.Cid{c}, nil)
	}*/
}

func setChan(ch chan int) <-chan int {
	go func(c chan int) {
		defer close(ch)
		c <- 1
		time.Sleep(10 * time.Second)
		//close(ch)
	}(ch)
	return ch
}

func TestSpec(t *testing.T) {
	ch := make(chan int)

	for p := range setChan(ch) {
		fmt.Printf("helo p: %v\n", p)
		//time.Sleep(1 * time.Second)
	}
	return
	for p := range setChan(ch) {
		fmt.Printf("helo p: %v\n", p)
	}

	p := <-setChan(ch)
	fmt.Printf("helo p: %v\n", p)

	fmt.Printf("----------end---: \n")
}

type student struct {
	name string
	id   int
}

func TestCp(t *testing.T) {
	ms := make(map[string]*student, 2)
	ms2 := make(map[string]*student, 2)
	ms["1"] = &student{
		"zs",
		1,
	}

	ms["2"] = &student{
		"ls",
		2,
	}

	val, _ := ms["1"]

	s := *val
	//ms2["1"] = &s

	ms2["1"] = val

	fmt.Printf("s :%v \n", s)
	s.name = "ww"
	s.id = 2
	fmt.Printf("ms1 :%v \n", ms["1"])
	fmt.Printf("s :%v \n", s)
	//fmt.Printf("ms2 :%v \n", ms2["1"])
	//ms["1"].name = "ms1 name"
	ms["1"] = nil
	fmt.Printf("s :%v \n", s)
	fmt.Printf("ms1 :%v \n", ms["1"])
	fmt.Printf("ms2 :%v \n", ms2["1"])
}

func UpdateMap(st map[string]*student) {
	st["1"] = &student{"zhangs", 1}

}

func TestMap2(t *testing.T) {
	st := make(map[string]*student, 1)
	UpdateMap(st)

	fmt.Printf("m1: %v \n", st["1"])
}

func TestChan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan int, 2)

	go f1(ctx, ch)
	ch <- 1
	time.Sleep(1 * time.Second)
	cancel()
	cancel()
	fmt.Printf("cancel\n")

	time.Sleep(10 * time.Second)

}

//10.
func f1(ctx context.Context, ch chan int) {
	for {
		select {
		case a := <-ch:
			fmt.Printf("a :%d\n", a)
			time.Sleep(1 * time.Second)
		case <-ctx.Done():
			fmt.Printf("exit")
			return
		}
	}
}

func TestLog(t *testing.T) {
	ding := new(montior.DingTalk)
	ding.Init("task_tracker", "b9319d14e6795244645968ae4a83f0903518af3836c48f6d0ce9460bc815bc16", "SEC319c35ae4bdc94d2c54752aeba9d3a8476786c265d83f348249b406558284df6")
	logger.InitLoggerWithDingTalk("warn", ding)
	utils.Log(utils.ERROR, "servername", "errinfo", "event")
}

func TestSplit(t *testing.T) {
	sz := strings.Split("helo123", "")
	if len(sz) < 1 {
		fmt.Printf(" sz < 1\n")
	}

	orderId := sz[0]
	fmt.Printf(" orderId: %v\n", orderId)
}

var fidEventChan chan *e.FidEvent

func InitThreadPool(size int) {
	fidEventChan = make(chan *e.FidEvent, size)

	for i := 0; i < size; i++ {
		go handler()
	}
}

func Proc(num, nthread int, orderId string) {
	var err error
	rets := make([]chan *e.FidRet, num, num)

	count := 0
	div := num / nthread

	t1 := time.Now().UnixMilli()

	for j := 0; j < div; j++ {
		for i := j * nthread; i < (j+1)*nthread; i++ {
			ret := make(chan *e.FidRet)
			fidEventChan <- &e.FidEvent{
				Fid:   fmt.Sprintf("fid_%v", i),
				Group: "chengdu",
				Ret:   ret,
			}
			//fmt.Printf("helo : %v\n", i)
			rets[i] = ret
		}

		//index := 0

		for i := j * nthread; i < (j+1)*nthread; i++ {
			fidRet := <-rets[i]

			if fidRet.Err != nil {
				err = fidRet.Err
			}

			count++
		}
	}

	//mod
	for i := div * nthread; i < num; i++ {
		ret := make(chan *e.FidRet)
		fidEventChan <- &e.FidEvent{
			Fid:   fmt.Sprintf("fid_%v", i),
			Group: "chengdu",
			Ret:   ret,
		}
		//fmt.Printf("helo : %v\n", i)
		rets[i] = ret
	}

	//index := 0
	for i := div * nthread; i < num; i++ {
		fidRet := <-rets[i]

		if fidRet.Err != nil {
			err = fidRet.Err
		}
		count++
	}

	if err != nil {
		fmt.Printf(" proc fail : %v\n", err.Error())
	}

	t2 := time.Now().UnixMilli()
	logger.Infof("-------------------helo orderId: %v generateUploadOrderEvent task end, cost_time: %v, count: %v", orderId, t2-t1, count)
}

var g_index int

func handler() {
	for {
		g_index++
		orderEvent := <-fidEventChan
		fidInfo, err := getValue(orderEvent)

		//fmt.Printf(" --------getvalue: %v \n", g_index)
		orderEvent.Ret <- &e.FidRet{FidInfo: fidInfo, Err: err}
	}
}

func getValue(event *e.FidEvent) (*e.FidInfo, error) {
	//time.Sleep(1 * time.Millisecond)
	return &e.FidInfo{
		Fid:     event.Fid,
		Cid:     "",
		Repeate: 0,
		Origins: "",
		Status:  1,
	}, nil

}

func Proc2(num, nthread int, orderId string, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	rets := make([]chan *e.FidRet, 0, num)

	count := 0

	total := 0
	t1 := time.Now().UnixMilli()

	for i := 0; i < num; i++ {
		ret := make(chan *e.FidRet)

		SetFidEvent(fidEventChan, ret)

		rets = append(rets, ret)

		total++
		count++

		if count >= nthread-1 { //等待接收完成。
			for _, ret := range rets {
				fidRet := <-ret
				err = fidRet.Err
			}

			count = 0
			rets = rets[0:0] //清空数组

			if err != nil {
				fmt.Printf(" return err: %v \n", err.Error())
			}
		}
	}

	for _, ret := range rets {
		fidRet := <-ret
		err = fidRet.Err
	}

	if err != nil {
		fmt.Printf(" return err: %v \n", err.Error())
	}

	t2 := time.Now().UnixMilli()
	logger.Infof("-------------------helo orderId: %v generateUploadOrderEvent task end, cost_time: %v, count: %v", orderId, t2-t1, total)
}

func SetFidEvent(fidEventChan chan *e.FidEvent, ret chan *e.FidRet) {
	fidEventChan <- &e.FidEvent{
		Fid:   fmt.Sprintf("fid_%v", 0),
		Group: "chengdu",
		Ret:   ret,
	}
}

func TestThreadPool(t *testing.T) {
	//runtime.GOMAXPROCS(30 * runtime.NumCPU())
	nThread := 20
	InitThreadPool(nThread)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go Proc2(871, nThread, fmt.Sprintf("orderId_%v", i), &wg)
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	//ctx, _ := context.WithCancel(context.Background())
	req := &api.UploadTaskRequest{
		Ext: &api.Extend{Ctx: ctx},
	}

	rsp := &api.UploadTaskRequest{
		Ext: &api.Extend{Ctx: ctx},
	}
	bt, _ := json.Marshal(req)

	json.Unmarshal(bt, rsp)

	bt1, _ := json.Marshal(rsp)

	fmt.Printf("helo: ctx: %v, \n req: %v,\n rsp: %v\n", ctx, string(bt), string(bt1))
}

func TestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("timeout return\n")
			return
		}
	}
}
