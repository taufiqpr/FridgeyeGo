package controllers

import (
	"FridgeEye-Go/config"
	"FridgeEye-Go/models"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
		return
	}

	if err := validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", req.Email).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "db error"})
		return
	}
	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "email already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "hash error"})
		return
	}

	var id int
	var createdAt time.Time
	err = config.DB.QueryRow(
		`INSERT INTO users (name, email, password, created_at) 
		 VALUES ($1, $2, $3, NOW()) 
		 RETURNING id, created_at`,
		req.Name, req.Email, string(hash),
	).Scan(&id, &createdAt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "insert error"})
		return
	}

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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
		return
	}

	if err := validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var id int
	var name, email, password string
	err := config.DB.QueryRow(
		"SELECT id, name, email, password FROM users WHERE email=$1 AND deleted_at IS NULL",
		req.Email,
	).Scan(&id, &name, &email, &password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "email not registered"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "wrong password"})
		return
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := jwt.MapClaims{
		"sub":   id,
		"email": email,
		"name":  name,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "token error"})
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	_, err = config.DB.Exec(
		`INSERT INTO login_history (user_email, ip_address, user_agent, timestamp)
		 VALUES ($1, $2, $3, NOW())`,
		email,
		ip,
		r.UserAgent(),
	)
	if err != nil {

		fmt.Println("failed to insert login history:", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
	})
}

func GetLoginHistory(w http.ResponseWriter, r *http.Request) {
	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}
	currentUserEmail := emailCtx.(string)

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 AND deleted_at IS NULL)", currentUserEmail).Scan(&exists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "db error"})
		return
	}
	if !exists {

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}

	rows, err := config.DB.Query(
		`SELECT id, user_email, COALESCE(ip_address,''), 
		        COALESCE(user_agent,''), timestamp
		 FROM login_history
		 WHERE user_email=$1
		 ORDER BY timestamp DESC`,
		currentUserEmail,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "db error"})
		return
	}
	defer rows.Close()

	var history []models.LoginHistory
	for rows.Next() {
		var h models.LoginHistory
		if err := rows.Scan(&h.ID, &h.UserEmail, &h.IPAddress, &h.UserAgent, &h.Timestamp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "scan error"})
			return
		}
		history = append(history, h)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
