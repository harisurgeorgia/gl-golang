package models

import (
	"gl/db"
	"log"
	"time"
)

type PeriodStatus string

const (
	StatusClosed PeriodStatus = "closed"
	StatusActive PeriodStatus = "active"
)

type Period struct {
	Id         *int          `db:"id"`
	StartDate  time.Time     `db:"start_date"`
	EndDate    time.Time     `db:"end_date"`
	Status     *PeriodStatus `db:"status"`
	PeriodName string        `db:"period_name"`
}

func ClosePeriod(startDate, endDate time.Time) error {
	var id int64
	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`update general_ledger.periods set closed = true`)
	if err != nil {
		tx.Rollback()
	}

	period_name := startDate.Format("2025-03-01") + " to " + endDate.Format("2025-03-30")
	err = tx.QueryRow(
		`INSERT INTO general_ledger.periods (start_date, end_date, period_name) VALUES ($1, $2, $3) RETURNING id`,
		startDate, endDate, period_name).Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		`WITH last_period AS (SELECT id FROM general_ledger.periods WHERE id < $1 ORDER BY id DESC LIMIT 1)
		INSERT INTO general_ledger.account_balances (period_id, account_id, closing_debit, closing_credit)
		SELECT $1 AS new_period_id, a.id AS account_id,
  		COALESCE(ab.closing_debit, 0) AS debit_balance,
  		COALESCE(ab.closing_credit, 0) AS credit_balance
		FROM general_ledger.accounts a
		LEFT JOIN general_ledger.account_balances ab 
  		ON a.id = ab.account_id
  		AND ab.period_id = (SELECT id FROM last_period);`,
		id)
	if err != nil {
		tx.Rollback()
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type MonthYear struct {
	Month int
	Year  int
}

func GetCurrentPeriod() MonthYear {
	var month_year MonthYear

	// Query the database to get the current open period's month and year
	// Assuming the periods table has a start_date column and closed is a boolean indicating if the period is closed
	// The query retrieves the most recent open period based on start_date) {

	db.Conn.QueryRow(
		`SELECT EXTRACT(MONTH FROM start_date)+1, EXTRACT(YEAR FROM start_date)
		FROM general_ledger.periods
		WHERE closed = false
		ORDER BY start_date DESC
		LIMIT 1`).Scan(&month_year.Month, &month_year.Year)
	return month_year
}

func GetLastPeriodEndDate() (*Period, error) {
	var p Period
	err := db.Conn.QueryRow(`select end_date from general_ledger.periods order by id desc limit 1`).Scan(&p.EndDate)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func SavePeriod(period Period) int64 {
	var id int64
	err := db.Conn.QueryRow(
		`INSERT INTO general_ledger.periods (period_name, start_date, end_date) 
		 VALUES ($1, $2, $3) RETURNING id`,
		period.PeriodName,
		period.StartDate,
		period.EndDate,
	).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}
	return id
}
