package repository

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/cache"
	"github.com/jeheskielSunloy77/go-kickstart/internal/lib/utils"
	"gorm.io/gorm"
)

type ResourceRepository[T any] interface {
	Store(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*T, error)
	GetMany(ctx context.Context, opts GetManyOptions) ([]T, int64, error)
	Update(ctx context.Context, entity T, updates ...map[string]any) (*T, error)
	Destroy(ctx context.Context, id uuid.UUID) error
	Kill(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (*T, error)
}

type resourceRepository[T any] struct {
	db       *gorm.DB
	cache    cache.Cache
	cacheTTL time.Duration
}

func NewResourceRepository[T any](db *gorm.DB, cacheClient cache.Cache, cacheTTL time.Duration) ResourceRepository[T] {
	if cacheTTL <= 0 {
		cacheClient = nil
	}
	return &resourceRepository[T]{db: db, cache: cacheClient, cacheTTL: cacheTTL}
}

func (r *resourceRepository[T]) Store(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return err
	}
	r.setCacheFromEntity(ctx, entity)
	return nil
}

func (r *resourceRepository[T]) GetByID(ctx context.Context, id uuid.UUID, preloads []string) (*T, error) {
	if r.cacheEnabled() && len(preloads) == 0 {
		if cached, ok := r.getCachedByID(ctx, id); ok {
			return cached, nil
		}
	}

	var entity T
	query := applyPreloads(r.db.WithContext(ctx), preloads)
	if err := query.First(&entity, id).Error; err != nil {
		return nil, err
	}
	r.setCache(ctx, id, &entity, preloads)
	return &entity, nil
}

func (r *resourceRepository[T]) Update(ctx context.Context, entity T, updates ...map[string]any) (*T, error) {
	// if the updates are provided, use them to only update specific fields, if not replace the entire entity
	var err error
	if len(updates) > 0 {
		err = r.db.WithContext(ctx).Model(&entity).Updates(updates[0]).Error
	} else {
		err = r.db.WithContext(ctx).Save(&entity).Error
	}
	if err != nil {
		return nil, err
	}

	// updated entity
	r.evictCacheFromEntity(ctx, entity)
	return &entity, nil
}

func (r *resourceRepository[T]) Destroy(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(new(T), id).Error; err != nil {
		return err
	}
	r.evictCache(ctx, id)
	return nil
}

func (r *resourceRepository[T]) Kill(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(new(T), id).Error; err != nil {
		return err
	}
	r.evictCache(ctx, id)
	return nil
}

func (r *resourceRepository[T]) Restore(ctx context.Context, id uuid.UUID) (*T, error) {
	if err := r.db.WithContext(ctx).
		Unscoped().
		Model(new(T)).
		Where("id = ?", id).
		Update("deleted_at", nil).
		Error; err != nil {
		return nil, err
	}

	r.evictCache(ctx, id)
	return r.GetByID(ctx, id, nil)
}

func (r *resourceRepository[T]) GetMany(ctx context.Context, opts GetManyOptions) ([]T, int64, error) {
	opts.Normalize()

	var (
		entities []T
		total    int64
	)

	countQuery := r.db.WithContext(ctx).Model(new(T))
	countQuery = applyJoins(countQuery, opts.Joins)
	countQuery = applyFilters(countQuery, opts.Filters)
	countQuery = applyWheres(countQuery, opts.Wheres)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	listQuery := r.db.WithContext(ctx).Model(new(T))
	listQuery = applyJoins(listQuery, opts.Joins)
	listQuery = applyFilters(listQuery, opts.Filters)
	listQuery = applyWheres(listQuery, opts.Wheres)
	listQuery = applyPreloads(listQuery, opts.Preloads)
	if err := listQuery.Limit(opts.Limit).Offset(opts.Offset).Order(opts.OrderBy + " " + opts.OrderDirection).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

type JoinClause struct {
	Query string
	Args  []any
}

type WhereClause struct {
	Query string
	Args  []any
}

type GetManyOptions struct {
	Filters        map[string]any
	Joins          []JoinClause
	Wheres         []WhereClause
	Preloads       []string
	OrderBy        string
	OrderDirection string
	Limit          int
	Offset         int
}

func NewGetManyOptionsFromRequest(c *fiber.Ctx) GetManyOptions {
	opts := GetManyOptions{
		Limit:          utils.ParseQueryInt(c.Query("limit")),
		Offset:         utils.ParseQueryInt(c.Query("offset")),
		Preloads:       ParsePreloads(c.Query("preloads")),
		OrderBy:        c.Query("orderBy"),
		OrderDirection: c.Query("orderDirection"),
	}
	opts.Normalize()
	return opts
}

func (o *GetManyOptions) Normalize() {
	if o.Limit <= 0 {
		o.Limit = 20
	}

	o.OrderDirection = strings.ToLower(strings.TrimSpace(o.OrderDirection))
	if o.OrderDirection == "" || (o.OrderDirection != "asc" && o.OrderDirection != "desc") {
		o.OrderDirection = "desc"
	}

	if o.OrderBy == "" {
		o.OrderBy = "created_at"
	}
}

func ParsePreloads(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	preloads := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		preloads = append(preloads, name)
	}

	if len(preloads) == 0 {
		return nil
	}
	return preloads
}

