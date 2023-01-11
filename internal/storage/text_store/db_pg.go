package text_store

import (
	"GophKeeper/internal/model/text"
	"GophKeeper/pkg/errs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	queryInsertText = `INSERT INTO text_data (meta, text) 
                       VALUES ($1, $2)`
	queryDeleteText = `DELETE FROM text_data 
                        WHERE meta = $1`
	queryUpdateText = `UPDATE text_data
                       SET text = $1
                       WHERE meta = $2`
	queryGetText = `SELECT text
                    FROM text_data 
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

// Create Создание новых текстовых данных.
func (store *PostgresStorage) Create(data text.DataTextFull) error {

	if _, err := store.db.ExecContext(context.Background(), queryInsertText, data.MetaInfo, data.Text); err != nil {

		pqErr := err.(*pq.Error)
		if pqErr.Code == "23505" {
			return errs.ErrAlreadyExist
		}

		err = fmt.Errorf("pg error on INSERT: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed create text data", zap.Error(err))
		return err
	}
	return nil
}

// Delete Удаление текстовых данных.
func (store *PostgresStorage) Delete(in text.DataTextGet) error {

	res, err := store.db.ExecContext(context.Background(), queryDeleteText, in.MetaInfo)
	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on DELETE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed delete text data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Change Изменение текстовых данных.
func (store *PostgresStorage) Change(in text.DataTextFull) error {

	res, err := store.db.ExecContext(context.Background(), queryUpdateText, in.Text, in.MetaInfo)
	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on UPDATE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed update text data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Get Получение текстовых данных по метаинформации.
func (store *PostgresStorage) Get(in text.DataTextGet) (text.DataTextFull, error) {

	row := store.db.QueryRowContext(context.Background(), queryGetText, in.MetaInfo)

	var data string
	if err := row.Scan(&data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return text.DataTextFull{}, errs.ErrNotFound
		}

		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on GET: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed get text data", zap.Error(err))
		return text.DataTextFull{}, err
	}

	return text.DataTextFull{
		MetaInfo: in.MetaInfo,
		Text:     data,
	}, nil
}
