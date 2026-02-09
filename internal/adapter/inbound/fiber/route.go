package fiber_inbound_adapter

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"prabogo/internal/domain"
	inbound_port "prabogo/internal/port/inbound"
	outbound_port "prabogo/internal/port/outbound"
)

func InitRoute(
	ctx context.Context,
	app *fiber.App,
	port inbound_port.HttpPort,
	d domain.Domain,
	dbPort outbound_port.DatabasePort,
) {
	// Enable CORS for frontend access with dynamic subdomain support
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// 1. Allow development origins if not in release mode
			if os.Getenv("APP_MODE") != "release" {
				if origin == "http://localhost:5173" || origin == "http://localhost:3000" {
					return true
				}
			}

			// 2. Allow main production domains
			if origin == "https://eduvera.ve-lora.my.id" || origin == "https://api-eduvera.ve-lora.my.id" {
				return true
			}

			// 3. Allow any subdomain of eduvera.ve-lora.my.id (HTTPS)
			if strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, ".eduvera.ve-lora.my.id") {
				return true
			}

			// 4. Temporary: Allow HTTP for testing while SSL propagates
			if strings.HasPrefix(origin, "http://") && strings.HasSuffix(origin, ".eduvera.ve-lora.my.id") {
				return true
			}

			return false
		},
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Sandbox-Tenant",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Rate Limiting (Public API Protection)
	// 60 requests per minute per IP for general endpoints
	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
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
	onboarding.Post("/check-subdomain", func(c *fiber.Ctx) error {
		return port.Onboarding().CheckSubdomain(c)
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

	// Auth Routes with stricter rate limiting
	auth := api.Group("/auth")

	// Stricter rate limit for login (5 req/min per IP) - brute force protection
	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak percobaan. Silakan coba lagi dalam 1 menit.",
			})
		},
	})

	auth.Post("/login", authLimiter, func(c *fiber.Ctx) error {
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
	// Stricter rate limit for forgot password (3 req/min) - prevent abuse
	forgotLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "forgot:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak permintaan reset password. Silakan coba lagi nanti.",
			})
		},
	})
	auth.Post("/forgot-password", forgotLimiter, func(c *fiber.Ctx) error {
		return port.Auth().ForgotPassword(c)
	})
	auth.Post("/reset-password", authLimiter, func(c *fiber.Ctx) error {
		return port.Auth().ResetPassword(c)
	})

	// Owner Routes
	owner := api.Group("/owner")
	owner.Post("/login", func(c *fiber.Ctx) error {
		return port.Owner().Login(c)
	})

	// Public Content Route
	publicApi := api.Group("/public")
	publicApi.Get("/content/:key", func(c *fiber.Ctx) error {
		return port.Content().Get(c)
	})

	// Landing Content (Dynamic)
	landingHandler := NewLandingContentHandler(d.LandingContent())
	publicApi.Get("/landing/:key", landingHandler.Get)

	// Protected Owner Routes
	ownerProtected := owner.Group("/")
	ownerProtected.Use(func(c *fiber.Ctx) error {
		return port.Middleware().OwnerAuth(c)
	})
	ownerProtected.Put("/landing/:key", landingHandler.Set)
	ownerProtected.Get("/tenants", func(c *fiber.Ctx) error {
		return port.Owner().GetTenants(c)
	})
	ownerProtected.Get("/tenants/:id", func(c *fiber.Ctx) error {
		return port.Owner().GetTenantDetail(c)
	})
	ownerProtected.Put("/tenants/:id/status", func(c *fiber.Ctx) error {
		return port.Owner().UpdateTenantStatus(c)
	})
	ownerProtected.Get("/stats", func(c *fiber.Ctx) error {
		return port.Owner().GetStats(c)
	})
	ownerProtected.Post("/impersonate", func(c *fiber.Ctx) error {
		return port.Owner().Impersonate(c)
	})
	ownerProtected.Post("/content", func(c *fiber.Ctx) error {
		return port.Content().Upsert(c)
	})
	ownerProtected.Post("/invoices/generate", func(c *fiber.Ctx) error {
		return port.SPP().GenerateManual(c)
	})
	ownerProtected.Post("/invoices/broadcast", func(c *fiber.Ctx) error {
		return port.SPP().BroadcastOverdueManual(c)
	})

	// Registration logs
	ownerProtected.Get("/registrations", func(c *fiber.Ctx) error {
		return port.Owner().GetRegistrations(c)
	})

	// SPP Transactions
	ownerProtected.Get("/transactions", func(c *fiber.Ctx) error {
		return port.Owner().GetSPPTransactions(c)
	})

	// Disbursements
	ownerProtected.Get("/disbursements", func(c *fiber.Ctx) error {
		return port.Owner().GetDisbursements(c)
	})
	ownerProtected.Post("/disbursements/:id/approve", func(c *fiber.Ctx) error {
		return port.Owner().ApproveDisbursement(c)
	})
	ownerProtected.Post("/disbursements/:id/reject", func(c *fiber.Ctx) error {
		return port.Owner().RejectDisbursement(c)
	})

	// Pesantren / Tenant Routes
	// Feature Gating: Only allow pesantren and hybrid plans
	// Pesantren / Tenant Routes
	// Feature Gating: Only allow pesantren and hybrid plans
	pesantren := api.Group("/pesantren")
	// ADDED: Authentication required before checking plan
	pesantren.Use(func(c *fiber.Ctx) error {
		return port.Middleware().ClientAuth(c)
	})
	pesantren.Use(CheckTrialStatus(d)) // Check trial status before allowing write operations
	pesantren.Use(RequirePlan("pesantren", "hybrid"))
	pesantren.Get("/dashboard/stats", func(c *fiber.Ctx) error {
		return port.PesantrenDashboard().GetStats(c)
	})

	// Tenant SPP Routes
	tenantSPP := pesantren.Group("/spp")
	tenantSPP.Get("/", func(c *fiber.Ctx) error {
		return port.SPP().List(c)
	})
	tenantSPP.Post("/", func(c *fiber.Ctx) error {
		return port.SPP().Create(c)
	})
	tenantSPP.Post("/:id/pay", func(c *fiber.Ctx) error {
		return port.SPP().RecordPayment(c)
	})
	tenantSPP.Get("/stats", func(c *fiber.Ctx) error {
		return port.SPP().GetStats(c)
	})
	// Manual payment confirmation routes
	tenantSPP.Get("/overdue", func(c *fiber.Ctx) error {
		return port.SPP().ListOverdue(c)
	})
	tenantSPP.Put("/:id", func(c *fiber.Ctx) error {
		return port.SPP().Update(c)
	})
	tenantSPP.Delete("/:id", func(c *fiber.Ctx) error {
		return port.SPP().Delete(c)
	})
	tenantSPP.Post("/:id/upload-proof", func(c *fiber.Ctx) error {
		return port.SPP().UploadProof(c)
	})
	tenantSPP.Post("/:id/confirm", func(c *fiber.Ctx) error {
		return port.SPP().ConfirmPayment(c)
	})

	// Notification logs
	ownerProtected.Get("/notifications", func(c *fiber.Ctx) error {
		return port.Owner().GetNotificationLogs(c)
	})

	// Owner WhatsApp Routes (Evolution API)
	ownerWAAdapter := NewOwnerWhatsAppAdapter(d)
	ownerWA := ownerProtected.Group("/whatsapp")
	ownerWA.Post("/connect", func(c *fiber.Ctx) error {
		return ownerWAAdapter.Connect(c)
	})
	ownerWA.Get("/status", func(c *fiber.Ctx) error {
		return ownerWAAdapter.GetStatus(c)
	})
	ownerWA.Post("/disconnect", func(c *fiber.Ctx) error {
		return ownerWAAdapter.Disconnect(c)
	})
	ownerWA.Post("/test", func(c *fiber.Ctx) error {
		return ownerWAAdapter.TestSend(c)
	})

	// Notification Template Routes (Owner only)
	templateAdapter := NewNotificationTemplateAdapter(dbPort)
	templates := ownerProtected.Group("/notification-templates")
	templates.Get("/", func(c *fiber.Ctx) error {
		return templateAdapter.List(c)
	})
	templates.Get("/:id", func(c *fiber.Ctx) error {
		return templateAdapter.Get(c)
	})
	templates.Post("/", func(c *fiber.Ctx) error {
		return templateAdapter.Create(c)
	})
	templates.Put("/:id", func(c *fiber.Ctx) error {
		return templateAdapter.Update(c)
	})
	templates.Delete("/:id", func(c *fiber.Ctx) error {
		return templateAdapter.Delete(c)
	})
	templates.Post("/:id/test", func(c *fiber.Ctx) error {
		return templateAdapter.TestSend(c)
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
	// SPP Payment Routes (Premium tier only)
	payment.Post("/spp/create", func(c *fiber.Ctx) error {
		return port.Payment().CreateSPPPayment(c)
	})
	payment.Post("/spp/webhook", func(c *fiber.Ctx) error {
		return port.Payment().SPPWebhook(c)
	})

	// Tenant WhatsApp Routes (Premium only)
	tenantWA := api.Group("/tenant/whatsapp")
	tenantWA.Use(func(c *fiber.Ctx) error {
		return port.Middleware().ClientAuth(c)
	})
	tenantWAAdapter := NewTenantWhatsAppAdapter(d)
	tenantWA.Post("/connect", tenantWAAdapter.Connect)
	tenantWA.Get("/status", tenantWAAdapter.Status)
	tenantWA.Post("/disconnect", tenantWAAdapter.Disconnect)
	tenantWA.Post("/test", tenantWAAdapter.SendTest)

	// Sekolah Routes (Protected)
	sekolah := api.Group("/sekolah")
	sekolah.Use(func(c *fiber.Ctx) error {
		return port.Middleware().ClientAuth(c)
	})
	sekolah.Use(CheckTrialStatus(d)) // Check trial status

	// Dashboard Stats
	sekolah.Get("/dashboard/stats", func(c *fiber.Ctx) error {
		return port.Sekolah().GetDashboardStats(c)
	})

	// Analytics Charts
	sekolah.Get("/dashboard/analytics", func(c *fiber.Ctx) error {
		return port.Analytics().GetAnalytics(c)
	})

	// ========================================
	// STUDENTS - Unified Siswa + Santri (NEW)
	// ========================================
	students := sekolah.Group("/students")
	students.Get("/", func(c *fiber.Ctx) error {
		return port.Student().List(c)
	})
	students.Get("/count", func(c *fiber.Ctx) error {
		return port.Student().Count(c)
	})
	students.Get("/:id", func(c *fiber.Ctx) error {
		return port.Student().Get(c)
	})
	// POST with 10-student limit for Trial/Basic tier
	students.Post("/", CheckDataLimit(d, "students", 10), func(c *fiber.Ctx) error {
		return port.Student().Create(c)
	})
	students.Put("/:id", func(c *fiber.Ctx) error {
		return port.Student().Update(c)
	})
	students.Delete("/:id", func(c *fiber.Ctx) error {
		return port.Student().Delete(c)
	})

	// Akademik
	akademik := sekolah.Group("/akademik")
	akademik.Get("/siswa", func(c *fiber.Ctx) error {
		return port.Sekolah().GetSiswaList(c)
	})
	akademik.Post("/siswa", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateSiswa(c)
	})
	akademik.Get("/guru", func(c *fiber.Ctx) error {
		return port.Sekolah().GetGuruList(c)
	})
	akademik.Post("/guru", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateGuru(c)
	})
	akademik.Get("/mapel", func(c *fiber.Ctx) error {
		return port.Sekolah().GetMapelList(c)
	})
	akademik.Get("/kelas", func(c *fiber.Ctx) error {
		return port.Sekolah().GetKelasList(c)
	})
	akademik.Post("/kelas", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateKelas(c)
	})

	// Kepesantrenan
	kepesantrenan := sekolah.Group("/kepesantrenan")
	kepesantrenan.Get("/aturan", func(c *fiber.Ctx) error {
		return port.Sekolah().GetPelanggaranAturanList(c)
	})
	kepesantrenan.Post("/aturan", func(c *fiber.Ctx) error {
		return port.Sekolah().CreatePelanggaranAturan(c)
	})
	kepesantrenan.Get("/pelanggaran", func(c *fiber.Ctx) error {
		return port.Sekolah().GetPelanggaranSiswaList(c)
	})
	kepesantrenan.Post("/pelanggaran", func(c *fiber.Ctx) error {
		return port.Sekolah().CreatePelanggaranSiswa(c)
	})
	kepesantrenan.Get("/perizinan", func(c *fiber.Ctx) error {
		return port.Sekolah().GetPerizinanList(c)
	})
	kepesantrenan.Post("/perizinan", func(c *fiber.Ctx) error {
		return port.Sekolah().CreatePerizinan(c)
	})

	// Asrama
	asrama := sekolah.Group("/asrama")
	asrama.Get("/gedung", func(c *fiber.Ctx) error {
		return port.Sekolah().GetAsramaList(c)
	})
	asrama.Post("/gedung", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateAsrama(c)
	})
	asrama.Get("/kamar", func(c *fiber.Ctx) error {
		return port.Sekolah().GetKamarList(c)
	})
	asrama.Post("/kamar", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateKamar(c)
	})
	asrama.Get("/penempatan", func(c *fiber.Ctx) error {
		return port.Sekolah().GetPenempatanList(c)
	})
	asrama.Post("/penempatan", func(c *fiber.Ctx) error {
		return port.Sekolah().CreatePenempatan(c)
	})

	// Tahfidz
	tahfidz := sekolah.Group("/tahfidz")
	tahfidz.Get("/setoran", func(c *fiber.Ctx) error {
		return port.Sekolah().GetTahfidzSetoranList(c)
	})
	tahfidz.Post("/setoran", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateTahfidzSetoran(c)
	})

	// Diniyah
	diniyah := sekolah.Group("/diniyah")
	diniyah.Get("/kitab", func(c *fiber.Ctx) error {
		return port.Sekolah().GetDiniyahKitabList(c)
	})
	diniyah.Post("/kitab", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateDiniyahKitab(c)
	})

	// Rapor & E-Rapor Routes
	erapor := sekolah.Group("/erapor")
	erapor.Get("/", func(c *fiber.Ctx) error {
		return port.Sekolah().GetRaporList(c)
	})
	erapor.Post("/", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateRapor(c)
	})

	// Tabungan
	tabungan := sekolah.Group("/tabungan")
	tabungan.Get("/", func(c *fiber.Ctx) error {
		return port.Sekolah().GetTabunganList(c)
	})
	tabungan.Post("/mutasi", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateTabunganMutasi(c)
	})

	// Kalender
	kalender := sekolah.Group("/kalender")
	kalender.Get("/", func(c *fiber.Ctx) error {
		return port.Sekolah().GetKalenderEvents(c)
	})
	kalender.Post("/", func(c *fiber.Ctx) error {
		return port.Sekolah().CreateKalenderEvent(c)
	})

	// Profil
	profil := sekolah.Group("/profil")
	profil.Get("/", func(c *fiber.Ctx) error {
		return port.Sekolah().GetProfil(c)
	})
	profil.Put("/", func(c *fiber.Ctx) error {
		return port.Sekolah().UpdateProfil(c)
	})

	// Laporan
	laporan := sekolah.Group("/laporan")
	laporan.Get("/", func(c *fiber.Ctx) error {
		return port.Sekolah().GetReportData(c)
	})

	// Subject Management
	erapor.Get("/subjects", func(c *fiber.Ctx) error {
		return port.ERapor().GetSubjects(c)
	})
	erapor.Post("/subjects", func(c *fiber.Ctx) error {
		return port.ERapor().CreateSubject(c)
	})
	erapor.Put("/subjects/:id", func(c *fiber.Ctx) error {
		return port.ERapor().UpdateSubject(c)
	})
	erapor.Delete("/subjects/:id", func(c *fiber.Ctx) error {
		return port.ERapor().DeleteSubject(c)
	})

	// Grade Management
	erapor.Post("/grades", func(c *fiber.Ctx) error {
		return port.ERapor().SaveGrade(c)
	})
	erapor.Post("/grades/batch", func(c *fiber.Ctx) error {
		return port.ERapor().BatchSaveGrades(c)
	})
	erapor.Get("/grades/student/:student_id", func(c *fiber.Ctx) error {
		return port.ERapor().GetStudentGrades(c)
	})
	erapor.Get("/grades/subject/:subject_id", func(c *fiber.Ctx) error {
		return port.ERapor().GetSubjectGrades(c)
	})

	// Rapor
	erapor.Get("/rapor/:student_id/:semester", func(c *fiber.Ctx) error {
		return port.ERapor().GetStudentRapor(c)
	})
	erapor.Post("/generate", func(c *fiber.Ctx) error {
		return port.ERapor().GenerateRapor(c)
	})

	// Stats
	erapor.Get("/stats", func(c *fiber.Ctx) error {
		return port.ERapor().GetStats(c)
	})

	// Curriculum Settings
	erapor.Get("/curriculum", func(c *fiber.Ctx) error {
		return port.ERapor().GetCurriculum(c)
	})
	erapor.Put("/curriculum", func(c *fiber.Ctx) error {
		return port.ERapor().SetCurriculum(c)
	})

	// Rapor History
	erapor.Get("/rapor/history", func(c *fiber.Ctx) error {
		return port.ERapor().GetRaporHistory(c)
	})

	// SDM Routes (Employee, Payroll, Attendance)
	sdm := sekolah.Group("/sdm")

	// Employee Management
	sdm.Get("/employees", func(c *fiber.Ctx) error {
		return port.SDM().GetEmployees(c)
	})
	sdm.Post("/employees", func(c *fiber.Ctx) error {
		return port.SDM().CreateEmployee(c)
	})
	sdm.Put("/employees/:id", func(c *fiber.Ctx) error {
		return port.SDM().UpdateEmployee(c)
	})
	sdm.Delete("/employees/:id", func(c *fiber.Ctx) error {
		return port.SDM().DeleteEmployee(c)
	})

	// Payroll
	sdm.Get("/payroll", func(c *fiber.Ctx) error {
		return port.SDM().GetPayrollByPeriod(c)
	})
	sdm.Post("/payroll/generate", func(c *fiber.Ctx) error {
		return port.SDM().GeneratePayroll(c)
	})
	sdm.Post("/payroll/:id/pay", func(c *fiber.Ctx) error {
		return port.SDM().MarkPayrollPaid(c)
	})
	sdm.Get("/payroll/:id/slip", func(c *fiber.Ctx) error {
		return port.SDM().GetPaySlip(c)
	})
	sdm.Get("/payroll/:id/slip/download", func(c *fiber.Ctx) error {
		return port.SDM().DownloadPaySlip(c)
	})
	sdm.Get("/payroll/config", func(c *fiber.Ctx) error {
		return port.SDM().GetPayrollConfig(c)
	})
	sdm.Put("/payroll/config", func(c *fiber.Ctx) error {
		return port.SDM().SavePayrollConfig(c)
	})

	// Subscription & Billing Routes
	sub := api.Group("/subscription")
	sub.Use(func(c *fiber.Ctx) error {
		return port.Middleware().ClientAuth(c)
	})
	sub.Get("/", func(c *fiber.Ctx) error {
		return port.Subscription().GetSubscription(c)
	})
	sub.Get("/pricing", func(c *fiber.Ctx) error {
		return port.Subscription().GetPricing(c)
	})
	sub.Post("/calculate-upgrade", func(c *fiber.Ctx) error {
		return port.Subscription().CalculateUpgrade(c)
	})
	sub.Post("/upgrade", func(c *fiber.Ctx) error {
		return port.Subscription().UpgradePlan(c)
	})
	sub.Post("/downgrade", func(c *fiber.Ctx) error {
		return port.Subscription().DowngradePlan(c)
	})

	// Attendance
	sdm.Get("/attendance", func(c *fiber.Ctx) error {
		return port.SDM().GetAttendance(c)
	})
	sdm.Post("/attendance", func(c *fiber.Ctx) error {
		return port.SDM().RecordAttendance(c)
	})
	sdm.Get("/attendance/summary", func(c *fiber.Ctx) error {
		return port.SDM().GetAttendanceSummary(c)
	})

	// Export Routes (Protected)
	export := api.Group("/export")
	export.Use(func(c *fiber.Ctx) error {
		return port.Middleware().ClientAuth(c)
	})
	export.Get("/students", func(c *fiber.Ctx) error {
		return port.Export().ExportStudents(c)
	})
	export.Get("/payments", func(c *fiber.Ctx) error {
		return port.Export().ExportPayments(c)
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
