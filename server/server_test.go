package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_CacheManager_parseLogReader(t *testing.T) {
	// 读取的日志文件在 my_config/access_security.log 下
	file, err := os.Open("../my_config/access_security.log")
	assert.Nil(t, err)
	defer file.Close()

	cm := NewCacheManager("my_config")

	rows, err := cm.parseLogReader(file, "my_config/access_security.log")
	assert.Nil(t, err)

	for _, row := range rows {
		fmt.Printf("%+v\n", row)
	}
}
