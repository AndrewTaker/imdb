package api

import (
	"net/http"
)

type API struct {
	mux   *http.ServeMux
	users *UsersHandler
}

func NewApi(mux *http.ServeMux, users *UsersHandler) *API {
	api := &API{mux: mux, users: users}
	api.setUpRoutes()

	return api
}

func (a *API) setUpRoutes() {
	a.mux.HandleFunc("POST /auth/signin/", a.users.SignIn)
	a.mux.HandleFunc("POST /auth/signup/", a.users.SignUp)
}
