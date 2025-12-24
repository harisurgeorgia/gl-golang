package models

import (
	"database/sql"
)

type User struct {
	Id       *int64 `db:"id" form:"id"`
	Email    string `db:"email" form:"email" binding:"required,email"`
	Fullname string `db:"fullname" form:"fullname" binding:"required"`
	Password string `db:"password" form:"password" binding:"required,min=6"`
	IsActive *bool  `db:"is_active" form:"is_active"`
	Role     int    `db:"role" form:"role"`
}

type UserMenu struct {
	Id              int64
	MenuDescription string
	Url             string
	Icon            sql.NullString
	UserType        string
	ItemType        string
	Page            string
}

func GetUserMenu(db *sql.DB, page string) ([]UserMenu, error) {
	query := `SELECT id, menu_description, url, icon, role, item_type, page 
              FROM general_ledger.user_menu WHERE page = $1 order by ordered asc`

	rows, err := db.Query(query, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 👇 This is valid because UserMenu is already defined above
	var menus []UserMenu

	for rows.Next() {
		var um UserMenu
		err := rows.Scan(
			&um.Id,
			&um.MenuDescription,
			&um.Url,
			&um.Icon,
			&um.UserType,
			&um.ItemType,
			&um.Page,
		)
		if err != nil {
			return nil, err
		}
		menus = append(menus, um)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return menus, nil
}
