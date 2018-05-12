package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/goware/emailx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func UserAll(w http.ResponseWriter, r *http.Request) {
	users, err := env.DataStore.UserAll()
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := env.DataStore.UserGet(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	err = emailx.Validate(user.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "email: "+err.Error())
		return
	}
	byteEmail := []byte(user.Email)
	hash, err := bcrypt.GenerateFromPassword(byteEmail, bcrypt.MinCost)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	// TODO replace with smtp
	fmt.Println(string(hash))

	newUser, err := env.DataStore.UserCreate(&user)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, newUser)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := env.DataStore.UserDelete(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	response := map[string]bool{
		"success": true,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := env.DataStore.UserGet(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	json.NewDecoder(r.Body).Decode(user)
	user, err = env.DataStore.UserUpdate(user)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func EmailConfirm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	id := vars["id"]
	err := env.DataStore.UserActivate(id, token)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	response := map[string]bool{
		"success": true,
	}
	respondWithJSON(w, http.StatusOK, response)
}
