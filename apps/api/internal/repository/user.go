package repository

import (
	"time"

	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/cache"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	ResourceRepository[model.User]
}

type userRepository struct {
	ResourceRepository[model.User]
}

func NewUserRepository(db *gorm.DB, cacheClient cache.Cache, cacheTTL time.Duration) UserRepository {
	return &userRepository{
		ResourceRepository: NewResourceRepository[model.User](db, cacheClient, cacheTTL),
	}
}
