package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gofermart/internal/middleware"
	"github.com/moonicy/gofermart/internal/storage"
)

func NewRoute(uh *UsersHandler, oh *OrdersHandler, us *storage.UsersStorage, mh *MovementHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Post("/user/register", uh.PostUserReqister)
		r.Post("/user/login", uh.PostUserLogin)
		r.Route("/user", func(r chi.Router) {
			r.Use(middleware.Auth(us))
			r.Post("/orders", oh.PostUserOrders)
			r.Get("/orders", oh.GetUserOrders)
			r.Route("/balance", func(r chi.Router) {
				r.Get("/", uh.GetUserBalance)
				r.Post("/withdraw", mh.PostUserBalanceWithdraw)
			})
			r.Get("/withdrawals", mh.GetUserWithdrawals)
		})
	})

	return router
}
