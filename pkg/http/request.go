package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	COUNT         = 3
	TIME_INTERNAL = 100
)

var HttpStateCode = map[int]bool{403: true, 408: true, 409: true}

func DoRequestNew(method, url string, args map[string]string, bt []byte) (rsp []byte, err error) {
	stateCode := 0
	var body *bytes.Reader
	ctx := context.Background()
	for i := 0; i < COUNT; i++ {
		if bt != nil {
			body = bytes.NewReader(bt)
		}

		if rsp, stateCode, err = doRequest(ctx, method, url, args, body); err == nil {
			return
		}

		if _, ok := HttpStateCode[stateCode]; !ok {
			return
		}

		time.Sleep(TIME_INTERNAL * time.Millisecond)
	}

	if err != nil {
		err = errors.New(fmt.Sprintf("DoRequest Fail: %v, url:%v", err.Error(), url))
	}
	return
}

func DoRequest(ctx context.Context, method, url string, args map[string]string, body io.Reader) (rsp []byte, err error) {
	stateCode := 0

	for i := 0; i < COUNT; i++ {
		if rsp, stateCode, err = doRequest(ctx, method, url, args, body); err == nil {
			return
		}

		if _, ok := HttpStateCode[stateCode]; !ok {
			return
		}

		time.Sleep(TIME_INTERNAL * time.Millisecond)
	}

	if err != nil {
		err = errors.New(fmt.Sprintf("DoRequest Fail: %v, url:%v", err.Error(), url))
	}
	return
}

func doRequest(ctx context.Context, method, url string, args map[string]string, body io.Reader) ([]byte, int, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, -1, errors.New(fmt.Sprintf(" http.NewRequest failed: %v", err.Error()))
	}

	name := "doRequest"
	szPath := strings.Split(url, "/")
	nLen := len(szPath)
	if nLen > 0 {
		name = szPath[nLen-1]
	}

	_, span := trace.StartSpan(ctx, name)
	defer span.End()

	f := &propagation.HTTPFormat{}
	//request = request.WithContext(nctx)
	f.SpanContextToRequest(span.SpanContext(), request)

	if len(args) > 0 {
		param := request.URL.Query()
		for k, v := range args {
			param.Add(k, v)
		}
		request.URL.RawQuery = param.Encode()
		//logger.Infof("do request: Method: %v,  url: %v,    query: %v,  args: %v", method, url, request.URL.RawQuery, args)
	}

	client := http.Client{}
	resp, err1 := client.Do(request)
	if err1 != nil {
		return nil, -1, err1
	}

	result, err2 := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err2 != nil {
		return nil, -1, errors.New(fmt.Sprintf(" ioutil.ReadAll(resp.Body)  failed: %v", err2.Error()))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New(fmt.Sprintf("doRequest fail, statuscode: %d, msg: %v,  url:%v", resp.StatusCode, string(result), url))
	}

	return result, resp.StatusCode, nil
}

func DoRequest2(method, url string, args, headers map[string]string, body io.Reader) (rsp []byte, err error) {
	for i := 0; i < COUNT; i++ {
		if rsp, err = doRequest2(method, url, args, headers, body); err == nil {
			return
		}
		time.Sleep(TIME_INTERNAL * time.Millisecond)
	}

	if err != nil {
		err = errors.New(fmt.Sprintf("DoRequest Fail: %v, url:%v", err.Error(), url))
	}
	return
}

func doRequest2(method, url string, args, headers map[string]string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(" http.NewRequest failed: %v", err.Error()))
	}

	if len(args) > 0 {
		param := request.URL.Query()
		for k, v := range args {
			param.Add(k, v)
		}
		request.URL.RawQuery = param.Encode()
		//logger.Infof("do request: Method: %v,  url: %v,    query: %v,  args: %v", method, url, request.URL.RawQuery, args)
	}

	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	client := http.Client{}
	resp, err1 := client.Do(request)
	if err1 != nil {
		return nil, err1
	}

	result, err2 := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err2 != nil {
		return nil, errors.New(fmt.Sprintf(" ioutil.ReadAll(resp.Body)  failed: %v", err2.Error()))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("doRequest fail, statuscode: %d, msg: %v", resp.StatusCode, string(result)))
	}

	return result, nil
}
