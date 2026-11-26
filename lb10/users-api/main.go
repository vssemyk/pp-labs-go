package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
			return
		}
		loginHandler(w, r)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
		switch r.Method {
		case http.MethodGet:
			getUsers(w, r)
		case http.MethodPost:
			createUser(w, r)
		default:
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			if !requireAdmin(w, r) {
				return
			}
		} else {
			if !requireAuth(w, r) {
				return
			}
		}
		switch r.Method {
		case http.MethodGet:
			getUser(w, r)
		case http.MethodPut:
			updateUser(w, r)
		case http.MethodDelete:
			deleteUser(w, r)
		default:
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
		}
	})

	http.HandleFunc("/session/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
			return
		}
		email := r.FormValue("email")
		var id int
		err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
		if err != nil {
			sendError(w, http.StatusUnauthorized, "auth", "Не найден")
			return
		}
		sid := createSession(id)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sid,
			Path:     "/",
			HttpOnly: true,
		})
		w.Write([]byte("Сессия создана"))
	})

	http.HandleFunc("/session/whoami", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			sendError(w, http.StatusUnauthorized, "auth", "Нет сессии")
			return
		}
		userID, ok := getSession(cookie.Value)
		if !ok {
			sendError(w, http.StatusUnauthorized, "auth", "Сессия истекла")
			return
		}
		w.Write([]byte("Ваш user_id: " + strconv.Itoa(userID)))
	})

	http.HandleFunc("/csrf-demo", csrfMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("CSRF-проверка пройдена"))
	}))

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
