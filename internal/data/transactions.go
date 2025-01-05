package data

import "time"

type TransactionQ interface {
	Get() (*ReturnTransaction, error)
	Select() ([]ReturnTransaction, error)
	Insert(value InsertTransaction) (*InsertTransaction, error)

	FilterByAddress(addresses ...string) TransactionQ
	FilterExpired() TransactionQ
}

type InsertTransaction struct {
	BalanceNumeric   string `db:"balance_wei" structs:"balance_wei" json:"value"`
	Sender           string `db:"sender" structs:"sender" json:"from"`
	Recipient        string `db:"recipient" structs:"recipient" json:"to"`
	TransactionHash  string `db:"transaction_hash" structs:"transaction_hash" json:"hash"`
	TransactionIndex string `db:"transaction_index" structs:"transaction_index" json:"transactionIndex"`
}

type ReturnTransaction struct {
	ID int64 `db:"id" structs:"-" json:"id"`

	BalanceNumeric   string    `db:"balance_wei" structs:"balance_wei" json:"value"`
	Sender           string    `db:"sender" structs:"sender" json:"from"`
	Recipient        string    `db:"recipient" structs:"recipient" json:"to"`
	TransactionHash  string    `db:"transaction_hash" structs:"transaction_hash" json:"hash"`
	TransactionIndex string    `db:"transaction_index" structs:"transaction_index" json:"transactionIndex"`
	Timestamp        time.Time `db:"timestamp" structs:"timestamp" json:"timestamp"`

	// Message string `db:"message" structs:"message"`
	// Expires int64  `db:"expiresat" structs:"expiresat"`
	// Address string `db:"address" structs:"address"`
}
