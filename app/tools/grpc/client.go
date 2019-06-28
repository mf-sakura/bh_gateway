package main

import (
	gpb "github.com/mf-sakura/bh_gateway/app/proto/gateway"

	"context"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50003", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := gpb.NewBookHotelServiceClient(conn)
	_, err = c.BookHotel(context.Background(), &gpb.BookHotelMessage{
		UserId: 1,
		PlanId: 2,
	})
	if err != nil {
		fmt.Printf("error:%v", err)
	}
}
