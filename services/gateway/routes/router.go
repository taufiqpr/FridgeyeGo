package routes

import (
	"FridgeEye-Go/services/gateway/config"
	"FridgeEye-Go/services/gateway/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

func newProxy(target string) *httputil.ReverseProxy {
	u, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(u)
}

func Router() *mux.Router {
	r := mux.NewRouter()

	// Middlewares
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimit(60, time.Minute))

	// Upstreams
	authProxy := newProxy(config.AppConfig.AuthURL)
	profileProxy := newProxy(config.AppConfig.ProfileURL)
	recipeProxy := newProxy(config.AppConfig.RecipeURL)

	// Auth routes (public)
	r.PathPrefix("/auth/").Handler(http.StripPrefix("", authProxy))
	r.PathPrefix("/register").Handler(http.StripPrefix("", authProxy))
	r.PathPrefix("/login").Handler(http.StripPrefix("", authProxy))

	// Protected routes via gateway JWT check
	protected := r.NewRoute().Subrouter()
	protected.Use(middleware.JWT)

	protected.PathPrefix("/profile").Handler(http.StripPrefix("", profileProxy))
	protected.PathPrefix("/get_recipes").Handler(http.StripPrefix("", recipeProxy))
	protected.PathPrefix("/get_recipe_detail").Handler(http.StripPrefix("", recipeProxy))

	// Health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

	return r
}
