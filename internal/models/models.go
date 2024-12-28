package models

import "math/big"

type TransactionData struct {
	BalanceNumeric   big.Int `json:"value"`
	Sender           string  `json:"from"`
	Recipient        string  `json:"to"`
	TransactionHash  string  `json:"hash"`
	TransactionIndex int     `json:"transactionIndex"`
}
