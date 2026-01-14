package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"loan/internal/database"
	"loan/internal/model"
)

func newLoanRepaymentTransactionsCache() *gotest.Cache {
	record1 := &model.LoanRepaymentTransactions{}
	record1.ID = 1
	record2 := &model.LoanRepaymentTransactions{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanRepaymentTransactionsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanRepaymentTransactionsCache_Set(t *testing.T) {
	c := newLoanRepaymentTransactionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRepaymentTransactions)
	err := c.ICache.(LoanRepaymentTransactionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanRepaymentTransactionsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanRepaymentTransactionsCache_Get(t *testing.T) {
	c := newLoanRepaymentTransactionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRepaymentTransactions)
	err := c.ICache.(LoanRepaymentTransactionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRepaymentTransactionsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanRepaymentTransactionsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanRepaymentTransactionsCache_MultiGet(t *testing.T) {
	c := newLoanRepaymentTransactionsCache()
	defer c.Close()

	var testData []*model.LoanRepaymentTransactions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRepaymentTransactions))
	}

	err := c.ICache.(LoanRepaymentTransactionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRepaymentTransactionsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanRepaymentTransactions))
	}
}

func Test_loanRepaymentTransactionsCache_MultiSet(t *testing.T) {
	c := newLoanRepaymentTransactionsCache()
	defer c.Close()

	var testData []*model.LoanRepaymentTransactions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRepaymentTransactions))
	}

	err := c.ICache.(LoanRepaymentTransactionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRepaymentTransactionsCache_Del(t *testing.T) {
	c := newLoanRepaymentTransactionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRepaymentTransactions)
	err := c.ICache.(LoanRepaymentTransactionsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRepaymentTransactionsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanRepaymentTransactionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRepaymentTransactions)
	err := c.ICache.(LoanRepaymentTransactionsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanRepaymentTransactionsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanRepaymentTransactionsCache(t *testing.T) {
	c := NewLoanRepaymentTransactionsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanRepaymentTransactionsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanRepaymentTransactionsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
