package model

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/validation"
)

type BaseModel interface {
	GetID() uuid.UUID
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
