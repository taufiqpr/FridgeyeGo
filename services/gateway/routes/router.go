package routes

import (
	"FridgeEye-Go/services/gateway/config"
	"FridgeEye-Go/services/gateway/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

func newProxy(target string) *httputil.ReverseProxy {
	u, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(u)
}

func Router() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.CORS)

	authProxy := newProxy(config.AppConfig.AuthURL)
	r.PathPrefix("/register").Handler(http.StripPrefix("", authProxy))
	r.PathPrefix("/login").Handler(http.StripPrefix("", authProxy))
	r.PathPrefix("/login-history").Handler(
		middleware.JWT(http.StripPrefix("", authProxy)),
	)

	profileProxy := newProxy(config.AppConfig.ProfileURL)
	r.PathPrefix("/profile").Handler(
		middleware.JWT(http.StripPrefix("", profileProxy)),
	)

	RecipeProxy := newProxy(config.AppConfig.RecipeURL)
	r.PathPrefix("/get_recipes").Handler(
		middleware.JWT(http.StripPrefix("", RecipeProxy)),
	)
	r.PathPrefix("/get_recipe_detail").Handler(
		middleware.JWT(http.StripPrefix("", RecipeProxy)),
	)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	return r
}
