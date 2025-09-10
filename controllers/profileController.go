package controllers

import (
	"FridgeEye-Go/helper"
	userrepo "FridgeEye-Go/repository"
	"encoding/json"
	"net/http"
	"strconv"

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
		json.NewEncoder(w).Encode(helper.ErrInvalidPayload)
		return
	}

	emailCtx := r.Context().Value("email")
	if emailCtx == nil {
		helper.Error("UpdateProfile unauthorized request")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(helper.ErrUnauthorized)
		return
	}
	currentUserEmail := emailCtx.(string)

	owner, err := userrepo.GetUserID(userID)
	if err != nil {
		helper.Error("DB error on UpdateProfile: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}
	if owner == nil {
		helper.Error("UpdateProfile user not found id=" + userIDStr)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(helper.ErrUserNotFound)
		return
	}
	if owner.Email != currentUserEmail {
		helper.Error("UpdateProfile access denied for email: " + currentUserEmail)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(helper.ErrUnauthorized)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.Error("UpdateProfile invalid payload: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(helper.ErrInvalidPayload)
		return
	}

	name := req["name"]
	if name == "" {
		helper.Error("UpdateProfile no fields to update for user id=" + userIDStr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "no fields to update"})
		return
	}

	err = userrepo.UpdateUserName(userID, name)
	if err != nil {
		helper.Error("UpdateProfile failed DB update: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
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

	owner, err := userrepo.GetUserID(userID)
	if err != nil {
		helper.Error("DB error on DeleteProfile: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(helper.ErrDB)
		return
	}
	if owner == nil {
		helper.Error("DeleteProfile user not found id=" + userIDStr)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(helper.ErrUserNotFound)
		return
	}
	if owner.Email != currentUserEmail {
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
