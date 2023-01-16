package credential_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"GophKeeper/internal/model/cred"
	"GophKeeper/pkg/errs"
)

var (
	queryInsert = `INSERT INTO cred_data (meta, email, password_hash) 
                   VALUES ($1, $2, $3)`
	queryDelete = `DELETE FROM cred_data 
                   WHERE meta = $1`
	queryUpdate = `UPDATE cred_data
                   SET email = $1, password_hash = $2
                   WHERE meta = $3`
	queryGet = `SELECT email, password_hash
                FROM cred_data 
                WHERE meta = $1`
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

// Create Создание новых данных.
func (store *PostgresStorage) Create(data cred.CredentialFull) error {

	if _, err := store.db.ExecContext(context.Background(), queryInsert, data.MetaInfo, data.Email, data.Password); err != nil {

		pqErr := err.(*pq.Error)
		if pqErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrAlreadyExist
		}

		err = fmt.Errorf("pg error on INSERT: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed create cred data", zap.Error(err))
		return err
	}
	return nil
}

// Delete Удаление данных.
func (store *PostgresStorage) Delete(in cred.CredentialGet) error {

	res, err := store.db.ExecContext(context.Background(), queryDelete, in.MetaInfo)
	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on DELETE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed delete cred data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Change Изменение текстовых данных.
func (store *PostgresStorage) Change(in cred.CredentialFull) error {

	res, err := store.db.ExecContext(context.Background(), queryUpdate, in.Email, in.Password, in.MetaInfo)
	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on UPDATE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed update cred data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Get Получение текстовых данных по метаинформации.
func (store *PostgresStorage) Get(in cred.CredentialGet) (cred.CredentialFull, error) {

	row := store.db.QueryRowContext(context.Background(), queryGet, in.MetaInfo)

	var email string
	var pwd string

	if err := row.Scan(&email, &pwd); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cred.CredentialFull{}, errs.ErrNotFound
		}

		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on GET: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed get cred data", zap.Error(err))
		return cred.CredentialFull{}, err
	}

	return cred.CredentialFull{
		MetaInfo: in.MetaInfo,
		Email:    email,
		Password: pwd,
	}, nil
}
