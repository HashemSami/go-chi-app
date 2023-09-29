package main

import (
	"fmt"

	"github.com/HashemSami/go-chi-app/models"
)

func main() {
	gs := models.GalleryService{}

	fmt.Println(gs.Images(2))
}
