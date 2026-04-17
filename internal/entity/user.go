package entity

import "time"

type UserFilterQuery struct {
	PageParam
	IsLocked *bool
	Search   string
}

type UserShort struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	IsActive bool   `json:"isActive"`
}

type User struct {
	ID         string     `json:"id"`
	Username   string     `json:"username"`
	Firstname  string     `json:"firstname"`
	Lastname   string     `json:"lastname"`
	Middlename string     `json:"middlename"`
	Email      string     `json:"email"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	IsActive   bool       `json:"isActive"`
	LockedAt   *time.Time `json:"lockedAt"`
}

type UserEdit struct {
	Username   string `json:"username"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Middlename string `json:"middlename"`
	Email      string `json:"email"`
}

// UsersPaged is for swagger only
type UsersPaged struct {
	PageInfo
	Items []*UserShort `json:"items"`
}
