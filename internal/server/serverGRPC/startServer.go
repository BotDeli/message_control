package serverGRPC

import (
	"context"
	"google.golang.org/grpc"
	"message_control/internal/config"
	"message_control/internal/message"
	"message_control/internal/server/middleware"
	"message_control/internal/server/serverGRPC/pb"
	"message_control/internal/storage"
	"net"
	"time"
)

type ServiceMessageControl interface {
	PostNewMessage(ctx context.Context, req *pb.BodyMessage) (*pb.Response, error)
	GetChatMessages(ctx context.Context, req *pb.Users) (*pb.ChatUsers, error)
	GetFriendsList(ctx context.Context, req *pb.User) (*pb.FriendList, error)
}

func StartServer(cfg config.GrpcServerConfig, controller storage.MessageControl) error {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.CheckUUIDRequest(cfg.UUID)),
	)
	var srv ServiceMessageControl = &ServerMessageController{
		Controller: controller,
	}
	pb.RegisterMessageControllerServer(server, srv)
	l, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		return err
	}
	return server.Serve(l)
}

type ServerMessageController struct {
	Controller storage.MessageControl
}

func (ms *ServerMessageController) PostNewMessage(_ context.Context, req *pb.BodyMessage) (*pb.Response, error) {
	msg := message.Message{
		From: req.From,
		To:   req.To,
		Text: req.Text,
		Date: time.Now(),
	}

	added, err := ms.Controller.AddNewMessage(msg)
	return &pb.Response{Success: added}, err
}

func (ms *ServerMessageController) GetChatMessages(_ context.Context, req *pb.Users) (*pb.ChatUsers, error) {
	messages, err := ms.Controller.GetMessagesChat(req.Username, req.Friend)
	return &pb.ChatUsers{Messages: messages}, err
}

func (ms *ServerMessageController) GetFriendsList(_ context.Context, req *pb.User) (*pb.FriendList, error) {
	friends, err := ms.Controller.GetFriendsList(req.Username)
	return &pb.FriendList{Friends: friends}, err
}
