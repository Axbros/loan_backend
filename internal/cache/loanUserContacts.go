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
	loanUserContactsCachePrefixKey = "loanUserContacts:"
	// LoanUserContactsExpireTime expire time
	LoanUserContactsExpireTime = 5 * time.Minute
)

var _ LoanUserContactsCache = (*loanUserContactsCache)(nil)

// LoanUserContactsCache cache interface
type LoanUserContactsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanUserContacts, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanUserContacts, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserContacts, error)
	MultiSet(ctx context.Context, data []*model.LoanUserContacts, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanUserContactsCache define a cache struct
type loanUserContactsCache struct {
	cache cache.Cache
}

// NewLoanUserContactsCache new a cache
func NewLoanUserContactsCache(cacheType *database.CacheType) LoanUserContactsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserContacts{}
		})
		return &loanUserContactsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanUserContacts{}
		})
		return &loanUserContactsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanUserContactsCacheKey cache key
func (c *loanUserContactsCache) GetLoanUserContactsCacheKey(id uint64) string {
	return loanUserContactsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanUserContactsCache) Set(ctx context.Context, id uint64, data *model.LoanUserContacts, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanUserContactsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanUserContactsCache) Get(ctx context.Context, id uint64) (*model.LoanUserContacts, error) {
	var data *model.LoanUserContacts
	cacheKey := c.GetLoanUserContactsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanUserContactsCache) MultiSet(ctx context.Context, data []*model.LoanUserContacts, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanUserContactsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanUserContactsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanUserContacts, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanUserContactsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanUserContacts)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanUserContacts)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanUserContactsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanUserContactsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserContactsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanUserContactsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanUserContactsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanUserContactsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
