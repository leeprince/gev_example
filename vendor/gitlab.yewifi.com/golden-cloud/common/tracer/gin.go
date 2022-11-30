/*
 * @Date: 2021-02-07 19:19:11
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2021-09-03 16:30:03
 */
package tracer

import (
	"context"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const mdGinOutgoingKey = "mdGinOutgoingKey"

const tracerContextKey = "tracer:context"

// Gin start point, defer span.Finish()
func GinStart(c *gin.Context, operationName string) opentracing.Span {
	var err error
	var spanCtx opentracing.SpanContext
	md, _ := c.Get(mdGinOutgoingKey)
	if md != nil {
		spanCtx, _ = GetTracer().Extract(opentracing.TextMap, MDReaderWriter{md.(metadata.MD)})
	} else {
		spanCtx, _ = GetTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	}

	span := opentracing.StartSpan(
		operationName,
		ext.RPCServerOption(spanCtx),
		opentracing.Tag{Key: string(ext.Component), Value: "Gin-Http"},
		ext.SpanKindRPCClient,
	)

	carrier := opentracing.TextMapCarrier{}
	err = GetTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		fmt.Println("inject-error", err.Error())
	}

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	c.Set(mdGinOutgoingKey, metadata.New(carrier))
	c.Set(tracerContextKey, ctx)
	return span
}

// CtxFromGinContext get context.Context with span from *gin.Context.
// Usage:
// 		ctx := tracer.CtxFromGinContext(c)
// 		logID := common.LogIdByCtx(ctx)  // this logID equal to opentracing tracer id
//
// This function should used with [gin middleware](./middleware/tracer.go)
func CtxFromGinContext(c *gin.Context) context.Context {
	val, ok := c.Get(tracerContextKey)
	if !ok {
		return context.Background()
	}
	ctx, ok := val.(context.Context)
	if !ok {
		return context.Background()
	}
	return ctx
}

// Gin client
func GinGrpcDialOption() grpc.DialOption {
	return grpc.WithUnaryInterceptor(ginClientInterceptor(GetTracer()))
}

// Gin ginClientInterceptor
func ginClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// md := ctx.Value(mdGinOutgoingKey).(metadata.MD)
		md := metadata.New(nil)
		spanCtx, _ := tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
		span := tracer.StartSpan(
			method,
			ext.RPCServerOption(spanCtx),
			opentracing.Tag{Key: string(ext.Component), Value: "Gin-gRPC"},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()
		mdWriter := MDReaderWriter{md}
		err := tracer.Inject(span.Context(), opentracing.TextMap, mdWriter)
		if err != nil {
			fmt.Println("inject-error", err.Error())
		}
		parentCtx := metadata.NewOutgoingContext(ctx, md)
		err = invoker(parentCtx, method, req, reply, cc, opts...)
		if err != nil && err != io.EOF {
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
			fmt.Println("call-error", err.Error())
		}
		return err
	}
}
