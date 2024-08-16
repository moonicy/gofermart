package demon

import (
	"context"
	"errors"
	"fmt"
	"github.com/moonicy/gofermart/internal/accrual"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/internal/storage"
	"github.com/moonicy/gofermart/pkg/workerpool"
	"sync/atomic"
	"time"
)

var ErrRateLimit = errors.New("rate limit exceeded")

type SyncOrders struct {
	isReachRateLimit atomic.Bool
	accrualCl        *accrual.Client
	us               *storage.UsersStorage
	os               *storage.OrdersStorage
	uos              *storage.UserOrderStorage
}

func NewSyncOrders(orderStorage *storage.OrdersStorage, accrualCl *accrual.Client, us *storage.UsersStorage, uos *storage.UserOrderStorage) *SyncOrders {
	return &SyncOrders{os: orderStorage, accrualCl: accrualCl, us: us, uos: uos}
}

func (so *SyncOrders) Run(ctx context.Context) func() {
	wp := workerpool.NewWorkerPool(so.Worker, 10)
	wp.Run()

	ticker := time.NewTicker(1 * time.Second)

	handle := func() {
		batch, err := so.os.GetBatch(ctx)
		if err != nil {
			return
		}
		for _, order := range batch {
			wp.AddJob(so.MakeJob(ctx, order))
		}
	}

	go func() {
		handle()
		for {
			select {
			case <-ticker.C:
				handle()
			}
		}
	}()
	return func() {
		ticker.Stop()
		wp.Close()
	}
}

func (so *SyncOrders) MakeJob(ctx context.Context, order models.Order) func() error {
	return func() error {
		info, err := so.accrualCl.GetOrderInfo(order.Number)
		if err != nil {
			if errors.Is(err, accrual.ErrTooManyRequests) {
				return ErrRateLimit
			}
			return err
		}

		info.UserID = order.UserID
		info.ID = order.ID

		err = so.uos.UpdateAccrual(ctx, info)
		if err != nil {
			return err
		}
		return nil
	}
}

func (so *SyncOrders) Worker(ch <-chan workerpool.Job) {
	for {
		if so.isReachRateLimit.Load() {
			time.Sleep(1 * time.Second)
			so.isReachRateLimit.Store(false)
		}
		job, has := <-ch
		if !has {
			return
		}
		err := job()
		if err == nil {
			continue
		}
		if errors.Is(err, ErrRateLimit) {
			so.isReachRateLimit.Store(true)
			continue
		}
		fmt.Printf("Error in job %v\n", err)
	}
}
