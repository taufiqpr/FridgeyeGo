package controllers

import (
	"FridgeEye-Go/config"
	"FridgeEye-Go/helper"
	"FridgeEye-Go/models"
	"FridgeEye-Go/repository/db"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.Error("Register payload invalid: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(helper.ErrInvalidPayload)
		return
	}

	if err := validate.Struct(req); err != nil {
		helper.Error("Register validation failed: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var count int
	err := config.DB.QueryRow(db.QueryCheckEmailExists, req.Email).Scan(&count)
	if err != nil {
		helper.Error("DB error on check email: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}
	if count > 0 {
		helper.Info("Register failed: email already exists (" + req.Email + ")")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(helper.ErrEmailAlreadyExist)
		return
	}

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		helper.Error("Failed to hash password: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrHashPassword)
		return
	}

	var id int
	var createdAt time.Time
	err = config.DB.QueryRow(db.QueryRegister, req.Name, req.Email, hash).
		Scan(&id, &createdAt)
	if err != nil {
		helper.Error("Failed to insert user: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrInsertUser)
		return
	}

	helper.Info(fmt.Sprintf("User registered: ID=%d, Email=%s", id, req.Email))

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"name":    req.Name,
		"email":   req.Email,
		"created": createdAt,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.Error("Login payload invalid: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(helper.ErrInvalidPayload)
		return
	}

	if err := validate.Struct(req); err != nil {
		helper.Error("Login validation failed: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var id int
	var name, email, password string
	err := config.DB.QueryRow(db.QueryLogin, req.Email).Scan(&id, &name, &email, &password)
	if err != nil {
		helper.Info("Login failed: email not found (" + req.Email + ")")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrEmailNotFound)
		return
	}

	if !helper.CheckPasswordHash(req.Password, password) {
		helper.Info("Login failed: wrong password (" + req.Email + ")")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrWrongPassword)
		return
	}

	tokenString, err := helper.GenerateToken(config.AppConfig.JWTSecret, id, name, email, 30*time.Minute)
	if err != nil {
		helper.Error("Failed to generate token for user " + email + ": " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrToken)
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	_, err = config.DB.Exec(db.QueryInsertLoginHistory, email, ip, r.UserAgent())
	if err != nil {
		helper.Error("Failed to insert login history for " + email + ": " + err.Error())
	} else {
		helper.Info("Login history recorded for " + email)
	}

	helper.Info(fmt.Sprintf("User logged in successfully: ID=%d, Email=%s, IP=%s", id, email, ip))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
	})
}

func GetLoginHistory(w http.ResponseWriter, r *http.Request) {
	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		helper.Info("Unauthorized access to login history (no email in context)")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrUnauthorized)
		return
	}
	currentUserEmail := emailCtx.(string)

	var exists bool
	err := config.DB.QueryRow(db.QueryUserExists, currentUserEmail).Scan(&exists)
	if err != nil {
		helper.Error("DB error while checking user exists for " + currentUserEmail + ": " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}
	if !exists {
		helper.Info("GetLoginHistory failed: user not found (" + currentUserEmail + ")")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(helper.ErrUserNotFound)
		return
	}

	rows, err := config.DB.Query(db.QueryLoginHistory, currentUserEmail)
	if err != nil {
		helper.Error("DB error while fetching login history for " + currentUserEmail + ": " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}
	defer rows.Close()

	var history []models.LoginHistory
	for rows.Next() {
		var h models.LoginHistory
		if err := rows.Scan(&h.ID, &h.UserEmail, &h.IPAddress, &h.UserAgent, &h.Timestamp); err != nil {
			helper.Error("Error scanning login history row for " + currentUserEmail + ": " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(helper.ErrScan)
			return
		}
		history = append(history, h)
	}

	helper.Info("Fetched login history for " + currentUserEmail)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
