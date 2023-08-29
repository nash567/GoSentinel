package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	"github.com/nash567/GoSentinel/api/v1/pb/goSentinel"
	"github.com/nash567/GoSentinel/api/v1/rpc"
	"github.com/nash567/GoSentinel/internal/config"
	"github.com/nash567/GoSentinel/pkg/cache"
	"github.com/nash567/GoSentinel/pkg/cache/redis"
	"github.com/nash567/GoSentinel/pkg/db"
	"github.com/nash567/GoSentinel/pkg/logger"
	logModel "github.com/nash567/GoSentinel/pkg/logger/model"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const timeout = 5 * time.Second

type Application struct {
	db         *sqlx.DB
	cache      cache.Cache
	grpcServer *grpc.Server
	httpServer *http.Server
	cfg        *config.Config
	router     *mux.Router
	log        logModel.Logger
	services   *services
}

func (a *Application) Init(ctx context.Context, configFile string, migrationPath string, seedDataPath string) {
	config, err := config.Load(configFile)
	if err != nil {
		log.Fatal("failed to read config")
		return
	}
	a.cfg = config
	a.log = logger.NewSLogger(&a.cfg.Log)
	if err != nil {
		panic(err)
	}

	a.log = a.log.WithFields(logModel.Fields{
		"appName": a.cfg.AppName,
		"env":     a.cfg.Env,
	})

	db, err := db.NewConnection(&config.DB)
	if err != nil {
		a.log.Fatal(err.Error())
		return
	}
	a.db = db
	a.log.WithField("host", a.cfg.DB.Host).WithField("port", a.cfg.DB.Port).Info("created database connection successfully")

	a.router = mux.NewRouter()
	a.cache, err = redis.New(ctx, &a.cfg.Redis)
	if err != nil {
		a.log.Fatal(err.Error())
		return
	}
	a.services = a.buildServices(a.db, a.log, a.cfg, a.cache)
	a.setupHandlers()
}

func (a *Application) Start(ctx context.Context) {
	a.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"accept", "Authorization", "content-type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)

	a.httpServer = &http.Server{
		Addr:              ":" + fmt.Sprintf("%v", a.cfg.HttpServer.Port),
		Handler:           a.registerHTTPEndpoints(ctx, a.cfg, a.log),
		ReadHeaderTimeout: timeout,
	}
	go func() {
		defer a.log.Infof("server stopped listening")
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.log.Errorf("failed to listen and serve: %v ", err)
			return
		}

	}()
	a.log.Infof("http server started on %d ...", a.cfg.HttpServer.Port)

	go func() {
		// grpc server
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GrpcServer.Host, a.cfg.GrpcServer.Port))
		if err != nil {
			a.log.Fatalf("failed to listen: %v", err)
		}

		a.grpcServer = a.registerGRPCEndpoints(a.services, a.log)
		if err := a.grpcServer.Serve(lis); err != nil {
			a.log.Fatalf("failed to serve: %v", err)
		}

	}()

	a.log.Infof("grpc server started on %d ...", a.cfg.GrpcServer.Port)
}

func (a *Application) Stop(ctx context.Context) {
	a.grpcServer.GracefulStop()

	a.log.Warn("shutting down....")
}

func (a *Application) setupHandlers() {

}
func (a *Application) registerHTTPEndpoints(ctx context.Context, cfg *config.Config, log logModel.Logger) *runtime.ServeMux {
	mux := runtime.NewServeMux()
	err := goSentinel.RegisterGoSentinelServiceHandlerFromEndpoint(ctx,
		mux,
		fmt.Sprintf("%v:%v", cfg.GrpcServer.Host, cfg.GrpcServer.Port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Errorf("Error while registering application service: %w", err)
	}
	return mux
}

func (a *Application) registerGRPCEndpoints(services *services, log logModel.Logger) *grpc.Server {
	// interceptor := authInterceptor.NewAuthInterceptor(log, true, *services.AuthSvc)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(a.services.authSvc.AuthenticationInterceptor),
		// intercept the request to check the token
		// grpc.UnaryInterceptor(interceptor.Auth),
	}
	grpcServer := grpc.NewServer(opts...)

	goSentinel.RegisterGoSentinelServiceServer(grpcServer, rpc.NewServer(a.services.applicationSvc, a.services.authSvc, a.cfg.AuthConfig.EncryptionKey))

	return grpcServer
}
