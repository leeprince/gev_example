package common

import "github.com/grpc-ecosystem/grpc-gateway/runtime"

func Marshal(msg interface{}) ([]byte, error) {
	marshaler := &runtime.JSONPb{OrigName: true, EmitDefaults: true}
	return marshaler.Marshal(msg)
}
