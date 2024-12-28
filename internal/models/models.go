package models

import (
	"math/big"
	"time"
)

type TransactionData struct {
	BalanceNumeric   big.Int `json:"value"`
	Sender           string  `json:"from"`
	Recipient        string  `json:"to"`
	TransactionHash  string  `json:"hash"`
	TransactionIndex string  `json:"transactionIndex"`
}

type TransactionDataWithTimestamp struct {
	BalanceNumeric   big.Int   `json:"value"`
	Sender           string    `json:"from"`
	Recipient        string    `json:"to"`
	TransactionHash  string    `json:"hash"`
	TransactionIndex string    `json:"transactionIndex"`
	Timestamp        time.Time `json:"timestamp"`
}
