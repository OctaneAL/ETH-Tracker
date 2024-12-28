package handlers

import (
	"net/http"

	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	database := DB(r.Context())

	transactions, err := database.GetAllTransactions(r.Context())

	if err != nil {
		ape.RenderErr(w, []*jsonapi.ErrorObject{problems.InternalError()}...)
		return
	}
	ape.Render(w, transactions)
}
