package models

import (
	"math/big"
	"time"
)

type BlockHash struct {
	BlockNumber *big.Int
	Timestamp   *time.Time
}
