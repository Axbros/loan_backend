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

func newLoanDepartmentsCache() *gotest.Cache {
	record1 := &model.LoanDepartments{}
	record1.ID = 1
	record2 := &model.LoanDepartments{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanDepartmentsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanDepartmentsCache_Set(t *testing.T) {
	c := newLoanDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartments)
	err := c.ICache.(LoanDepartmentsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanDepartmentsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanDepartmentsCache_Get(t *testing.T) {
	c := newLoanDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartments)
	err := c.ICache.(LoanDepartmentsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanDepartmentsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanDepartmentsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanDepartmentsCache_MultiGet(t *testing.T) {
	c := newLoanDepartmentsCache()
	defer c.Close()

	var testData []*model.LoanDepartments
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanDepartments))
	}

	err := c.ICache.(LoanDepartmentsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanDepartmentsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanDepartments))
	}
}

func Test_loanDepartmentsCache_MultiSet(t *testing.T) {
	c := newLoanDepartmentsCache()
	defer c.Close()

	var testData []*model.LoanDepartments
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanDepartments))
	}

	err := c.ICache.(LoanDepartmentsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanDepartmentsCache_Del(t *testing.T) {
	c := newLoanDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartments)
	err := c.ICache.(LoanDepartmentsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanDepartmentsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanDepartmentsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDepartments)
	err := c.ICache.(LoanDepartmentsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanDepartmentsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanDepartmentsCache(t *testing.T) {
	c := NewLoanDepartmentsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanDepartmentsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanDepartmentsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
