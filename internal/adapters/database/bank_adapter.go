package database

import "log"

func (adapter DatabaseAdapter) GetBankAccountByAccountNumber(acct string) (BankAccountORM, error) {

	var bankAccountORM BankAccountORM
	if err := adapter.db.First(&bankAccountORM, "account_number = ?", acct).Error; err != nil {
		log.Printf("Can't find bank account with number %v", acct)
		return bankAccountORM, err
	}
	return bankAccountORM, nil

}
