package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	Uniqueviolation     = "23505"
)

var ErrorRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: Uniqueviolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
