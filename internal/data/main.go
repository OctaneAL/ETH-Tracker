package data

type MasterQ interface {
	New() MasterQ

	Trans() TransactionQ

	Transaction(fn func(db MasterQ) error) error
}
