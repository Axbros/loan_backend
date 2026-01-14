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

func newLoanUserCallRecordsCache() *gotest.Cache {
	record1 := &model.LoanUserCallRecords{}
	record1.ID = 1
	record2 := &model.LoanUserCallRecords{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanUserCallRecordsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanUserCallRecordsCache_Set(t *testing.T) {
	c := newLoanUserCallRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserCallRecords)
	err := c.ICache.(LoanUserCallRecordsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanUserCallRecordsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanUserCallRecordsCache_Get(t *testing.T) {
	c := newLoanUserCallRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserCallRecords)
	err := c.ICache.(LoanUserCallRecordsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserCallRecordsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanUserCallRecordsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanUserCallRecordsCache_MultiGet(t *testing.T) {
	c := newLoanUserCallRecordsCache()
	defer c.Close()

	var testData []*model.LoanUserCallRecords
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserCallRecords))
	}

	err := c.ICache.(LoanUserCallRecordsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserCallRecordsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanUserCallRecords))
	}
}

func Test_loanUserCallRecordsCache_MultiSet(t *testing.T) {
	c := newLoanUserCallRecordsCache()
	defer c.Close()

	var testData []*model.LoanUserCallRecords
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserCallRecords))
	}

	err := c.ICache.(LoanUserCallRecordsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserCallRecordsCache_Del(t *testing.T) {
	c := newLoanUserCallRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserCallRecords)
	err := c.ICache.(LoanUserCallRecordsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserCallRecordsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanUserCallRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserCallRecords)
	err := c.ICache.(LoanUserCallRecordsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanUserCallRecordsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanUserCallRecordsCache(t *testing.T) {
	c := NewLoanUserCallRecordsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanUserCallRecordsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanUserCallRecordsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
