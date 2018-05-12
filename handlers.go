package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"encoding/json"
	"github.com/gorilla/mux"
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
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	users, err := env.DataStore.UserAll()
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)
	user, err := env.DataStore.UserGet(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	newUser, _ := env.DataStore.UserCreate(&user)
	respondWithJSON(w, http.StatusOK, newUser)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.Header().Set("Allow", "DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
	if r.Method != "PUT" {
		w.Header().Set("Allow", "PUT")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)
	user, err := env.DataStore.UserGet(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	json.NewDecoder(r.Body).Decode(user)
	err = env.DataStore.UserUpdate(user)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
