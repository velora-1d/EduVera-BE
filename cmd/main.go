package main

import (
	"os"

	"prabogo/internal"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	app := internal.NewApp()
	option := "http"
	if len(os.Args) > 1 {
		option = os.Args[1]
	}

	app.Run(option)
}
