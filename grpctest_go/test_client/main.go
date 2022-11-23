package main

import (
	"fmt"
	testservice "grpctest/test_client/test_service"
)

func main() {
	for {
		fmt.Println("开始连接服务器...")
		s := testservice.NewClientService()
		s.StartService()
		fmt.Println("与服务器连接断开...")
	}
}
