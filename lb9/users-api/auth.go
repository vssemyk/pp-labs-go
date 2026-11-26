package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
)

var (
	tokensMu sync.Mutex
	tokens   = map[string]int{}
)

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "bad_json", "Некорректный JSON")
		return
	}
	if req.Email == "" {
		sendError(w, http.StatusBadRequest, "validation", "Email обязателен")
		return
	}

	var id int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&id)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "auth", "Пользователь с таким email не найден")
		return
	}

	token := generateToken()

	tokensMu.Lock()
	tokens[token] = id
	tokensMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func checkAuth(r *http.Request) (int, bool) {
	h := r.Header.Get("Authorization")
	if len(h) < 8 || h[:7] != "Bearer " {
		return 0, false
	}
	token := h[7:]

	tokensMu.Lock()
	id, ok := tokens[token]
	tokensMu.Unlock()

	return id, ok
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if _, ok := checkAuth(r); !ok {
		sendError(w, http.StatusUnauthorized, "auth", "Требуется авторизация")
		return false
	}
	return true
}
