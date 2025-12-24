package repository

import (
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	*ResourceRepository[model.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		ResourceRepository: NewResourceRepository[model.User](db),
	}
}
