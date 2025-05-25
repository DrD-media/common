package auth

import (
	"common/errors"
	"database/sql"
)

type UserRepositoryIMPL struct {
	db *sql.DB
}

type UserRepository interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
	GetByID(id int) (*User, error)
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryIMPL{db: db}
}

func (r *UserRepositoryIMPL) Create(user *User) error {
	query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

func (r *UserRepositoryIMPL) GetByUsername(username string) (*User, error) {
	query := "SELECT id, username, password, email, role FROM users WHERE username = ?"
	row := r.db.QueryRow(query, username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by username")
	}
	return user, nil
}

func (r *UserRepositoryIMPL) GetByID(id int) (*User, error) {
	query := "SELECT id, username, password, email, role FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	return user, nil
}
