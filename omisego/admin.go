package omisego

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	// "io/ioutil"
	"net/http"
)

type AdminAPI struct {
	auth      Authorization
	c         *http.Client
	baseUrl   string
	serverUrl string
}

type LoginArgs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// {
//   "version": "1",
//   "success": false,
//   "data": {
//     "object": "error",
//     "code": "server:internal_server_error",
//     "description": "Something went wrong on the server",
//     "messages": {
//       "error_key": "error_reason"
//     }
//   }
// }

// {
//   "version": "1",
//   "success": true,
//   "data": {
//     "object": "list",
//     "data": [
//       {
//         "object": "minted_token",
//         "id": "ABC:ce3982f5-4a27-498d-a91b-7bb2e2a8d3d1",
//         "symbol": "ABC",
//         "name": "ABC Point",
//         "subunit_to_unit": 100,
//         "created_at": "2018-01-01T00:00:00Z",
//         "updated_at": "2018-01-01T10:00:00Z"
//       }
//     ],
//     "pagination": {
//       "per_page": 10,
//       "current_page": 1,
//       "is_first_page": true,
//       "is_last_page": true
//     }
//   }
// }

type Response struct {
	Version string                 `json:"version"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

func (a *AdminAPI) Login(reqBody LoginArgs) Response {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"login", b)

	req.Header.Set("Authorization", a.auth.CreateAuthorizationHeader())
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")

	res, err := a.c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resBody := Response{}
	json.NewDecoder(res.Body).Decode(&resBody)

	if !resBody.Success {
		log.Errorf("Unsuccesful request with error message: %s", resBody.Data["description"])
	}
	return resBody
}

func (a *AdminAPI) Logout() Response {
	req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"logout", nil)

	req.Header.Set("Authorization", a.auth.CreateAuthorizationHeader())
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")

	res, err := a.c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resBody := Response{}
	json.NewDecoder(res.Body).Decode(&resBody)

	if !resBody.Success {
		log.Errorf("Unsuccesful request with error message: %s", resBody.Data["description"])
	}
	return resBody
}
