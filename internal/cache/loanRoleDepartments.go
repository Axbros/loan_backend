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
	loanRoleDepartmentsCachePrefixKey = "loanRoleDepartments:"
	// LoanRoleDepartmentsExpireTime expire time
	LoanRoleDepartmentsExpireTime = 5 * time.Minute
)

var _ LoanRoleDepartmentsCache = (*loanRoleDepartmentsCache)(nil)

// LoanRoleDepartmentsCache cache interface
type LoanRoleDepartmentsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanRoleDepartments, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanRoleDepartments, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRoleDepartments, error)
	MultiSet(ctx context.Context, data []*model.LoanRoleDepartments, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanRoleDepartmentsCache define a cache struct
type loanRoleDepartmentsCache struct {
	cache cache.Cache
}

// NewLoanRoleDepartmentsCache new a cache
func NewLoanRoleDepartmentsCache(cacheType *database.CacheType) LoanRoleDepartmentsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRoleDepartments{}
		})
		return &loanRoleDepartmentsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRoleDepartments{}
		})
		return &loanRoleDepartmentsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanRoleDepartmentsCacheKey cache key
func (c *loanRoleDepartmentsCache) GetLoanRoleDepartmentsCacheKey(id uint64) string {
	return loanRoleDepartmentsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanRoleDepartmentsCache) Set(ctx context.Context, id uint64, data *model.LoanRoleDepartments, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanRoleDepartmentsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanRoleDepartmentsCache) Get(ctx context.Context, id uint64) (*model.LoanRoleDepartments, error) {
	var data *model.LoanRoleDepartments
	cacheKey := c.GetLoanRoleDepartmentsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanRoleDepartmentsCache) MultiSet(ctx context.Context, data []*model.LoanRoleDepartments, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanRoleDepartmentsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanRoleDepartmentsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRoleDepartments, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanRoleDepartmentsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanRoleDepartments)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanRoleDepartments)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanRoleDepartmentsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanRoleDepartmentsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRoleDepartmentsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanRoleDepartmentsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRoleDepartmentsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanRoleDepartmentsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
