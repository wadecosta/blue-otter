package main

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
	)

func authenticateVault(AESHash, AESKey string) (error) {
	fmt.Println("AESHash:", AESHash)
	fmt.Println("AESKey:", AESKey)
	err := bcrypt.CompareHashAndPassword([]byte(AESHash), []byte(AESKey))
	return err
}
