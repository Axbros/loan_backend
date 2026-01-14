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
	loanRepaymentSchedulesCachePrefixKey = "loanRepaymentSchedules:"
	// LoanRepaymentSchedulesExpireTime expire time
	LoanRepaymentSchedulesExpireTime = 5 * time.Minute
)

var _ LoanRepaymentSchedulesCache = (*loanRepaymentSchedulesCache)(nil)

// LoanRepaymentSchedulesCache cache interface
type LoanRepaymentSchedulesCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanRepaymentSchedules, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanRepaymentSchedules, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRepaymentSchedules, error)
	MultiSet(ctx context.Context, data []*model.LoanRepaymentSchedules, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanRepaymentSchedulesCache define a cache struct
type loanRepaymentSchedulesCache struct {
	cache cache.Cache
}

// NewLoanRepaymentSchedulesCache new a cache
func NewLoanRepaymentSchedulesCache(cacheType *database.CacheType) LoanRepaymentSchedulesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRepaymentSchedules{}
		})
		return &loanRepaymentSchedulesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRepaymentSchedules{}
		})
		return &loanRepaymentSchedulesCache{cache: c}
	}

	return nil // no cache
}

// GetLoanRepaymentSchedulesCacheKey cache key
func (c *loanRepaymentSchedulesCache) GetLoanRepaymentSchedulesCacheKey(id uint64) string {
	return loanRepaymentSchedulesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanRepaymentSchedulesCache) Set(ctx context.Context, id uint64, data *model.LoanRepaymentSchedules, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanRepaymentSchedulesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanRepaymentSchedulesCache) Get(ctx context.Context, id uint64) (*model.LoanRepaymentSchedules, error) {
	var data *model.LoanRepaymentSchedules
	cacheKey := c.GetLoanRepaymentSchedulesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanRepaymentSchedulesCache) MultiSet(ctx context.Context, data []*model.LoanRepaymentSchedules, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanRepaymentSchedulesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanRepaymentSchedulesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRepaymentSchedules, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanRepaymentSchedulesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanRepaymentSchedules)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanRepaymentSchedules)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanRepaymentSchedulesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanRepaymentSchedulesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRepaymentSchedulesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanRepaymentSchedulesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRepaymentSchedulesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanRepaymentSchedulesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
