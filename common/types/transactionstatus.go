package types

type TransactionStatus string

const (
	Success TransactionStatus = "success"
	Failed  TransactionStatus = "failed"
	Pending TransactionStatus = "pending"
	Unknown TransactionStatus = "unknown"
)
