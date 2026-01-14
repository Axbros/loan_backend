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
	loanUsersCachePrefixKey = "loanUsers:"
	// LoanUsersExpireTime expire time
	LoanUsersExpireTime = 5 * time.Minute
)

var _ LoanUsersCache = (*loanUsersCache)(nil)

// LoanUsersCache cache interface
type LoanUsersCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanUsers, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanUsers, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUsers, error)
	MultiSet(ctx context.Context, data []*model.LoanUsers, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanUsersCache define a cache struct
type loanUsersCache struct {
	cache cache.Cache
}

// NewLoanUsersCache new a cache
func NewLoanUsersCache(cacheType *database.CacheType) LoanUsersCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUsers{}
		})
		return &loanUsersCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUsers{}
		})
		return &loanUsersCache{cache: c}
	}

	return nil // no cache
}

// GetLoanUsersCacheKey cache key
func (c *loanUsersCache) GetLoanUsersCacheKey(id uint64) string {
	return loanUsersCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanUsersCache) Set(ctx context.Context, id uint64, data *model.LoanUsers, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanUsersCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanUsersCache) Get(ctx context.Context, id uint64) (*model.LoanUsers, error) {
	var data *model.LoanUsers
	cacheKey := c.GetLoanUsersCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanUsersCache) MultiSet(ctx context.Context, data []*model.LoanUsers, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanUsersCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanUsersCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUsers, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanUsersCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanUsers)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanUsers)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanUsersCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanUsersCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUsersCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanUsersCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUsersCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanUsersCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
