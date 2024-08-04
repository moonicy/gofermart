package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/moonicy/gofermart/internal/config"
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

	uh := handlers.NewUsersHandler(us)
	oh := handlers.NewOrdersHandler(os)

	route := handlers.NewRoute(uh, oh, us)

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: route,
	}

	err = server.ListenAndServe()
	if err != nil {
		sugar.Fatal(err)
	}
}
