package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	ageStr := query.Get("age")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page := 1
	limit := 10
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	offset := (page - 1) * limit

	sqlQuery := "SELECT id, name, email, age FROM users WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if name != "" {
		sqlQuery += fmt.Sprintf(" AND name ILIKE $%d", argNum)
		args = append(args, "%"+name+"%")
		argNum++
	}
	if ageStr != "" {
		if age, err := strconv.Atoi(ageStr); err == nil {
			sqlQuery += fmt.Sprintf(" AND age = $%d", argNum)
			args = append(args, age)
			argNum++
		}
	}

	sqlQuery += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age)
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "bad_id", "Некорректный ID")
		return
	}

	var u User
	err = db.QueryRow("SELECT id, name, email, age FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Age)

	if err == sql.ErrNoRows {
		sendError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendError(w, http.StatusBadRequest, "bad_json", "Некорректный JSON")
		return
	}

	if msg := validateUser(u); msg != "" {
		sendError(w, http.StatusBadRequest, "validation", msg)
		return
	}

	err := db.QueryRow(
		"INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id",
		u.Name, u.Email, u.Age,
	).Scan(&u.ID)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "bad_id", "Некорректный ID")
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendError(w, http.StatusBadRequest, "bad_json", "Некорректный JSON")
		return
	}

	if msg := validateUser(u); msg != "" {
		sendError(w, http.StatusBadRequest, "validation", msg)
		return
	}

	result, err := db.Exec(
		"UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4",
		u.Name, u.Email, u.Age, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		sendError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
		return
	}

	u.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "bad_id", "Некорректный ID")
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		sendError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
