package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/krup4/goPetProject/web/backend/db"
	requests "github.com/krup4/goPetProject/web/backend/request"
	"github.com/krup4/goPetProject/web/backend/response"
)

func UserSignUP(w http.ResponseWriter, r *http.Request) {
	var userRequest requests.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user db.User
	err := db.Pool.QueryRow(r.Context(), `
		SELECT (id, login, password, name) FROM users WHERE login = $1
	`, userRequest.Login).Scan(&user)

	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == nil {
		http.Error(w, "User is already exists", http.StatusBadRequest)
		return
	}

	_, err = db.Pool.Exec(r.Context(), `
		INSERT INTO users (login, password, name) 
		VALUES ($1, $2, $3) ON CONFLICT DO NOTHING
	`, userRequest.Login, userRequest.Password, userRequest.Name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.StatusResponse{Message: "ok"}
	json.NewEncoder(w).Encode(resp)
}
