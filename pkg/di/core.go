package di

import (
	"agentos/internal/domain"
	"agentos/internal/domain/ports/out"
	"agentos/pkg/db"
)

type Container struct {
	Logger out.Logger
	Config domain.Configuration
	DB     *db.Adapter
}

func New() Container {
	return Container{}
}
