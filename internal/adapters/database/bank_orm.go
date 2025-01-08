package database

import (
	"time"

	"github.com/google/uuid"
)

type BankAccountORM struct {
	AccountUUID    uuid.UUID `gorm:"primaryKey"`
	AccountNumber  string
	AccountName    string
	Currency       string
	CurrentBalance float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Transactions   []BankTransactionORM `gorm:"foreignKey:AccountUUID"`
}

func (bankAccORM BankAccountORM) TableName() string {
	return "bank_accounts"
}

type BankTransactionORM struct {
	TransactionUUID      uuid.UUID `gorm:"primaryKey"`
	AccountUUID          uuid.UUID
	TransactionType      string
	Amount               float64
	TransactionTimestamp time.Time
	Notes                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (BankTransactionORM BankTransactionORM) TableName() string {
	return "bank_transactions"
}

type BankExchangeRateORM struct {
	ExchangeRateUUID   uuid.UUID `gorm:"primaryKey"`
	FromCurrency       string
	ToCurrency         string
	Rate               float64
	ValidFromTimestamp time.Time
	ValidToTimestamp   time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (BankExchangeRateORM) TableName() string {
	return "bank_exchange_rates"
}

type BankTransferORM struct {
	TransferUUID      uuid.UUID `gorm:"primaryKey"`
	FromAccountUUID   uuid.UUID
	ToAccountUUID     uuid.UUID
	Currency          string
	Amount            float64
	TransferTimestamp time.Time
	TransferSuccess   bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (BankTransferORM) TableName() string {
	return "bank_transfers"
}
