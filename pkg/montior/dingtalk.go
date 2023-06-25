package montior

import (
	"bytes"
	ctl "controller/pkg/http"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strconv"
	"time"
)

// 钉钉机器人接入官方文档：
// https://open.dingtalk.com/document/robots/custom-robot-access#section-e4x-4y8-9k0

type DingTalk struct {
	ServiceName string      // 服务名称
	AccessToken string      // 机器人access token
	Secret      string      // 机器人加签秘钥
	EnableAt    bool        // 指定@群成员
	AtAll       bool        // @所有人
	BaseURL     string      // 机器人webhook地址
	pipeLine    chan *Event // 机器人限流消息队列
	queue       chan *Event // 机器人待发送消息队列
}

const (
	TEXT = iota + 1
	LINK
	MARKDOWN
)

type Event struct {
	MsgType int32
	Data    interface{}
}

type Message struct {
	At      AtPerson `json:"at"`
	Text    Text     `json:"text"`
	Msgtype string   `json:"msgtype"`
}

type Text struct {
	Content string `json:"content"`
}

type AtPerson struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}

type Rsp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// 初始化
func (p *DingTalk) Init(serviceName, accessToken, secret string) {
	p.ServiceName = serviceName
	p.BaseURL = "https://oapi.dingtalk.com/robot/send"
	p.AccessToken = accessToken
	if secret != "" {
		p.Secret = secret
	}

	p.pipeLine = make(chan *Event, 20)
	p.queue = make(chan *Event, 1000)
	go p.Handle()
	go p.PopMsgEventQueue()
}

func (p *DingTalk) PushMsgEventQueue(logLevel zapcore.Level, msgType int32, template string, args ...interface{}) {
	request := new(Message)

	switch msgType {
	case TEXT:
		request.Msgtype = "text"
		request.Text = Text{Content: p.buildMsgContent(p.getLevel(logLevel), template, args)}
	case LINK:
		return
	case MARKDOWN:
		return
	default:
		return
	}

	e := &Event{MsgType: msgType, Data: request}
	p.queue <- e
}

func (p *DingTalk) Handle() {
	for {
		e := <-p.pipeLine
		p.SendMessage(e)
	}
}

func (p *DingTalk) PopMsgEventQueue() {
	for {
		e := <-p.queue
		p.pipeLine <- e
	}
}

// 发送消息
func (p *DingTalk) SendMessage(e *Event) {
	bt, err := json.Marshal(e.Data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	queryParams := make(map[string]string, 0)
	queryParams["access_token"] = p.AccessToken
	if p.Secret != "" {
		queryParams["timestamp"], queryParams["sign"] = p.genSign()
	}

	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"

	ret, err1 := ctl.DoRequest2(http.MethodPost, p.BaseURL, queryParams, headers, bytes.NewReader(bt))
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	rsp := new(Rsp)
	if err := json.Unmarshal(ret, rsp); err != nil {
		fmt.Println(err.Error())
		return
	}
	if rsp.ErrCode != 0 {
		fmt.Printf("request err: %+v", rsp)
	}
}

func (p *DingTalk) buildLogContent(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func (p *DingTalk) buildMsgContent(logLevel, template string, fmtArgs []interface{}) string {
	content := orderedmap.New()
	content.Set("发生时间", time.Now().Format("2006-01-02 15:04:05"))
	content.Set("日志等级", logLevel)
	content.Set("服务名称", p.ServiceName)
	content.Set("日志内容", p.buildLogContent(template, fmtArgs))

	ret := ""
	keys := content.Keys()
	for _, k := range keys {
		v, _ := content.Get(k)
		ret += fmt.Sprintf("%s:%s\n", k, v)
	}
	return ret
}

func (p *DingTalk) genSign() (ts, sign string) {
	ts = strconv.FormatInt(time.Now().UnixMilli(), 10)
	stringToSign := fmt.Sprintf("%s\n%s", ts, p.Secret)
	sign = HmacSha256(stringToSign, p.Secret)
	return
}

func (p *DingTalk) getLevel(level zapcore.Level) string {
	switch level {
	case zap.InfoLevel:
		return "info"
	case zap.DebugLevel:
		return "debug"
	case zap.ErrorLevel:
		return "ctlerror"
	case zap.FatalLevel:
		return "fatal"
	case zap.DPanicLevel:
		return "panic"
	default:
		return "info"
	}
}

func HmacSha256(source, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(source))
	signedBytes := h.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}
