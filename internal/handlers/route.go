package handlers

import (
	"github.com/go-chi/chi/v5"
)

func NewRoute(uh *UsersHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", uh.PostUserReqister)
			r.Post("/login", uh.PostUserLogin)
			//r.Post("/orders", PostUserOrders)
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
