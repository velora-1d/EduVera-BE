package internal

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	joonix "github.com/joonix/log"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"

	command_inbound_adapter "prabogo/internal/adapter/inbound/command"
	fiber_inbound_adapter "prabogo/internal/adapter/inbound/fiber"
	rabbitmq_inbound_adapter "prabogo/internal/adapter/inbound/rabbitmq"
	temporal_inbound_adapter "prabogo/internal/adapter/inbound/temporal"
	gibrun_outbound_adapter "prabogo/internal/adapter/outbound/gibrun"
	postgres_outbound_adapter "prabogo/internal/adapter/outbound/postgres"
	rabbitmq_outbound_adapter "prabogo/internal/adapter/outbound/rabbitmq"
	redis_outbound_adapter "prabogo/internal/adapter/outbound/redis"
	temporal_outbound_adapter "prabogo/internal/adapter/outbound/temporal"
	whatsapp_outbound_adapter "prabogo/internal/adapter/outbound/whatsapp"
	"prabogo/internal/adapter/outbound/whatsback"
	"prabogo/internal/domain"
	_ "prabogo/internal/migration/postgres"
	outbound_port "prabogo/internal/port/outbound"
	"prabogo/internal/scheduler"
	service_notification "prabogo/internal/service/notification"
	"prabogo/utils"
	"prabogo/utils/activity"
	"prabogo/utils/database"
	"prabogo/utils/gibrun"
	"prabogo/utils/log"
	"prabogo/utils/rabbitmq"
	"prabogo/utils/redis"
)

var databaseDriverList = []string{"postgres"}
var httpDriverList = []string{"fiber"}
var messageDriverList = []string{"rabbitmq", "whatsapp", "none"}
var workflowDriverList = []string{"temporal", "none"}
var outboundDatabaseDriver string
var outboundMessageDriver string
var outboundCacheDriver string
var outboundWorkflowDriver string
var inboundHttpDriver string
var inboundMessageDriver string
var inboundWorkflowDriver string

type App struct {
	ctx       context.Context
	domain    domain.Domain
	db        outbound_port.DatabasePort
	message   outbound_port.MessagePort
	scheduler *scheduler.Scheduler
	notif     *service_notification.NotificationService
}

func NewApp() *App {
	ctx := activity.NewContext("init")
	ctx = activity.WithClientID(ctx, "system")
	_ = godotenv.Load(".env")
	configureLogging()
	outboundDatabaseDriver = os.Getenv("OUTBOUND_DATABASE_DRIVER")
	outboundMessageDriver = os.Getenv("OUTBOUND_MESSAGE_DRIVER")
	outboundCacheDriver = os.Getenv("OUTBOUND_CACHE_DRIVER")
	outboundWorkflowDriver = os.Getenv("OUTBOUND_WORKFLOW_DRIVER")
	inboundHttpDriver = os.Getenv("INBOUND_HTTP_DRIVER")
	inboundMessageDriver = os.Getenv("INBOUND_MESSAGE_DRIVER")
	inboundWorkflowDriver = os.Getenv("INBOUND_WORKFLOW_DRIVER")

	dbPort := databaseOutbound(ctx)
	messagePort := messageOutbound(ctx)
	evolutionPort := whatsback.NewWhatsbackAdapter()
	fonnteAdapter := whatsapp_outbound_adapter.NewAdapter()

	notifService := service_notification.NewNotificationService(fonnteAdapter.WhatsApp(), evolutionPort, dbPort)

	dom := domain.NewDomain(
		dbPort,
		messagePort,
		cacheOutbound(ctx),
		workflowOutbound(ctx),
		evolutionPort,
		notifService,
	)

	return &App{
		ctx:     ctx,
		domain:  dom,
		db:      dbPort,
		message: messagePort,
		notif:   notifService,
	}
}

func (a *App) Run(option string) {
	switch option {
	case "http":
		a.httpInbound()
	case "message":
		a.messageInbound()
	case "workflow":
		a.workflowInbound()
	case "migrate":
		a.runMigrations()
	default:
		a.commandInbound()
	}
}

func (a *App) runMigrations() {
	if outboundDatabaseDriver != "postgres" {
		log.WithContext(a.ctx).Fatal("Migration only supports postgres driver")
	}

	db := database.InitDatabase(a.ctx, outboundDatabaseDriver)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.WithContext(a.ctx).Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.WithContext(a.ctx).Fatalf("failed to run migrations: %v", err)
	}

	log.WithContext(a.ctx).Info("Migrations completed successfully")
}

func databaseOutbound(ctx context.Context) outbound_port.DatabasePort {
	if !utils.IsInList(databaseDriverList, outboundDatabaseDriver) {
		log.WithContext(ctx).Fatal("database driver is not supported")
		os.Exit(1)
	}
	db := database.InitDatabase(ctx, outboundDatabaseDriver)

	switch outboundDatabaseDriver {
	case "postgres":
		return postgres_outbound_adapter.NewAdapter(db)
	}
	return nil
}

