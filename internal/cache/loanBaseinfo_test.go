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

func newLoanBaseinfoCache() *gotest.Cache {
	record1 := &model.LoanBaseinfo{}
	record1.ID = 1
	record2 := &model.LoanBaseinfo{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanBaseinfoCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanBaseinfoCache_Set(t *testing.T) {
	c := newLoanBaseinfoCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfo)
	err := c.ICache.(LoanBaseinfoCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanBaseinfoCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanBaseinfoCache_Get(t *testing.T) {
	c := newLoanBaseinfoCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfo)
	err := c.ICache.(LoanBaseinfoCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanBaseinfoCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanBaseinfoCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanBaseinfoCache_MultiGet(t *testing.T) {
	c := newLoanBaseinfoCache()
	defer c.Close()

	var testData []*model.LoanBaseinfo
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanBaseinfo))
	}

	err := c.ICache.(LoanBaseinfoCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanBaseinfoCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanBaseinfo))
	}
}

func Test_loanBaseinfoCache_MultiSet(t *testing.T) {
	c := newLoanBaseinfoCache()
	defer c.Close()

	var testData []*model.LoanBaseinfo
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanBaseinfo))
	}

	err := c.ICache.(LoanBaseinfoCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanBaseinfoCache_Del(t *testing.T) {
	c := newLoanBaseinfoCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfo)
	err := c.ICache.(LoanBaseinfoCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanBaseinfoCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanBaseinfoCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfo)
	err := c.ICache.(LoanBaseinfoCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanBaseinfoCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanBaseinfoCache(t *testing.T) {
	c := NewLoanBaseinfoCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanBaseinfoCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanBaseinfoCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
