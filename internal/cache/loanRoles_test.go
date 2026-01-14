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

func newLoanRolesCache() *gotest.Cache {
	record1 := &model.LoanRoles{}
	record1.ID = 1
	record2 := &model.LoanRoles{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanRolesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanRolesCache_Set(t *testing.T) {
	c := newLoanRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoles)
	err := c.ICache.(LoanRolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanRolesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanRolesCache_Get(t *testing.T) {
	c := newLoanRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoles)
	err := c.ICache.(LoanRolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRolesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanRolesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanRolesCache_MultiGet(t *testing.T) {
	c := newLoanRolesCache()
	defer c.Close()

	var testData []*model.LoanRoles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRoles))
	}

	err := c.ICache.(LoanRolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRolesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanRoles))
	}
}

func Test_loanRolesCache_MultiSet(t *testing.T) {
	c := newLoanRolesCache()
	defer c.Close()

	var testData []*model.LoanRoles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRoles))
	}

	err := c.ICache.(LoanRolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRolesCache_Del(t *testing.T) {
	c := newLoanRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoles)
	err := c.ICache.(LoanRolesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRolesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoles)
	err := c.ICache.(LoanRolesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanRolesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanRolesCache(t *testing.T) {
	c := NewLoanRolesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanRolesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanRolesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
