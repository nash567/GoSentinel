package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/nash567/GoSentinel/internal/config"
	logModel "github.com/nash567/GoSentinel/pkg/logger/model"
)

type services struct {
	cfg *config.Config
}
type repos struct {
}

func buildServices(db *sqlx.DB, log logModel.Logger, cfg *config.Config) *services {
	svc := &services{
		cfg: cfg,
	}
	repo := repos{}
	repo.buildRepos(db)
	svc.buildServies(repo, log)
	return svc

}
func (s *services) buildServies(repo repos, log logModel.Logger) {

}

func (r *repos) buildRepos(db *sqlx.DB) {

}
