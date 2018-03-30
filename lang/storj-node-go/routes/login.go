package routes

import (
	"crypto/ecdsa"
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

type Config struct {
	key *ecdsa.PublicKey
}

const (
	TEST_KEY_PATH = "./key/test_ecdsa"
)

var (
	ErrTokenInvalid = errors.New("Token is not valid")
	ErrParsingKey   = errors.New("Unable to parse key")
	ErrParsingToken = errors.New("Unable to parse token")
	ErrSigningToken = errors.New("Error signing token")
)

func (t Token) GenerateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		Subject:   user.Username,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(1 * time.Hour.Unix()),
	})

	tokenString, err := token.SignedString([]byte(user.PubKey))
	if err != nil {
		return ErrSigningToken
	}

	return tokenString
}

func ValidateToken(encryptedToken string) (string, error) {
	tokData := regexp.MustCompile(`/s*$`).ReplaceAll([]byte(encryptedToken), []byte{})
	currentToken, err := jwt.Parse(string(tokData), func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return ErrParsingToken
	}

	fmt.Println("Current token ", currentToken)

	if !currentToken.Valid {
		return ErrTokenInvalid
	}

	formatToken, err := json.Marshal(currentToken.Claims)
	return string(tokData), err
}

func (c *Config) LoadKey(path string) error {
	pubKey, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	parsedKey, err := jtw.ParseECPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		return ErrParsingKey
	}

	c.key = parsedKey

	return nil
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
