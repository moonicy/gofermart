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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	cfg := config.NewConfig()
	sugar := logger.NewLogger()
	db, err := storage.NewDB(cfg)
	if err != nil {
		sugar.Fatal(err)
	}
	migrator := storage.NewMigrator(db)
	err = migrator.Migrate(ctx)
	if err != nil {
		sugar.Fatal(err)
	}
	userStorage := storage.NewUsersStorage(db)
	ordersStorage := storage.NewOrdersStorage(db)
	userOrderStorage := storage.NewUserOrderStorage(db, ordersStorage, userStorage)
	movementsStorage := storage.NewMovementsStorage(db, userStorage, ordersStorage)

	client := accrual.NewClient(cfg.AccrualSystemAddress)
	syncOrders := demon.NewSyncOrders(ordersStorage, client, userOrderStorage)
	syncOrders.Run(ctx)

	usersHandler := handlers.NewUsersHandler(userStorage)
	ordersHandler := handlers.NewOrdersHandler(ordersStorage)
	movementHandler := handlers.NewMovementHandler(movementsStorage)

	route := handlers.NewRoute(usersHandler, ordersHandler, userStorage, movementHandler)

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: route,
	}

	err = server.ListenAndServe()
	if err != nil {
		sugar.Fatal(err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		sugar.Fatalw("Server shutdown error", "error", err)
	}

	cancelFn()
	time.Sleep(1 * time.Second)
}
