package gclog

/*
 * @Date: 2020-10-29 15:52:25
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 15:54:43
 */

import "context"

type ctxKey struct{}

var logIDCtxKey ctxKey

func SetLogIDToCtx(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, logIDCtxKey, logID)
}

func GetLogIDFromCtx(ctx context.Context) string {
	logIDInterface := ctx.Value(logIDCtxKey)
	if logIDInterface == nil {
		return ""
	}
	logID, ok := logIDInterface.(string)
	if !ok {
		return ""
	}
	return logID
}
