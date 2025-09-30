package models

type Recipe struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type RecipeResponse struct {
	Results []Recipe `json:"results"`
}

type RecipeDetailResponse struct {
	ID                  int    `json:"id"`
	Title               string `json:"title"`
	Image               string `json:"image"`
	ReadyInMinutes      int    `json:"readyInMinutes"`
	Servings            int    `json:"servings"`
	Instructions        string `json:"instructions"`
	ExtendedIngredients []struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	} `json:"extendedIngredients"`
}
