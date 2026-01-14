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

func newLoanMfaDevicesCache() *gotest.Cache {
	record1 := &model.LoanMfaDevices{}
	record1.ID = 1
	record2 := &model.LoanMfaDevices{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanMfaDevicesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanMfaDevicesCache_Set(t *testing.T) {
	c := newLoanMfaDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaDevices)
	err := c.ICache.(LoanMfaDevicesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanMfaDevicesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanMfaDevicesCache_Get(t *testing.T) {
	c := newLoanMfaDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaDevices)
	err := c.ICache.(LoanMfaDevicesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanMfaDevicesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanMfaDevicesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanMfaDevicesCache_MultiGet(t *testing.T) {
	c := newLoanMfaDevicesCache()
	defer c.Close()

	var testData []*model.LoanMfaDevices
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanMfaDevices))
	}

	err := c.ICache.(LoanMfaDevicesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanMfaDevicesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanMfaDevices))
	}
}

func Test_loanMfaDevicesCache_MultiSet(t *testing.T) {
	c := newLoanMfaDevicesCache()
	defer c.Close()

	var testData []*model.LoanMfaDevices
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanMfaDevices))
	}

	err := c.ICache.(LoanMfaDevicesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanMfaDevicesCache_Del(t *testing.T) {
	c := newLoanMfaDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaDevices)
	err := c.ICache.(LoanMfaDevicesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanMfaDevicesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanMfaDevicesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaDevices)
	err := c.ICache.(LoanMfaDevicesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanMfaDevicesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanMfaDevicesCache(t *testing.T) {
	c := NewLoanMfaDevicesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanMfaDevicesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanMfaDevicesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
