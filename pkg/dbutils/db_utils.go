package dbutils

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// txKey is an unexported type for storing the transaction in the context.
type txKey struct{}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func WithinTransactionReturning[T any](
	tm TxManager,
	ctx context.Context,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var result T

	err := tm.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		result, err = fn(txCtx)
		return err
	})

	return result, err
}

// ExtractTx is a helper to retrieve the pgx.Tx from the context.
func ExtractTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)

	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return tx, nil
}

func ConvertPgtypeTimestampToTime(ts pgtype.Timestamp) (time.Time, error) {
	if !ts.Valid {
		return time.Time{}, errors.New("timestamp is not valid")
	}
	return ts.Time, nil
}

func ConvertTimeToPgtypeTimestamp(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}
