package main

import (
	// "bitbucket.org/onyxium/onyxium-strategy-worker/email"
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

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

type NewUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	var userBody NewUserBody
	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := models.NewUser(userBody.Email, userBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = env.DataStore.UserCreate(user)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	// eventUserSignUp() UNCOMMENT ON MASTER

	// Email confirmation by user
	// err = email.EmailActivateUser(user.Email, user.Id.Hex())
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// TEST FAIL IF WE UNCOMMENT THIS BECAUSE OMGProvider is not authenticated during tests
	// Need to look at dependancy injection
	// For our local ledger we need to create a user and credit him starting money
	err = NewUser(user)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := env.DataStore.UserDelete(params["id"])
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
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
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	json.NewDecoder(r.Body).Decode(user)
	err = env.DataStore.UserUpdate(user)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func EmailConfirm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	id := vars["id"]
	id, _ = url.QueryUnescape(id)
	token, _ = url.QueryUnescape(token)
	err := env.DataStore.UserActivate(id, token)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	// eventUserActivated() UNCOMMENT ON MASTER
	response := map[string]bool{
		"success": true,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func StrategyAll(w http.ResponseWriter, r *http.Request) {
	strategies, err := env.DataStore.StrategyAll()
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, strategies)
}

func StrategyGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strategy, err := env.DataStore.StrategyGet(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, strategy)
}

func StrategyUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strategy, err := env.DataStore.StrategyGet(params["id"])
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	json.NewDecoder(r.Body).Decode(strategy)
	err = env.DataStore.StrategyUpdate(strategy)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, strategy)

}

func StrategyDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := env.DataStore.StrategyDelete(params["id"])
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	response := map[string]bool{
		"success": true,
	}
	respondWithJSON(w, http.StatusOK, response)
}
