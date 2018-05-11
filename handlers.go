package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func UserAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	users, _ := env.DataStore.UserAll()

	json.NewEncoder(w).Encode(users)
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)
	user, _ := env.DataStore.UserGet(params["id"])
	json.NewEncoder(w).Encode(user)
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	newUser, _ := env.DataStore.UserCreate(&user)
	json.NewEncoder(w).Encode(newUser)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.Header().Set("Allow", "DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)
	_ = env.DataStore.UserDelete(params["id"])
	response := map[string]string{
		"status": "deleted",
	}
	json.NewEncoder(w).Encode(response)
}
