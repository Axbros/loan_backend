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
	loanBaseinfoFilesCachePrefixKey = "loanBaseinfoFiles:"
	// LoanBaseinfoFilesExpireTime expire time
	LoanBaseinfoFilesExpireTime = 5 * time.Minute
)

var _ LoanBaseinfoFilesCache = (*loanBaseinfoFilesCache)(nil)

// LoanBaseinfoFilesCache cache interface
type LoanBaseinfoFilesCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanBaseinfoFiles, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanBaseinfoFiles, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfoFiles, error)
	MultiSet(ctx context.Context, data []*model.LoanBaseinfoFiles, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanBaseinfoFilesCache define a cache struct
type loanBaseinfoFilesCache struct {
	cache cache.Cache
}

// NewLoanBaseinfoFilesCache new a cache
func NewLoanBaseinfoFilesCache(cacheType *database.CacheType) LoanBaseinfoFilesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanBaseinfoFiles{}
		})
		return &loanBaseinfoFilesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanBaseinfoFiles{}
		})
		return &loanBaseinfoFilesCache{cache: c}
	}

	return nil // no cache
}

// GetLoanBaseinfoFilesCacheKey cache key
func (c *loanBaseinfoFilesCache) GetLoanBaseinfoFilesCacheKey(id uint64) string {
	return loanBaseinfoFilesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanBaseinfoFilesCache) Set(ctx context.Context, id uint64, data *model.LoanBaseinfoFiles, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanBaseinfoFilesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanBaseinfoFilesCache) Get(ctx context.Context, id uint64) (*model.LoanBaseinfoFiles, error) {
	var data *model.LoanBaseinfoFiles
	cacheKey := c.GetLoanBaseinfoFilesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanBaseinfoFilesCache) MultiSet(ctx context.Context, data []*model.LoanBaseinfoFiles, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanBaseinfoFilesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanBaseinfoFilesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanBaseinfoFiles, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanBaseinfoFilesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanBaseinfoFiles)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanBaseinfoFiles)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanBaseinfoFilesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanBaseinfoFilesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanBaseinfoFilesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanBaseinfoFilesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanBaseinfoFilesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanBaseinfoFilesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
