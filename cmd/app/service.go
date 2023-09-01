package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/nash567/GoSentinel/internal/config"
	mailSvc "github.com/nash567/GoSentinel/internal/notifications/email"
	mailModel "github.com/nash567/GoSentinel/internal/notifications/email/model"
	"github.com/nash567/GoSentinel/pkg/cache"

	applicationSvc "github.com/nash567/GoSentinel/internal/service/application"
	authSvc "github.com/nash567/GoSentinel/internal/service/auth"

	applicationModel "github.com/nash567/GoSentinel/internal/service/application/model"
	application "github.com/nash567/GoSentinel/internal/service/application/repo"
	authModel "github.com/nash567/GoSentinel/internal/service/auth/model"
	userSvc "github.com/nash567/GoSentinel/internal/service/user"
	userModel "github.com/nash567/GoSentinel/internal/service/user/model"
	userRepo "github.com/nash567/GoSentinel/internal/service/user/repo"

	authRepo "github.com/nash567/GoSentinel/internal/service/auth/repo"
	logModel "github.com/nash567/GoSentinel/pkg/logger/model"
)

type services struct {
	cfg            *config.Config
	applicationSvc applicationModel.Service
	mailSvc        mailModel.Service
	authSvc        authModel.Service
	userSvc        userModel.Service
}
type repos struct {
	applicationRepo applicationModel.Repository
	authRepo        authModel.Repository
	userRepo        userModel.Repository
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
	s.userSvc = userSvc.NewService(repo.userRepo, s.authSvc, s.cfg.AuthConfig)
	s.applicationSvc = applicationSvc.NewService(&s.cfg.ApplicationConfig, s.mailSvc, repo.applicationRepo, cacheSvc, s.cfg.AuthConfig, s.authSvc)
}

func (r *repos) buildRepos(db *sqlx.DB) {
	r.applicationRepo = application.NewRepository(db.DB)
	r.authRepo = authRepo.NewRepository(db.DB)
	r.userRepo = userRepo.NewRepository(db.DB)

}
