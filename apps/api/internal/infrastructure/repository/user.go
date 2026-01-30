package repository

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/application/port"
	"github.com/jeheskielSunloy77/go-kickstart/internal/domain"
	"github.com/jeheskielSunloy77/go-kickstart/internal/infrastructure/config"
	"github.com/jeheskielSunloy77/go-kickstart/internal/infrastructure/lib/cache"
	"gorm.io/gorm"
)

type UserRepository = port.UserRepository

type userRepository struct {
	ResourceRepository[domain.User]
}

func NewUserRepository(cfg *config.Config, db *gorm.DB, cacheClient cache.Cache) UserRepository {
	return &userRepository{
		ResourceRepository: NewResourceRepository[domain.User](cfg, db, cacheClient),
	}
}
