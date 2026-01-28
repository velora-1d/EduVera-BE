package fiber_inbound_adapter

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	inbound_port "prabogo/internal/port/inbound"
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

	// Protected Owner Routes
	ownerProtected := owner.Group("/")
	ownerProtected.Use(func(c *fiber.Ctx) error {
		return port.Middleware().OwnerAuth(c)
	})
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
	ownerProtected.Post("/content", func(c *fiber.Ctx) error {
		return port.Content().Upsert(c)
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
	pesantren := api.Group("/pesantren")
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

	// Notification logs
	ownerProtected.Get("/notifications", func(c *fiber.Ctx) error {
		return port.Owner().GetNotificationLogs(c)
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

	// Sekolah Routes (Protected)
	sekolah := api.Group("/sekolah")
	sekolah.Use(func(c *fiber.Ctx) error {
		return port.Middleware().ClientAuth(c)
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

	// Stats
	erapor.Get("/stats", func(c *fiber.Ctx) error {
		return port.ERapor().GetStats(c)
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
