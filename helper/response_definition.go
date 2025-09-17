package helper

var (
	ErrInvalidPayload    = map[string]string{"error": "invalid payload"}
	ErrValidationFailed  = map[string]string{"error": "validation failed"}
	ErrDB                = map[string]string{"error": "db error"}
	ErrEmailAlreadyExist = map[string]string{"error": "email already registered"}
	ErrEmailNotFound     = map[string]string{"error": "email not registered"}
	ErrWrongPassword     = map[string]string{"error": "wrong password"}
	ErrHashPassword      = map[string]string{"error": "hash error"}
	ErrInsertUser        = map[string]string{"error": "insert error"}
	ErrToken             = map[string]string{"error": "token error"}
	ErrUnauthorized      = map[string]string{"error": "unauthorized"}
	ErrUserNotFound      = map[string]string{"error": "user not found"}
	ErrScan              = map[string]string{"error": "scan error"}
	ErrRecipeQueryRequired = map[string]string{"error": "at least one query param is required"}
	ErrRecipeIDRequired    = map[string]string{"error": "recipe ID is required"}
	ErrRecipeFetchFailed   = map[string]string{"error": "failed to fetch data"}
	ErrRecipeGetFailed     = map[string]string{"error": "failed to get recipes"}
	ErrRecipeParseFailed   = map[string]string{"error": "failed to parse response"}
	ErrRecipeDetailFailed  = map[string]string{"error": "failed to get recipe detail"}
)
