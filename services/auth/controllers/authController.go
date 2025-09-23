package controllers

import (
	"FridgeEye-Go/services/auth/config"
	"FridgeEye-Go/services/auth/helper"
	"FridgeEye-Go/services/auth/models"
	q "FridgeEye-Go/services/auth/repository/db"
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
	err := config.DB.QueryRow(q.QueryCheckEmailExists, req.Email).Scan(&count)
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
	var user models.User
	err = config.DB.QueryRow(q.QueryRegister, req.Name, req.Email, hash).Scan(&user.ID, &user.CreatedAt)
	user.Name = req.Name
	user.Email = req.Email
	if err != nil {
		helper.Error("Failed to insert user: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrInsertUser)
		return
	}
	helper.Info(fmt.Sprintf("User registered: ID=%d, Email=%s", user.ID, user.Email))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
	var user models.User
	err := config.DB.QueryRow(q.QueryLogin, req.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		helper.Info("Login failed: email not found (" + req.Email + ")")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrEmailNotFound)
		return
	}
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		helper.Info("Login failed: wrong password (" + req.Email + ")")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrWrongPassword)
		return
	}
	tokenString, err := helper.GenerateToken(config.AppConfig.JWTSecret, user.ID, user.Name, user.Email, 30*time.Minute)
	if err != nil {
		helper.Error("Failed to generate token for user " + user.Email + ": " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrToken)
		return
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		helper.Error("Failed to parse IP from RemoteAddr")
	}
	_, err = config.DB.Exec(q.QueryInsertLoginHistory, user.Email, ip, r.UserAgent())
	if err != nil {
		helper.Error("Failed to insert login history for " + user.Email + ": " + err.Error())
	} else {
		helper.Info("Login history recorded for " + user.Email)
	}
	helper.Info(fmt.Sprintf("User logged in successfully: ID=%d, Email=%s, IP=%s", user.ID, user.Email, ip))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"token": tokenString})
}
