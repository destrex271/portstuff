package main

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	store AuthStore
}

func NewService(store AuthStore) *service {
	return &service{store}
}

func (s *service) CreateUser(ctx context.Context, user User) error {
	if user.Username == "" || user.Password == "" || user.EmailID == "" {
		return errors.New("username, password, and email_id are required")
	}

	err := s.store.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DeleteUser(ctx context.Context, userId int) error {
	if userId <= 0 {
		return errors.New("invalid user ID")
	}

	err := s.store.Delete(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, user User) error {
	if user.Id <= 0 || user.Username == "" || user.Password == "" || user.EmailID == "" {
		return errors.New("valid user ID, username, password, and email_id are required")
	}

	err := s.store.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func generateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"id":         user.Id,
		"username":   user.Username,
		"user_group": user.UserGroup,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key")) // Use a strong, secure secret key
}

func (s *service) LoginUser(ctx context.Context, username, password string) (string, error) {
	// Retrieve the user from the database by username
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT token upon successful authentication
	token, err := generateToken(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *service) HealthCheck(ctx context.Context) string {
	return "OK"
}
