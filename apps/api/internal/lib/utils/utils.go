package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/errs"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
)

func PrintJSON(v any) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println("JSON:", string(json))
}

func ParseUUIDParam(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, errs.NewBadRequestError("invalid id provided", true, []errs.FieldError{{Field: "id", Error: "must be a valid uuid"}}, nil)
	}
	return id, nil
}

func ParseQueryInt(raw string, maxAndDefaultVal ...int) int {
	var (
		defaultVal int
		max        *int
	)

	if len(maxAndDefaultVal) > 0 {
		max = &maxAndDefaultVal[0]
	}
	if len(maxAndDefaultVal) > 1 {
		defaultVal = maxAndDefaultVal[1]
	}

	if raw == "" {
		return defaultVal
	}

	if v, err := strconv.Atoi(raw); err == nil {
		if v < 1 {
			return defaultVal
		}
		if max != nil && v > *max {
			return *max
		}
		return v
	}

	return defaultVal
}

func GetModelName[T model.BaseModel]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}

func GetModelNameLower[T model.BaseModel]() string {
	return strings.ToLower(GetModelName[T]())
}

func GetModelSemanticName[T model.BaseModel]() string {
	var name strings.Builder

	for i, r := range GetModelName[T]() {
		if i > 0 && r >= 'A' && r <= 'Z' {
			name.WriteRune(' ')
		}
		name.WriteRune(r)
	}
	return name.String()

}
func GetModelCacheKey[T model.BaseModel](id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return "resource:" + GetModelNameLower[T]() + ":id:" + id.String()
}
