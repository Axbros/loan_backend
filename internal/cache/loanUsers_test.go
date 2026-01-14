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

func newLoanUsersCache() *gotest.Cache {
	record1 := &model.LoanUsers{}
	record1.ID = 1
	record2 := &model.LoanUsers{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanUsersCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanUsersCache_Set(t *testing.T) {
	c := newLoanUsersCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUsers)
	err := c.ICache.(LoanUsersCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanUsersCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanUsersCache_Get(t *testing.T) {
	c := newLoanUsersCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUsers)
	err := c.ICache.(LoanUsersCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUsersCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanUsersCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanUsersCache_MultiGet(t *testing.T) {
	c := newLoanUsersCache()
	defer c.Close()

	var testData []*model.LoanUsers
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUsers))
	}

	err := c.ICache.(LoanUsersCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanUsersCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanUsers))
	}
}

func Test_loanUsersCache_MultiSet(t *testing.T) {
	c := newLoanUsersCache()
	defer c.Close()

	var testData []*model.LoanUsers
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanUsers))
	}

	err := c.ICache.(LoanUsersCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUsersCache_Del(t *testing.T) {
	c := newLoanUsersCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUsers)
	err := c.ICache.(LoanUsersCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanUsersCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanUsersCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanUsers)
	err := c.ICache.(LoanUsersCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanUsersCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanUsersCache(t *testing.T) {
	c := NewLoanUsersCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanUsersCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanUsersCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
