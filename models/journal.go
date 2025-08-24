package models

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Journal struct {
	ID            *int64     `db:"id"`
	JournalNumber *string    `db:"journal_number"`
	JournalDate   time.Time  `db:"journal_date"`
	Description   string     `db:"description"` // nullable
	PeriodID      int64      `db:"period_id"`   // nullable FK
	Posted        bool       `db:"posted"`
	PostedBy      *string    `db:"posted_by"` // nullable
	PostedAt      *time.Time `db:"posted_at"` // nullable
	CreatedAt     time.Time  `db:"created_at"`
	VerifiedAt    *time.Time `db:"verified_at"`
	VerifiedBy    *string    `db:"verified_by"`
	Verified      bool       `db:"verified"`
	Lines         []JournalLine
}

type JournalLine struct {
	ID          *int64          `db:"id"`
	JournalID   int64           `db:"journal_id"`
	AccountID   int64           `db:"account_id"`
	Debit       decimal.Decimal `db:"debit"`
	Credit      decimal.Decimal `db:"credit"`
	Description string          `db:"line_description"`
	LineNumber  int             `db:"line_number"`
}

func JournalSave(journal Journal, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		`insert into general_ledger.journals (journal_date, description, period_id, posted, posted_by, posted_at) values ($1, $2, $3, $4, $5, $6) RETURNING id
		`, journal.JournalDate, journal.Description, 1, false, 1, nil).Scan(&journal.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, line := range journal.Lines {
		_, err = tx.Exec(
			`insert into general_ledger.journal_lines (journal_id, account_id, debit, credit, line_description, line_number) values ($1, $2, $3, $4, $5, $6)
			`, journal.ID, line.AccountID, line.Debit, line.Credit, line.Description, line.LineNumber)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// List all journals not posted
func GetPendingJournals(db *sql.DB) ([]Journal, error) {
	rows, err := db.Query(`SELECT id, journal_number, journal_date, description, period_id, posted, posted_by, posted_at, created_at FROM general_ledger.journals WHERE posted = false`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var journal Journal
		if err := rows.Scan(&journal.ID, &journal.JournalNumber, &journal.JournalDate, &journal.Description, &journal.PeriodID, &journal.Posted, &journal.PostedBy, &journal.PostedAt, &journal.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch journal lines for each journal

		journals = append(journals, journal)
	}

	return journals, nil
}
