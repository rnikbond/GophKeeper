package card_store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"GophKeeper/internal/model/card"
	"GophKeeper/pkg/errs"
)

var (
	queryInsert = `INSERT INTO card_data (meta, num, period_dt, cvv, full_name) 
                   VALUES ($1, $2, $3, $4, $5)`
	queryDelete = `DELETE FROM card_data 
                   WHERE meta = $1`
	queryUpdate = `UPDATE card_data
                   SET num = $1, period_dt = $2, cvv = $3, full_name = $4
                   WHERE meta = $5`
	queryGet = `SELECT num, period_dt, cvv, full_name
                FROM card_data 
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

// Create Создание новых данных банковской карты.
func (store *PostgresStorage) Create(data card.DataCardFull) error {

	if _, err := store.db.ExecContext(context.Background(),
		queryInsert,
		data.MetaInfo,
		data.Number,
		data.Period,
		data.CVV,
		data.FullName); err != nil {

		pqErr := err.(*pq.Error)
		if pqErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrAlreadyExist
		}

		err = fmt.Errorf("pg error on INSERT: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed create card data", zap.Error(err))
		return err
	}
	return nil
}

// Delete Удаление данных банковской карты.
func (store *PostgresStorage) Delete(in card.DataCardGet) error {

	res, err := store.db.ExecContext(context.Background(), queryDelete, in.MetaInfo)
	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on DELETE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed delete card data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Change Изменение текстовых данных.
func (store *PostgresStorage) Change(in card.DataCardFull) error {

	res, err := store.db.ExecContext(
		context.Background(),
		queryUpdate,
		in.Number,
		in.Period,
		in.CVV,
		in.FullName,
		in.MetaInfo)

	if err != nil {
		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on UPDATE: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed update card data", zap.Error(err))
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

// Get Получение данных анковской карты по метаинформации.
func (store *PostgresStorage) Get(in card.DataCardGet) (card.DataCardFull, error) {

	row := store.db.QueryRowContext(context.Background(), queryGet, in.MetaInfo)
	data := card.DataCardFull{
		MetaInfo: in.MetaInfo,
	}

	if err := row.Scan(&data.Number, &data.Period, &data.CVV, &data.FullName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return card.DataCardFull{}, errs.ErrNotFound
		}

		pqErr := err.(*pq.Error)
		err = fmt.Errorf("pg error on GET: %s. %v", pqErr.Code.Name(), err)
		store.logger.Error("failed get card data", zap.Error(err))
		return card.DataCardFull{}, err
	}

	return data, nil
}
