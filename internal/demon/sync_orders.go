package demon

import (
	"context"
	"errors"
	"fmt"
	"github.com/moonicy/gofermart/internal/accrual"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/pkg/workerpool"
	"sync/atomic"
	"time"
)

var ErrRateLimit = errors.New("rate limit exceeded")

type Client interface {
	GetOrderInfo(number string) (models.Order, error)
}

type OrdersStorage interface {
	GetBatch(ctx context.Context) ([]models.Order, error)
}

type UserOrderStorage interface {
	UpdateAccrual(ctx context.Context, order models.Order) error
}

type SyncOrders struct {
	isReachRateLimit atomic.Bool
	accrualCl        Client
	ordersStorage    OrdersStorage
	userOrderStorage UserOrderStorage
}

func NewSyncOrders(orderStorage OrdersStorage, accrualCl Client, userOrderStorage UserOrderStorage) *SyncOrders {
	return &SyncOrders{ordersStorage: orderStorage, accrualCl: accrualCl, userOrderStorage: userOrderStorage}
}

// Run запускает процесс синхронизации заказов с accrual.
// Работает на основе worker pool, для того чтобы распараллелить синхронизацию.
// Такая синхронизация выбрана потому что в accrual не имеется ручки для сихронизации батчем.
// В worker pool 10 workers поэтому мы синхронизируем в момент сразу 10 заказов.
func (so *SyncOrders) Run(ctx context.Context) {
	wp := workerpool.NewWorkerPool(so.Worker, 10)
	wp.Run()

	ticker := time.NewTicker(1 * time.Second)

	handle := func() {
		batch, err := so.ordersStorage.GetBatch(ctx)
		if err != nil {
			return
		}
		for _, order := range batch {
			wp.AddJob(so.MakeJob(ctx, order))
		}
	}

	go func() {
		defer func() {
			ticker.Stop()
			wp.Close()
		}()
		handle()
		for {
			select {
			case <-ticker.C:
				handle()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// MakeJob синхронизирует один заказ.
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

		err = so.userOrderStorage.UpdateAccrual(ctx, info)
		if err != nil {
			return err
		}
		return nil
	}
}

// Worker единица работы с механизмом лимитирования по количеству запросов.
// Вынесен отдельно т.к. имеет уникальный механизм лимитирования, который необходим только для синхронизации заказов.
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
