package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type AccessLog struct {
	Method   string `json:"method"`
	Url      string `json:"url"`
	ReqBody  string `json:"req_body"`
	RespBody string `json:"resp_body"`
	Duration string `json:"duration"`
	Status   int    `json:"status"`
}

type LoggerMiddlewareBuilder struct {
	logFunc func(ctx context.Context, al *AccessLog)
}

func NewLoggerMiddlewareBuilder(fn func(ctx context.Context, al *AccessLog)) *LoggerMiddlewareBuilder {
	return &LoggerMiddlewareBuilder{logFunc: fn}
}

func (l *LoggerMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}

		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}

		if ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewBuffer(body))
			ctx.Request.Body = reader

			if len(body) > 1024 {
				body = body[:1024]
			}

			al.ReqBody = string(body)
		}

		ctx.Writer = responseWriter{
			al:             al,
			ResponseWriter: ctx.Writer,
		}

		defer func() {
			al.Duration = time.Since(start).String()
			l.logFunc(ctx, al)
		}()

		ctx.Next()
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}
