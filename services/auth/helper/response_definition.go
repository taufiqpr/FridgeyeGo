package helper

var (
	ErrInvalidPayload = map[string]string{"error": "invalid payload"}
	ErrDB             = map[string]string{"error": "db error"}
	ErrToken          = map[string]string{"error": "token error"}
	ErrUnauthorized   = map[string]string{"error": "unauthorized"}
	ErrEmailAlreadyExist = map[string]string{"error": "email already exists"}
	ErrHashPassword      = map[string]string{"error": "failed to hash password"}
	ErrInsertUser        = map[string]string{"error": "failed to insert user"}
	ErrEmailNotFound     = map[string]string{"error": "email not found"}
	ErrWrongPassword     = map[string]string{"error": "wrong password"}
	ErrUserNotFound      = map[string]string{"error": "user not found"}
	ErrScan              = map[string]string{"error": "scan error"}
)
