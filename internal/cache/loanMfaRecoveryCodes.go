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
	loanMfaRecoveryCodesCachePrefixKey = "loanMfaRecoveryCodes:"
	// LoanMfaRecoveryCodesExpireTime expire time
	LoanMfaRecoveryCodesExpireTime = 5 * time.Minute
)

var _ LoanMfaRecoveryCodesCache = (*loanMfaRecoveryCodesCache)(nil)

// LoanMfaRecoveryCodesCache cache interface
type LoanMfaRecoveryCodesCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanMfaRecoveryCodes, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanMfaRecoveryCodes, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanMfaRecoveryCodes, error)
	MultiSet(ctx context.Context, data []*model.LoanMfaRecoveryCodes, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanMfaRecoveryCodesCache define a cache struct
type loanMfaRecoveryCodesCache struct {
	cache cache.Cache
}

// NewLoanMfaRecoveryCodesCache new a cache
func NewLoanMfaRecoveryCodesCache(cacheType *database.CacheType) LoanMfaRecoveryCodesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanMfaRecoveryCodes{}
		})
		return &loanMfaRecoveryCodesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanMfaRecoveryCodes{}
		})
		return &loanMfaRecoveryCodesCache{cache: c}
	}

	return nil // no cache
}

// GetLoanMfaRecoveryCodesCacheKey cache key
func (c *loanMfaRecoveryCodesCache) GetLoanMfaRecoveryCodesCacheKey(id uint64) string {
	return loanMfaRecoveryCodesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanMfaRecoveryCodesCache) Set(ctx context.Context, id uint64, data *model.LoanMfaRecoveryCodes, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanMfaRecoveryCodesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanMfaRecoveryCodesCache) Get(ctx context.Context, id uint64) (*model.LoanMfaRecoveryCodes, error) {
	var data *model.LoanMfaRecoveryCodes
	cacheKey := c.GetLoanMfaRecoveryCodesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanMfaRecoveryCodesCache) MultiSet(ctx context.Context, data []*model.LoanMfaRecoveryCodes, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanMfaRecoveryCodesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanMfaRecoveryCodesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanMfaRecoveryCodes, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanMfaRecoveryCodesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanMfaRecoveryCodes)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanMfaRecoveryCodes)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanMfaRecoveryCodesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanMfaRecoveryCodesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanMfaRecoveryCodesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanMfaRecoveryCodesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanMfaRecoveryCodesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanMfaRecoveryCodesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
