package pg

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/OctaneAL/ETH-Tracker/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const transactionsTableName = "transactions"

func newTransactionQ(db *pgdb.DB) data.TransactionQ {
	return &transactionQ{
		db:  db,
		sql: sq.StatementBuilder,
	}
}

type transactionQ struct {
	db  *pgdb.DB
	sql sq.StatementBuilderType
}

func (q *transactionQ) Get() (*data.ReturnTransaction, error) {
	var result data.ReturnTransaction
	err := q.db.Get(&result, q.sql.Select("*").From(transactionsTableName))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce from db")
	}
	return &result, nil
}

func (q *transactionQ) Select() ([]data.ReturnTransaction, error) {
	var result []data.ReturnTransaction
	err := q.db.Select(&result, q.sql.Select("*").From(transactionsTableName).OrderBy("id DESC").Limit(100))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to select nonces from db")
	}
	return result, nil
}

func (q *transactionQ) Insert(value data.InsertTransaction) (*data.InsertTransaction, error) {
	clauses := structs.Map(value)

	var result data.InsertTransaction
	stmt := sq.Insert(transactionsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert nonce to db")
	}
	return &result, nil
}

// func (q *transactionQ) Update(value data.Transaction) (*data.Transaction, error) {
// 	clauses := structs.Map(value)

// 	var result data.Transaction
// 	stmt := q.sql.Update(transactionsTableName).SetMap(clauses).Suffix("returning *")
// 	err := q.db.Get(&result, stmt)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to update nonce in db")
// 	}
// 	return &result, nil
// }

// func (q *transactionQ) Delete() error {
// 	err := q.db.Exec(q.sql.Delete(transactionsTableName))
// 	if err != nil {
// 		return errors.Wrap(err, "failed to delete nonces from db")
// 	}
// 	return nil
// }

// func (q *transactionQ) FilterByAddress(addresses ...string) data.TransactionQ {
// 	pred := sq.Eq{"address": addresses}
// 	q.sql = q.sql.Where(pred)
// 	return q
// }

func (q *transactionQ) FilterBySenderRecipientHash(sender, recipient, transactionHash string) data.TransactionQ {
	pred := sq.Eq{}
	if sender != "" {
		pred["sender"] = sender
	}
	if recipient != "" {
		pred["recipient"] = recipient
	}
	if transactionHash != "" {
		pred["transaction_hash"] = transactionHash
	}

	q.sql = q.sql.Where(pred)
	return q
}

// func (q *transactionQ) FilterExpired() data.TransactionQ {
// 	q.sql = sq.StatementBuilder.Where("expiresat < ?", time.Now().Unix())
// 	return q
// }
