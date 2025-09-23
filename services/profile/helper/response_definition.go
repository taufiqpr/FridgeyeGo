package helper

var (
	ErrInvalidPayload   = map[string]string{"error": "invalid payload"}
	ErrValidationFailed = map[string]string{"error": "validation failed"}
	ErrDB               = map[string]string{"error": "db error"}
	ErrUnauthorized     = map[string]string{"error": "unauthorized"}
	ErrUserNotFound     = map[string]string{"error": "user not found"}
)
