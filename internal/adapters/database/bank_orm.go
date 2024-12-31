package database

import (
	"time"

	uuid "github.com/gofrs/uuid"
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
