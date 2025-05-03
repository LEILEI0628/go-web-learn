package grpc

import (
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	// grpc的server
	server := grpc.NewServer()
	defer func() {
		server.GracefulStop() // 优雅退出
	}()
	// 业务server
	userServer := &UserServer{}
	RegisterUserServiceServer(server, userServer)
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	t.Log(err)
}
