package application

import (
	"log"

	port "github.com/VallabhSLEPAM/grpc-server/internal/ports.go"
)

type BankService struct {
	db port.BankDatabasePort
}

func NewBankService(dbPort port.BankDatabasePort) *BankService {
	return &BankService{
		db: dbPort,
	}
}

func (service *BankService) FindCurrentBalance(acct string) float64 {
	bankAccount, err := service.db.GetBankAccountByAccountNumber(acct)
	if err != nil {
		log.Println("Error in FindCurrentBalance:", err)
	}
	return bankAccount.CurrentBalance
}
