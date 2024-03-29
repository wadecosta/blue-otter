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
}

func authenticateUser(username, password string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password, email FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.AESKey)
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
	err := db.QueryRow("SELECT id, username, email, password, AESKey FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.AESKey)
	if err != nil {
        	return nil, err
	}
	return &user, nil
}
