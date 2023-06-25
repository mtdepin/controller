package statemachine

import (
	e "controller/task_tracker/event"
	"fmt"
)

const (
	INIT           = iota + 1
	UPLOAD_PROCEED //上传进行中
	UPLOAD_FINISH  //上传完成
	REP_PROCEED    //备份进行中
	//DELETE_PROCEED  //删除进行中
	REP_FINISH      //备份完成
	CHARGE_PROCEED  //计费完成
	DOWNLOAD_FINISH //下载完成
)

type OrderState struct {
	Status uint64
	Event  *e.Event
}

const (
	SUCCESS = 1
	FAIL    = 0
)

type Key struct {
	orderId string
}

func (p *Key) String() string {
	return fmt.Sprintf("%s", p.orderId)
}
