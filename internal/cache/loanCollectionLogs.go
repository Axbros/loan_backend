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
	loanCollectionLogsCachePrefixKey = "loanCollectionLogs:"
	// LoanCollectionLogsExpireTime expire time
	LoanCollectionLogsExpireTime = 5 * time.Minute
)

var _ LoanCollectionLogsCache = (*loanCollectionLogsCache)(nil)

// LoanCollectionLogsCache cache interface
type LoanCollectionLogsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanCollectionLogs, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanCollectionLogs, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanCollectionLogs, error)
	MultiSet(ctx context.Context, data []*model.LoanCollectionLogs, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanCollectionLogsCache define a cache struct
type loanCollectionLogsCache struct {
	cache cache.Cache
}

// NewLoanCollectionLogsCache new a cache
func NewLoanCollectionLogsCache(cacheType *database.CacheType) LoanCollectionLogsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanCollectionLogs{}
		})
		return &loanCollectionLogsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanCollectionLogs{}
		})
		return &loanCollectionLogsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanCollectionLogsCacheKey cache key
func (c *loanCollectionLogsCache) GetLoanCollectionLogsCacheKey(id uint64) string {
	return loanCollectionLogsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanCollectionLogsCache) Set(ctx context.Context, id uint64, data *model.LoanCollectionLogs, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanCollectionLogsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanCollectionLogsCache) Get(ctx context.Context, id uint64) (*model.LoanCollectionLogs, error) {
	var data *model.LoanCollectionLogs
	cacheKey := c.GetLoanCollectionLogsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanCollectionLogsCache) MultiSet(ctx context.Context, data []*model.LoanCollectionLogs, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanCollectionLogsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanCollectionLogsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanCollectionLogs, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanCollectionLogsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanCollectionLogs)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanCollectionLogs)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanCollectionLogsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanCollectionLogsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanCollectionLogsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanCollectionLogsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanCollectionLogsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanCollectionLogsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
