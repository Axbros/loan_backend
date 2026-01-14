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
	loanBaseinfoCachePrefixKey = "loanBaseinfo:"
	// LoanBaseinfoExpireTime expire time
	LoanBaseinfoExpireTime = 5 * time.Minute
)

var _ LoanBaseinfoCache = (*loanBaseinfoCache)(nil)

// LoanBaseinfoCache cache interface
type LoanBaseinfoCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanBaseinfo, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanBaseinfo, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfo, error)
	MultiSet(ctx context.Context, data []*model.LoanBaseinfo, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanBaseinfoCache define a cache struct
type loanBaseinfoCache struct {
	cache cache.Cache
}

// NewLoanBaseinfoCache new a cache
func NewLoanBaseinfoCache(cacheType *database.CacheType) LoanBaseinfoCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanBaseinfo{}
		})
		return &loanBaseinfoCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanBaseinfo{}
		})
		return &loanBaseinfoCache{cache: c}
	}

	return nil // no cache
}

// GetLoanBaseinfoCacheKey cache key
func (c *loanBaseinfoCache) GetLoanBaseinfoCacheKey(id uint64) string {
	return loanBaseinfoCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanBaseinfoCache) Set(ctx context.Context, id uint64, data *model.LoanBaseinfo, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanBaseinfoCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanBaseinfoCache) Get(ctx context.Context, id uint64) (*model.LoanBaseinfo, error) {
	var data *model.LoanBaseinfo
	cacheKey := c.GetLoanBaseinfoCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanBaseinfoCache) MultiSet(ctx context.Context, data []*model.LoanBaseinfo, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanBaseinfoCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanBaseinfoCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfo, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanBaseinfoCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanBaseinfo)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanBaseinfo)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanBaseinfoCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanBaseinfoCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanBaseinfoCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanBaseinfoCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanBaseinfoCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanBaseinfoCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