func applyFilters(db *gorm.DB, filters map[string]any) *gorm.DB {
	if len(filters) > 0 {
		return db.Where(filters)
	}
	return db
}

func applyWheres(db *gorm.DB, wheres []WhereClause) *gorm.DB {
	for _, where := range wheres {
		query := strings.TrimSpace(where.Query)
		if query == "" {
			continue
		}
		if len(where.Args) > 0 {
			db = db.Where(query, where.Args...)
			continue
		}
		db = db.Where(query)
	}
	return db
}

func applyJoins(db *gorm.DB, joins []JoinClause) *gorm.DB {
	for _, join := range joins {
		query := strings.TrimSpace(join.Query)
		if query == "" {
			continue
		}
		if len(join.Args) > 0 {
			db = db.Joins(query, join.Args...)
			continue
		}
		db = db.Joins(query)
	}
	return db
}

func applyPreloads(db *gorm.DB, preloads []string) *gorm.DB {
	for _, preload := range preloads {
		name := strings.TrimSpace(preload)
		if name == "" {
			continue
		}
		db = db.Preload(name)
	}
	return db
}

func (r *resourceRepository[T]) cacheEnabled() bool {
	return r.cache != nil && r.cacheTTL > 0
}

func (r *resourceRepository[T]) getCachedByID(ctx context.Context, id uuid.UUID) (*T, bool) {
	if !r.cacheEnabled() {
		return nil, false
	}
	key := r.cacheKey(id)
	data, err := r.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return nil, false
		}
		return nil, false
	}

	var entity T
	if err := json.Unmarshal(data, &entity); err != nil {
		_ = r.cache.Delete(ctx, key)
		return nil, false
	}
	return &entity, true
}

func (r *resourceRepository[T]) setCache(ctx context.Context, id uuid.UUID, entity *T, preloads []string) {
	if !r.cacheEnabled() || entity == nil || len(preloads) != 0 {
		return
	}

	payload, err := json.Marshal(entity)
	if err != nil {
		return
	}
	_ = r.cache.Set(ctx, r.cacheKey(id), payload, r.cacheTTL)
}

func (r *resourceRepository[T]) setCacheFromEntity(ctx context.Context, entity *T) {
	if !r.cacheEnabled() || entity == nil {
		return
	}
	id, ok := extractEntityID(entity)
	if !ok {
		return
	}
	r.setCache(ctx, id, entity, nil)
}

func (r *resourceRepository[T]) evictCacheFromEntity(ctx context.Context, entity T) {
	if !r.cacheEnabled() {
		return
	}
	id, ok := extractEntityID(entity)
	if !ok {
		return
	}
	r.evictCache(ctx, id)
}

func (r *resourceRepository[T]) evictCache(ctx context.Context, id uuid.UUID) {
	if !r.cacheEnabled() {
		return
	}
	_ = r.cache.Delete(ctx, r.cacheKey(id))
}

func (r *resourceRepository[T]) cacheKey(id uuid.UUID) string {
	return "resource:" + resourceTypeName[T]() + ":id:" + id.String()
}

func resourceTypeName[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	name := t.Name()
	pkg := strings.ReplaceAll(t.PkgPath(), "/", ".")
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

func extractEntityID(entity any) (uuid.UUID, bool) {
	if entity == nil {
		return uuid.Nil, false
	}

	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return uuid.Nil, false
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return uuid.Nil, false
	}

	field := val.FieldByName("ID")
	if !field.IsValid() {
		return uuid.Nil, false
	}
	if field.Type() != reflect.TypeOf(uuid.UUID{}) {
		return uuid.Nil, false
	}

	id, ok := field.Interface().(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, false
	}
	return id, true
}
