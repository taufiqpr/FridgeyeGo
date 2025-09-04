package controllers

import (
	"FridgeEye-Go/config"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}
	currentUserEmail := emailCtx.(string)

	var id int
	var name, email string
	err := config.DB.QueryRow(
		"SELECT id, name, email FROM users WHERE email=$1 AND deleted_at IS NULL",
		currentUserEmail,
	).Scan(&id, &name, &email)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    id,
		"name":  name,
		"email": email,
	})
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
		return
	}

	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}
	currentUserEmail := emailCtx.(string)

	var ownerEmail string
	err = config.DB.QueryRow("SELECT email FROM users WHERE id=$1 AND deleted_at IS NULL", userID).Scan(&ownerEmail)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	if ownerEmail != currentUserEmail {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "access denied"})
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
		return
	}

	name := req["name"]

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "no fields to update"})
		return
	}

	_, err = config.DB.Exec(
		"UPDATE users SET name=$1, updated_at=$2 WHERE id=$3",
		name, time.Now(), userID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "update error"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "profile updated"})
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
		return
	}

	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}
	currentUserEmail := emailCtx.(string)

	var ownerEmail string
	err = config.DB.QueryRow("SELECT email FROM users WHERE id=$1", userID).Scan(&ownerEmail)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	if ownerEmail != currentUserEmail {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "access denied"})
		return
	}

	_, err = config.DB.Exec("UPDATE users SET deleted_at=$1 WHERE id=$2", time.Now(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "delete error"})
		return
	}

	_, _ = config.DB.Exec("DELETE FROM login_history WHERE user_email=$1", currentUserEmail)

	json.NewEncoder(w).Encode(map[string]string{"message": "account deleted"})
}
