package errors

const (
	DUPLICATE_DB_NAME        = "Database name already exists"
	UNKNOWN_QUERY_IDENTIFIER = "Invalid param in query: "
)

type DbError struct {
	message string
}

func (e *DbError) Error() string {
	return e.message
}

func NewDbError(message string) *DbError {
	return &DbError{message}
}
