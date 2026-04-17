package postgres

import "time"

type user struct {
	ID         string     `db:"id"`
	Username   string     `db:"username"`
	FirstName  string     `db:"firstname"`
	LastName   string     `db:"lastname"`
	MiddleName string     `db:"middlename"`
	Email      string     `db:"email"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	IsActive   bool       `db:"is_active"`
	LockedAt   *time.Time `db:"locked_at"`
}

type userShort struct {
	ID         string `db:"id"`
	Username   string `db:"username"`
	FirstName  string `db:"firstname"`
	LastName   string `db:"lastname"`
	MiddleName string `db:"middlename"`
	IsActive   bool   `db:"is_active"`
}
