package context

import (
	"context"

	"github.com/HashemSami/go-chi-app/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// user the point allows us to return a nil value when its not defined
func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)

	user, ok := val.(*models.User)
	if !ok {
		// the most likely case is that nothing was ever stored in the context,
		// so it doesn't have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the user key
		return nil
	}
	return user
}
