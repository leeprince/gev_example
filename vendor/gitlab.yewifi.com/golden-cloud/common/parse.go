package common

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseYaml(file string, configRaw interface{}) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic("加载配置文件错误" + file + "错误原因" + err.Error())
	}

	err = yaml.Unmarshal(content, configRaw)
	if err != nil {
		panic("解析配置文件错误" + file + "错误原因" + err.Error())
	}
}
