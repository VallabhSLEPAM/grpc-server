package bank

import (
	"errors"
	"time"
)

const (
	TransactionTypeUnknown string = "UNKNOWN"
	TransactionTypeIn      string = "IN"
	TransactionTypeOut     string = "OUT"
)

type ExchangeRate struct {
	FromCurrency       string
	ToCurrency         string
	Rate               float64
	ValidFromTimeStamp time.Time
	ValidToTimeStamp   time.Time
}

type Transaction struct {
	TransactionType string
	Amount          float64
	Timestamp       time.Time
	Notes           string
}

type TransactionSummary struct {
	SumIn         float64
	SumOut        float64
	SumTotal      float64
	SummaryOnDate time.Time
}

type TransferTransaction struct {
	FromAccount string
	ToAccount   string
	Currency    string
	Amount      float64
}

var ErrTransferSourceAccountNotFound = errors.New("source account not found")
var ErrTransferDestincationAccountNotFound = errors.New("destination account not found")
var ErrTransferRecordCreationFailed = errors.New("can't create transfer record")
var ErrTransferTransactionPair = errors.New("can't create transfer transaction pair possible due to insufficient balance at source account")
