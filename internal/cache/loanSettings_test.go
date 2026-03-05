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

func newLoanSettingsCache() *gotest.Cache {
	record1 := &model.LoanSettings{}
	record1.ID = 1
	record2 := &model.LoanSettings{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanSettingsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanSettingsCache_Set(t *testing.T) {
	c := newLoanSettingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanSettings)
	err := c.ICache.(LoanSettingsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanSettingsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanSettingsCache_Get(t *testing.T) {
	c := newLoanSettingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanSettings)
	err := c.ICache.(LoanSettingsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanSettingsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanSettingsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanSettingsCache_MultiGet(t *testing.T) {
	c := newLoanSettingsCache()
	defer c.Close()

	var testData []*model.LoanSettings
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanSettings))
	}

	err := c.ICache.(LoanSettingsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanSettingsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanSettings))
	}
}

func Test_loanSettingsCache_MultiSet(t *testing.T) {
	c := newLoanSettingsCache()
	defer c.Close()

	var testData []*model.LoanSettings
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanSettings))
	}

	err := c.ICache.(LoanSettingsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanSettingsCache_Del(t *testing.T) {
	c := newLoanSettingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanSettings)
	err := c.ICache.(LoanSettingsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanSettingsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanSettingsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanSettings)
	err := c.ICache.(LoanSettingsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanSettingsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanSettingsCache(t *testing.T) {
	c := NewLoanSettingsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanSettingsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanSettingsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
