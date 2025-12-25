package model

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/validation"
)

type BaseWithId struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type BaseWithCreatedAt struct {
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type BaseWithUpdatedAt struct {
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Base struct {
	BaseWithId
	BaseWithCreatedAt
	BaseWithUpdatedAt
}

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type BaseDTO[T any] interface {
	ToModel() *T
}

type StoreDTO[T any] interface {
	BaseDTO[T]
	validation.Validatable
}

type UpdateDTO[T any] interface {
	BaseDTO[T]
	validation.Validatable
	ToMap() map[string]any
}

type EmptyDTO struct{}

func (d *EmptyDTO) Validate() error { return nil }

func NewDTO[T any]() T {
	var dto T
	t := reflect.TypeOf(dto)
	if t != nil && t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem()).Interface().(T)
	}
	return dto
}
