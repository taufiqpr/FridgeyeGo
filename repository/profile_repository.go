package userrepo

import (
	"FridgeEye-Go/config"
	"FridgeEye-Go/models"
	"FridgeEye-Go/repository/db"
	"database/sql"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.QueryRow(db.QueryGetUserByEmail, email).
		Scan(&user.ID, &user.Name, &user.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserID(id int) (*models.User, error) {
	var user models.User
	err := config.DB.QueryRow(db.QueryGetUserByID, id).
		Scan(&user.ID, &user.Name, &user.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserName(id int, name string) error {
	_, err := config.DB.Exec(db.QueryUpdateUser, name, id)
	return err
}

func SoftDeleteUser(id int) error {
	_, err := config.DB.Exec(db.QuerySoftDeleteUser, id)
	return err
}

func DeleteLoginHistoryByEmail(email string) error {
	_, err := config.DB.Exec(db.QueryDeleteLoginHistoryByEmail, email)
	return err
}
