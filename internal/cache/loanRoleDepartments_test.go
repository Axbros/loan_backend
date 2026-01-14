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

func newLoanRoleDepartmentsCache() *gotest.Cache {
	record1 := &model.LoanRoleDepartments{}
	record1.ID = 1
	record2 := &model.LoanRoleDepartments{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanRoleDepartmentsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanRoleDepartmentsCache_Set(t *testing.T) {
	c := newLoanRoleDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoleDepartments)
	err := c.ICache.(LoanRoleDepartmentsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanRoleDepartmentsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanRoleDepartmentsCache_Get(t *testing.T) {
	c := newLoanRoleDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoleDepartments)
	err := c.ICache.(LoanRoleDepartmentsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRoleDepartmentsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanRoleDepartmentsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanRoleDepartmentsCache_MultiGet(t *testing.T) {
	c := newLoanRoleDepartmentsCache()
	defer c.Close()

	var testData []*model.LoanRoleDepartments
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRoleDepartments))
	}

	err := c.ICache.(LoanRoleDepartmentsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanRoleDepartmentsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanRoleDepartments))
	}
}

func Test_loanRoleDepartmentsCache_MultiSet(t *testing.T) {
	c := newLoanRoleDepartmentsCache()
	defer c.Close()

	var testData []*model.LoanRoleDepartments
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanRoleDepartments))
	}

	err := c.ICache.(LoanRoleDepartmentsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRoleDepartmentsCache_Del(t *testing.T) {
	c := newLoanRoleDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoleDepartments)
	err := c.ICache.(LoanRoleDepartmentsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanRoleDepartmentsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanRoleDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanRoleDepartments)
	err := c.ICache.(LoanRoleDepartmentsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanRoleDepartmentsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanRoleDepartmentsCache(t *testing.T) {
	c := NewLoanRoleDepartmentsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanRoleDepartmentsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanRoleDepartmentsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
