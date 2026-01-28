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
	"prabogo/internal/domain"
	_ "prabogo/internal/migration/postgres"
	outbound_port "prabogo/internal/port/outbound"
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
	ctx    context.Context
	domain domain.Domain
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
	domain := domain.NewDomain(
		databaseOutbound(ctx),
		messageOutbound(ctx),
		cacheOutbound(ctx),
		workflowOutbound(ctx),
	)

	return &App{
		ctx:    ctx,
		domain: domain,
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
		inboundHttpAdapter := fiber_inbound_adapter.NewAdapter(a.domain)
		fiber_inbound_adapter.InitRoute(ctx, app, inboundHttpAdapter)
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
	logrus.SetLevel(logrus.DebugLevel)
	logrus.AddHook(utils.LogrusSourceContextHook{})

	if os.Getenv("APP_MODE") != "release" {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	} else {
		logrus.SetFormatter(&joonix.FluentdFormatter{})
	}
}
