package models

import (
	"database/sql"
	"log"
)

type Account struct {
	Id          int64  `db:"id"`
	AccountCode string `db:"account_code"`
	AccountName string `db:"account_name"`
	AccountType string `db:"account_type"`
	IsActive    *bool  `db:"is_active"`
}

func GetAllAccounts(db *sql.DB) []Account {
	var accounts []Account
	rows, err := db.Query("SELECT id, account_code, account_name, account_type, is_active FROM general_ledger.accounts")
	if err != nil {
		log.Println("Query error:", err)
		return accounts
	}
	defer rows.Close()

	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.Id, &acc.AccountCode, &acc.AccountName, &acc.AccountType, &acc.IsActive)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		//log.Printf("Account Code %s", acc.AccountCode)
		accounts = append(accounts, acc)

	}
	return accounts
}
