package auth_store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"GophKeeper/internal/server/model/auth"
	"GophKeeper/pkg/errs"
)

var (
	queryCreate = `INSERT INTO users (email, password_hash) 
                   VALUES ($1, $2)`
	queryDelete = `DELETE FROM users 
                   WHERE id = $1`
	queryUpdate = `UPDATE users
                   SET password_hash = $1
                   WHERE id = $2`
	queryFind = `SELECT password_hash
                 FROM users 
                 WHERE email = $1`
	queryGetID = `SELECT id
                  FROM users 
                  WHERE email = $1`
)

type PostgresStorage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresStorage - Создание хранилища в БД Postgres
func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{
		db:     db,
		logger: zap.L(),
	}
}

// Create Создание нового пользователя.
func (store *PostgresStorage) Create(cred auth.Credential) error {

	if _, ok := store.userID(cred.Email); ok {
		return errs.ErrAlreadyExist
	}

	if _, err := store.db.ExecContext(context.Background(), queryCreate, cred.Email, cred.Password); err != nil {
		return err
	}
	return nil
}

// Delete Удаление пользователя
func (store *PostgresStorage) Delete(email string) error {

	userID, ok := store.userID(email)
	if !ok {
		return errs.ErrNotFound
	}

	if _, err := store.db.ExecContext(context.Background(), queryDelete, userID); err != nil {
		return err
	}
	return nil
}

// Update - Обновление пароля пользователя
func (store *PostgresStorage) Update(email, password string) error {

	userID, ok := store.userID(email)
	if !ok {
		return errs.ErrNotFound
	}

	if _, err := store.db.ExecContext(context.Background(), queryUpdate, password, userID); err != nil {
		return err
	}

	return nil
}

// Find - Поиск пользователя по имени и паролю
func (store *PostgresStorage) Find(cred auth.Credential) error {

	row := store.db.QueryRowContext(context.Background(), queryFind, cred.Email)

	var pwd string
	if err := row.Scan(&pwd); err != nil {
		return errs.ErrNotFound
	}

	if pwd != cred.Password {
		return ErrInvalidPassword
	}

	return nil
}

func (store *PostgresStorage) userID(email string) (int64, bool) {

	row := store.db.QueryRowContext(context.Background(), queryGetID, email)

	var userID int64
	if err := row.Scan(&userID); err != nil {
		return 0, false
	}

	return userID, true
}
