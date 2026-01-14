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

func newLoanReferralVisitsCache() *gotest.Cache {
	record1 := &model.LoanReferralVisits{}
	record1.ID = 1
	record2 := &model.LoanReferralVisits{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanReferralVisitsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanReferralVisitsCache_Set(t *testing.T) {
	c := newLoanReferralVisitsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanReferralVisits)
	err := c.ICache.(LoanReferralVisitsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanReferralVisitsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanReferralVisitsCache_Get(t *testing.T) {
	c := newLoanReferralVisitsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanReferralVisits)
	err := c.ICache.(LoanReferralVisitsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanReferralVisitsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanReferralVisitsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanReferralVisitsCache_MultiGet(t *testing.T) {
	c := newLoanReferralVisitsCache()
	defer c.Close()

	var testData []*model.LoanReferralVisits
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanReferralVisits))
	}

	err := c.ICache.(LoanReferralVisitsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanReferralVisitsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanReferralVisits))
	}
}

func Test_loanReferralVisitsCache_MultiSet(t *testing.T) {
	c := newLoanReferralVisitsCache()
	defer c.Close()

	var testData []*model.LoanReferralVisits
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanReferralVisits))
	}

	err := c.ICache.(LoanReferralVisitsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanReferralVisitsCache_Del(t *testing.T) {
	c := newLoanReferralVisitsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanReferralVisits)
	err := c.ICache.(LoanReferralVisitsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanReferralVisitsCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanReferralVisitsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanReferralVisits)
	err := c.ICache.(LoanReferralVisitsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanReferralVisitsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanReferralVisitsCache(t *testing.T) {
	c := NewLoanReferralVisitsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanReferralVisitsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanReferralVisitsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
