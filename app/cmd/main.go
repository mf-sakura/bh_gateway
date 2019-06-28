package main

import (
	"github.com/mf-sakura/bh_gateway/app/config"
	gpb "github.com/mf-sakura/bh_gateway/app/proto/gateway"
	"github.com/mf-sakura/bh_gateway/app/server"

	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	open_config "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	port = ":5003"
)

func main() {
	fmt.Println("Process Started.")
	cfg := open_config.Configuration{
		Sampler: &open_config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &open_config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "jaeger:6831",
		},
	}
	tracer, closer, err := cfg.New(
		"booking_hotel",
		open_config.Logger(jaeger.StdLogger),
	)
	defer closer.Close()
	if err != nil {
		fmt.Println(err)
	}
	opentracing.SetGlobalTracer(tracer)
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
