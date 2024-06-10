package main

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       	int
	Username 	string
	Password	string
	Email    	string
	AESKey		string
	isLoggedIn	bool
	isAdmin		bool
}

func authenticateUser(username, password string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password, email, AESKey FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.AESKey)
	if (err != nil) {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if (err != nil) {
		return nil, err
	}
	return &user, nil
}

func getUserByID(userID int) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password, email, AESKey, is_admin FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.AESKey, &user.isAdmin)
	if err != nil {
        	return nil, err
	}
	return &user, nil
}
