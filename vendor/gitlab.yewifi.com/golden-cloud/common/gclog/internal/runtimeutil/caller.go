package runtimeutil

/*
 * @Date: 2020-10-29 13:33:23
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 13:34:26
 */

import (
	"path/filepath"
	"runtime"
)

func GetCaller(i int) (string, int) {
	_, path, lineNo, ok := runtime.Caller(i)
	if !ok {
		return "", 0
	}
	dirName := filepath.Base(filepath.Dir(path))
	filename := filepath.Base(path)
	callerName := filepath.Join(dirName, filename)
	return callerName, lineNo
}
