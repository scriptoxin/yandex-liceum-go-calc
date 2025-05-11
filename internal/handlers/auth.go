package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/scriptoxin/yandex-liceum-go-calc/pkg/db"
	"github.com/scriptoxin/yandex-liceum-go-calc/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Register — POST /api/v1/register
func Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	_, err = db.Conn.Exec(
		"INSERT INTO users(login, password) VALUES(?, ?)",
		req.Login, string(hash),
	)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Login — POST /api/v1/login
func Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var id int
	var hash string
	err := db.Conn.QueryRow(
		"SELECT id, password FROM users WHERE login = ?",
		req.Login,
	).Scan(&id, &hash)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := jwt.Generate(id)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
