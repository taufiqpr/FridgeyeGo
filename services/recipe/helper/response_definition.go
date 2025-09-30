package helper

var (
	ErrRecipeQueryRequired = map[string]string{"error": "at least one query param is required"}
	ErrRecipeIDRequired    = map[string]string{"error": "recipe ID is required"}
	ErrRecipeFetchFailed   = map[string]string{"error": "failed to fetch data"}
	ErrRecipeGetFailed     = map[string]string{"error": "failed to get recipes"}
	ErrRecipeParseFailed   = map[string]string{"error": "failed to parse response"}
	ErrRecipeDetailFailed  = map[string]string{"error": "failed to get recipe detail"}
)
