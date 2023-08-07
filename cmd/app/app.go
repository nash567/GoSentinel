package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/nash567/GoSentinel/internal/config"
	"github.com/nash567/GoSentinel/pkg/db"
	"github.com/nash567/GoSentinel/pkg/logger"
	logModel "github.com/nash567/GoSentinel/pkg/logger/model"
)

const timeout = 5 * time.Second

type Application struct {
	db         *sqlx.DB
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

	a.services = buildServices(a.db, a.log, a.cfg)
	a.setupHandlers()
}

func (a *Application) Start(ctx context.Context) {
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type"}),
	)
	a.router.Use(corsHandler)

	a.httpServer = &http.Server{
		Addr:              ":" + fmt.Sprintf("%v", a.cfg.Server.Port),
		Handler:           a.router,
		ReadHeaderTimeout: timeout,
	}
	go func() {
		defer a.log.Error("server stopped listening...")

		if err := a.httpServer.ListenAndServe(); err != nil {
			a.log.Fatal("failed to listen and serve")
			return
		}
	}()
	a.log.Infof("http server started on port %d..", a.cfg.Server.Port)
}
func (a *Application) Stop(ctx context.Context) {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}

	a.log.Warn("shutting down....")
}

func (a *Application) setupHandlers() {

}
