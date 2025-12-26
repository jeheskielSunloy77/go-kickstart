package repository

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	ResourceRepository[model.User]
}

type userRepository struct {
	ResourceRepository[model.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		ResourceRepository: NewResourceRepository[model.User](db),
	}
}
