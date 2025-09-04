package controllers

// import (
// 	"FridgeEye-Go/config"
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"net"
// 	"net/http"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/golang-jwt/jwt"
// 	"golang.org/x/crypto/bcrypt"
// 	"google.golang.org/api/idtoken"
// )

// type GoogleAuthRequest struct {
// 	IDToken string `json:"id_token"`
// }

// type GoogleAuthResponse struct {
// 	Status string `json:"status"`
// 	Token  string `json:"token"`
// 	User   struct {
// 		Email string `json:"email"`
// 		Name  string `json:"name"`
// 	} `json:"user"`
// }

// func GoogleAuth(w http.ResponseWriter, r *http.Request) {
// 	// Decode request
// 	var req GoogleAuthRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Validate Google ID token
// 	payload, err := idtoken.Validate(context.Background(), req.IDToken, config.AppConfig.GoogleClientID)
// 	if err != nil {
// 		http.Error(w, `{"error":"invalid google token"}`, http.StatusUnauthorized)
// 		return
// 	}

// 	// Extract user info
// 	email, _ := payload.Claims["email"].(string)
// 	name, _ := payload.Claims["name"].(string)
// 	if email == "" {
// 		http.Error(w, `{"error":"email not found"}`, http.StatusBadRequest)
// 		return
// 	}
// 	if name == "" {
// 		name = "No Name"
// 	}

// 	// Check if user exists
// 	var (
// 		id      int
// 		dbEmail string
// 		dbName  string
// 	)
// 	err = config.DB.QueryRow(
// 		`SELECT id, email, name 
// 		 FROM users 
// 		 WHERE email=$1 AND deleted_at IS NULL`,
// 		email,
// 	).Scan(&id, &dbEmail, &dbName)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// Create new user
// 			randomPass := "GOOGLE_" + uuid.NewString()
// 			hashedPass, _ := bcrypt.GenerateFromPassword([]byte(randomPass), bcrypt.DefaultCost)

// 			err = config.DB.QueryRow(
// 				`INSERT INTO users (email, name, password, created_at, is_verified)
// 				 VALUES ($1, $2, $3, NOW(), true)
// 				 RETURNING id, email, name`,
// 				email, name, string(hashedPass),
// 			).Scan(&id, &dbEmail, &dbName)
// 			if err != nil {
// 				http.Error(w, `{"error":"create user failed"}`, http.StatusInternalServerError)
// 				return
// 			}
// 		} else {
// 			http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// Save login history
// 	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
// 	_, _ = config.DB.Exec(
// 		`INSERT INTO login_history (user_email, ip_address, user_agent, timestamp, status)
// 		 VALUES ($1, $2, $3, NOW(), 'success')`,
// 		email, ip, r.UserAgent(),
// 	)

// 	// Generate JWT
// 	expirationTime := time.Now().Add(30 * time.Minute)
// 	claims := jwt.MapClaims{
// 		"sub":   id,
// 		"email": dbEmail,
// 		"name":  dbName,
// 		"exp":   expirationTime.Unix(),
// 		"iat":   time.Now().Unix(),
// 		"auth":  "google",
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
// 	if err != nil {
// 		http.Error(w, `{"error":"jwt error"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	// Response
// 	var resp GoogleAuthResponse
// 	resp.Status = "success"
// 	resp.Token = tokenString
// 	resp.User.Email = dbEmail
// 	resp.User.Name = dbName

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	_ = json.NewEncoder(w).Encode(resp)
// }
