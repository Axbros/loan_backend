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
	loanDisbursementsCachePrefixKey = "loanDisbursements:"
	// LoanDisbursementsExpireTime expire time
	LoanDisbursementsExpireTime = 5 * time.Minute
)

var _ LoanDisbursementsCache = (*loanDisbursementsCache)(nil)

// LoanDisbursementsCache cache interface
type LoanDisbursementsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanDisbursements, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanDisbursements, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDisbursements, error)
	MultiSet(ctx context.Context, data []*model.LoanDisbursements, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanDisbursementsCache define a cache struct
type loanDisbursementsCache struct {
	cache cache.Cache
}

// NewLoanDisbursementsCache new a cache
func NewLoanDisbursementsCache(cacheType *database.CacheType) LoanDisbursementsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanDisbursements{}
		})
		return &loanDisbursementsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanDisbursements{}
		})
		return &loanDisbursementsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanDisbursementsCacheKey cache key
func (c *loanDisbursementsCache) GetLoanDisbursementsCacheKey(id uint64) string {
	return loanDisbursementsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanDisbursementsCache) Set(ctx context.Context, id uint64, data *model.LoanDisbursements, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanDisbursementsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanDisbursementsCache) Get(ctx context.Context, id uint64) (*model.LoanDisbursements, error) {
	var data *model.LoanDisbursements
	cacheKey := c.GetLoanDisbursementsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanDisbursementsCache) MultiSet(ctx context.Context, data []*model.LoanDisbursements, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanDisbursementsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanDisbursementsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanDisbursements, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanDisbursementsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanDisbursements)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanDisbursements)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanDisbursementsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanDisbursementsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanDisbursementsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanDisbursementsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanDisbursementsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanDisbursementsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
