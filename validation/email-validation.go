package validation

import (
	"fmt"
	"gl/db"
	"net/mail"

	_ "github.com/lib/pq"
)

func EmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%q is invalid", email)
	}
	var exists string
	query := `SELECT email FROM general_ledger.users WHERE email = $1`
	row := db.Conn.QueryRow(query, email)
	err = row.Scan(&exists)
	if err != nil {
		return fmt.Errorf("no user found with email %q", email)
	}
	return nil
}
