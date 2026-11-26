package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my-secret-key-for-lab")

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
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

	role := "user"
	if id == 4 {
		role = "admin"
	}

	claims := Claims{
		UserID: id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "jwt", "Не удалось создать токен")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

func checkJWT(r *http.Request) (*Claims, bool) {
	h := r.Header.Get("Authorization")
	if len(h) < 8 || h[:7] != "Bearer " {
		return nil, false
	}
	tokenStr := h[7:]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	return claims, true
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if _, ok := checkJWT(r); !ok {
		sendError(w, http.StatusUnauthorized, "auth", "Требуется авторизация")
		return false
	}
	return true
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := checkJWT(r)
	if !ok {
		sendError(w, http.StatusUnauthorized, "auth", "Требуется авторизация")
		return false
	}
	if claims.Role != "admin" {
		sendError(w, http.StatusForbidden, "forbidden", "Только для администратора")
		return false
	}
	return true
}
