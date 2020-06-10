package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
)

const salt = "ymOcixH29JWua0684iEd"

func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s + salt))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RandomBase64() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func OTP(digit int) (string, error) {
	b := make([]byte, digit)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = '0' + byte(math.Floor(float64(b[i])/25.6))
	}
	return string(b), nil
}

const (
	ErrorWrongOTPCode           = "Wrong OTP Code"
	ErrorWrongLoginInfo         = "Wrong Email/Phonenumber or Password"
	ErrorIdentifierAlreadyExist = "Email/Phonenumber Already Registered on Your Company"
	ErrorNotUpdateRowNotExist   = "Cannot Update, Row is Not Existed"
)
