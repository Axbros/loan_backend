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
	loanPaymentChannelsCachePrefixKey = "loanPaymentChannels:"
	// LoanPaymentChannelsExpireTime expire time
	LoanPaymentChannelsExpireTime = 5 * time.Minute
)

var _ LoanPaymentChannelsCache = (*loanPaymentChannelsCache)(nil)

// LoanPaymentChannelsCache cache interface
type LoanPaymentChannelsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanPaymentChannels, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanPaymentChannels, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanPaymentChannels, error)
	MultiSet(ctx context.Context, data []*model.LoanPaymentChannels, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanPaymentChannelsCache define a cache struct
type loanPaymentChannelsCache struct {
	cache cache.Cache
}

// NewLoanPaymentChannelsCache new a cache
func NewLoanPaymentChannelsCache(cacheType *database.CacheType) LoanPaymentChannelsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanPaymentChannels{}
		})
		return &loanPaymentChannelsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanPaymentChannels{}
		})
		return &loanPaymentChannelsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanPaymentChannelsCacheKey cache key
func (c *loanPaymentChannelsCache) GetLoanPaymentChannelsCacheKey(id uint64) string {
	return loanPaymentChannelsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanPaymentChannelsCache) Set(ctx context.Context, id uint64, data *model.LoanPaymentChannels, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanPaymentChannelsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanPaymentChannelsCache) Get(ctx context.Context, id uint64) (*model.LoanPaymentChannels, error) {
	var data *model.LoanPaymentChannels
	cacheKey := c.GetLoanPaymentChannelsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanPaymentChannelsCache) MultiSet(ctx context.Context, data []*model.LoanPaymentChannels, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanPaymentChannelsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanPaymentChannelsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanPaymentChannels, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanPaymentChannelsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanPaymentChannels)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanPaymentChannels)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanPaymentChannelsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanPaymentChannelsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanPaymentChannelsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanPaymentChannelsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanPaymentChannelsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanPaymentChannelsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
