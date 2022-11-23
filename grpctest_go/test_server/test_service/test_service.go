package testservice

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	proto "grpctest/protos"
	"net"
	"sync"
)

var Service *TestService

type ClientInfo struct {
	UUID       string
	chWrite    chan string
	chRead     chan string
	chWriteErr chan struct{}
	chReadErr  chan struct{}
	chFinish   chan struct{}
}

type TestService struct {
	mu      sync.Mutex
	Clients map[string]*ClientInfo
}

func (s *TestService) CallEachOther(t proto.TestService_CallEachOtherServer) error {
	//加入Client
	client := &ClientInfo{
		UUID:       uuid.NewV4().String(),
		chWrite:    make(chan string),
		chRead:     make(chan string),
		chWriteErr: make(chan struct{}),
		chReadErr:  make(chan struct{}),
		chFinish:   make(chan struct{}, 2),
	}
	s.Clients[client.UUID] = client
	fmt.Printf("客户端：%s已连接\n", client.UUID)
	//写
	go func() {
		for {
			select {
			case str, ok := <-client.chWrite:
				if !ok {
					client.chFinish <- struct{}{}
					client.chWriteErr <- struct{}{}
					return //通道关闭，调用完成
				}
				err := t.Send(&proto.CallResponse{Data: str})
				if err != nil {
					client.chWriteErr <- struct{}{}
					client.chFinish <- struct{}{}
					return
				}
			case <-client.chReadErr: //读错误时，写要停止
				return
			}
		}
	}()
	go func() {
		for {
			r, err := t.Recv()
			if err != nil {
				client.chReadErr <- struct{}{}
				client.chFinish <- struct{}{}
				return
			}
			fmt.Printf("[%s-接收]:%s\n", client.UUID, r.Data)
		}
	}()
	select {
	case <-client.chFinish:
		delete(s.Clients, client.UUID)
		fmt.Printf("客户端：%s已断开\n", client.UUID)
		return nil
	}
}

func Send(msg string) {
	for _, v := range Service.Clients {
		v.chWrite <- msg
	}
	fmt.Printf("已向%d个客户端发送了消息：%s", len(Service.Clients), msg)
}

func StartService() {
	lis, err := net.Listen("tcp", "0.0.0.0:9091")
	if err != nil {
		fmt.Printf("listen error:%v\n", err)
		return
	}
	s := grpc.NewServer()
	Service = &TestService{
		Clients: map[string]*ClientInfo{},
	}
	proto.RegisterTestServiceServer(s, Service)
	if err := s.Serve(lis); err != nil {
		fmt.Printf("serve error:%v\n", err)
		return
	}
}
