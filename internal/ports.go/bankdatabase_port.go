package port

import (
	"time"

	"github.com/VallabhSLEPAM/grpc-server/internal/adapters/database"
	"github.com/google/uuid"
)

// This is like DAO layer interface and its implementation would be in database folder
type BankDatabasePort interface {
	GetBankAccountByAccountNumber(acct string) (database.BankAccountORM, error)
	CreateExchangeRate(r database.BankExchangeRateORM) (uuid.UUID, error)
	GetExchangeRateAtTimestamp(fromCurr, toCurr string, ts time.Time) (database.BankExchangeRateORM, error)
	CreateTransactions(database.BankAccountORM, database.BankTransactionORM) (uuid.UUID, error)

	CreateTransfer(database.BankTransferORM) (uuid.UUID, error)
	CreateTransferTransactionPair(fromAccount, toAccount database.BankAccountORM, fromTransactionORM, toTransactionORM database.BankTransactionORM) (bool, error)
	UpdateTransferStatus(bankTransferORM database.BankTransferORM, status bool) error
}
