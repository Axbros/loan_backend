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

func newLoanUserDeviceAppsCache() *gotest.Cache {
	record1 := &model.LoanUserDeviceApps{}
	record1.ID = 1
	record2 := &model.LoanUserDeviceApps{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanUserDeviceAppsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanUserDeviceAppsCache_Set(t *testing.T) {
	c := newLoanUserDeviceAppsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserDeviceApps)
	err := c.ICache.(LoanUserDeviceAppsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanUserDeviceAppsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanUserDeviceAppsCache_Get(t *testing.T) {
	c := newLoanUserDeviceAppsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserDeviceApps)
	err := c.ICache.(LoanUserDeviceAppsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserDeviceAppsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanUserDeviceAppsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanUserDeviceAppsCache_MultiGet(t *testing.T) {
	c := newLoanUserDeviceAppsCache()
	defer c.Close()

	var testData []*model.LoanUserDeviceApps
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserDeviceApps))
	}

	err := c.ICache.(LoanUserDeviceAppsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserDeviceAppsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanUserDeviceApps))
	}
}

func Test_loanUserDeviceAppsCache_MultiSet(t *testing.T) {
	c := newLoanUserDeviceAppsCache()
	defer c.Close()

	var testData []*model.LoanUserDeviceApps
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserDeviceApps))
	}

	err := c.ICache.(LoanUserDeviceAppsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserDeviceAppsCache_Del(t *testing.T) {
	c := newLoanUserDeviceAppsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserDeviceApps)
	err := c.ICache.(LoanUserDeviceAppsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserDeviceAppsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanUserDeviceAppsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserDeviceApps)
	err := c.ICache.(LoanUserDeviceAppsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanUserDeviceAppsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanUserDeviceAppsCache(t *testing.T) {
	c := NewLoanUserDeviceAppsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanUserDeviceAppsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanUserDeviceAppsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
