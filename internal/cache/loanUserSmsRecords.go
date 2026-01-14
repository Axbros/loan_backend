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
	loanUserSmsRecordsCachePrefixKey = "loanUserSmsRecords:"
	// LoanUserSmsRecordsExpireTime expire time
	LoanUserSmsRecordsExpireTime = 5 * time.Minute
)

var _ LoanUserSmsRecordsCache = (*loanUserSmsRecordsCache)(nil)

// LoanUserSmsRecordsCache cache interface
type LoanUserSmsRecordsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanUserSmsRecords, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanUserSmsRecords, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserSmsRecords, error)
	MultiSet(ctx context.Context, data []*model.LoanUserSmsRecords, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanUserSmsRecordsCache define a cache struct
type loanUserSmsRecordsCache struct {
	cache cache.Cache
}

// NewLoanUserSmsRecordsCache new a cache
func NewLoanUserSmsRecordsCache(cacheType *database.CacheType) LoanUserSmsRecordsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserSmsRecords{}
		})
		return &loanUserSmsRecordsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserSmsRecords{}
		})
		return &loanUserSmsRecordsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanUserSmsRecordsCacheKey cache key
func (c *loanUserSmsRecordsCache) GetLoanUserSmsRecordsCacheKey(id uint64) string {
	return loanUserSmsRecordsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanUserSmsRecordsCache) Set(ctx context.Context, id uint64, data *model.LoanUserSmsRecords, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanUserSmsRecordsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanUserSmsRecordsCache) Get(ctx context.Context, id uint64) (*model.LoanUserSmsRecords, error) {
	var data *model.LoanUserSmsRecords
	cacheKey := c.GetLoanUserSmsRecordsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanUserSmsRecordsCache) MultiSet(ctx context.Context, data []*model.LoanUserSmsRecords, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanUserSmsRecordsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanUserSmsRecordsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserSmsRecords, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanUserSmsRecordsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanUserSmsRecords)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanUserSmsRecords)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanUserSmsRecordsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanUserSmsRecordsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserSmsRecordsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanUserSmsRecordsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserSmsRecordsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanUserSmsRecordsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
