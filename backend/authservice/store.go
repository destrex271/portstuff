package main

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type store struct {
	Conn *sql.DB
}

func NewStore(connString string) (*store, error) {

	conn, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, err
	}

	return &store{
		Conn: conn,
	}, nil
}

func (s *store) Create(ctx context.Context, req User) error {
	s.Conn.Exec("INSERT INTO Users(username, passowrd, dob, email_id, mobile, userGroup) VALUES($1, $2, $3, $4, $5, $6)",
		req.Username, req.Password, req.Dob, req.EmailID, req.Mobile, req.UserGroup)
	return nil
}

func (s *store) Delete(ctx context.Context, userId int) error {
	_, err := s.Conn.Exec("DELETE FROM Users WHERE id = $1", userId)
	return err // Return the error if Exec fails
}

func (s *store) Update(ctx context.Context, req User) error {
	_, err := s.Conn.Exec("UPDATE Users SET username = $1, password = $2, dob = $3, email_id = $4, mobile = $5, userGroup = $6 WHERE id = $7",
		req.Username, req.Password, req.Dob, req.EmailID, req.Mobile, req.UserGroup, req.Id)
	return err // Return the error if Exec fails
}

func (s *store) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var user User

	// Query the database to find the user by username
	err := s.Conn.QueryRowContext(ctx, "SELECT id, username, password, dob, email_id, mobile, userGroup FROM Users WHERE username = $1", username).
		Scan(&user.Id, &user.Username, &user.Password, &user.Dob, &user.EmailID, &user.Mobile, &user.UserGroup)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errors.New("user not found")
		}
		return User{}, err // Return the error if there is a different failure
	}

	return user, nil
}
