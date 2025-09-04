package controllers

import (
	"FridgeEye-Go/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRecipes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Query is required"})
		return
	}

	url := fmt.Sprintf("%s/recipes/complexSearch?query=%s&number=10&addRecipeInformation=true&apiKey=%s",
		config.AppConfig.BaseURL, query, config.AppConfig.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch data"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get recipes"})
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse response"})
		return
	}

	results, ok := data["results"].([]interface{})
	if ok {
		for _, item := range results {
			m := item.(map[string]interface{})
			if m["image"] == nil || m["image"] == "" {
				if id, ok := m["id"].(float64); ok {
					m["image"] = fmt.Sprintf("https://spoonacular.com/recipeImages/%.0f-556x370.jpg", id)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func GetRecipeDetail(w http.ResponseWriter, r *http.Request) {
	recipeID := r.URL.Query().Get("id")
	if recipeID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Recipe ID is required"})
		return
	}

	url := fmt.Sprintf("%s/recipes/%s/information?apiKey=%s",
		config.AppConfig.BaseURL, recipeID, config.AppConfig.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch recipe detail"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get recipe detail"})
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse response"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}