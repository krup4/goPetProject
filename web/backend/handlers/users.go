package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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
		log.Fatalln(err.Error())
		return
	} else if err == nil {
		http.Error(w, "User is already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := others.GenerateFromPassword(userRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err.Error())
		return
	}

	_, err = queries.CreateNewUser(userRequest.Login, hashedPassword, userRequest.Name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := response.StatusResponse{Message: "ok"}
	json.NewEncoder(w).Encode(resp)
}
