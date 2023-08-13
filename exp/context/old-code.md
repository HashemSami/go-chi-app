```go
package main

import (
	"context"
	"fmt"
	"strings"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, favoriteColorKey, "blue")

	value := ctx.Value(favoriteColorKey)

	// type assertion in go is similar to the generics feature
	intValue, ok := value.(int)
	if !ok {
		fmt.Println("it isn't an int")
	} else {
		fmt.Println(intValue + 4)
	}

	strValue, ok := value.(string)
	if !ok {
		fmt.Println("it isn't a string")
	} else {
		fmt.Println(strings.HasPrefix(strValue, "b"))
	}
}

```
