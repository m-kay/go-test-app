package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

//go:generate protoc --proto_path=../../resources/grpc/ --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative --descriptor_set_out=../../resources/grpc/helloService.dsc helloService.proto
type HelloService struct {
	conn *grpc.ClientConn
}

func NewHelloService(port int) (*HelloService, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &HelloService{conn: conn}, nil
}

func (s *HelloService) GetHello() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := NewHelloServiceClient(s.conn).GetHello(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	return response.GetMessage(), nil
}