func messageOutbound(ctx context.Context) outbound_port.MessagePort {
	if !utils.IsInList(messageDriverList, outboundMessageDriver) {
		log.WithContext(ctx).Fatal("message driver is not supported")
		os.Exit(1)
	}

	switch outboundMessageDriver {
	case "rabbitmq":
		if err := rabbitmq.InitMessage(); err != nil {
			log.WithContext(ctx).Fatalf("failed to init rabbitmq: %v", err)
		}
		return rabbitmq_outbound_adapter.NewAdapter()
	case "whatsapp":
		return whatsapp_outbound_adapter.NewAdapter()
	case "none":
		return nil
	}
	return nil
}

func cacheOutbound(ctx context.Context) outbound_port.CachePort {
	if !utils.IsInList([]string{"redis", "gibrun"}, outboundCacheDriver) {
		log.WithContext(ctx).Fatal("cache driver is not supported")
		os.Exit(1)
	}

	switch outboundCacheDriver {
	case "redis":
		redis.InitDatabase()
		return redis_outbound_adapter.NewAdapter()
	case "gibrun":
		gibrun.Init()
		return gibrun_outbound_adapter.NewAdapter()
	}
	return nil
}

func workflowOutbound(ctx context.Context) outbound_port.WorkflowPort {
	if !utils.IsInList([]string{"temporal", "none"}, outboundWorkflowDriver) {
		log.WithContext(ctx).Fatal("workflow driver is not supported")
		os.Exit(1)
	}

	switch outboundWorkflowDriver {
	case "temporal":
		return temporal_outbound_adapter.NewAdapter()
	case "none":
		return nil
	}
	return nil
}

func (a *App) httpInbound() {
	ctx := a.ctx
	if !utils.IsInList(httpDriverList, inboundHttpDriver) {
		log.WithContext(ctx).Fatal("http driver is not supported")
		os.Exit(1)
	}

	switch inboundHttpDriver {
	case "fiber":
		engine := html.New("./web/templates", ".html")
		app := fiber.New(fiber.Config{
			Views: engine,
		})
		inboundHttpAdapter := fiber_inbound_adapter.NewAdapter(a.domain, a.message)
		fiber_inbound_adapter.InitRoute(ctx, app, inboundHttpAdapter, a.domain, a.db)
		go func() {
			port := os.Getenv("SERVER_PORT")
			if port == "" {
				port = os.Getenv("PORT")
			}
			if port == "" {
				port = "8081" // Default fallback
			}
			if err := app.Listen(":" + port); err != nil {
				log.WithContext(ctx).Fatalf("failed to listen fiber: %v", err)
			}
		}()

		// Start scheduler for subscription reminders
		a.scheduler = scheduler.NewScheduler(ctx, a.domain, a.db)
		a.scheduler.Start()

		// Start WhatsApp Notification Consumer (RabbitMQ -> Fonnte/Evolution)
		if inboundMessageDriver == "rabbitmq" {
			go func() {
				if err := a.notif.StartWhatsAppConsumer(ctx); err != nil {
					log.WithContext(ctx).Errorf("failed to start whatsapp consumer: %v", err)
				}
			}()
		}
	}

	ctx, shutdown := context.WithTimeout(ctx, 5*time.Second)
	defer shutdown()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	log.WithContext(ctx).Info("http server stopped")
}

func (a *App) messageInbound() {
	ctx := a.ctx
	if !utils.IsInList(messageDriverList, inboundMessageDriver) {
		log.WithContext(ctx).Fatal("message driver is not supported")
		os.Exit(1)
	}

	switch inboundMessageDriver {
	case "rabbitmq":
		inboundMessageAdapter := rabbitmq_inbound_adapter.NewAdapter(a.domain)
		rabbitmq_inbound_adapter.InitRoute(ctx, os.Args, inboundMessageAdapter)
	}
}

func (a *App) commandInbound() {
	ctx := a.ctx
	inboundCommandAdapter := command_inbound_adapter.NewAdapter(a.domain)
	command_inbound_adapter.InitRoute(ctx, os.Args, inboundCommandAdapter)
}

func (a *App) workflowInbound() {
	ctx := a.ctx
	if !utils.IsInList(workflowDriverList, inboundWorkflowDriver) {
		log.WithContext(ctx).Fatal("workflow driver is not supported")
		os.Exit(1)
	}

	switch inboundWorkflowDriver {
	case "temporal":
		inboundWorkflowAdapter := temporal_inbound_adapter.NewAdapter(a.domain)
		temporal_inbound_adapter.InitRoute(ctx, os.Args, inboundWorkflowAdapter)
	}
}

func configureLogging() {
	// SECURITY: Use INFO level in production to avoid leaking sensitive data
	if os.Getenv("APP_MODE") == "release" {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&joonix.FluentdFormatter{})
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}
	logrus.AddHook(utils.LogrusSourceContextHook{})
}
