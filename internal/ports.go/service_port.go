package port

import (
	"time"

	"github.com/VallabhSLEPAM/grpc-server/internal/application/domain/bank"
	"github.com/google/uuid"
)

// This acts as a service layer interface and implementation would be in application folder
type HelloServicePort interface {
	GenerateHello(string) string
}

type BankServicePort interface {
	FindCurrentBalance(acct string) float64
	CreateExchangeRate(r bank.ExchangeRate) (uuid.UUID, error)
	FindExchangeRate(fromCurr, toCurr string, ts time.Time) float64
	CreateTransaction(acct string, t bank.Transaction) (uuid.UUID, error)
	CalculateTransactionSummary(tcur *bank.TransactionSummary, transaction bank.Transaction) error
	Transfer(bank.TransferTransaction) (uuid.UUID, bool, error)
}
