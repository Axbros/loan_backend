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

func newLoanDepartmentRolesCache() *gotest.Cache {
	record1 := &model.LoanDepartmentRoles{}
	record1.ID = 1
	record2 := &model.LoanDepartmentRoles{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanDepartmentRolesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanDepartmentRolesCache_Set(t *testing.T) {
	c := newLoanDepartmentRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartmentRoles)
	err := c.ICache.(LoanDepartmentRolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanDepartmentRolesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanDepartmentRolesCache_Get(t *testing.T) {
	c := newLoanDepartmentRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartmentRoles)
	err := c.ICache.(LoanDepartmentRolesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanDepartmentRolesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanDepartmentRolesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanDepartmentRolesCache_MultiGet(t *testing.T) {
	c := newLoanDepartmentRolesCache()
	defer c.Close()

	var testData []*model.LoanDepartmentRoles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanDepartmentRoles))
	}

	err := c.ICache.(LoanDepartmentRolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanDepartmentRolesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanDepartmentRoles))
	}
}

func Test_loanDepartmentRolesCache_MultiSet(t *testing.T) {
	c := newLoanDepartmentRolesCache()
	defer c.Close()

	var testData []*model.LoanDepartmentRoles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanDepartmentRoles))
	}

	err := c.ICache.(LoanDepartmentRolesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanDepartmentRolesCache_Del(t *testing.T) {
	c := newLoanDepartmentRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartmentRoles)
	err := c.ICache.(LoanDepartmentRolesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanDepartmentRolesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanDepartmentRolesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartmentRoles)
	err := c.ICache.(LoanDepartmentRolesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanDepartmentRolesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanDepartmentRolesCache(t *testing.T) {
	c := NewLoanDepartmentRolesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanDepartmentRolesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanDepartmentRolesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
