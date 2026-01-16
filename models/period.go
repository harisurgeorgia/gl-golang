package models

import (
	"database/sql"
	"errors"
	"gl/db"
	"log"
	"time"
)

type PeriodStatus string

const (
	StatusClosed  PeriodStatus = "closed"
	StatusActive  PeriodStatus = "active"
	StatusPending PeriodStatus = "pending"
)

type Period struct {
	Id         *int64       `db:"id"`
	StartDate  time.Time    `db:"start_date"`
	EndDate    time.Time    `db:"end_date"`
	Status     PeriodStatus `db:"status"`
	PeriodName string       `db:"period_name"`
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
	)
	if err != nil {
		tx.Rollback()
		return err
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

func GetLastPeriodEndDate() (*time.Time, error) {
	var dt *time.Time
	err := db.Conn.QueryRow(`select end_date from general_ledger.periods order by id desc limit 1`).Scan(&dt)
	if err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil // not found
			}
			return nil, err
		}
	}
	return dt, nil
}

func SavePeriod(period Period) int64 {
	if period.Id != nil {
		return UpdatePeriod(period)
	}
	var id int64
	err := db.Conn.QueryRow(
		`INSERT INTO general_ledger.periods (period_name, start_date, end_date, status) 
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		period.PeriodName,
		period.StartDate,
		period.EndDate,
		period.Status,
	).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}
	return id
}

func UpdatePeriod(period Period) int64 {
	var id int64
	err := db.Conn.QueryRow(
		`UPDATE general_ledger.periods SET period_name = $1, start_date = $2, end_date = $3, status = $4 RETURNING id`,
		period.PeriodName,
		period.StartDate,
		period.EndDate,
		period.Status,
	).Scan(&id)

	if err != nil {
		log.Fatal(err)
	}
	return id
}

func GetPeriod(id int64) (*Period, error) {
	var period Period
	err := db.Conn.QueryRow(
		`SELECT id, period_name, start_date, end_date, status
		FROM general_ledger.periods
		WHERE id = $1`,
		id,
	).Scan(&period.Id, &period.PeriodName, &period.StartDate, &period.EndDate, &period.Status)
	if err != nil {
		return nil, err
	}
	return &period, nil
}

func GetPeriods() ([]Period, error) {
	var periods []Period
	rows, err := db.Conn.Query(`SELECT id, period_name, start_date, end_date, status FROM general_ledger.periods`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var period Period
		err := rows.Scan(&period.Id, &period.PeriodName, &period.StartDate, &period.EndDate, &period.Status)
		if err != nil {
			return nil, err
		}
		periods = append(periods, period)
	}
	return periods, nil
}

func CheckStatus(period Period) error {
	var count int
	if period.Status == StatusActive {

		err := db.Conn.QueryRow(
			`SELECT COUNT(*) FROM general_ledger.periods WHERE status = $1 and id < $2`,
			StatusActive, *period.Id,
		).Scan(&count)

		if err != nil {
			return err
		}
		tx, err := db.Conn.Begin()
		if err != nil {
			return err
		}
		if count > 0 {
			_, err = tx.Exec(`update general_ledger.periods set status = $1 where status = $2 and id < $3`, StatusClosed, StatusActive, *period.Id)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		_, err = tx.Exec(`update general_ledger.periods set status = $1 where status = $2 and id = $3`, StatusActive, StatusPending, *period.Id)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			`WITH last_period AS (
    SELECT id
    FROM general_ledger.periods
    WHERE id < $1
    ORDER BY id DESC
    LIMIT 1
	)INSERT INTO general_ledger.account_balances (period_id, account_id, balance)
		SELECT $1 AS new_period_id, a.id AS account_id,
  		COALESCE(ab.balance, 0) AS balance
		FROM general_ledger.accounts a
		LEFT JOIN general_ledger.account_balances ab 
  		ON a.id = ab.account_id
  		AND ab.period_id = (SELECT id FROM last_period);`,
			*period.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}

	}

	return nil
}
