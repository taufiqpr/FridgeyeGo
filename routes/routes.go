package routes

import (
	"FridgeEye-Go/controllers"
	"FridgeEye-Go/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.Hello).Methods(http.MethodGet)
	r.HandleFunc("/register", controllers.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	r.Handle("/login-history", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetLoginHistory))).Methods(http.MethodGet)
	r.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetProfile))).Methods(http.MethodGet)
	r.Handle("/profile/{id}", middleware.AuthMiddleware(http.HandlerFunc(controllers.UpdateProfile))).Methods(http.MethodPut)
	r.Handle("/profile/{id}", middleware.AuthMiddleware(http.HandlerFunc(controllers.DeleteProfile))).Methods(http.MethodDelete)
	r.Handle("/get_recipes", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetRecipes))).Methods(http.MethodGet)
	r.Handle("/get_recipe_detail", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetRecipeDetail))).Methods(http.MethodGet)
	return r
}
