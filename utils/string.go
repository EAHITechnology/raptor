package utils

import (
	"fmt"
	"strings"
)

/*
直接拼接字符串,如果缓冲变得太大，Write会采用错误值 ErrTooLarge 引发panic
*/
func Write(strs ...string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Write 捕获到的错误：%s\n", r)
		}
	}()
	var strBuild strings.Builder
	for _, str := range strs {
		strBuild.WriteString(str)
	}
	return strBuild.String(), nil
}
