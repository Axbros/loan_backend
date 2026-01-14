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

func newLoanMfaRecoveryCodesCache() *gotest.Cache {
	record1 := &model.LoanMfaRecoveryCodes{}
	record1.ID = 1
	record2 := &model.LoanMfaRecoveryCodes{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewLoanMfaRecoveryCodesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_loanMfaRecoveryCodesCache_Set(t *testing.T) {
	c := newLoanMfaRecoveryCodesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaRecoveryCodes)
	err := c.ICache.(LoanMfaRecoveryCodesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(LoanMfaRecoveryCodesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_loanMfaRecoveryCodesCache_Get(t *testing.T) {
	c := newLoanMfaRecoveryCodesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaRecoveryCodes)
	err := c.ICache.(LoanMfaRecoveryCodesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanMfaRecoveryCodesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(LoanMfaRecoveryCodesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_loanMfaRecoveryCodesCache_MultiGet(t *testing.T) {
	c := newLoanMfaRecoveryCodesCache()
	defer c.Close()

	var testData []*model.LoanMfaRecoveryCodes
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanMfaRecoveryCodes))
	}

	err := c.ICache.(LoanMfaRecoveryCodesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(LoanMfaRecoveryCodesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.LoanMfaRecoveryCodes))
	}
}

func Test_loanMfaRecoveryCodesCache_MultiSet(t *testing.T) {
	c := newLoanMfaRecoveryCodesCache()
	defer c.Close()

	var testData []*model.LoanMfaRecoveryCodes
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.LoanMfaRecoveryCodes))
	}

	err := c.ICache.(LoanMfaRecoveryCodesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanMfaRecoveryCodesCache_Del(t *testing.T) {
	c := newLoanMfaRecoveryCodesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaRecoveryCodes)
	err := c.ICache.(LoanMfaRecoveryCodesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_loanMfaRecoveryCodesCache_SetCacheWithNotFound(t *testing.T) {
	c := newLoanMfaRecoveryCodesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.LoanMfaRecoveryCodes)
	err := c.ICache.(LoanMfaRecoveryCodesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(LoanMfaRecoveryCodesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewLoanMfaRecoveryCodesCache(t *testing.T) {
	c := NewLoanMfaRecoveryCodesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewLoanMfaRecoveryCodesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewLoanMfaRecoveryCodesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
