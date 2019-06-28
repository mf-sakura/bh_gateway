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

	// if _, err := userClient.IncrUserCounter(ctx, &upb.IncrUserCounterMessage{
	// 	UserId: req.UserId,
	// }); err != nil {
	// 	log.Printf("IncrUserCounter failed:%v\n", err)
	// 	return nil, grpc.Errorf(codes.Internal, "予約に失敗しました")
	// }
	// res, err := hotelClient.ReserveHotel(ctx, &hpb.ReserveHotellMessage{
	// 	PlanId:     req.PlanId,
	// 	UserId:     req.UserId,
	// 	SequenceId: req.SequenceId,
	// })
	// if err != nil {
	// 	log.Printf("ReserveHotel failed:%v\n", err)
	// 	if _, err := userClient.DecrUserCounter(ctx, &upb.DecrUserCounterMessage{UserId: req.UserId}); err != nil {
	// 		// ここで失敗するとデータ整合性が保てないので、リトライやエラー通知の機構が必要
	// 		log.Printf("Compensating Transaction: DecrUserCounter failed:%v\n", err)
	// 	}
	// 	return nil, grpc.Errorf(codes.Internal, "予約に失敗しました")
	// }
	// return &gpb.BookHotelResponse{ReservtionId: res.ReservationId}, nil
	bs := NewBookStateMachine(req.UserId, req.PlanId)
	bs.FSM.Event("incr")
	if bs.Err != nil {
		return nil, grpc.Errorf(codes.Internal, "予約に失敗しました")
	}
	return &gpb.BookHotelResponse{ReservtionId: *bs.ReservationID}, nil
}

type BookHotelMachine struct {
	UserID        int64
	PlanID        int64
	ReservationID *int64
	FSM           *fsm.FSM
	// ErrChan chan error
	Err error
}

func NewBookStateMachine(userID, planID int64) *BookHotelMachine {
	// ch := make(chan error)
	b := &BookHotelMachine{
		UserID: userID,
		PlanID: planID,
	}
	b.FSM = fsm.NewFSM(
		"start",
		fsm.Events{
			{Name: "incr", Src: []string{"start"}, Dst: "incr"},
			{Name: "reserve", Src: []string{"incr"}, Dst: "reserve"},
			{Name: "success", Src: []string{"reserve"}, Dst: "end"},
			{Name: "incr_err", Src: []string{"incr"}, Dst: "end"},
			{Name: "reserve_err", Src: []string{"reserve"}, Dst: "incr_rollback"},
			{Name: "rollbackend", Src: []string{"incr_rollback"}, Dst: "end"},
		},
		fsm.Callbacks{
			"before_incr":        b.Increment(),
			"before_reserve":     b.Reserve(),
			"before_reserve_err": b.RollBackIncrement(),
		},
	)
	return b
}

func (b *BookHotelMachine) Increment() func(*fsm.Event) {
	return func(event *fsm.Event) {
		ctx := event.Args[0].(context.Context)

		if _, err := userClient.IncrUserCounter(ctx, &upb.IncrUserCounterMessage{
			UserId: b.UserID,
		}); err != nil {
			// b.ErrChan <- err
			b.Err = err
			b.FSM.Event("incr_err", ctx, err)
		}
		b.FSM.Event("reserve", ctx)
	}
}

func (b *BookHotelMachine) Reserve() func(*fsm.Event) {
	return func(event *fsm.Event) {
		ctx := event.Args[0].(context.Context)
		res, err := hotelClient.ReserveHotel(ctx, &hpb.ReserveHotellMessage{
			PlanId: b.PlanID,
			UserId: b.UserID,
		})
		if err != nil {
			log.Printf("reserve_err:%v\n", err)
			b.FSM.Event("reserve_err", ctx)
		}
		b.ReservationID = &res.ReservationId
	}
}

func (b *BookHotelMachine) RollBackIncrement() func(*fsm.Event) {
	return func(event *fsm.Event) {
		ctx := event.Args[0].(context.Context)
		if _, err := userClient.DecrUserCounter(ctx, &upb.DecrUserCounterMessage{UserId: b.UserID}); err != nil {
			// ここで失敗するとデータ整合性が保てないので、リトライやエラー通知の機構が必要
			log.Printf("Compensating Transaction: DecrUserCounter failed:%v\n", err)
			// b.ErrChan <- err
			b.Err = err
		}
	}
}
