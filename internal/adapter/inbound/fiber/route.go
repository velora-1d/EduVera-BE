package fiber_inbound_adapter

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	inbound_port "eduvera/internal/port/inbound"
)

func InitRoute(
	ctx context.Context,
	app *fiber.App,
	port inbound_port.HttpPort,
) {
	// Enable CORS for frontend access
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://eduvera.ve-lora.my.id",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Internal routes (API key protected)
	internal := app.Group("/internal")
	internal.Use(func(c *fiber.Ctx) error {
		return port.Middleware().InternalAuth(c)
	})
	internal.Post("/client-upsert", func(c *fiber.Ctx) error {
		return port.Client().Upsert(c)
	})
	internal.Post("/client-find", func(c *fiber.Ctx) error {
		return port.Client().Find(c)
	})
	internal.Delete("/client-delete", func(c *fiber.Ctx) error {
		return port.Client().Delete(c)
	})

	// Protected API routes (JWT protected)
	v1Protected := app.Group("/v1")
	v1Protected.Use(func(c *fiber.Ctx) error {
		return port.Middleware().ClientAuth(c)
	})
	v1Protected.Get("/ping", func(c *fiber.Ctx) error {
		return port.Ping().GetResource(c)
	})

	// ========================================
	// PUBLIC API ROUTES (No Auth Required)
	// ========================================

	// API v1 - Onboarding (Public)
	api := app.Group("/api/v1")

	// Onboarding Routes
	onboarding := api.Group("/onboarding")
	onboarding.Post("/register", func(c *fiber.Ctx) error {
		return port.Onboarding().Register(c)
	})
	onboarding.Post("/institution", func(c *fiber.Ctx) error {
		return port.Onboarding().Institution(c)
	})
	onboarding.Post("/subdomain", func(c *fiber.Ctx) error {
		return port.Onboarding().Subdomain(c)
	})
	onboarding.Post("/bank-account", func(c *fiber.Ctx) error {
		return port.Onboarding().BankAccount(c)
	})
	onboarding.Post("/confirm", func(c *fiber.Ctx) error {
		return port.Onboarding().Confirm(c)
	})
	onboarding.Get("/status/:id", func(c *fiber.Ctx) error {
		return port.Onboarding().Status(c)
	})

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error {
		return port.Auth().Login(c)
	})
	auth.Get("/me", func(c *fiber.Ctx) error {
		return port.Auth().Me(c)
	})
	auth.Post("/refresh", func(c *fiber.Ctx) error {
		return port.Auth().Refresh(c)
	})
	auth.Post("/logout", func(c *fiber.Ctx) error {
		return port.Auth().Logout(c)
	})

	// Payment Routes (Midtrans)
	payment := api.Group("/payment")
	payment.Post("/create", func(c *fiber.Ctx) error {
		return port.Payment().CreateTransaction(c)
	})
	payment.Post("/webhook", func(c *fiber.Ctx) error {
		return port.Payment().Webhook(c)
	})
	payment.Get("/status/:order_id", func(c *fiber.Ctx) error {
		return port.Payment().GetStatus(c)
	})

	// ========================================
	// LEGACY ROUTES (Keep for backward compatibility)
	// ========================================

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return port.Landing().Home(c)
	})

	// Legacy onboarding (HTML-based - to be deprecated)
	app.Get("/register", func(c *fiber.Ctx) error {
		return port.Landing().Register(c)
	})
	app.Post("/register", func(c *fiber.Ctx) error {
		return port.Landing().RegisterProcess(c)
	})
	app.Get("/onboarding/step-2", func(c *fiber.Ctx) error {
		return port.Landing().Step2(c)
	})
	app.Post("/onboarding/step-2", func(c *fiber.Ctx) error {
		return port.Landing().Step2Process(c)
	})
	app.Get("/onboarding/step-3", func(c *fiber.Ctx) error {
		return port.Landing().Step3(c)
	})
	app.Post("/onboarding/step-3", func(c *fiber.Ctx) error {
		return port.Landing().Step3Process(c)
	})
	app.Get("/onboarding/step-4", func(c *fiber.Ctx) error {
		return port.Landing().Step4(c)
	})
	app.Post("/onboarding/step-4", func(c *fiber.Ctx) error {
		return port.Landing().Step4Process(c)
	})
	app.Get("/onboarding/step-5", func(c *fiber.Ctx) error {
		return port.Landing().Step5(c)
	})
	app.Post("/onboarding/step-5", func(c *fiber.Ctx) error {
		return port.Landing().Step5Process(c)
	})
	app.Get("/onboarding/step-6", func(c *fiber.Ctx) error {
		return port.Landing().Step6(c)
	})
	app.Post("/onboarding/step-6", func(c *fiber.Ctx) error {
		return port.Landing().Step6Process(c)
	})
	app.Get("/onboarding/step-7", func(c *fiber.Ctx) error {
		return port.Landing().Step7(c)
	})
	app.Post("/onboarding/step-7", func(c *fiber.Ctx) error {
		return port.Landing().Step7Process(c)
	})
}
