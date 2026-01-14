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

func newLoanPermissionsCache() *gotest.Cache {
	record1 := &model.LoanPermissions{}
	record1.ID = 1
	record2 := &model.LoanPermissions{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanPermissionsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanPermissionsCache_Set(t *testing.T) {
	c := newLoanPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPermissions)
	err := c.ICache.(LoanPermissionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanPermissionsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanPermissionsCache_Get(t *testing.T) {
	c := newLoanPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPermissions)
	err := c.ICache.(LoanPermissionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanPermissionsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanPermissionsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanPermissionsCache_MultiGet(t *testing.T) {
	c := newLoanPermissionsCache()
	defer c.Close()

	var testData []*model.LoanPermissions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanPermissions))
	}

	err := c.ICache.(LoanPermissionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanPermissionsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanPermissions))
	}
}

func Test_loanPermissionsCache_MultiSet(t *testing.T) {
	c := newLoanPermissionsCache()
	defer c.Close()

	var testData []*model.LoanPermissions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanPermissions))
	}

	err := c.ICache.(LoanPermissionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanPermissionsCache_Del(t *testing.T) {
	c := newLoanPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPermissions)
	err := c.ICache.(LoanPermissionsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanPermissionsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanPermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanPermissions)
	err := c.ICache.(LoanPermissionsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanPermissionsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanPermissionsCache(t *testing.T) {
	c := NewLoanPermissionsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanPermissionsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanPermissionsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
