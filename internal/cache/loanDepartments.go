package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"loan/internal/database"
	"loan/internal/model"
)

const (
	// cache prefix key, must end with a colon
	loanDepartmentsCachePrefixKey = "loanDepartments:"
	// LoanDepartmentsExpireTime expire time
	LoanDepartmentsExpireTime = 5 * time.Minute
)

var _ LoanDepartmentsCache = (*loanDepartmentsCache)(nil)

// LoanDepartmentsCache cache interface
type LoanDepartmentsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanDepartments, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanDepartments, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDepartments, error)
	MultiSet(ctx context.Context, data []*model.LoanDepartments, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanDepartmentsCache define a cache struct
type loanDepartmentsCache struct {
	cache cache.Cache
}

// NewLoanDepartmentsCache new a cache
func NewLoanDepartmentsCache(cacheType *database.CacheType) LoanDepartmentsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanDepartments{}
		})
		return &loanDepartmentsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanDepartments{}
		})
		return &loanDepartmentsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanDepartmentsCacheKey cache key
func (c *loanDepartmentsCache) GetLoanDepartmentsCacheKey(id uint64) string {
	return loanDepartmentsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanDepartmentsCache) Set(ctx context.Context, id uint64, data *model.LoanDepartments, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanDepartmentsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanDepartmentsCache) Get(ctx context.Context, id uint64) (*model.LoanDepartments, error) {
	var data *model.LoanDepartments
	cacheKey := c.GetLoanDepartmentsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanDepartmentsCache) MultiSet(ctx context.Context, data []*model.LoanDepartments, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanDepartmentsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanDepartmentsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDepartments, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanDepartmentsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanDepartments)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanDepartments)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanDepartmentsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanDepartmentsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanDepartmentsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanDepartmentsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanDepartmentsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanDepartmentsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
