package port

import "github.com/VallabhSLEPAM/grpc-server/internal/adapters/database"

type BankDatabasePort interface {
	GetBankAccountByAccountNumber(acct string) (database.BankAccountORM, error)
}
