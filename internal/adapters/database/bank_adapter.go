package database

import (
	"log"
	"time"

	"github.com/VallabhSLEPAM/grpc-server/internal/application/domain/bank"
	uuid "github.com/google/uuid"
)

//DAO layer implementation

func (adapter DatabaseAdapter) GetBankAccountByAccountNumber(acct string) (BankAccountORM, error) {

	var bankAccountORM BankAccountORM
	if err := adapter.db.First(&bankAccountORM, "account_number = ?", acct).Error; err != nil {
		log.Printf("Can't find bank account with number %v", acct)
		return bankAccountORM, err
	}
	return bankAccountORM, nil

}

func (adapter DatabaseAdapter) CreateExchangeRate(r BankExchangeRateORM) (uuid.UUID, error) {
	if err := adapter.db.Create(r).Error; err != nil {
		return uuid.Nil, err
	}
	return r.ExchangeRateUUID, nil
}

func (adapter DatabaseAdapter) GetExchangeRateAtTimestamp(fromCurr, toCurr string, ts time.Time) (BankExchangeRateORM, error) {
	var bankExchangeRateORM BankExchangeRateORM
	err := adapter.db.First(&bankExchangeRateORM, "from_currency = ? AND to_currency = ? AND (? BETWEEN valid_from_timestamp and valid_to_timestamp)", fromCurr, toCurr, ts).Error
	return bankExchangeRateORM, err
}

func (adapter DatabaseAdapter) CreateTransactions(bankAccount BankAccountORM, bankTransaction BankTransactionORM) (uuid.UUID, error) {

	tx := adapter.db.Begin()
	if err := tx.Create(bankTransaction).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	//recalculate current balance
	newAmount := bankTransaction.Amount
	if bankTransaction.TransactionType == bank.TransactionTypeOut {
		newAmount = -1 * bankTransaction.Amount
	}

	newAccountBalance := bankAccount.CurrentBalance + newAmount
	if err := tx.Model(&bankAccount).Updates(
		map[string]interface{}{
			"current_balance": newAccountBalance,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	tx.Commit()
	return bankTransaction.TransactionUUID, nil
}

func (adapter DatabaseAdapter) CreateTransfer(bankTransferORM BankTransferORM) (uuid.UUID, error) {
	if err := adapter.db.Create(&bankTransferORM).Error; err != nil {
		return uuid.Nil, err
	}
	return bankTransferORM.TransferUUID, nil
}

func (adapter DatabaseAdapter) CreateTransferTransactionPair(fromAccount, toAccount BankAccountORM, fromTransactionORM, toTransactionORM BankTransactionORM) (bool, error) {
	tx := adapter.db.Begin()

	if err := tx.Create(&fromAccount).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Create(&toAccount).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	//recalculate current balance (fromAccount)
	fromAccountBalanceNew := fromAccount.CurrentBalance - fromTransactionORM.Amount
	if err := tx.Model(&fromAccount).Updates(
		map[string]interface{}{
			"current_balance": fromAccountBalanceNew,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	//recalculate current balance (toAccount)
	toAccountBalanceNew := toAccount.CurrentBalance + toTransactionORM.Amount
	if err := tx.Model(&toAccount).Updates(
		map[string]interface{}{
			"current_balance": toAccountBalanceNew,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (adapter DatabaseAdapter) UpdateTransferStatus(bankTransferORM BankTransferORM, status bool) error {

	if err := adapter.db.Model(&bankTransferORM).Updates(
		map[string]interface{}{
			"transfer_success": status,
			"updated_at":       time.Now(),
		},
	).Error; err != nil {
		return err
	}
	return nil

}
