package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/rigbyel/auth-server/internal/models"
)

type Storage struct {
	db *sql.DB
}

// creates new instance of Storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return &Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return &Storage{}, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// stops storage work
func (s *Storage) Stop() error {
	return s.db.Close()
}

// saves user in storage
func (s *Storage) SaveUser(u *models.User) (*models.User, error) {
	const op = "storage.sqlite.SaveUser"

	// prepare query
	stmt, err := s.db.Prepare(
		"INSERT INTO users (email, passHash) VALUES ($1, $2)",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// execute query
	res, err := stmt.Exec(u.Email, u.PassHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// get id of created user
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u.Id = id

	return u, nil
}

// gets user with the given email from storage
func (s *Storage) User(email string) (*models.User, error) {
	const op = "storage.sqlite.User"

	// get user's passhash from database
	row := s.db.QueryRow(
		"SELECT passHash FROM users WHERE email = $1",
		email,
	)

	var passHash []byte
	if err := row.Scan(&passHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// get user id from database
	row = s.db.QueryRow(
		"SELECT id FROM users WHERE email = $1",
		email,
	)

	var id int64
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// create models.User struct with the given email, id and passHash
	user := &models.User{
		Id:       id,
		Email:    email,
		PassHash: passHash,
	}

	return user, nil
}
