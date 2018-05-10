package omisego

import ()

type AdminAPI struct {
	Client
}

// func (a *AdminAPI) addDefaultHeaders(req *http.Request) {
// 	req.Header.Set("Authorization", a.auth.CreateAuthorizationHeader())
// 	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
// 	req.Header.Set("accept", "application/vnd.omisego.v1+json")
// }

/////////////////
// Session
/////////////////
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AdminAPI) Login(reqBody LoginParams) (Response, error) {
	req, err := a.newRequest("POST", "/login", reqBody)
	if err != nil {
		return Response{}, err
	}

	var res Response
	_, err = a.do(req, &res)
	return res, err

	// b := new(bytes.Buffer)
	// json.NewEncoder(b).Encode(reqBody)
	// req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"login", b)
	// a.addDefaultHeaders(req)
	// res, err := a.c.Do(req)
	// if err != nil {
	// 	return Response{}, err
	// }

	// resBody := Response{}
	// json.NewDecoder(res.Body).Decode(&resBody)
	// if !resBody.Success {
	// 	return resBody, errors.New("Unsuccessful request.")
	// }
	// return resBody, nil
}

func (a *AdminAPI) Logout() (Response, error) {
	req, err := a.newRequest("POST", "/logout", nil)
	if err != nil {
		return Response{}, err
	}

	var res Response
	_, err = a.do(req, &res)
	return res, err

	// req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"logout", nil)
	// a.addDefaultHeaders(req)
	// res, err := a.c.Do(req)
	// if err != nil {
	// 	return Response{}, err
	// }

	// resBody := Response{}
	// json.NewDecoder(res.Body).Decode(&resBody)
	// if !resBody.Success {
	// 	return resBody, errors.New("Unsuccessful request.")
	// }
	// return resBody, nil
}

/////////////////
// API Access
/////////////////
func (a *AdminAPI) AccessKeyCreate() (Response, error) {
	req, err := a.newRequest("POST", "/access_key.create", nil)
	if err != nil {
		return Response{}, err
	}

	var res Response
	_, err = a.do(req, &res)
	return res, err

	// req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"access_key.create", nil)
	// a.addDefaultHeaders(req)
	// res, err := a.c.Do(req)
	// if err != nil {
	// 	return Response{}, errors.New("Unsuccesful request")
	// }

	// resBody := Response{}
	// json.NewDecoder(res.Body).Decode(&resBody)
	// if !resBody.Success {
	// 	return resBody, errors.New("Unsuccesful request")
	// }
	// return resBody, nil
}
