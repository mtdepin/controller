package processor

import (
	"controller/pkg/logger"
	"encoding/json"
)

func log(level int, name, errInfo string, event interface{}) {
	bt, _ := json.Marshal(event)
	logger.Warnf("%s fail: %v | event: %v", name, errInfo, string(bt))
}
