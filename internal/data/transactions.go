package data

import "time"

type TransactionQ interface {
	Get() (*Transaction, error)
	Select() ([]Transaction, error)
	Insert(value Transaction) (*Transaction, error)
	GetLastRecord() (*Transaction, error)

	FilterBySenderRecipientHash(sender, recipient, transactionHash string) TransactionQ
}

type Transaction struct {
	BalanceNumeric   string    `db:"balance_wei" structs:"balance_wei" json:"value"`
	Sender           string    `db:"sender" structs:"sender" json:"from"`
	Recipient        string    `db:"recipient" structs:"recipient" json:"to"`
	TransactionHash  string    `db:"transaction_hash" structs:"transaction_hash" json:"hash"`
	TransactionIndex int64     `db:"transaction_index" structs:"transaction_index" json:"transactionIndex"`
	BlockNumber      int64     `db:"block_number" structs:"block_number" json:"blockNumber"`
	Timestamp        time.Time `db:"timestamp" structs:"timestamp" json:"timestamp"`
}

// type InsertTransaction struct {
// 	BalanceNumeric   string `db:"balance_wei" structs:"balance_wei" json:"value"`
// 	Sender           string `db:"sender" structs:"sender" json:"from"`
// 	Recipient        string `db:"recipient" structs:"recipient" json:"to"`
// 	TransactionHash  string `db:"transaction_hash" structs:"transaction_hash" json:"hash"`
// 	TransactionIndex int64  `db:"transaction_index" structs:"transaction_index" json:"transactionIndex"`
// 	BlockNumber      int64  `db:"block_number" structs:"block_number" json:"blockNumber"`
// }

// type ReturnTransaction struct {
// 	ID int64 `db:"id" structs:"-" json:"id"`

// 	BalanceNumeric   string    `db:"balance_wei" structs:"balance_wei" json:"value"`
// 	Sender           string    `db:"sender" structs:"sender" json:"from"`
// 	Recipient        string    `db:"recipient" structs:"recipient" json:"to"`
// 	TransactionHash  string    `db:"transaction_hash" structs:"transaction_hash" json:"hash"`
// 	TransactionIndex int64     `db:"transaction_index" structs:"transaction_index" json:"transactionIndex"`
// 	BlockNumber      int64     `db:"block_number" structs:"block_number" json:"blockNumber"`
// 	Timestamp        time.Time `db:"timestamp" structs:"timestamp" json:"timestamp"`
// }
