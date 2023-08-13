package main

import (
	stdCtx "context"
	"fmt"

	"github.com/HashemSami/go-chi-app/context"
	"github.com/HashemSami/go-chi-app/models"
)

func main() {
	ctx := stdCtx.Background()

	user := models.User{
		Email: "hash@hash.io",
	}

	ctx = context.WithUser(ctx, &user)

	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
}
