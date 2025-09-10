package userrepo

import (
	"FridgeEye-Go/config"
	"FridgeEye-Go/repository/db"
	"database/sql"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := config.DB.QueryRow(db.QueryGetUserByEmail, email).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserID(id int) (*User, error) {
	var user User
	err := config.DB.QueryRow(db.QueryGetUserByID, id).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
