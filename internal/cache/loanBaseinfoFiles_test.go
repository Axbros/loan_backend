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

func newLoanBaseinfoFilesCache() *gotest.Cache {
	record1 := &model.LoanBaseinfoFiles{}
	record1.ID = 1
	record2 := &model.LoanBaseinfoFiles{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanBaseinfoFilesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanBaseinfoFilesCache_Set(t *testing.T) {
	c := newLoanBaseinfoFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfoFiles)
	err := c.ICache.(LoanBaseinfoFilesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanBaseinfoFilesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanBaseinfoFilesCache_Get(t *testing.T) {
	c := newLoanBaseinfoFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfoFiles)
	err := c.ICache.(LoanBaseinfoFilesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanBaseinfoFilesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanBaseinfoFilesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanBaseinfoFilesCache_MultiGet(t *testing.T) {
	c := newLoanBaseinfoFilesCache()
	defer c.Close()

	var testData []*model.LoanBaseinfoFiles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanBaseinfoFiles))
	}

	err := c.ICache.(LoanBaseinfoFilesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanBaseinfoFilesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanBaseinfoFiles))
	}
}

func Test_loanBaseinfoFilesCache_MultiSet(t *testing.T) {
	c := newLoanBaseinfoFilesCache()
	defer c.Close()

	var testData []*model.LoanBaseinfoFiles
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanBaseinfoFiles))
	}

	err := c.ICache.(LoanBaseinfoFilesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanBaseinfoFilesCache_Del(t *testing.T) {
	c := newLoanBaseinfoFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfoFiles)
	err := c.ICache.(LoanBaseinfoFilesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanBaseinfoFilesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanBaseinfoFilesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanBaseinfoFiles)
	err := c.ICache.(LoanBaseinfoFilesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanBaseinfoFilesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanBaseinfoFilesCache(t *testing.T) {
	c := NewLoanBaseinfoFilesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanBaseinfoFilesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanBaseinfoFilesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
