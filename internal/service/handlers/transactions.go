package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetFilteredTransactions(w http.ResponseWriter, r *http.Request) {
	database := DB(r.Context())

	sender := strings.TrimSpace(r.URL.Query().Get("sender"))
	recipient := strings.TrimSpace(r.URL.Query().Get("recipient"))
	transactionHash := strings.TrimSpace(r.URL.Query().Get("transactionHash"))

	transactions, err := database.GetTransactionsWithFilters(r.Context(), sender, recipient, transactionHash)
	if err != nil {
		log.Printf("%v", err)
		ape.RenderErr(w, []*jsonapi.ErrorObject{problems.InternalError()}...)
		return
	}

	ape.Render(w, transactions)
}
