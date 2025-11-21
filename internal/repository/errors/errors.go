package errors

import "github.com/lib/pq"

const (
	UniqueViolation     pq.ErrorCode = "23505"
	ForeignKeyViolation pq.ErrorCode = "23503"
	NotNullViolation    pq.ErrorCode = "23502"
	CheckViolation      pq.ErrorCode = "23514"
)
