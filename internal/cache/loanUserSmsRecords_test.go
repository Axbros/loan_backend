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

func newLoanUserSmsRecordsCache() *gotest.Cache {
	record1 := &model.LoanUserSmsRecords{}
	record1.ID = 1
	record2 := &model.LoanUserSmsRecords{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanUserSmsRecordsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanUserSmsRecordsCache_Set(t *testing.T) {
	c := newLoanUserSmsRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserSmsRecords)
	err := c.ICache.(LoanUserSmsRecordsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanUserSmsRecordsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanUserSmsRecordsCache_Get(t *testing.T) {
	c := newLoanUserSmsRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserSmsRecords)
	err := c.ICache.(LoanUserSmsRecordsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserSmsRecordsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanUserSmsRecordsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanUserSmsRecordsCache_MultiGet(t *testing.T) {
	c := newLoanUserSmsRecordsCache()
	defer c.Close()

	var testData []*model.LoanUserSmsRecords
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserSmsRecords))
	}

	err := c.ICache.(LoanUserSmsRecordsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserSmsRecordsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanUserSmsRecords))
	}
}

func Test_loanUserSmsRecordsCache_MultiSet(t *testing.T) {
	c := newLoanUserSmsRecordsCache()
	defer c.Close()

	var testData []*model.LoanUserSmsRecords
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserSmsRecords))
	}

	err := c.ICache.(LoanUserSmsRecordsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserSmsRecordsCache_Del(t *testing.T) {
	c := newLoanUserSmsRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserSmsRecords)
	err := c.ICache.(LoanUserSmsRecordsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserSmsRecordsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanUserSmsRecordsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserSmsRecords)
	err := c.ICache.(LoanUserSmsRecordsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanUserSmsRecordsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanUserSmsRecordsCache(t *testing.T) {
	c := NewLoanUserSmsRecordsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanUserSmsRecordsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanUserSmsRecordsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
