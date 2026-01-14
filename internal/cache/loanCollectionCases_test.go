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

func newLoanCollectionCasesCache() *gotest.Cache {
	record1 := &model.LoanCollectionCases{}
	record1.ID = 1
	record2 := &model.LoanCollectionCases{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanCollectionCasesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanCollectionCasesCache_Set(t *testing.T) {
	c := newLoanCollectionCasesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionCases)
	err := c.ICache.(LoanCollectionCasesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanCollectionCasesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanCollectionCasesCache_Get(t *testing.T) {
	c := newLoanCollectionCasesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionCases)
	err := c.ICache.(LoanCollectionCasesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanCollectionCasesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanCollectionCasesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanCollectionCasesCache_MultiGet(t *testing.T) {
	c := newLoanCollectionCasesCache()
	defer c.Close()

	var testData []*model.LoanCollectionCases
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanCollectionCases))
	}

	err := c.ICache.(LoanCollectionCasesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanCollectionCasesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanCollectionCases))
	}
}

func Test_loanCollectionCasesCache_MultiSet(t *testing.T) {
	c := newLoanCollectionCasesCache()
	defer c.Close()

	var testData []*model.LoanCollectionCases
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanCollectionCases))
	}

	err := c.ICache.(LoanCollectionCasesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanCollectionCasesCache_Del(t *testing.T) {
	c := newLoanCollectionCasesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionCases)
	err := c.ICache.(LoanCollectionCasesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanCollectionCasesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanCollectionCasesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionCases)
	err := c.ICache.(LoanCollectionCasesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanCollectionCasesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanCollectionCasesCache(t *testing.T) {
	c := NewLoanCollectionCasesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanCollectionCasesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanCollectionCasesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
