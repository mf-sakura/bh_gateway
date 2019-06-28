package server

import (
	"fmt"
	"github.com/mf-sakura/bh_gateway/app/config"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	hpb "github.com/mf-sakura/bh_gateway/app/proto/hotel"
	upb "github.com/mf-sakura/bh_gateway/app/proto/user"
)

var userClient upb.UserServiceClient
var hotelClient hpb.HotelServiceClient

func CreateConnections(conf *config.GRPCConfig) error {
	userTarget := fmt.Sprintf("%s:%s", conf.UserHost, conf.UserPort)
	hotelTarget := fmt.Sprintf("%s:%s", conf.HotelHost, conf.HotelPort)
	uconn, err := grpc.Dial(userTarget, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()))
	if err != nil {
		return err
	}
	hconn, err := grpc.Dial(hotelTarget, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()))
	if err != nil {
		return err
	}
	uc := upb.NewUserServiceClient(uconn)
	hc := hpb.NewHotelServiceClient(hconn)
	userClient = uc
	hotelClient = hc

	return nil
}
