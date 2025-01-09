package application

import (
	"fmt"
	"log"
	"time"

	"github.com/VallabhSLEPAM/grpc-server/internal/adapters/database"
	"github.com/VallabhSLEPAM/grpc-server/internal/application/domain/bank"
	port "github.com/VallabhSLEPAM/grpc-server/internal/ports.go"
	"github.com/google/uuid"
)

type BankService struct {
	db port.BankDatabasePort
}

func NewBankService(dbPort port.BankDatabasePort) *BankService {
	return &BankService{
		db: dbPort,
	}
}

func (service *BankService) FindCurrentBalance(acct string) (float64, error) {
	bankAccount, err := service.db.GetBankAccountByAccountNumber(acct)
	if err != nil {
		log.Println("Error in FindCurrentBalance:", err)
		return 0, err
	}
	return bankAccount.CurrentBalance, nil
}

func (service *BankService) CreateExchangeRate(r bank.ExchangeRate) (uuid.UUID, error) {
	now := time.Now()
	bankExchangeRateORM := database.BankExchangeRateORM{
		ExchangeRateUUID:   uuid.New(),
		FromCurrency:       r.FromCurrency,
		ToCurrency:         r.ToCurrency,
		Rate:               r.Rate,
		ValidFromTimestamp: r.ValidFromTimeStamp,
		ValidToTimestamp:   r.ValidToTimeStamp,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	return service.db.CreateExchangeRate(bankExchangeRateORM)
}

func (service *BankService) FindExchangeRate(fromCurr, toCurr string, ts time.Time) (float64, error) {
	exchangeRate, err := service.db.GetExchangeRateAtTimestamp(fromCurr, toCurr, ts)
	if err != nil {
		return 0, err
	}
	return float64(exchangeRate.Rate), nil
}

func (service *BankService) CreateTransaction(acct string, t bank.Transaction) (uuid.UUID, error) {
	now := time.Now()
	newUUID := uuid.New()

	bankAccountORM, err := service.db.GetBankAccountByAccountNumber(acct)
	if err != nil {
		log.Printf("Can't create transaction for %v: %v\n", acct, err)
		return uuid.Nil, fmt.Errorf("can't create transaction for %v: %v", acct, err)
	}

	if t.TransactionType == bank.TransactionTypeOut && bankAccountORM.CurrentBalance < t.Amount {
		return bankAccountORM.AccountUUID, fmt.Errorf(
			"insufficient balance",
		)
	}

	bankTransactionORM := database.BankTransactionORM{
		TransactionUUID:      newUUID,
		AccountUUID:          bankAccountORM.AccountUUID,
		TransactionType:      t.TransactionType,
		TransactionTimestamp: now,
		Amount:               t.Amount,
		Notes:                t.Notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	return service.db.CreateTransactions(bankAccountORM, bankTransactionORM)

}

func (service *BankService) CalculateTransactionSummary(tcur *bank.TransactionSummary, transaction bank.Transaction) error {

	switch transaction.TransactionType {
	case bank.TransactionTypeIn:
		tcur.SumIn += transaction.Amount
	case bank.TransactionTypeOut:
		tcur.SumOut += transaction.Amount
	default:
		return fmt.Errorf("unknown transaction type %v", transaction.TransactionType)
	}
	tcur.SumTotal = tcur.SumIn - tcur.SumOut
	return nil
}

func (service *BankService) Transfer(tt bank.TransferTransaction) (uuid.UUID, bool, error) {
	now := time.Now()

	fromAccountORM, err := service.db.GetBankAccountByAccountNumber(tt.FromAccount)
	if err != nil {
		log.Printf("Can't find transfer from account %v: %v\n", tt.FromAccount, err)
		return uuid.Nil, false, bank.ErrTransferSourceAccountNotFound
	}

	if fromAccountORM.CurrentBalance < tt.Amount {
		return uuid.Nil, false, bank.ErrTransferTransactionPair
	}
	toAccountORM, err := service.db.GetBankAccountByAccountNumber(tt.ToAccount)
	if err != nil {
		log.Printf("Can't find transfer to account %v: %v\n", tt.ToAccount, err)
		return uuid.Nil, false, bank.ErrTransferDestincationAccountNotFound
	}

	fromTransactionORM := database.BankTransactionORM{
		TransactionUUID:      uuid.New(),
		AccountUUID:          fromAccountORM.AccountUUID,
		Amount:               tt.Amount,
		CreatedAt:            now,
		UpdatedAt:            now,
		TransactionTimestamp: now,
		TransactionType:      bank.TransactionTypeOut,
		Notes:                "Transfer out to " + tt.ToAccount,
	}

	toTransactionORM := database.BankTransactionORM{
		TransactionUUID:      uuid.New(),
		AccountUUID:          toAccountORM.AccountUUID,
		Amount:               tt.Amount,
		CreatedAt:            now,
		UpdatedAt:            now,
		TransactionTimestamp: now,
		TransactionType:      bank.TransactionTypeIn,
		Notes:                "Transfer in from " + tt.FromAccount,
	}

	//create transfer request
	newTransferUUID := uuid.New()
	transferORM := database.BankTransferORM{
		TransferUUID:    newTransferUUID,
		FromAccountUUID: fromAccountORM.AccountUUID,
		ToAccountUUID:   toAccountORM.AccountUUID,
		Currency:        tt.Currency,
		Amount:          tt.Amount,
		// TransferSuccess: ,
		TransferTimestamp: now,
		CreatedAt:         now,
		UpdatedAt:         now,
		TransferSuccess:   false,
	}

	if _, err := service.db.CreateTransfer(transferORM); err != nil {
		log.Printf("Can't create transfer from %v to %v\n", fromAccountORM.AccountNumber, toAccountORM.AccountNumber)
		return uuid.Nil, false, bank.ErrTransferRecordCreationFailed
	}

	if transferPairSuccess, err := service.db.CreateTransferTransactionPair(fromAccountORM, toAccountORM, fromTransactionORM, toTransactionORM); transferPairSuccess {
		service.db.UpdateTransferStatus(transferORM, true)
		return newTransferUUID, true, bank.ErrTransferTransactionPair
	} else {
		return newTransferUUID, false, err
	}

}
