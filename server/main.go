package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	proto "github.com/utevo/gRPC-Chat/proto"
)

type Connection struct {
	stream proto.Broadcast_CreateStreamServer
	user   proto.User
	active bool
	errors chan error
}

type Server struct {
	connections []*Connection
}

func (server *Server) CreateStream(user *proto.User, stream proto.Broadcast_CreateStreamServer) error {
	connection := Connection{
		stream: stream,
		user:   *user,
		active: true,
		errors: make(chan error),
	}

	server.connections = append(server.connections, &connection)

	fmt.Println("New Steam Created")
	return nil
}

func (server *Server) BroadcastMessage(context context.Context, message *proto.Message) (*proto.Close, error) {
	fmt.Println("New Message: ", message)
	return &proto.Close{}, nil
}

func main() {
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", "localhost:5050")
	if err != nil {
		panic(err)
	}
	server := &Server{}

	proto.RegisterBroadcastServer(grpcServer, server)

	fmt.Println("Start serving")
	grpcServer.Serve(listener)
}
