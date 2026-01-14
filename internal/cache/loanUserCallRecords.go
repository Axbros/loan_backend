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
	loanUserCallRecordsCachePrefixKey = "loanUserCallRecords:"
	// LoanUserCallRecordsExpireTime expire time
	LoanUserCallRecordsExpireTime = 5 * time.Minute
)

var _ LoanUserCallRecordsCache = (*loanUserCallRecordsCache)(nil)

// LoanUserCallRecordsCache cache interface
type LoanUserCallRecordsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanUserCallRecords, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanUserCallRecords, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserCallRecords, error)
	MultiSet(ctx context.Context, data []*model.LoanUserCallRecords, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanUserCallRecordsCache define a cache struct
type loanUserCallRecordsCache struct {
	cache cache.Cache
}

// NewLoanUserCallRecordsCache new a cache
func NewLoanUserCallRecordsCache(cacheType *database.CacheType) LoanUserCallRecordsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserCallRecords{}
		})
		return &loanUserCallRecordsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserCallRecords{}
		})
		return &loanUserCallRecordsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanUserCallRecordsCacheKey cache key
func (c *loanUserCallRecordsCache) GetLoanUserCallRecordsCacheKey(id uint64) string {
	return loanUserCallRecordsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanUserCallRecordsCache) Set(ctx context.Context, id uint64, data *model.LoanUserCallRecords, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanUserCallRecordsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanUserCallRecordsCache) Get(ctx context.Context, id uint64) (*model.LoanUserCallRecords, error) {
	var data *model.LoanUserCallRecords
	cacheKey := c.GetLoanUserCallRecordsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanUserCallRecordsCache) MultiSet(ctx context.Context, data []*model.LoanUserCallRecords, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanUserCallRecordsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanUserCallRecordsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserCallRecords, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanUserCallRecordsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanUserCallRecords)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanUserCallRecords)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanUserCallRecordsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanUserCallRecordsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserCallRecordsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanUserCallRecordsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserCallRecordsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanUserCallRecordsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
