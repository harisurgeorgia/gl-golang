package models

import (
	"gl/db"
	"log"

	"github.com/shopspring/decimal"
)

type GeneralLedgerRow struct {
	JournalNumber string
	Description   string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
}

type ClosingBalance struct {
	AccountCode string
	AccountName string
	Description string
	Debit       decimal.Decimal
	Credit      decimal.Decimal
}

func GetGeneralLedger(period, account_id int) ([]GeneralLedgerRow, *ClosingBalance, error) {
	var generalLedgerRows []GeneralLedgerRow
	rows, err := db.Conn.Query(`select je.journal_number, jl.line_description, jl.debit, jl.credit
	from general_ledger.journal_lines jl 
	left join general_ledger.accounts a on a.id = jl.account_id 
	left join general_ledger.journals je on je.id = jl.journal_id
	where jl.account_id = $1 and je.period_id = $2`, account_id, period)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row GeneralLedgerRow
		err := rows.Scan(&row.JournalNumber, &row.Description, &row.Debit, &row.Credit)
		if err != nil {
			return nil, nil, err
		}
		generalLedgerRows = append(generalLedgerRows, row)
	}
	closingBalance, err := getLastPeriodClosingBalance(period, account_id)
	if err != nil {
		return nil, nil, err
	}

	return generalLedgerRows, closingBalance, nil
}

func getLastPeriodClosingBalance(period int, accountID int) (*ClosingBalance, error) {
	var row ClosingBalance

	sql := `
	WITH last_period AS (
		SELECT id
		FROM general_ledger.periods
		WHERE id < $1
		ORDER BY id DESC
		LIMIT 1
	)
	SELECT
		a.account_code,
		a.account_name,
		'Closing Balance' AS description,
		CASE 
			WHEN COALESCE(ab.balance, 0) > 0 THEN ab.balance
			ELSE 0::numeric
		END AS debit,
		CASE 
			WHEN COALESCE(ab.balance, 0) < 0 THEN ABS(ab.balance)
			ELSE 0::numeric
		END AS credit
	FROM general_ledger.accounts a
	LEFT JOIN general_ledger.account_balances ab
		ON ab.account_id = a.id
		AND ab.period_id = (SELECT id FROM last_period)
	WHERE a.id = $2;
	`

	err := db.Conn.QueryRow(sql, period, accountID).Scan(
		&row.AccountCode,
		&row.AccountName,
		&row.Description,
		&row.Debit,
		&row.Credit,
	)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &row, nil
}
