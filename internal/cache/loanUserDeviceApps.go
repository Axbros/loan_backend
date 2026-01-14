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
	loanUserDeviceAppsCachePrefixKey = "loanUserDeviceApps:"
	// LoanUserDeviceAppsExpireTime expire time
	LoanUserDeviceAppsExpireTime = 5 * time.Minute
)

var _ LoanUserDeviceAppsCache = (*loanUserDeviceAppsCache)(nil)

// LoanUserDeviceAppsCache cache interface
type LoanUserDeviceAppsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanUserDeviceApps, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanUserDeviceApps, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserDeviceApps, error)
	MultiSet(ctx context.Context, data []*model.LoanUserDeviceApps, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanUserDeviceAppsCache define a cache struct
type loanUserDeviceAppsCache struct {
	cache cache.Cache
}

// NewLoanUserDeviceAppsCache new a cache
func NewLoanUserDeviceAppsCache(cacheType *database.CacheType) LoanUserDeviceAppsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserDeviceApps{}
		})
		return &loanUserDeviceAppsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserDeviceApps{}
		})
		return &loanUserDeviceAppsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanUserDeviceAppsCacheKey cache key
func (c *loanUserDeviceAppsCache) GetLoanUserDeviceAppsCacheKey(id uint64) string {
	return loanUserDeviceAppsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanUserDeviceAppsCache) Set(ctx context.Context, id uint64, data *model.LoanUserDeviceApps, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanUserDeviceAppsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanUserDeviceAppsCache) Get(ctx context.Context, id uint64) (*model.LoanUserDeviceApps, error) {
	var data *model.LoanUserDeviceApps
	cacheKey := c.GetLoanUserDeviceAppsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanUserDeviceAppsCache) MultiSet(ctx context.Context, data []*model.LoanUserDeviceApps, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanUserDeviceAppsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanUserDeviceAppsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserDeviceApps, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanUserDeviceAppsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanUserDeviceApps)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanUserDeviceApps)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanUserDeviceAppsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanUserDeviceAppsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserDeviceAppsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanUserDeviceAppsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserDeviceAppsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanUserDeviceAppsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
