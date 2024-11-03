package service

import "codebase-app/internal/module/z_template/ports"

var _ ports.XxxService = &xxxService{}

type xxxService struct {
	repo ports.XxxRepository
}

func NewXxxService(repo ports.XxxRepository) *xxxService {
	return &xxxService{
		repo: repo,
	}
}
