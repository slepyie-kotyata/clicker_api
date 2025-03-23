package service

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	var password_byte = []byte(password)
	hashed_password, _ := bcrypt.GenerateFromPassword(password_byte, bcrypt.DefaultCost)
	return string(hashed_password)
}

func DoPasswordsMatch(hashed_password, current_password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(current_password))
	return err == nil
}