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

func newLoanCollectionLogsCache() *gotest.Cache {
	record1 := &model.LoanCollectionLogs{}
	record1.ID = 1
	record2 := &model.LoanCollectionLogs{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanCollectionLogsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanCollectionLogsCache_Set(t *testing.T) {
	c := newLoanCollectionLogsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionLogs)
	err := c.ICache.(LoanCollectionLogsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanCollectionLogsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanCollectionLogsCache_Get(t *testing.T) {
	c := newLoanCollectionLogsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionLogs)
	err := c.ICache.(LoanCollectionLogsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanCollectionLogsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanCollectionLogsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanCollectionLogsCache_MultiGet(t *testing.T) {
	c := newLoanCollectionLogsCache()
	defer c.Close()

	var testData []*model.LoanCollectionLogs
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanCollectionLogs))
	}

	err := c.ICache.(LoanCollectionLogsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanCollectionLogsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanCollectionLogs))
	}
}

func Test_loanCollectionLogsCache_MultiSet(t *testing.T) {
	c := newLoanCollectionLogsCache()
	defer c.Close()

	var testData []*model.LoanCollectionLogs
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanCollectionLogs))
	}

	err := c.ICache.(LoanCollectionLogsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanCollectionLogsCache_Del(t *testing.T) {
	c := newLoanCollectionLogsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionLogs)
	err := c.ICache.(LoanCollectionLogsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanCollectionLogsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanCollectionLogsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanCollectionLogs)
	err := c.ICache.(LoanCollectionLogsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanCollectionLogsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanCollectionLogsCache(t *testing.T) {
	c := NewLoanCollectionLogsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanCollectionLogsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanCollectionLogsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
