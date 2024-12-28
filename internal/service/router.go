package service

import (
	"context"

	"github.com/OctaneAL/ETH-Tracker/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxDB(context.Background(), s.db),
		),
	)

	r.Route("/integrations/ETH-Tracker", func(r chi.Router) {
		r.Get("/transactions", handlers.GetFilteredTransactions)
	})

	return r
}
