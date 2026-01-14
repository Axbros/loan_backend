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

func newLoanUserRolesCache() *gotest.Cache {
	record1 := &model.LoanUserRoles{}
	record1.ID = 1
	record2 := &model.LoanUserRoles{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanUserRolesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanUserRolesCache_Set(t *testing.T) {
	c := newLoanUserRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserRoles)
	err := c.ICache.(LoanUserRolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanUserRolesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanUserRolesCache_Get(t *testing.T) {
	c := newLoanUserRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserRoles)
	err := c.ICache.(LoanUserRolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserRolesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanUserRolesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanUserRolesCache_MultiGet(t *testing.T) {
	c := newLoanUserRolesCache()
	defer c.Close()

	var testData []*model.LoanUserRoles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserRoles))
	}

	err := c.ICache.(LoanUserRolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserRolesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanUserRoles))
	}
}

func Test_loanUserRolesCache_MultiSet(t *testing.T) {
	c := newLoanUserRolesCache()
	defer c.Close()

	var testData []*model.LoanUserRoles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserRoles))
	}

	err := c.ICache.(LoanUserRolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserRolesCache_Del(t *testing.T) {
	c := newLoanUserRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserRoles)
	err := c.ICache.(LoanUserRolesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserRolesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanUserRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserRoles)
	err := c.ICache.(LoanUserRolesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanUserRolesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanUserRolesCache(t *testing.T) {
	c := NewLoanUserRolesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanUserRolesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanUserRolesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
