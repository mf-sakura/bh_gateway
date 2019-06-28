package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/looplab/fsm"
	gpb "github.com/mf-sakura/bh_gateway/app/proto/gateway"
	hpb "github.com/mf-sakura/bh_gateway/app/proto/hotel"
	upb "github.com/mf-sakura/bh_gateway/app/proto/user"
	"log"
)

type BookHotelServiceServerImpl struct {
}

func (b *BookHotelServiceServerImpl) BookHotel(ctx context.Context, req *gpb.BookHotelMessage) (*gpb.BookHotelResponse, error) {

	if _, err := userClient.IncrUserCounter(ctx, &upb.IncrUserCounterMessage{
		UserId: req.UserId,
	}); err != nil {
		log.Printf("IncrUserCounter failed:%v\n", err)
		return nil, grpc.Errorf(codes.Internal, "予約に失敗しました")
	}
	res, err := hotelClient.ReserveHotel(ctx, &hpb.ReserveHotellMessage{
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
