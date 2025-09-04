package controllers

import (
	"encoding/json"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "LANGSUNG SAJA HIT API NYA!!!"})
}