package main

import (
	"strings"
)

func validateUser(u User) string {
	if strings.TrimSpace(u.Name) == "" {
		return "Имя не может быть пустым"
	}
	if !strings.Contains(u.Email, "@") {
		return "Некорректный email"
	}
	if u.Age < 0 || u.Age > 150 {
		return "Возраст должен быть от 0 до 150"
	}
	return ""
}
