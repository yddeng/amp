package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

//读取json文件并反序列化
func DecodeJsonFromFile(i interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, i)
}

// 序列化对象并保存
func EncodeJsonToFile(i interface{}, filename string) error {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}
	_ = os.MkdirAll(path.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, data, os.ModePerm)
}
