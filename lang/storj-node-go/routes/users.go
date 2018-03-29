package routes

import (
	"github.com/Storj/research/lang/storj-node-go/storage/boltdb"
	"github.com/google/uuid"
	"github.com/kataras/iris"
	"log"
)

// Users contains items needed to process requests to the user namespace
type Users struct {
	DB *boltdb.Client
}

func (u *Users) CreateUser(ctx iris.Context) {
	user := boltdb.User{
		Id:       uuid.New(),
		Username: ctx.Params().Get("id"),
		Email:    `dece@trali.zzd`,
	}

	if err := ctx.ReadJSON(user); err != nil {
		ctx.JSON(iris.StatusNotAcceptable)
	}

	if err := u.DB.CreateUser(user); err != nil {
		log.Println(err)
	} else {
		ctx.StatusCode(iris.StatusOK)
		ctx.HTML("<h1>User successfully created!</h1>")
	}
}

func (u *Users) GetUser(ctx iris.Context) {
	userId := ctx.Params().Get("id")
	userInfo, err := u.DB.GetUser(userId)
	if err != nil {
		log.Println(err)
	}

	ctx.Writef("%s's info is: %s", userId, userInfo)
}

// Updates only email for now
// Uses two db queries now, can refactor
func (u *Users) UpdateUser(ctx iris.Context) {
	userId := ctx.Params().Get("id")
	userInfo, err := u.DB.GetUser(userId)
	if err != nil {
		log.Println(err)
	}

	updated := boltdb.User{
		Id:       userInfo.Id,
		Username: userInfo.Username,
		Email:    ctx.Params().Get("email"),
	}

	err1 := u.DB.UpdateUser(updated)
	if err1 != nil {
		log.Println(err)
	}
}

func (u *Users) DeleteUser(ctx iris.Context) {
	userId := ctx.Params().Get("id")
	u.DB.DeleteUser(userId)
}
