package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/auth"
	"github.com/AlexWilliam12/silent-signal/internal/dtos"
	"github.com/AlexWilliam12/silent-signal/internal/models"
	"github.com/AlexWilliam12/silent-signal/internal/repositories"
	"github.com/AlexWilliam12/silent-signal/internal/services"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	var request dtos.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repositories.FindUserByCredentials(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auth.JWTToken{Token: token})
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {

	var request dtos.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := services.GenerateHash(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:        request.Username,
		Password:        request.Password,
		CredentialsHash: hash,
	}

	if _, err := repositories.CreateUser(&user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auth.JWTToken{Token: token})
}

func HandleValidateToken(w http.ResponseWriter, r *http.Request) {
	_, err := HandleAuthorization(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleValidateHash(w http.ResponseWriter, r *http.Request) {

	var request dtos.CredentialsHashRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repositories.FindUserByHash(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auth.JWTToken{Token: token})
}
