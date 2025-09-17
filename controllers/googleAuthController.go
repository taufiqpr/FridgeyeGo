package controllers

import (
	"FridgeEye-Go/config"
	"FridgeEye-Go/helper"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

func InitGoogleOAuth() {
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  config.AppConfig.GoogleRedirectURL,
		ClientID:     config.AppConfig.GoogleClientID,
		ClientSecret: config.AppConfig.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if googleOAuthConfig == nil {
		http.Error(w, "Google OAuth not initialized", http.StatusInternalServerError)
		return
	}
	url := googleOAuthConfig.AuthCodeURL("random-state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if googleOAuthConfig == nil {
		http.Error(w, "Google OAuth not initialized", http.StatusInternalServerError)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}

	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	email := userData["email"].(string)
	name := userData["name"].(string)

	var id int
	err = config.DB.QueryRow("SELECT id FROM users WHERE email=$1", email).Scan(&id)

	if err == sql.ErrNoRows {
		err = config.DB.QueryRow(
			"INSERT INTO users (name, email, created_at) VALUES ($1, $2, NOW()) RETURNING id",
			name, email,
		).Scan(&id)
		if err != nil {
			http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString, err := helper.GenerateToken(
		config.AppConfig.JWTSecret,
		id,
		name,
		email,
		30*time.Minute, 
	)
	if err != nil {
		http.Error(w, "Failed to sign JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user":  userData,
	})

}
