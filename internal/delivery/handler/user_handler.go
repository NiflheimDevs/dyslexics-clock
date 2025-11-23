package handler

import "net/http"


type UserHandler struct {

}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) Login(w http.Response, r *http.Request) {

}

func (uh *UserHandler) GetAlarms(w http.Response, r *http.Request) {

}

func (uh *UserHandler) GetColor(w http.Response, r *http.Request) {

}

func (uh *UserHandler) UpdateColor(w http.Response, r *http.Request) {
}

func (uh *UserHandler) UpdateAlarm(w http.Response, r *http.Request) {
}
