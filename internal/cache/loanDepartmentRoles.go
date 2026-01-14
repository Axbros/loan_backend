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
	loanDepartmentRolesCachePrefixKey = "loanDepartmentRoles:"
	// LoanDepartmentRolesExpireTime expire time
	LoanDepartmentRolesExpireTime = 5 * time.Minute
)

var _ LoanDepartmentRolesCache = (*loanDepartmentRolesCache)(nil)

// LoanDepartmentRolesCache cache interface
type LoanDepartmentRolesCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanDepartmentRoles, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanDepartmentRoles, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDepartmentRoles, error)
	MultiSet(ctx context.Context, data []*model.LoanDepartmentRoles, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanDepartmentRolesCache define a cache struct
type loanDepartmentRolesCache struct {
	cache cache.Cache
}

// NewLoanDepartmentRolesCache new a cache
func NewLoanDepartmentRolesCache(cacheType *database.CacheType) LoanDepartmentRolesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanDepartmentRoles{}
		})
		return &loanDepartmentRolesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanDepartmentRoles{}
		})
		return &loanDepartmentRolesCache{cache: c}
	}

	return nil // no cache
}

// GetLoanDepartmentRolesCacheKey cache key
func (c *loanDepartmentRolesCache) GetLoanDepartmentRolesCacheKey(id uint64) string {
	return loanDepartmentRolesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanDepartmentRolesCache) Set(ctx context.Context, id uint64, data *model.LoanDepartmentRoles, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanDepartmentRolesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanDepartmentRolesCache) Get(ctx context.Context, id uint64) (*model.LoanDepartmentRoles, error) {
	var data *model.LoanDepartmentRoles
	cacheKey := c.GetLoanDepartmentRolesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanDepartmentRolesCache) MultiSet(ctx context.Context, data []*model.LoanDepartmentRoles, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanDepartmentRolesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanDepartmentRolesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDepartmentRoles, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanDepartmentRolesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanDepartmentRoles)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanDepartmentRoles)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanDepartmentRolesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanDepartmentRolesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanDepartmentRolesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanDepartmentRolesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanDepartmentRolesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanDepartmentRolesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
