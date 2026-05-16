package analytics

import (
	"context"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Overview(ctx context.Context) (Overview, error) {
	return s.repo.Overview(ctx)
}

func (s *Service) Revenue(ctx context.Context) (Revenue, error) {
	return s.repo.Revenue(ctx)
}

func (s *Service) ExamStats(ctx context.Context) (ExamStats, error) {
	return s.repo.ExamStats(ctx)
}
