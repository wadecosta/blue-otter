package main

import (
	"golang.org/x/crypto/bcrypt"
	)

func authenticateVault(AESHash, AESKey string) (error) {
	err := bcrypt.CompareHashAndPassword([]byte(AESHash), []byte(AESKey))
	return err
}
