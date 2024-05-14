package user

import (
	"database/sql"
	"time"

	"github.com/mamenzul/go-rest/types"
	"github.com/nrednav/cuid2"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	id := cuid2.Generate()
	_, err := s.db.Exec("INSERT INTO users ( id, email, password, created_at) VALUES (?, ?, ?, ?)", id, user.Email, user.Password, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) StoreResetToken(email string) (string, error) {
	token := cuid2.Generate()
	_, err := s.db.Exec("INSERT INTO password_resets (email, token, created_at, expires_at) VALUES (?, ?, ?, ?)", email, token, time.Now(), time.Now().Add(time.Minute*10))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Store) CheckResetToken(token string) (bool, error) {
	var email string
	err := s.db.QueryRow("SELECT email FROM password_resets WHERE token = ? AND expires_at > ?", token, time.Now()).Scan(&email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Store) DeleteResetToken(token string) error {
	_, err := s.db.Exec("DELETE FROM password_resets WHERE token = ?", token)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdatePassword(email string, password string) error {
	_, err := s.db.Exec("UPDATE users SET password = ? WHERE email = ?", password, email)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUsers() ([]types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users ")
	if err != nil {
		return nil, err
	}

	users := []types.User{}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}

	return users, nil
}
