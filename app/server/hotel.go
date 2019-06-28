package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	gpb "github.com/mf-sakura/bh_gateway/app/proto/gateway"
	hpb "github.com/mf-sakura/bh_gateway/app/proto/hotel"
	upb "github.com/mf-sakura/bh_gateway/app/proto/user"
	"log"

	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"time"
)

type BookHotelServiceServerImpl struct {
}

func (b *BookHotelServiceServerImpl) BookHotel(ctx context.Context, req *gpb.BookHotelMessage) (*gpb.BookHotelResponse, error) {

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "jaeger:6831",
		},
	}
	tracer, closer, err := cfg.New(
		"booking_hotel",
		config.Logger(jaeger.StdLogger),
	)
	defer closer.Close()
	if err != nil {
		fmt.Println(err)
	}
	opentracing.SetGlobalTracer(tracer)
	span := opentracing.GlobalTracer().StartSpan("gateway")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx, span)

	if _, err := userClient.IncrUserCounter(spanCtx, &upb.IncrUserCounterMessage{
		UserId: req.UserId,
	}); err != nil {
		log.Printf("IncrUserCounter failed:%v\n", err)
		return nil, grpc.Errorf(codes.Internal, "予約に失敗しました")
	}
	res, err := hotelClient.ReserveHotel(spanCtx, &hpb.ReserveHotellMessage{
		PlanId:     req.PlanId,
		UserId:     req.UserId,
		SequenceId: req.SequenceId,
	})
	if err != nil {
		log.Printf("ReserveHotel failed:%v\n", err)
		if _, err := userClient.DecrUserCounter(ctx, &upb.DecrUserCounterMessage{UserId: req.UserId}); err != nil {
			// ここで失敗するとデータ整合性が保てないので、リトライやエラー通知の機構が必要
			log.Printf("Compensating Transaction: DecrUserCounter failed:%v\n", err)
		}
		return nil, grpc.Errorf(codes.Internal, "予約に失敗しました")
	}
	return &gpb.BookHotelResponse{ReservtionId: res.ReservationId}, nil

}
