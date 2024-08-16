package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/moonicy/gofermart/internal/accrual"
	"github.com/moonicy/gofermart/internal/config"
	"github.com/moonicy/gofermart/internal/demon"
	"github.com/moonicy/gofermart/internal/handlers"
	"github.com/moonicy/gofermart/internal/storage"
	"github.com/moonicy/gofermart/pkg/logger"
	"net/http"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	sugar := logger.NewLogger()
	db, err := storage.NewDB(cfg)
	if err != nil {
		sugar.Fatal(err)
	}
	mg := storage.NewMigrator(db)
	err = mg.Migrate(ctx)
	if err != nil {
		sugar.Fatal(err)
	}
	us := storage.NewUsersStorage(db)
	os := storage.NewOrdersStorage(db)
	uos := storage.NewUserOrderStorage(db, os, us)
	ms := storage.NewMovementsStorage(db, us, os)

	cl := accrual.NewClient(cfg.AccrualSystemAddress)
	syncOrders := demon.NewSyncOrders(os, cl, us, uos)
	cancelFn := syncOrders.Run(ctx)

	uh := handlers.NewUsersHandler(us)
	oh := handlers.NewOrdersHandler(os)
	mh := handlers.NewMovementHandler(ms)

	route := handlers.NewRoute(uh, oh, us, mh)

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: route,
	}

	err = server.ListenAndServe()
	if err != nil {
		sugar.Fatal(err)
	}
	cancelFn()
}
