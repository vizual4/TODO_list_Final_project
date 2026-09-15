package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "secret key for TODO-list 123"

type ReqPassword struct {
	Password string `json:"password"`
}

type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var reqPassword ReqPassword
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &reqPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	if len(pass) == 0 || reqPassword.Password != pass {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "wrong password"}, http.StatusBadRequest)
		return
	}

	token, err := generateToken(pass)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"token": token}, http.StatusOK)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			var valid bool
			if len(jwt) > 0 {
				valid = parseValidateToken(jwt)
			}

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func getPassHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateToken(password string) (string, error) {
	passHash := getPassHash(password)

	claims := Claims{
		Hash: passHash,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

func parseValidateToken(tokenStr string) bool {
	password := os.Getenv("TODO_PASSWORD")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	expectedHash := getPassHash(password)
	return claims.Hash == expectedHash
}
