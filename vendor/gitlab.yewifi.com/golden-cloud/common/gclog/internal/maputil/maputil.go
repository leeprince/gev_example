package maputil

/*
 * @Date: 2020-09-27 12:49:04
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 11:55:36
 */

import (
	"fmt"
	"strings"
)

// GetOrCreatePath 根据给定的path，在map中获取相对应的子map。如果路径不存在，则创建并返回。
// 详细用例参考 ./maputil_test.go
func GetOrCreatePath(obj map[string]interface{}, path string) (map[string]interface{}, error) {
	if path == "" {
		return obj, nil
	}

	pathNames := strings.Split(path, ".")

	for i := 0; i < len(pathNames); i++ {
		pathName := pathNames[i]
		subObjInterface, ok := obj[pathName]
		if !ok {
			subObj := make(map[string]interface{})
			obj[pathName] = subObj
			obj = subObj
		} else {
			subObj, ok := subObjInterface.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("path %s is not a map", strings.Join(pathNames[:i+1], "."))
			}
			obj = subObj
		}
	}

	return obj, nil
}
