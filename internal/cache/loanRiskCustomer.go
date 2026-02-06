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
	loanRiskCustomerCachePrefixKey = "loanRiskCustomer:"
	// LoanRiskCustomerExpireTime expire time
	LoanRiskCustomerExpireTime = 5 * time.Minute
)

var _ LoanRiskCustomerCache = (*loanRiskCustomerCache)(nil)

// LoanRiskCustomerCache cache interface
type LoanRiskCustomerCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanRiskCustomer, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanRiskCustomer, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRiskCustomer, error)
	MultiSet(ctx context.Context, data []*model.LoanRiskCustomer, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanRiskCustomerCache define a cache struct
type loanRiskCustomerCache struct {
	cache cache.Cache
}

// NewLoanRiskCustomerCache new a cache
func NewLoanRiskCustomerCache(cacheType *database.CacheType) LoanRiskCustomerCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRiskCustomer{}
		})
		return &loanRiskCustomerCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRiskCustomer{}
		})
		return &loanRiskCustomerCache{cache: c}
	}

	return nil // no cache
}

// GetLoanRiskCustomerCacheKey cache key
func (c *loanRiskCustomerCache) GetLoanRiskCustomerCacheKey(id uint64) string {
	return loanRiskCustomerCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanRiskCustomerCache) Set(ctx context.Context, id uint64, data *model.LoanRiskCustomer, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanRiskCustomerCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanRiskCustomerCache) Get(ctx context.Context, id uint64) (*model.LoanRiskCustomer, error) {
	var data *model.LoanRiskCustomer
	cacheKey := c.GetLoanRiskCustomerCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanRiskCustomerCache) MultiSet(ctx context.Context, data []*model.LoanRiskCustomer, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanRiskCustomerCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanRiskCustomerCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRiskCustomer, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanRiskCustomerCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanRiskCustomer)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanRiskCustomer)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanRiskCustomerCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanRiskCustomerCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRiskCustomerCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanRiskCustomerCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRiskCustomerCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanRiskCustomerCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
