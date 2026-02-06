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

func newLoanRiskCustomerCache() *gotest.Cache {
	record1 := &model.LoanRiskCustomer{}
	record1.ID = 1
	record2 := &model.LoanRiskCustomer{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanRiskCustomerCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanRiskCustomerCache_Set(t *testing.T) {
	c := newLoanRiskCustomerCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRiskCustomer)
	err := c.ICache.(LoanRiskCustomerCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanRiskCustomerCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanRiskCustomerCache_Get(t *testing.T) {
	c := newLoanRiskCustomerCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRiskCustomer)
	err := c.ICache.(LoanRiskCustomerCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRiskCustomerCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanRiskCustomerCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanRiskCustomerCache_MultiGet(t *testing.T) {
	c := newLoanRiskCustomerCache()
	defer c.Close()

	var testData []*model.LoanRiskCustomer
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRiskCustomer))
	}

	err := c.ICache.(LoanRiskCustomerCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRiskCustomerCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanRiskCustomer))
	}
}

func Test_loanRiskCustomerCache_MultiSet(t *testing.T) {
	c := newLoanRiskCustomerCache()
	defer c.Close()

	var testData []*model.LoanRiskCustomer
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRiskCustomer))
	}

	err := c.ICache.(LoanRiskCustomerCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRiskCustomerCache_Del(t *testing.T) {
	c := newLoanRiskCustomerCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRiskCustomer)
	err := c.ICache.(LoanRiskCustomerCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRiskCustomerCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanRiskCustomerCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRiskCustomer)
	err := c.ICache.(LoanRiskCustomerCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanRiskCustomerCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanRiskCustomerCache(t *testing.T) {
	c := NewLoanRiskCustomerCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanRiskCustomerCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanRiskCustomerCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
