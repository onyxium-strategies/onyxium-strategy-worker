package omisego

import ()

type AdminAPI struct {
	Client
}

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
}

func (a *AdminAPI) Logout() (Response, error) {
	req, err := a.newRequest("POST", "/logout", nil)
	if err != nil {
		return Response{}, err
	}

	var res Response
	_, err = a.do(req, &res)
	return res, err
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
}
