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
	loanSettingsCachePrefixKey = "loanSettings:"
	// LoanSettingsExpireTime expire time
	LoanSettingsExpireTime = 5 * time.Minute
)

var _ LoanSettingsCache = (*loanSettingsCache)(nil)

// LoanSettingsCache cache interface
type LoanSettingsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanSettings, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanSettings, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanSettings, error)
	MultiSet(ctx context.Context, data []*model.LoanSettings, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanSettingsCache define a cache struct
type loanSettingsCache struct {
	cache cache.Cache
}

// NewLoanSettingsCache new a cache
func NewLoanSettingsCache(cacheType *database.CacheType) LoanSettingsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanSettings{}
		})
		return &loanSettingsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanSettings{}
		})
		return &loanSettingsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanSettingsCacheKey cache key
func (c *loanSettingsCache) GetLoanSettingsCacheKey(id uint64) string {
	return loanSettingsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanSettingsCache) Set(ctx context.Context, id uint64, data *model.LoanSettings, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanSettingsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanSettingsCache) Get(ctx context.Context, id uint64) (*model.LoanSettings, error) {
	var data *model.LoanSettings
	cacheKey := c.GetLoanSettingsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanSettingsCache) MultiSet(ctx context.Context, data []*model.LoanSettings, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanSettingsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanSettingsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanSettings, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanSettingsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanSettings)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanSettings)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanSettingsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanSettingsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanSettingsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanSettingsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanSettingsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanSettingsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
