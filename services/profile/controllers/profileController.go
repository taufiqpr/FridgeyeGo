package controllers

import (
	"FridgeEye-Go/services/profile/helper"
	"FridgeEye-Go/services/profile/models"
	userrepo "FridgeEye-Go/services/profile/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		helper.Error("GetProfile unauthorized request")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrUnauthorized)
		return
	}

	currentUserEmail := emailCtx.(string)
	user, err := userrepo.GetUserByEmail(currentUserEmail)
	if err != nil {
		helper.Error("GetProfile DB error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}
	if user == nil {
		helper.Error("GetProfile user not found: " + currentUserEmail)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(helper.ErrUserNotFound)
		return
	}

	helper.Info("Profile fetched for email: " + currentUserEmail)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		helper.Error("UpdateProfile invalid user id: " + userIDStr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
		return
	}

	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		helper.Error("UpdateProfile unauthorized request")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}
	currentUserEmail := emailCtx.(string)

	targetUser, err := userrepo.GetUserID(userID)
	if err != nil {
		helper.Error("DB error on UpdateProfile: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}
	if targetUser == nil {
		helper.Error("UpdateProfile user not found id=" + userIDStr)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
		return
	}
	if targetUser.Email != currentUserEmail {
		helper.Error("UpdateProfile access denied for email: " + currentUserEmail)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "access denied"})
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.Error("UpdateProfile invalid payload: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		helper.Error("UpdateProfile validation failed: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "validation failed"})
		return
	}

	if err := userrepo.UpdateUserName(userID, req.Name); err != nil {
		helper.Error("UpdateProfile failed DB update: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	helper.Info("Profile updated for user id=" + userIDStr)
	json.NewEncoder(w).Encode(map[string]string{"message": "profile updated"})
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		helper.Error("DeleteProfile invalid user id: " + userIDStr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(helper.ErrInvalidPayload)
		return
	}

	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		helper.Error("DeleteProfile unauthorized request")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrUnauthorized)
		return
	}
	currentUserEmail := emailCtx.(string)

	targetUser, err := userrepo.GetUserID(userID)
	if err != nil {
		helper.Error("DB error on DeleteProfile: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}
	if targetUser == nil {
		helper.Error("DeleteProfile user not found id=" + userIDStr)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(helper.ErrUserNotFound)
		return
	}
	if targetUser.Email != currentUserEmail {
		helper.Error("DeleteProfile access denied for email: " + currentUserEmail)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(helper.ErrUnauthorized)
		return
	}

	err = userrepo.SoftDeleteUser(userID)
	if err != nil {
		helper.Error("DeleteProfile failed soft delete: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}

	_ = userrepo.DeleteLoginHistoryByEmail(currentUserEmail)

	helper.Info("Account deleted for email: " + currentUserEmail)

	json.NewEncoder(w).Encode(map[string]string{"message": "account deleted"})
}
