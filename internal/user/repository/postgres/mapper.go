package postgres

import (
	"avanpost-show/internal/entity"
	"strings"
)

func toUser(u user) *entity.User {
	return &entity.User{
		ID:         u.ID,
		Username:   u.Username,
		Firstname:  u.FirstName,
		Lastname:   u.LastName,
		Middlename: u.MiddleName,
		Email:      u.Email,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		IsActive:   u.IsActive,
		LockedAt:   u.LockedAt,
	}
}

func toUserShort(u userShort) *entity.UserShort {
	usr := &entity.UserShort{
		ID:       u.ID,
		Username: u.Username,
		IsActive: u.IsActive,
	}
	sb := strings.Builder{}
	sb.WriteString(u.FirstName)
	sb.WriteString(" ")
	if u.MiddleName != "" {
		sb.WriteString(u.MiddleName)
		sb.WriteString(" ")
	}
	sb.WriteString(u.LastName)
	usr.FullName = sb.String()
	return usr
}

func toUserShortList(items []userShort) []*entity.UserShort {
	users := make([]*entity.UserShort, 0, len(items))
	for _, u := range items {
		users = append(users, toUserShort(u))
	}
	return users
}
