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
	loanReferralVisitsCachePrefixKey = "loanReferralVisits:"
	// LoanReferralVisitsExpireTime expire time
	LoanReferralVisitsExpireTime = 5 * time.Minute
)

var _ LoanReferralVisitsCache = (*loanReferralVisitsCache)(nil)

// LoanReferralVisitsCache cache interface
type LoanReferralVisitsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanReferralVisits, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanReferralVisits, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanReferralVisits, error)
	MultiSet(ctx context.Context, data []*model.LoanReferralVisits, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanReferralVisitsCache define a cache struct
type loanReferralVisitsCache struct {
	cache cache.Cache
}

// NewLoanReferralVisitsCache new a cache
func NewLoanReferralVisitsCache(cacheType *database.CacheType) LoanReferralVisitsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanReferralVisits{}
		})
		return &loanReferralVisitsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanReferralVisits{}
		})
		return &loanReferralVisitsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanReferralVisitsCacheKey cache key
func (c *loanReferralVisitsCache) GetLoanReferralVisitsCacheKey(id uint64) string {
	return loanReferralVisitsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanReferralVisitsCache) Set(ctx context.Context, id uint64, data *model.LoanReferralVisits, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanReferralVisitsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanReferralVisitsCache) Get(ctx context.Context, id uint64) (*model.LoanReferralVisits, error) {
	var data *model.LoanReferralVisits
	cacheKey := c.GetLoanReferralVisitsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanReferralVisitsCache) MultiSet(ctx context.Context, data []*model.LoanReferralVisits, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanReferralVisitsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanReferralVisitsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanReferralVisits, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanReferralVisitsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanReferralVisits)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanReferralVisits)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanReferralVisitsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanReferralVisitsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanReferralVisitsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanReferralVisitsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanReferralVisitsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanReferralVisitsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
