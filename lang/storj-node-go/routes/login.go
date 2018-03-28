package routes

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type BoltDB interface {
	CreateUser(user *boltdb.User)
	GetUser(key []byte) (*boltdb.User, error)
	UpdateUser(user User)
	DeleteUser(key []byte)
}

// needs to be in req body
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Token struct {
	Token string
}

func (u *Users) HandleLogin(ctx iris.Context) {
	if ctx.Request().Body == nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	userRequest := &LoginUserRequest{}
	err := json.Unmarshal(ctx.Request().Body, userRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	savedUser, err := u.DB.GetUser(userRequest.Username)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	if userRequest.Password != savedUser.Password {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject:   userRequest.Username,
		NotBefore: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})

	tokenString, err := token.SignedString(server.signkey)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
}
