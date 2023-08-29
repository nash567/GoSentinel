package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/nash567/GoSentinel/internal/config"
	mailSvc "github.com/nash567/GoSentinel/internal/notifications/email"
	mailModel "github.com/nash567/GoSentinel/internal/notifications/email/model"

	applicationSvc "github.com/nash567/GoSentinel/internal/service/application"
	authSvc "github.com/nash567/GoSentinel/internal/service/auth"

	applicationModel "github.com/nash567/GoSentinel/internal/service/application/model"
	application "github.com/nash567/GoSentinel/internal/service/application/repo"
	"github.com/nash567/GoSentinel/internal/service/auth"
	authModel "github.com/nash567/GoSentinel/internal/service/auth/model"
	authRepo "github.com/nash567/GoSentinel/internal/service/auth/repo"
	"github.com/nash567/GoSentinel/pkg/cache"
	logModel "github.com/nash567/GoSentinel/pkg/logger/model"
)

type services struct {
	cfg            *config.Config
	applicationSvc applicationModel.Service
	mailSvc        mailModel.Service
	authSvc        *auth.Service
}
type repos struct {
	applicationRepo applicationModel.Repository
	authRepo        authModel.Repository
}

func (a *Application) buildServices(db *sqlx.DB, log logModel.Logger, cfg *config.Config, cacheSvc cache.Cache) *services {
	svc := &services{
		cfg: cfg,
	}
	repo := repos{}
	repo.buildRepos(db)
	svc.buildServies(repo, a.cache, log)
	return svc

}
func (s *services) buildServies(repo repos, cacheSvc cache.Cache, log logModel.Logger) {
	s.mailSvc = mailSvc.NewService(&s.cfg.Mailer)
	s.authSvc = authSvc.NewService(&s.cfg.AuthConfig, repo.authRepo)
	s.applicationSvc = applicationSvc.NewService(&s.cfg.ApplicationConfig, s.mailSvc, repo.applicationRepo, cacheSvc, s.cfg.AuthConfig, s.authSvc)
}

func (r *repos) buildRepos(db *sqlx.DB) {
	r.applicationRepo = application.NewRepository(db.DB)
	r.authRepo = authRepo.NewRepository(db.DB)

}
