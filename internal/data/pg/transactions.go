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

func (q *transactionQ) Get() (*data.Transaction, error) {
	var result data.Transaction
	err := q.db.Get(&result, q.sql.Select("*").From(transactionsTableName))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transaction from db")
	}
	return &result, nil
}

func (q *transactionQ) Select() ([]data.Transaction, error) {
	var result []data.Transaction

	// err := q.db.Select(&result, q.sql.Select("*").From(transactionsTableName).OrderBy("id DESC").Limit(100))

	columns := []string{"balance_wei", "sender", "recipient", "transaction_hash", "transaction_index", "block_number", "timestamp"}
	err := q.db.Select(&result, q.sql.Select(columns...).From(transactionsTableName).OrderBy("timestamp DESC").Limit(100))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to select transactions from db")
	}
	return result, nil
}

func (q *transactionQ) Insert(value data.Transaction) (*data.Transaction, error) {
	clauses := structs.Map(value)

	var result data.Transaction
	stmt := sq.Insert(transactionsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert transaction to db")
	}
	return &result, nil
}

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

func (q *transactionQ) GetLastRecord() (*data.Transaction, error) {
	var result data.Transaction

	columns := []string{"balance_wei", "sender", "recipient", "transaction_hash", "transaction_index", "block_number", "timestamp"}
	err := q.db.Get(&result, q.sql.Select(columns...).From(transactionsTableName).OrderBy("timestamp DESC").Limit(1))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last transaction from db")
	}
	return &result, nil
}
