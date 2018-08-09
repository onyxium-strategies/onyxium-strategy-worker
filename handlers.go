package main

import (
	// "bitbucket.org/onyxium/onyxium-strategy-worker/email"
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
	"encoding/json"
	omg "github.com/Alainy/OmiseGo-Go-SDK"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"os"
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

func Login(w http.ResponseWriter, r *http.Request) {
	var userBody NewUserBody
	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	mgo := env.DataStore.(*models.MGO)

	c := mgo.DB(models.DatabaseName).C(models.UserCollection)
	user := &models.User{}
	err = c.Find(bson.M{"email": userBody.Email}).One(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if ok, err := models.ComparePasswords(user.Password, []byte(userBody.Password)); ok && err == nil {
		payload := map[string]string{
			"userId": user.Id.Hex(),
		}
		respondWithJSON(w, http.StatusOK, payload)
	} else {
		respondWithError(w, http.StatusBadRequest, "Incorrect password")
	}
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
	err = NewOMGUser(user)
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
	v := r.URL.Query()
	userId := v.Get("userId")
	var err error
	var strategies []models.Strategy
	if userId != "" {
		mgo := env.DataStore.(*models.MGO)
		c := mgo.DB(models.DatabaseName).C(models.StrategyCollection)
		err := c.Find(bson.M{"userId": bson.ObjectIdHex(userId)}).All(&strategies)
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
	} else {
		strategies, err = env.DataStore.StrategyAll()
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
	}

	payload := make([]map[string]interface{}, 0)
	for _, strategy := range strategies {
		payload = append(payload, map[string]interface{}{
			"id":        strategy.Id,
			"name":      strategy.Name,
			"status":    strategy.Status,
			"state":     strategy.State,
			"createdAt": strategy.CreatedAt.Unix(),
			"updatedAt": strategy.UpdatedAt.Unix(),
			"tree":      strategy.Tree.ToKaryArray(),
		})
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func StrategyGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	strategy, err := env.DataStore.StrategyGet(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	payload := map[string]interface{}{
		"id":        strategy.Id,
		"name":      strategy.Name,
		"status":    strategy.Status,
		"state":     strategy.State,
		"createdAt": strategy.CreatedAt.Unix(),
		"updatedAt": strategy.UpdatedAt.Unix(),
		"tree":      strategy.Tree.ToKaryArray(),
	}
	respondWithJSON(w, http.StatusOK, payload)
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

func BalancesGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	body := omg.ListByProviderUserIdParams{
		params["userId"],
		omg.ListParams{},
	}
	walletList, err := env.Ledger.UserGetWalletsByProviderUserId(body)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	var payload []omg.Balance
	for _, wallet := range walletList.Data {
		if wallet.Identifier == "primary" {
			// send initial funds of 10 BTC
			payload = wallet.Balances
		}
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func TransactionsGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var lastPage bool
	page := 1
	var payload []map[string]interface{}

	for !lastPage {
		body := omg.ListByProviderUserIdParams{
			params["userId"],
			omg.ListParams{
				PerPage: 100,
				Page:    page,
			},
		}
		transactionList, err := env.Ledger.UserGetTransactionsByProviderUserId(body)
		if err != nil {
			respondWithError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		lastPage = transactionList.IsLastPage

		for _, transaction := range transactionList.Data {
			// BUY
			if transaction.From.Address == os.Getenv("primaryWalletAddress") {
				payload = append(payload, map[string]interface{}{
					"Amount":        transaction.From.Amount,
					"Symbol":        transaction.From.Symbol,
					"SubunitToUnit": transaction.From.SubunitToUnit,
					"CreatedAt":     transaction.CreatedAt,
				})
			} else { // SELL
				payload = append(payload, map[string]interface{}{
					"Amount":        -transaction.To.Amount,
					"Symbol":        transaction.To.Symbol,
					"SubunitToUnit": transaction.To.SubunitToUnit,
					"CreatedAt":     transaction.CreatedAt,
				})
			}
		}
		page += 1
	}

	respondWithJSON(w, http.StatusOK, payload)
}
