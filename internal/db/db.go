package db

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/OctaneAL/ETH-Tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(connectionString string) *DB {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return &DB{Pool: pool}
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) SaveTransaction(ctx context.Context, transaction *models.TransactionData) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO transactions (balance_wei, sender, recipient, transaction_hash, transaction_index) 
		VALUES ($1, $2, $3, $4, $5)
	`, transaction.BalanceNumeric.String(), transaction.Sender, transaction.Recipient, transaction.TransactionHash, transaction.TransactionIndex)
	return err
}

func (db *DB) GetTransactionsWithFilters(ctx context.Context, sender, recipient, transactionHash string) ([]models.TransactionData, error) {
	var filters []string
	var args []interface{}

	if sender != "" {
		filters = append(filters, "sender = $"+fmt.Sprint(len(args)+1))
		args = append(args, sender)
	}
	if recipient != "" {
		filters = append(filters, "recipient = $"+fmt.Sprint(len(args)+1))
		args = append(args, recipient)
	}
	if transactionHash != "" {
		filters = append(filters, "transaction_hash = $"+fmt.Sprint(len(args)+1))
		args = append(args, transactionHash)
	}

	whereClause := ""
	if len(filters) > 0 {
		whereClause = "WHERE " + strings.Join(filters, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT balance_wei, sender, recipient, transaction_hash, transaction_index 
		FROM transactions 
		%s
		ORDER BY id
		LIMIT 100
	`, whereClause)

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.TransactionData

	for rows.Next() {
		var transaction models.TransactionData
		var balanceStr string
		err := rows.Scan(&balanceStr, &transaction.Sender, &transaction.Recipient, &transaction.TransactionHash, &transaction.TransactionIndex)
		if err != nil {
			log.Printf("Error scanning field: %v \n", err)
			return nil, err
		}
		balanceInt, ok := new(big.Int).SetString(balanceStr, 10)
		if !ok {
			log.Printf("Error converting balance to big.Int: %v \n", balanceStr)
			return nil, fmt.Errorf("error converting balance to big.Int: %v", balanceStr)
		}
		transaction.BalanceNumeric = *balanceInt
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
