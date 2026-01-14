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
	loanCollectionCasesCachePrefixKey = "loanCollectionCases:"
	// LoanCollectionCasesExpireTime expire time
	LoanCollectionCasesExpireTime = 5 * time.Minute
)

var _ LoanCollectionCasesCache = (*loanCollectionCasesCache)(nil)

// LoanCollectionCasesCache cache interface
type LoanCollectionCasesCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanCollectionCases, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanCollectionCases, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanCollectionCases, error)
	MultiSet(ctx context.Context, data []*model.LoanCollectionCases, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanCollectionCasesCache define a cache struct
type loanCollectionCasesCache struct {
	cache cache.Cache
}

// NewLoanCollectionCasesCache new a cache
func NewLoanCollectionCasesCache(cacheType *database.CacheType) LoanCollectionCasesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanCollectionCases{}
		})
		return &loanCollectionCasesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanCollectionCases{}
		})
		return &loanCollectionCasesCache{cache: c}
	}

	return nil // no cache
}

// GetLoanCollectionCasesCacheKey cache key
func (c *loanCollectionCasesCache) GetLoanCollectionCasesCacheKey(id uint64) string {
	return loanCollectionCasesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanCollectionCasesCache) Set(ctx context.Context, id uint64, data *model.LoanCollectionCases, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanCollectionCasesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanCollectionCasesCache) Get(ctx context.Context, id uint64) (*model.LoanCollectionCases, error) {
	var data *model.LoanCollectionCases
	cacheKey := c.GetLoanCollectionCasesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanCollectionCasesCache) MultiSet(ctx context.Context, data []*model.LoanCollectionCases, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanCollectionCasesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanCollectionCasesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanCollectionCases, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanCollectionCasesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanCollectionCases)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanCollectionCases)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanCollectionCasesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanCollectionCasesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanCollectionCasesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanCollectionCasesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanCollectionCasesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanCollectionCasesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
