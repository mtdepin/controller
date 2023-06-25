package api

import (
	"context"
	"controller/pkg/logger"
	utilruntime "controller/pkg/runtime"
	"encoding/json"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
	"net/http"
	"strings"
)

const (
	copyDirective    = "COPY"
	replaceDirective = "REPLACE"
)

// HttpTraceAll Log headers and body.
func HttpTraceAll(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer utilruntime.HandleCrash()
		format := &propagation.HTTPFormat{}

		name := "HttpTraceAll"
		szDir := strings.Split(r.URL.Path, "/")
		if len(szDir) > 0 {
			name = szDir[len(szDir)-1]
		}

		sc, ok := format.SpanContextFromRequest(r)
		if ok {
			ctx, span := trace.StartSpanWithRemoteParent(r.Context(), name, sc)
			span.AddAttributes(trace.StringAttribute("Host", r.Host))
			vars := r.URL.Query()
			args, _ := json.Marshal(vars)
			span.AddAttributes(trace.StringAttribute("args", string(args)))
			span.AddAttributes(trace.StringAttribute("method", r.Method))
			span.AddAttributes(trace.StringAttribute("uri", r.RequestURI))

			logger.Debug("------------------------------------------------")
			logger.Debugf("request: %s, method: %s", r.RequestURI, r.Method)
			logger.Debugf("body length: %d", r.ContentLength)
			logger.Debugf("query args: %s", args)
			//for k, _ := range vars {
			//	logger.Debugf("args :%s, %s", k, vars.Get(k))
			//}
			logger.Debug()
			defer span.End()

			r = r.WithContext(ctx)

			//logger.Infof("-----------------helo HttpTraceAll receive ctx:", ctx)
		} else { //不存在，新增一个
			r = r.WithContext(context.Background())
			ctx, span := trace.StartSpan(r.Context(), name)
			defer span.End()

			r = r.WithContext(ctx)
			//logger.Infof("-----------------helo HttpTraceAll not receive ctx:", ctx)
		}

		//todo : do something
		f.ServeHTTP(w, r)
	}
}
