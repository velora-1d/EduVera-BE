package main

import (
	"os"

	"prabogo/internal"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// @title EduVera API
// @version 1.0
// @description REST API for EduVera School Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@eduvera.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api-eduvera.ve-lora.my.id
// @BasePath /api/v1
// @schemes https http
func main() {
	app := internal.NewApp()
	option := "http"
	if len(os.Args) > 1 {
		option = os.Args[1]
	}

	app.Run(option)
}
