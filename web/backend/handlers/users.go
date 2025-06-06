package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/krup4/goPetProject/web/backend/db/queries"
	"github.com/krup4/goPetProject/web/backend/others"
	requests "github.com/krup4/goPetProject/web/backend/request"
	"github.com/krup4/goPetProject/web/backend/response"
)

func UserSignUP(w http.ResponseWriter, r *http.Request) {
	var userRequest requests.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := queries.GetUserByLogin(userRequest.Login)

	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	} else if err == nil {
		http.Error(w, "User is already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := others.GenerateFromPassword(userRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	_, err = queries.CreateNewUser(userRequest.Login, hashedPassword, userRequest.Name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.StatusResponse{Message: "ok"}
	json.NewEncoder(w).Encode(resp)
}

func UserSignIn(w http.ResponseWriter, r *http.Request) {
	var userRequest requests.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := queries.GetUserByLogin(userRequest.Login)

	if err == pgx.ErrNoRows {
		http.Error(w, "User was not found", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	valid, err := others.ValidatePassword(userRequest.Password, user.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	if !valid {
		http.Error(w, "Password or login is incorrect", http.StatusBadRequest)
		return
	}

	token, err := others.GenerateToken(user.Login, os.Getenv("SECRET_KEY"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.AuthResponse{Token: token}
	json.NewEncoder(w).Encode(resp)
}
