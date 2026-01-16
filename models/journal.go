package models

import (
	"database/sql"
	"gl/db"
	"time"

	"github.com/shopspring/decimal"
)

type Journal struct {
	ID            *int64          `db:"id"`
	JournalNumber *string         `db:"journal_number"`
	JournalDate   time.Time       `db:"journal_date"`
	Description   string          `db:"description"` // nullable
	PeriodID      int64           `db:"period_id"`   // nullable FK
	Status        string          `db:"posted"`
	PostedBy      *int64          `db:"posted_by"` // nullable
	PostedAt      *time.Time      `db:"posted_at"` // nullable
	CreatedBy     int64           `db:"created_by"`
	CreatedAt     time.Time       `db:"created_at"`
	VerifiedAt    *time.Time      `db:"verified_at"`
	VerifiedBy    *string         `db:"verified_by"`
	Verified      bool            `db:"verified"`
	TtlDebit      decimal.Decimal `db:"ttl_debit"`
	TtlCredit     decimal.Decimal `db:"ttl_credit"`
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

func JournalSave(journal Journal, db *sql.DB) (*int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	// if journal.ID not nil delete journal form journal table
	if journal.ID != nil {
		_, err = tx.Exec(`DELETE FROM general_ledger.journals WHERE id = $1`, *journal.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = tx.Exec(`DELETE FROM general_ledger.journal_lines WHERE journal_id = $1`, *journal.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if journal.JournalNumber != nil {
		// use the provided journal_number
		err = db.QueryRow(`
        INSERT INTO general_ledger.journals
        (journal_date, journal_number, description, period_id, status, created_by)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, journal.JournalDate, *journal.JournalNumber, journal.Description, 1, "pending", journal.CreatedBy).Scan(&journal.ID)
	} else {
		// let the trigger generate it
		err = db.QueryRow(`
        INSERT INTO general_ledger.journals
        (journal_date, description, period_id, status, created_by)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, journal_number
    `, journal.JournalDate, journal.Description, 1, "pending", journal.CreatedBy).Scan(&journal.ID, &journal.JournalNumber)
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, line := range journal.Lines {
		_, err = tx.Exec(
			`insert into general_ledger.journal_lines (journal_id, account_id, debit, credit, line_description, line_number) values ($1, $2, $3, $4, $5, $6)
			`, journal.ID, line.AccountID, line.Debit, line.Credit, line.Description, line.LineNumber)
		if err != nil {
			tx.Rollback()
			return nil, sql.ErrConnDone
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return journal.ID, nil
}

// List all journals not posted
func GetJournals(db *sql.DB, filter string) ([]Journal, error) {
	rows, err := db.Query(`SELECT id, journal_number, journal_date, description, period_id, status, posted_by, posted_at, created_at FROM general_ledger.journals WHERE status = $1`, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var journal Journal
		if err := rows.Scan(&journal.ID, &journal.JournalNumber, &journal.JournalDate, &journal.Description, &journal.PeriodID, &journal.Status, &journal.PostedBy, &journal.PostedAt, &journal.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch journal lines for each journal

		journals = append(journals, journal)
	}

	return journals, nil
}

func GetJournalLines(journalID int64, db *sql.DB) ([]JournalLine, decimal.Decimal, decimal.Decimal, error) {
	rows, err := db.Query(`SELECT id, journal_id, account_id, debit, credit, line_description, line_number FROM general_ledger.journal_lines WHERE journal_id = $1`, journalID)
	if err != nil {
		return nil, decimal.Zero, decimal.Zero, err
	}
	defer rows.Close()

	var lines []JournalLine

	totalDebit := decimal.Zero
	totalCredit := decimal.Zero
	for rows.Next() {
		var line JournalLine
		if err := rows.Scan(&line.ID, &line.JournalID, &line.AccountID, &line.Debit, &line.Credit, &line.Description, &line.LineNumber); err != nil {
			return nil, decimal.Zero, decimal.Zero, err
		}
		totalDebit = totalDebit.Add(line.Debit)
		totalCredit = totalCredit.Add(line.Credit)
		lines = append(lines, line)
	}

	return lines, totalDebit, totalCredit, nil
}

// Get a journal by ID
func GetJournalByID(journalID int64, db *sql.DB) (*Journal, error) {
	var journal Journal
	err := db.QueryRow(`SELECT id, journal_number, journal_date, description, period_id, status, created_at, created_by FROM general_ledger.journals WHERE id = $1`, journalID).
		Scan(&journal.ID, &journal.JournalNumber, &journal.JournalDate, &journal.Description, &journal.PeriodID, &journal.Status, &journal.CreatedAt, &journal.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No journal found
		}
		return nil, err
	}

	lines, totalDebit, totalCredit, err := GetJournalLines(journalID, db)
	if err != nil {
		return nil, err
	}
	journal.Lines = lines
	journal.TtlDebit = totalDebit
	journal.TtlCredit = totalCredit

	return &journal, nil
}

func ListAllJournals(db *sql.DB) ([]Journal, error) {
	rows, err := db.Query(`SELECT id, journal_number, journal_date, description, period_id, status, posted_by, posted_at, created_at FROM general_ledger.journals ORDER BY journal_date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []Journal
	for rows.Next() {
		var journal Journal
		if err := rows.Scan(&journal.ID, &journal.JournalNumber, &journal.JournalDate, &journal.Description, &journal.PeriodID, &journal.Status, &journal.PostedBy, &journal.PostedAt, &journal.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch journal lines for each journal

		journals = append(journals, journal)
	}

	return journals, nil
}

func FindJournalByJournalNumber(s string) (int64, error) {
	var id int64
	err := db.Conn.QueryRow("SELECT id FROM general_ledger.journals WHERE journal_number = $1", s).Scan(&id)
	return id, err
}

func JournalUpdate(journal Journal, db *sql.DB) (*int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		`UPDATE general_ledger.journals
		 SET posted_at = $1,
		     posted_by = $2,
		     status    = $3
		 WHERE id = $4`,
		journal.PostedAt,
		journal.PostedBy,
		journal.Status,
		journal.ID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return journal.ID, nil
}
