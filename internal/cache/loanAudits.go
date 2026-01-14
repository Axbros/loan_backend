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
	loanAuditsCachePrefixKey = "loanAudits:"
	// LoanAuditsExpireTime expire time
	LoanAuditsExpireTime = 5 * time.Minute
)

var _ LoanAuditsCache = (*loanAuditsCache)(nil)

// LoanAuditsCache cache interface
type LoanAuditsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanAudits, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanAudits, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanAudits, error)
	MultiSet(ctx context.Context, data []*model.LoanAudits, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanAuditsCache define a cache struct
type loanAuditsCache struct {
	cache cache.Cache
}

// NewLoanAuditsCache new a cache
func NewLoanAuditsCache(cacheType *database.CacheType) LoanAuditsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanAudits{}
		})
		return &loanAuditsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanAudits{}
		})
		return &loanAuditsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanAuditsCacheKey cache key
func (c *loanAuditsCache) GetLoanAuditsCacheKey(id uint64) string {
	return loanAuditsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanAuditsCache) Set(ctx context.Context, id uint64, data *model.LoanAudits, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanAuditsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanAuditsCache) Get(ctx context.Context, id uint64) (*model.LoanAudits, error) {
	var data *model.LoanAudits
	cacheKey := c.GetLoanAuditsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanAuditsCache) MultiSet(ctx context.Context, data []*model.LoanAudits, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanAuditsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanAuditsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanAudits, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanAuditsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanAudits)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanAudits)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanAuditsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanAuditsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanAuditsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanAuditsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanAuditsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanAuditsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
