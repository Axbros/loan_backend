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

func newLoanPaymentChannelsCache() *gotest.Cache {
	record1 := &model.LoanPaymentChannels{}
	record1.ID = 1
	record2 := &model.LoanPaymentChannels{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanPaymentChannelsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanPaymentChannelsCache_Set(t *testing.T) {
	c := newLoanPaymentChannelsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPaymentChannels)
	err := c.ICache.(LoanPaymentChannelsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanPaymentChannelsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanPaymentChannelsCache_Get(t *testing.T) {
	c := newLoanPaymentChannelsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPaymentChannels)
	err := c.ICache.(LoanPaymentChannelsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanPaymentChannelsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanPaymentChannelsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanPaymentChannelsCache_MultiGet(t *testing.T) {
	c := newLoanPaymentChannelsCache()
	defer c.Close()

	var testData []*model.LoanPaymentChannels
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanPaymentChannels))
	}

	err := c.ICache.(LoanPaymentChannelsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanPaymentChannelsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanPaymentChannels))
	}
}

func Test_loanPaymentChannelsCache_MultiSet(t *testing.T) {
	c := newLoanPaymentChannelsCache()
	defer c.Close()

	var testData []*model.LoanPaymentChannels
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanPaymentChannels))
	}

	err := c.ICache.(LoanPaymentChannelsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanPaymentChannelsCache_Del(t *testing.T) {
	c := newLoanPaymentChannelsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPaymentChannels)
	err := c.ICache.(LoanPaymentChannelsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanPaymentChannelsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanPaymentChannelsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPaymentChannels)
	err := c.ICache.(LoanPaymentChannelsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanPaymentChannelsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanPaymentChannelsCache(t *testing.T) {
	c := NewLoanPaymentChannelsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanPaymentChannelsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanPaymentChannelsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
