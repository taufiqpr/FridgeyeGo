package routes

import (
	"FridgeEye-Go/services/profile/controllers"
	"FridgeEye-Go/services/profile/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(controllers.GetProfile))).Methods(http.MethodGet)
	r.Handle("/profile/{id}", middleware.AuthMiddleware(http.HandlerFunc(controllers.UpdateProfile))).Methods(http.MethodPut)
	r.Handle("/profile/{id}", middleware.AuthMiddleware(http.HandlerFunc(controllers.DeleteProfile))).Methods(http.MethodDelete)
	return r
}
