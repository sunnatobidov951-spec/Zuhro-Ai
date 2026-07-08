package domain

import (
	"time"
	"github.com/google/uuid"
)

type Role string
type Status string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	StatusActive Status = "active"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         Role
	Status       Status
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) CanLogin() bool {
	return u.Status == StatusActive
}

func (u *User) UpdateLastLogin() {
	now := time.Now().UTC()
	u.LastLogin = &now
}

