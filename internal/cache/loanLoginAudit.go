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
	loanLoginAuditCachePrefixKey = "loanLoginAudit:"
	// LoanLoginAuditExpireTime expire time
	LoanLoginAuditExpireTime = 5 * time.Minute
)

var _ LoanLoginAuditCache = (*loanLoginAuditCache)(nil)

// LoanLoginAuditCache cache interface
type LoanLoginAuditCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanLoginAudit, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanLoginAudit, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanLoginAudit, error)
	MultiSet(ctx context.Context, data []*model.LoanLoginAudit, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanLoginAuditCache define a cache struct
type loanLoginAuditCache struct {
	cache cache.Cache
}

// NewLoanLoginAuditCache new a cache
func NewLoanLoginAuditCache(cacheType *database.CacheType) LoanLoginAuditCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanLoginAudit{}
		})
		return &loanLoginAuditCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanLoginAudit{}
		})
		return &loanLoginAuditCache{cache: c}
	}

	return nil // no cache
}

// GetLoanLoginAuditCacheKey cache key
func (c *loanLoginAuditCache) GetLoanLoginAuditCacheKey(id uint64) string {
	return loanLoginAuditCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanLoginAuditCache) Set(ctx context.Context, id uint64, data *model.LoanLoginAudit, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanLoginAuditCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanLoginAuditCache) Get(ctx context.Context, id uint64) (*model.LoanLoginAudit, error) {
	var data *model.LoanLoginAudit
	cacheKey := c.GetLoanLoginAuditCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanLoginAuditCache) MultiSet(ctx context.Context, data []*model.LoanLoginAudit, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanLoginAuditCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanLoginAuditCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanLoginAudit, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanLoginAuditCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanLoginAudit)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanLoginAudit)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanLoginAuditCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanLoginAuditCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanLoginAuditCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanLoginAuditCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanLoginAuditCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanLoginAuditCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
