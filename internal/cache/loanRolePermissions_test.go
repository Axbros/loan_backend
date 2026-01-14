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

func newLoanRolePermissionsCache() *gotest.Cache {
	record1 := &model.LoanRolePermissions{}
	record1.ID = 1
	record2 := &model.LoanRolePermissions{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanRolePermissionsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanRolePermissionsCache_Set(t *testing.T) {
	c := newLoanRolePermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRolePermissions)
	err := c.ICache.(LoanRolePermissionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanRolePermissionsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanRolePermissionsCache_Get(t *testing.T) {
	c := newLoanRolePermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRolePermissions)
	err := c.ICache.(LoanRolePermissionsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRolePermissionsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanRolePermissionsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanRolePermissionsCache_MultiGet(t *testing.T) {
	c := newLoanRolePermissionsCache()
	defer c.Close()

	var testData []*model.LoanRolePermissions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRolePermissions))
	}

	err := c.ICache.(LoanRolePermissionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRolePermissionsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanRolePermissions))
	}
}

func Test_loanRolePermissionsCache_MultiSet(t *testing.T) {
	c := newLoanRolePermissionsCache()
	defer c.Close()

	var testData []*model.LoanRolePermissions
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRolePermissions))
	}

	err := c.ICache.(LoanRolePermissionsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRolePermissionsCache_Del(t *testing.T) {
	c := newLoanRolePermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRolePermissions)
	err := c.ICache.(LoanRolePermissionsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRolePermissionsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanRolePermissionsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRolePermissions)
	err := c.ICache.(LoanRolePermissionsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanRolePermissionsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanRolePermissionsCache(t *testing.T) {
	c := NewLoanRolePermissionsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanRolePermissionsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanRolePermissionsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
