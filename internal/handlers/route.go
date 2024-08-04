package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gofermart/internal/middleware"
	"github.com/moonicy/gofermart/internal/storage"
)

func NewRoute(uh *UsersHandler, oh *OrdersHandler, us *storage.UsersStorage) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Post("/user/register", uh.PostUserReqister)
		r.Post("/user/login", uh.PostUserLogin)
		r.Route("/user", func(r chi.Router) {
			r.Use(middleware.Auth(us))
			r.Post("/orders", oh.PostUserOrders)
			//r.Get("/orders", GetUserOrders)
			//r.Route("/balance", func(r chi.Router) {
			//	r.Get("/", GetUserBalance)
			//	r.Post("/withdraw", PostUserBalanceWithdraw)
			//})
			//r.Get("/withdrawals", GetUserWithdrawals)
		})
	})

	return router
}
