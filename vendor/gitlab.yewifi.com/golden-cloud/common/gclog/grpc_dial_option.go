package gclog

/*
 * @Date: 2021-05-31 10:43:26
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2021-09-07 19:20:22
 */

import (
	"context"
	"time"

	"gitlab.yewifi.com/golden-cloud/common"
	"google.golang.org/grpc"
)

func grpcUnaryClientInterceptor(slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logID := common.LogIdByCtx(ctx)

		// access log
		WithField("method", method).WithField("req", req).Info(logID, "发起 GRPC 请求")

		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(startTime)
		millisecond := duration.Milliseconds()

		entry := WithField("method", method).
			WithField("req", req).
			WithField("rsp", reply).
			WithField("err", err).
			WithField("duration", duration).
			WithField("millisecond", millisecond)

		// access log
		entry.Info(logID, "接收 GRPC 响应")

		// error log
		if err != nil {
			entry.Error(logID, "上游服务报错")
		}

		// slow log
		if duration >= slowThreshold {
			entry.Warn(logID, "上游慢响应")
		}

		return err
	}
}

func GRPCDialOption(slowThreshold time.Duration) grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpcUnaryClientInterceptor(slowThreshold))
}

func GRPCUnaryClientInterceptor(slowThreshold time.Duration) grpc.UnaryClientInterceptor {
	return grpcUnaryClientInterceptor(slowThreshold)
}
