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

func newLoanUserContactsCache() *gotest.Cache {
	record1 := &model.LoanUserContacts{}
	record1.ID = 1
	record2 := &model.LoanUserContacts{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanUserContactsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanUserContactsCache_Set(t *testing.T) {
	c := newLoanUserContactsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserContacts)
	err := c.ICache.(LoanUserContactsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanUserContactsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanUserContactsCache_Get(t *testing.T) {
	c := newLoanUserContactsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserContacts)
	err := c.ICache.(LoanUserContactsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserContactsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanUserContactsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanUserContactsCache_MultiGet(t *testing.T) {
	c := newLoanUserContactsCache()
	defer c.Close()

	var testData []*model.LoanUserContacts
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserContacts))
	}

	err := c.ICache.(LoanUserContactsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUserContactsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanUserContacts))
	}
}

func Test_loanUserContactsCache_MultiSet(t *testing.T) {
	c := newLoanUserContactsCache()
	defer c.Close()

	var testData []*model.LoanUserContacts
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUserContacts))
	}

	err := c.ICache.(LoanUserContactsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserContactsCache_Del(t *testing.T) {
	c := newLoanUserContactsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserContacts)
	err := c.ICache.(LoanUserContactsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUserContactsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanUserContactsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUserContacts)
	err := c.ICache.(LoanUserContactsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanUserContactsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanUserContactsCache(t *testing.T) {
	c := NewLoanUserContactsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanUserContactsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanUserContactsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
