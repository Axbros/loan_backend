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
	loanRepaymentTransactionsCachePrefixKey = "loanRepaymentTransactions:"
	// LoanRepaymentTransactionsExpireTime expire time
	LoanRepaymentTransactionsExpireTime = 5 * time.Minute
)

var _ LoanRepaymentTransactionsCache = (*loanRepaymentTransactionsCache)(nil)

// LoanRepaymentTransactionsCache cache interface
type LoanRepaymentTransactionsCache interface {
	Set(ctx context.Context, id uint64, data *model.LoanRepaymentTransactions, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.LoanRepaymentTransactions, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRepaymentTransactions, error)
	MultiSet(ctx context.Context, data []*model.LoanRepaymentTransactions, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// loanRepaymentTransactionsCache define a cache struct
type loanRepaymentTransactionsCache struct {
	cache cache.Cache
}

// NewLoanRepaymentTransactionsCache new a cache
func NewLoanRepaymentTransactionsCache(cacheType *database.CacheType) LoanRepaymentTransactionsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRepaymentTransactions{}
		})
		return &loanRepaymentTransactionsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.LoanRepaymentTransactions{}
		})
		return &loanRepaymentTransactionsCache{cache: c}
	}

	return nil // no cache
}

// GetLoanRepaymentTransactionsCacheKey cache key
func (c *loanRepaymentTransactionsCache) GetLoanRepaymentTransactionsCacheKey(id uint64) string {
	return loanRepaymentTransactionsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *loanRepaymentTransactionsCache) Set(ctx context.Context, id uint64, data *model.LoanRepaymentTransactions, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLoanRepaymentTransactionsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *loanRepaymentTransactionsCache) Get(ctx context.Context, id uint64) (*model.LoanRepaymentTransactions, error) {
	var data *model.LoanRepaymentTransactions
	cacheKey := c.GetLoanRepaymentTransactionsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *loanRepaymentTransactionsCache) MultiSet(ctx context.Context, data []*model.LoanRepaymentTransactions, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLoanRepaymentTransactionsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *loanRepaymentTransactionsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.LoanRepaymentTransactions, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLoanRepaymentTransactionsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.LoanRepaymentTransactions)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.LoanRepaymentTransactions)
	for _, id := range ids {
		val, ok := itemMap[c.GetLoanRepaymentTransactionsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *loanRepaymentTransactionsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRepaymentTransactionsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *loanRepaymentTransactionsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetLoanRepaymentTransactionsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *loanRepaymentTransactionsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
