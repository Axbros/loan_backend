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
	loanRolesCachePrefixKey = "loanRoles:"
	// LoanRolesExpireTime expire time
	LoanRolesExpireTime = 5 * time.Minute
)

var _ LoanRolesCache = (*loanRolesCache)(nil)

// LoanRolesCache cache interface
type LoanRolesCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanRoles, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanRoles, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRoles, error)
	MultiSet(ctx context.Context, data []*model.LoanRoles, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanRolesCache define a cache struct
type loanRolesCache struct {
	cache cache.Cache
}

// NewLoanRolesCache new a cache
func NewLoanRolesCache(cacheType *database.CacheType) LoanRolesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRoles{}
		})
		return &loanRolesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRoles{}
		})
		return &loanRolesCache{cache: c}
	}

	return nil // no cache
}

// GetLoanRolesCacheKey cache key
func (c *loanRolesCache) GetLoanRolesCacheKey(id uint64) string {
	return loanRolesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanRolesCache) Set(ctx context.Context, id uint64, data *model.LoanRoles, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanRolesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanRolesCache) Get(ctx context.Context, id uint64) (*model.LoanRoles, error) {
	var data *model.LoanRoles
	cacheKey := c.GetLoanRolesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanRolesCache) MultiSet(ctx context.Context, data []*model.LoanRoles, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanRolesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanRolesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRoles, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanRolesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanRoles)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanRoles)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanRolesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanRolesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRolesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanRolesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRolesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanRolesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
