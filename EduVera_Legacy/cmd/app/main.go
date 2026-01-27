package main

import (
	"eduvera/internal/landing"
	"eduvera/internal/platform/database"
	"eduvera/internal/platform/prabogo"
	"log"
	"net/http"
)

func main() {
	// Koneksi Database (credentials from settings.json)
	database.Connect("postgres://eduvera:eduvera@localhost:5432/eduvera?sslmode=disable")

	// Inisialisasi Prabogo
	app := prabogo.New()

	// Serve Static Files
	app.Router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Route Landing Page
	app.GET("/{$}", landing.Home)
	app.GET("/register", landing.Register)
	app.GET("/onboarding/step-2", landing.OnboardingStep2)
	app.POST("/onboarding/step-2", landing.OnboardingStep2)
	app.GET("/onboarding/step-3", landing.OnboardingStep3)
	app.POST("/onboarding/step-3", landing.OnboardingStep3)
	app.GET("/onboarding/step-4", landing.OnboardingStep4)
	app.POST("/onboarding/step-4", landing.OnboardingStep4)
	app.GET("/onboarding/step-5", landing.OnboardingStep5)
	app.POST("/onboarding/step-5", landing.OnboardingStep5)

	// Jalankan Server
	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
