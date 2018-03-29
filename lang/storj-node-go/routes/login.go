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
	Pubkey   string `json:"pubkey"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Pubkey   string `json:"pubkey"`
}

type Token struct {
	Token string
}

type Claim struct {
	Id  string
	exp int64
}

var (
	ErrTokenInvalid = errors.New("Token is not valid")
	secret          = "dried mango"
)

func (t Token) GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Subject:   username,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(1 * time.Hour.Unix()),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	return tokenString
}

func ValidateToken(encryptedToken string) (string, error) {
	tokData := regexp.MustCompile(`/s*$`).ReplaceAll([]byte(encryptedToken), []byte{})

	currentToken, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Couldn't parse token: ", err)
		return
	}

	fmt.Println("Current token ", currentToken)

	if !currentToken.Valid {
		return ErrTokenInvalid
	}

	formatToken, err := json.Marshal(currentToken.Claims)

	return string(tokData), err
}

func (u *Users) HandleLogin(ctx iris.Context) {
	var tok Token

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
		if err == boltdb.ErrNotFound {
			u.DB.CreateUser(userRequest)
		}
		return
	}

	if userRequest.Password != savedUser.Password {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	tokenString, err := tok.GenerateToken(savedUser.Username)
}
