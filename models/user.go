package models

type User struct {
	Id       *int64   `db:"id" form:"id"`
	Email    string  `db:"email" form:"email" binding:"required,email"`
	Fullname string  `db:"fullname" form:"fullname" binding:"required"`
	Password string  `db:"password" form:"password" binding:"required,min=6"`
	IsActive *bool   `db:"is_active" form:"is_active"`
	Role     *string `db:"role" form:"role"`
}
