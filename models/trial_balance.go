package models

import "gl/db"

type TrialBalanceRow struct {
	AccountCode string
	AccountName string
	Debit       float64
	Credit      float64
}

func GetTrialBalance(period int) ([]TrialBalanceRow, error) {
	var trailBalanceRows []TrialBalanceRow
	rows, err := db.Conn.Query(`SELECT
    a.account_code,
    a.account_name ,
    CASE WHEN ab.balance > 0 THEN ab.balance ELSE 0 END AS debit,
    CASE WHEN ab.balance < 0 THEN ABS(ab.balance) ELSE 0 END AS credit
FROM general_ledger.account_balances ab
JOIN general_ledger.accounts a ON a.id = ab.account_id
WHERE ab.period_id = $1
ORDER BY a.account_code;`, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row TrialBalanceRow
		err := rows.Scan(&row.AccountCode, &row.AccountName, &row.Debit, &row.Credit)
		if err != nil {
			return nil, err
		}
		trailBalanceRows = append(trailBalanceRows, row)
	}
	return trailBalanceRows, nil
}
