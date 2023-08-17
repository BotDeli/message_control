package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"message_control/internal/config"
	"message_control/internal/server/middleware"
	"message_control/internal/storage"
	"net"
	"time"
)

//type serverS interface {
//	NewMessage(ctx context.Context, req *Message) (*Response, error)
//	ChatMessages(ctx context.Context, req *Users) (*ChatUsers, error)
//	FriendsList(ctx context.Context, req *User) (*List, error)
//	mustEmbedUnimplementedMessageControllerServer()
//}

func StartServer(cfg config.GrpcServerConfig, controller storage.MessageControl) error {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.CheckUUIDRequest(cfg.UUID)),
	)
	srv := &ServerMessageController{}
	RegisterMessageControllerServer(server, srv)
	l, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		return err
	}
	return server.Serve(l)
}

type ServerMessageController struct{}

func (ms *ServerMessageController) NewMessage(ctx context.Context, req *Message) (*Response, error) {
	fmt.Println(req)
	return &Response{Success: true}, nil
}

func (ms *ServerMessageController) ChatMessages(ctx context.Context, req *Users) (*ChatUsers, error) {
	fmt.Println(req)
	return &ChatUsers{Messages: []*Message{{From: "s", To: "SSS", Read: true, Text: "HiHiHi!!!!", Date: time.Now().String()}}}, nil
}

func (ms *ServerMessageController) FriendsList(ctx context.Context, req *User) (*List, error) {
	fmt.Println(req)
	return &List{}, nil
}

func (ms *ServerMessageController) mustEmbedUnimplementedMessageControllerServer() {
	fmt.Println("ooofff")
	return
}
