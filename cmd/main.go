package main

import (
	"os"

	"eduvera/internal"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := internal.NewApp()
	option := "http"
	if len(os.Args) > 1 {
		option = os.Args[1]
	}

	app.Run(option)
}
