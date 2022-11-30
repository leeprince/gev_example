/*
 * @Date: 2021-09-03 16:43:35
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2022-02-11 11:13:03
 */
package tracer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// InjectHTTPHeader inject tracing information from context.Context to http.Header.
//
// Usage:
// 		ctx := opentracing.ContextWithSpan(context.Background(), span)
// 		req := http.NewRequest(http.MethodPost, uri, reqBodyReader)
//  	tracer.InjectHTTPHeader(ctx, httpRequest.Header)
//
func InjectHTTPHeader(ctx context.Context, header http.Header) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	err := GetTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header),
	)
	if err != nil {
		fmt.Println("gitlab.yewifi.com/golden-cloud/common/tracer/http.go: inject header header failed: ", err.Error())
	}
}
