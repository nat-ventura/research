package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type ServerDependency func(*LoginServer) error

type BoltDB interface {
	CreateUser(user *boltdb.User)
	GetUser(key []byte) (*boltdb.User, error)
	UpdateUser(user User)
	DeleteUser(key []byte)
}

type LoginServer struct {
	signKey interface{}
	bolt    BoltDB
}

func (s *LoginServer) handleLogin(w http.ResponseWriter, req *http.Request) {
	
	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	type LoginUserRequest struct {
		Account string `json:"account"`
		Password string `json:"password"`
	}

	userRequest := &LoginUserRequest{}
	err := json.Unmarshal(body, userRequest)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}




