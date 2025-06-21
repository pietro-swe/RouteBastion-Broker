/*
Package shared provide shared constants and functionality between modules
*/
package shared

type ErrorCode string

const (
	ErrCodeInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeEncryptionFailure ErrorCode = "ENCRYPTION_FAILURE"
	ErrCodeDatabaseFailure ErrorCode = "DATABASE_FAILURE"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
)

type DomainError struct {
	Code ErrorCode
	Msg  string
}

func (e DomainError) Error() string {
	return e.Msg
}

type ApplicationError struct {
	Code ErrorCode
	Msg  string
}

func (e ApplicationError) Error() string {
	return e.Msg
}

type InfrastructureError struct {
	Code ErrorCode
	Msg  string
}

func (e InfrastructureError) Error() string {
	return e.Msg
}
