package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/ravilushqa/otelgqlgen"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"

	"github.com/kargotech/go-testapp/internal/handler/gql"
	order_handler "github.com/kargotech/go-testapp/internal/handler/rest/order"
	order_repo "github.com/kargotech/go-testapp/internal/repo/order"
	order_usecase "github.com/kargotech/go-testapp/internal/usecase/order"
	auditevent "github.com/kargotech/gokargo/audit/event"
	"github.com/kargotech/gokargo/feature_flag"
	"github.com/kargotech/gokargo/graceful"
	logger "github.com/kargotech/gokargo/logger"
	exporter_http_gin "github.com/kargotech/gokargo/metrics/exporter/http/gin"
	instrumenter_gorm "github.com/kargotech/gokargo/metrics/instrumenter/gorm"
	instrumenter_gql "github.com/kargotech/gokargo/metrics/instrumenter/gql"
	instrumenter_http_gin "github.com/kargotech/gokargo/metrics/instrumenter/http/gin"
	"github.com/kargotech/gokargo/metrics/recorder"
	panic_handler "github.com/kargotech/gokargo/panic_handler"
	panic_handler_gin "github.com/kargotech/gokargo/panic_handler/http/gin"
	error_gin "github.com/kargotech/gokargo/serror/gin"
	"github.com/kargotech/gokargo/unitofwork"

	"github.com/kargotech/go-testapp/config"
	"github.com/kargotech/go-testapp/gen/graph/generated"

	// 	"github.com/kargotech/go-testapp/internal/handler/gql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const gracefulPeriod = 10 * time.Second

func main() {

	mainCfg := config.InitializeConfig()

	// Initialize feature flag
	feature_flag.InitFeatureFlag(mainCfg.FeatureFlag.URL, mainCfg.FeatureFlag.Token, mainCfg.FeatureFlag.AppName)

	tp, err := initTracing(mainCfg.Service)
	if err != nil {
		logger.KargoLog.Error(err.Error())
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.KargoLog.InfoF("Error shutting down tracer provider: %v", err)
		}
	}()

	// Any setups concerning singletong goes here
	// ================ SINGLETON ========================
	panic_handler.SetOptions(panic_handler.Options{
		Logger: logger.KargoLog,
	})
	// ================ SINGLETON ========================

	// Uncomment to see the config during startup
	// fmt.Printf("%+v\n", mainCfg)

	// External resource initialization for Dependency Injection goes here
	// ================ DEPENDENCY INJECTION RESOURCE ======================
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		mainCfg.Database.Host,
		mainCfg.Database.Username,
		mainCfg.Database.Password,
		mainCfg.Database.Name,
		mainCfg.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.KargoLog.Fatal(err.Error())
	}

	if err := db.Use(otelgorm.NewPlugin(
		otelgorm.WithTracerProvider(tp),
		otelgorm.WithoutQueryVariables())); err != nil {
		logger.KargoLog.Fatal(err.Error())
	}

	// TODO: Change this to prometheus later, since we don't have NR licensing yet
	// nr, err := newrelic.NewApplication(
	// 	newrelic.ConfigAppName(mainCfg.NewRelic.AppName),
	// 	newrelic.ConfigLicense(mainCfg.NewRelic.LicenseKey),
	// 	newrelic.ConfigDistributedTracerEnabled(mainCfg.NewRelic.DistributedTracerEnabled),
	// )
	// if err != nil {
	// 	fmt.Println("String Error: ", err.Error())
	// }
	// ================ DEPENDENCY INJECTION RESOURCE ======================

	// ================ DEPENDENCY INJECTION ================
	// Repo Initialization
	orderRepo := order_repo.New(db)
	auditEventRepo := auditevent.NewAuditEvent(db)

	unitOfWork := unitofwork.NewUnitOfWorkGorm(db, &auditEventRepo)
	// Usecase initialization
	orderUsecase := order_usecase.New(orderRepo, &unitOfWork)

	// REST Handler initialization
	orderHandler := order_handler.New(orderUsecase)

	// Graphql Handler initialization
	// graphqlHandler := graphqlHandler(orderUsecase)

	// ================ DEPENDENCY INJECTION ================

	router := gin.Default()

	// Middleware Setup
	router.Use(error_gin.ErrorMiddleware())
	router.Use(panic_handler_gin.PanicHandlerMiddleware())
	router.Use(otelgin.Middleware(mainCfg.Service.Name, otelgin.WithTracerProvider(tp)))

	// Setup Standard metric
	recorder := recorder.NewRecorder()
	// Instantiate DB Gorm Metric
	instrumenter_gorm.New(recorder).Setup(db, mainCfg.Database.Name, "public")
	// Setup Http metric
	instrumenter_http_gin.New(recorder).Setup(router)
	// Setup Gql metric with Graphql Handler initialization
	normalEndpointServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &gql.Resolver{
		OrderUC: orderUsecase,
	}}))
	normalEndpointServer.Use(otelgqlgen.Middleware(otelgqlgen.WithTracerProvider(tp)))
	normalHandler := instrumenter_gql.GqlHandler{RequestPath: "/query", Handler: normalEndpointServer, TypeSchema: instrumenter_gql.Normal}
	instrumenter_gql.New(recorder).Setup(router, map[string]instrumenter_gql.GqlHandler{"/query": normalHandler})
	// Setup exporter

	exporter_http_gin.New(recorder).Setup(router)

	// HTTP endpoints
	router.GET("/health", func(ctx *gin.Context) {
		logger.KargoLog.Info("Going to be Healthy")

		ctx.JSON(200, gin.H{
			"message": "HEALTHY!",
		})
	})
	router.GET("health2", func(ctx *gin.Context) {
		logger.KargoLog.Info("Going to be Healthy 2")

		ctx.JSON(200, gin.H{
			"message": "PERFECTLY HEALTHY!",
		})
	})
	router.GET("trypanic", func(ctx *gin.Context) {
		panic("trypanic endpoint")
	})
	router.POST("/order", error_gin.GinHandler(orderHandler.CreateOrder))
	router.PUT("/order/:id", error_gin.GinHandler(orderHandler.UpdateOrder))
	router.GET("/order/:id", error_gin.GinHandler(orderHandler.GetOrderByID))

	// GQL Endpoints
	// router.POST("/query", graphqlHandler)
	router.GET("/", playgroundHandler())

	fmt.Printf("Server Running (Port = 8080), route: http://localhost:8080/health\n")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	graceful.RunHttpServer(context.Background(), server, gracefulPeriod)
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func initTracing(svcCfg config.ServiceConfig) (*tracesdk.TracerProvider, error) {

	logger.KargoLog.Info("Start initialize exporter: Blocking section")

	// Create OTLP exporter
	exp, err := otlptracegrpc.New(context.TODO(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		return nil, err
	}

	logger.KargoLog.Info("Start initialize resource")

	r, err := resource.New(context.TODO(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithProcessPID(),
		resource.WithProcessOwner(),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeDescription(),
		resource.WithProcessRuntimeVersion(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
	)
	if err != nil {
		return nil, err
	}

	r, err = resource.Merge(r,
		resource.NewSchemaless(
			semconv.ServiceNameKey.String(svcCfg.Name),
		))
	if err != nil {
		return nil, err
	}

	logger.KargoLog.Info("Start initialize tracer provider")

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(r),
	)

	logger.KargoLog.Info("Set trace provider & propagator globally")

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
