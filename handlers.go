package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/email"
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	omg "bitbucket.org/onyxium/onyxium-strategy-worker/omisego"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"os"
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

type NewUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	var newUser NewUserBody
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := models.NewUser(newUser.Email, newUser.Password)
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
	err = email.EmailActivateUser(user.Email, user.Id.Hex())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// For our local ledger we need to create a user and credit him starting money
	userCreateBody := omg.UserParams{
		ProviderUserId: user.Id.Hex(),
		Username:       user.Email,
	}
	_, err = env.ServerUser.UserCreate(userCreateBody)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	var tokenId string
	if baseTokenId != "" {
		tokenId = baseTokenId
	} else {
		tokenId = os.Getenv("baseTokenId")
	}

	creditBalanceBody := omg.BalanceAdjustmentParams{
		ProviderUserId: user.Id.Hex(),
		TokenId:        tokenId,
		Amount:         100,
	}

	_, err = env.ServerUser.UserCreditBalance(creditBalanceBody)
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
