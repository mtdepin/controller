package main

import (
	"controller/pkg/cache"
	ctl "controller/pkg/http"
	"controller/task_tracker/param"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
	"time"
)

func TestJson(t *testing.T) {
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
	str := `{"order_id":"ffe5b5aa-5678-4ea7-b7f3-8cdc27a6f6b8dfdsfdsdfdf","tasks":[{"fid":"177","cid":"Qmb3wsY43c8NmWLwiBY7a8zMP57U8uxG2EZscd75a74MjE","regions":["chengdu"]},{"fid":"178","cid":"Qmd7z6DMPGMwomU9MNbE3H7BdwCRyvDDpZjaMyUSfmcut6","regions":["chengdu"]},{"fid":"161","cid":"QmSDKo1bwWLXxvoXFCHr74FyMdwKsXRVZTkaCB2254WYtp","regions":["chengdu"]},{"fid":"162","cid":"QmUohPbYJHf4TXtbdLzqMZpmgaNjbqTQLnMSHNN6JV8Ae8","regions":["chengdu"]},{"fid":"167","cid":"QmcPfR6LfJAtxbtBVdNCZx2ppPY6dsBQSY78etqMEV6ngd","regions":["chengdu"]}]}`

	//http://192.168.2.98:8612/scheduler/v1/searchRep
	ctl.DoRequestNew(http.MethodGet, "http://192.168.2.99:8612/scheduler/v1/searchRep1", nil, []byte(str))
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

func TestChan(t *testing.T) {
	ch := make(chan int, 10)
	fmt.Printf("ch size : %v \n", len(ch))
	ch <- 1
	fmt.Printf("ch size : %v \n", len(ch))
}
