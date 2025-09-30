package routes

import (
	"FridgeEye-Go/services/recipe/controllers"
	"FridgeEye-Go/services/recipe/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/get_recipes", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetRecipes))).Methods(http.MethodGet)
	r.Handle("/get_recipe_detail", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetRecipeDetail))).Methods(http.MethodGet)
	return r
}
