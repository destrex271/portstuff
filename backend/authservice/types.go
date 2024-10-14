package main

import (
	"context"
	"time"
)

type AuthService interface {
	CreateUser(context.Context) error
}

type AuthStore interface {
	Create(context.Context, User) error
	Delete(context.Context, int) error
	Update(context.Context, User) error
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

type UserGroup int

const (
	ADMIN  = 0
	CLIENT = 1
	DRIVER = 2
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Dob       time.Time `json:"dob"`
	EmailID   string    `json:"email_id"`
	Mobile    string    `json:"mobile"`
	UserGroup UserGroup `json:"user_group"` // Can be Client, Driver OR Admin
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response returned after a successful login
type LoginResponse struct {
	Token string `json:"token"`
}
