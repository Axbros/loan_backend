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

func newLoanDisbursementsCache() *gotest.Cache {
	record1 := &model.LoanDisbursements{}
	record1.ID = 1
	record2 := &model.LoanDisbursements{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanDisbursementsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanDisbursementsCache_Set(t *testing.T) {
	c := newLoanDisbursementsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDisbursements)
	err := c.ICache.(LoanDisbursementsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanDisbursementsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanDisbursementsCache_Get(t *testing.T) {
	c := newLoanDisbursementsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDisbursements)
	err := c.ICache.(LoanDisbursementsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanDisbursementsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanDisbursementsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanDisbursementsCache_MultiGet(t *testing.T) {
	c := newLoanDisbursementsCache()
	defer c.Close()

	var testData []*model.LoanDisbursements
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanDisbursements))
	}

	err := c.ICache.(LoanDisbursementsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanDisbursementsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanDisbursements))
	}
}

func Test_loanDisbursementsCache_MultiSet(t *testing.T) {
	c := newLoanDisbursementsCache()
	defer c.Close()

	var testData []*model.LoanDisbursements
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanDisbursements))
	}

	err := c.ICache.(LoanDisbursementsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanDisbursementsCache_Del(t *testing.T) {
	c := newLoanDisbursementsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDisbursements)
	err := c.ICache.(LoanDisbursementsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanDisbursementsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanDisbursementsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanDisbursements)
	err := c.ICache.(LoanDisbursementsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanDisbursementsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanDisbursementsCache(t *testing.T) {
	c := NewLoanDisbursementsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanDisbursementsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanDisbursementsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
