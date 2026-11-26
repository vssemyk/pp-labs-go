package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)

// ===== CSRF =====

func csrfMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			token := generateRandomString(16)
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    token,
				Path:     "/",
				HttpOnly: false,
			})
			next(w, r)
			return
		}

		cookie, err := r.Cookie("csrf_token")
		header := r.Header.Get("X-CSRF-Token")
		if err != nil || header == "" || cookie.Value != header {
			sendError(w, http.StatusForbidden, "csrf", "CSRF-токен неверен")
			return
		}
		next(w, r)
	}
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ===== СЕССИИ =====

type Session struct {
	UserID    int
	ExpiresAt time.Time
}

var (
	sessionsMu sync.Mutex
	sessions   = map[string]Session{}
)

func createSession(userID int) string {
	id := generateRandomString(16)
	sessionsMu.Lock()
	sessions[id] = Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	sessionsMu.Unlock()
	return id
}

func getSession(id string) (int, bool) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s, ok := sessions[id]
	if !ok || time.Now().After(s.ExpiresAt) {
		delete(sessions, id)
		return 0, false
	}
	return s.UserID, true
}
