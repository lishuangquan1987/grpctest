package main

import (
	"fmt"
	testservice "grpctest/test_server/test_service"
)

func main() {
	fmt.Println("启动...")
	go testservice.StartService()

	for {
		fmt.Println("请输入要发送的字符串：")
		var str string
		fmt.Scanln(&str)
		testservice.Send(str)
	}
}
