package binary_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"GophKeeper/internal/server/model/binary"
	"GophKeeper/pkg/errs"
)

var (
	queryInsert = `INSERT INTO bin_data (meta, bytes) 
                   VALUES ($1, $2)`
	queryDelete = `DELETE FROM bin_data 
                   WHERE meta = $1`
	queryUpdate = `UPDATE bin_data
                   SET bytes = $1
                   WHERE meta = $2`
	queryGet = `SELECT bytes
                FROM bin_data 
                WHERE meta = $1`
)

type PostgresStorage struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresStorage - Создание хранилища в БД Postgres.
func NewPostgresStorage(db *sqlx.DB) *PostgresStorage {
	return &PostgresStorage{
		db:     db,
		logger: zap.L(),
	}
}

// Create Создание новых бинарных данных.
func (store *PostgresStorage) Create(data binary.DataFull) error {

	if _, err := store.db.ExecContext(context.Background(), queryInsert, data.MetaInfo, data.Bytes); err != nil {

		pqErr := err.(*pq.Error)
		if pqErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrAlreadyExist
		}

		err = fmt.Errorf("pg error on INSERT: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed create bin data", zap.Error(err))
		return err
	}
	return nil
}

// Delete Удаление бинарных данных.
func (store *PostgresStorage) Delete(in binary.DataGet) error {

	res, err := store.db.ExecContext(context.Background(), queryDelete, in.MetaInfo)
	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on DELETE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed delete bin data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Change Изменение бинарных данных.
func (store *PostgresStorage) Change(in binary.DataFull) error {

	res, err := store.db.ExecContext(context.Background(), queryUpdate, in.Bytes, in.MetaInfo)
	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on UPDATE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed update bin data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Get Получение бинарных данных по метаинформации.
func (store *PostgresStorage) Get(in binary.DataGet) (binary.DataFull, error) {

	row := store.db.QueryRowContext(context.Background(), queryGet, in.MetaInfo)

	var data []byte
	if err := row.Scan(&data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return binary.DataFull{}, errs.ErrNotFound
		}

		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on GET: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed get bin data", zap.Error(err))
		return binary.DataFull{}, err
	}

	return binary.DataFull{
		MetaInfo: in.MetaInfo,
		Bytes:    data,
	}, nil
}
