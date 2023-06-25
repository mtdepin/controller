package main

import (
	"controller/business/dict"
	"controller/business/param"
	"encoding/json"
	"fmt"
	"testing"
)

func TestF1(t *testing.T) {
	uploadTask := &param.UploadTaskRequest{
		RequestId:  "123",
		UserId:     "user123",
		UploadType: param.NoDevice,
		Group:      "chengdu",
		Tasks:      make([]*dict.UploadTask, 0, 2),
		NasList:    []string{},
	}

	task := &dict.UploadTask{
		Fid:        "ox12334354",
		Rep:        3,
		RepMin:     3,
		RepMax:     5,
		Encryption: 0,
		Expire:     100000000000,
		Level:      1,
		Size:       1024,
		Name:       "t1.txt",
		Ajust:      0,
	}

	uploadTask.Tasks = append(uploadTask.Tasks, task)
	if bt, err := json.Marshal(uploadTask); err == nil {
		fmt.Printf("%v\n", string(bt))
	} else {
		fmt.Printf("json marshal err: %v\n", err.Error())
	}
}
