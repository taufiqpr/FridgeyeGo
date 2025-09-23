package controllers

import (
	"FridgeEye-Go/services/recipe/config"
	"FridgeEye-Go/services/recipe/helper"
	"FridgeEye-Go/services/recipe/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetRecipes(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	apiQuery := queryParams.Encode()
	if apiQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(helper.ErrRecipeQueryRequired)
		return
	}

	url := fmt.Sprintf("%s?%s&number=10&addRecipeInformation=true&apiKey=%s",
		config.AppConfig.SpoonacularSearchEndpoint, apiQuery, config.AppConfig.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrRecipeFetchFailed)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrRecipeGetFailed)
		return
	}

	var data models.RecipeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrRecipeParseFailed)
		return
	}

	for i, recipe := range data.Results {
		if recipe.Image == "" {
			data.Results[i].Image = fmt.Sprintf("https://spoonacular.com/recipeImages/%d-556x370.jpg", recipe.ID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Results)
}

func GetRecipeDetail(w http.ResponseWriter, r *http.Request) {
	recipeID := r.URL.Query().Get("id")
	if recipeID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(helper.ErrRecipeIDRequired)
		return
	}

	url := fmt.Sprintf("%s/%s/information?apiKey=%s",
		config.AppConfig.SpoonacularDetailEndpoint, recipeID, config.AppConfig.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrRecipeFetchFailed)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrRecipeDetailFailed)
		return
	}

	var data models.RecipeDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrRecipeParseFailed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
