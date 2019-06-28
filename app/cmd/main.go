package main

import (
	"github.com/mf-sakura/bh_gateway/app/config"
	gpb "github.com/mf-sakura/bh_gateway/app/proto/gateway"
	"github.com/mf-sakura/bh_gateway/app/server"

	"fmt"
	"google.golang.org/grpc"
	"net"
)

const (
	port = ":5003"
)

func main() {
	fmt.Println("Process Started.")
	conf, err := config.LoadConifg()
	if err != nil {
		panic(err)
	}
	if err := server.CreateConnections(conf); err != nil {
		panic(err)
	}
	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	defer func() {
		err := recover()
		s.GracefulStop()
		if err != nil {
			panic(err)
		}
	}()
	gpb.RegisterBookHotelServiceServer(s, &server.BookHotelServiceServerImpl{})

	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
