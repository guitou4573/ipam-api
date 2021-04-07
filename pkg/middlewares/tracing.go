package middlewares

import (
	"log"
  "github.com/gin-gonic/gin"
  opentracing "github.com/opentracing/opentracing-go"
)

type tracingMiddleware struct {}

func (t tracingMiddleware) build() gin.HandlerFunc {
  return func(ctx *gin.Context) {
    tracer := opentracing.GlobalTracer()
		if tracer == nil {
			return
		}
    span := tracer.StartSpan("wire-api-req")
		ctx.Set("tracing-req-context", span)
		ctx.Next()
    span.Finish()
	}
}

func NewTracingMiddleware() (gin.HandlerFunc) {
  log.Println("NewTracingMiddleware")
	midTracing := tracingMiddleware{}.build()
  return midTracing
}
