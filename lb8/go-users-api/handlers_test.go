package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{"валидный", User{Name: "Иван", Email: "ivan@mail.ru", Age: 25}, false},
		{"пустое имя", User{Name: "", Email: "a@b.ru", Age: 20}, true},
		{"плохой email", User{Name: "A", Email: "bad", Age: 20}, true},
		{"отрицательный возраст", User{Name: "A", Email: "a@b.ru", Age: -1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := validateUser(tt.user)
			if (msg != "") != tt.wantErr {
				t.Errorf("validateUser() = %q, wantErr %v", msg, tt.wantErr)
			}
		})
	}
}

func TestCreateUserBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("не json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("код = %d, ожидался %d", w.Code, http.StatusBadRequest)
	}
}
