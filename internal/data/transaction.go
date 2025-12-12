package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ucok-man/tcsa/cmd/api/dto"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusFailed  TransactionStatus = "failed"
	TransactionStatusSucces  TransactionStatus = "success"
)

type Transaction struct {
	ID        int               `json:"id"`
	UserId    int               `json:"user_id"`
	Amount    int               `json:"amount"`
	Status    TransactionStatus `json:"status"`
	Version   int               `json:"-"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type TransactionModel struct {
	db *sql.DB
}

func (m TransactionModel) Insert(transaction *Transaction) error {
	query := `
        INSERT INTO transactions (user_id, amount, status)
        VALUES ($1, $2, $3)
        RETURNING id, version, created_at, updated_at`
	args := []any{transaction.UserId, transaction.Amount, transaction.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.db.QueryRowContext(ctx, query, args...).Scan(
		&transaction.ID,
		&transaction.Version,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
}

func (m TransactionModel) GetAll(dto dto.TransactionGetAllDTO) ([]*Transaction, Metadata, error) {
	if dto.Pagination.Limit == nil || dto.Pagination.Offset == nil {
		return nil, Metadata{}, fmt.Errorf("unset pagination limit or offset value")
	}

	if dto.Sort.RawValue != "" && (dto.Sort.Direction == nil || dto.Sort.ColumnValue == nil) {
		return nil, Metadata{}, fmt.Errorf("unset sort direction or columnValue")
	}

	query := fmt.Sprintf(`
	    SELECT count(*) OVER(), id, user_id, amount, status, version, created_at, updated_at
	    FROM transactions
	    WHERE status = $1 OR $1 = ''
	    ORDER BY %s %s, id ASC
	    LIMIT $2 OFFSET $3`, *dto.Sort.ColumnValue, *dto.Sort.Direction,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{dto.Filter.Status, dto.Pagination.Limit, dto.Pagination.Offset}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var totalRecords int
	var transactions []*Transaction

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&totalRecords, // count from window function
			&transaction.ID,
			&transaction.UserId,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Version,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, dto.Pagination.Page, dto.Pagination.PageSize)
	return transactions, metadata, nil
}

func (m TransactionModel) GetById(id int) (*Transaction, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, amount, status, version, created_at, updated_at
		FROM transactions
		WHERE id = $1`

	var transaction Transaction

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserId,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Version,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &transaction, nil
}

func (m TransactionModel) Update(transaction *Transaction) error {
	query := `
        UPDATE transactions
        SET amount = $1, status = $2, updated_at=$3, version = version + 1
        WHERE id = $4 AND version = $5
        RETURNING version`

	args := []any{
		&transaction.Amount,
		&transaction.Status,
		&transaction.UpdatedAt,
		&transaction.ID,
		&transaction.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, args...).Scan(&transaction.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m TransactionModel) DeleteOne(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM transactions WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
