package testservice

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "grpctest/protos"
)

func NewClientService() *ClientService {
	return &ClientService{
		chFinish:   make(chan struct{}, 2),
		chWriteErr: make(chan struct{}, 1),
		chReadErr:  make(chan struct{}, 1),
	}
}

type ClientService struct {
	chFinish   chan struct{}
	chReadErr  chan struct{}
	chWriteErr chan struct{}
}

func (s *ClientService) StartService() {
	conn, err := grpc.Dial("0.0.0.0:9091", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Printf("dial error:%v\n", err)
		return
	}
	defer conn.Close()
	c := proto.NewTestServiceClient(conn)
	callClient, err := c.CallEachOther(context.Background())
	if err != nil {
		fmt.Printf("call CallEachOther fail:%v\n", err)
		return
	}
	//读
	go func() {
		for {
			select {
			case <-s.chWriteErr:
				return
			default:
				break
			}
			data, err := callClient.Recv()
			if err != nil {
				fmt.Printf("recv error :%v\n", err)
				s.chReadErr <- struct{}{}
				s.chFinish <- struct{}{}
				return
			}
			fmt.Printf("[接受]：%s\n", data.Data)
		}
	}()
	//写
	go func() {
		for {
			select {
			case <-s.chReadErr:
				return
			default:
				break
			}
			fmt.Println("请输入要发送的字符串：")
			var str string
			fmt.Scanln(&str)
			err := callClient.Send(&proto.CallRequest{
				Data: str,
			})
			if err != nil {
				s.chWriteErr <- struct{}{}
				s.chFinish <- struct{}{}
				return
			}
		}
	}()
	select {
	case <-s.chFinish:
		break
	}
}
