package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanUserCallRecords business-level http error codes.
// the loanUserCallRecordsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanUserCallRecordsNO = 97
	loanUserCallRecordsName     = "loanUserCallRecords"
	loanUserCallRecordsBaseCode = errcode.HCode(loanUserCallRecordsNO)

	ErrCreateLoanUserCallRecords     = errcode.NewError(loanUserCallRecordsBaseCode+1, "failed to create "+loanUserCallRecordsName)
	ErrDeleteByIDLoanUserCallRecords = errcode.NewError(loanUserCallRecordsBaseCode+2, "failed to delete "+loanUserCallRecordsName)
	ErrUpdateByIDLoanUserCallRecords = errcode.NewError(loanUserCallRecordsBaseCode+3, "failed to update "+loanUserCallRecordsName)
	ErrGetByIDLoanUserCallRecords    = errcode.NewError(loanUserCallRecordsBaseCode+4, "failed to get "+loanUserCallRecordsName+" details")
	ErrListLoanUserCallRecords       = errcode.NewError(loanUserCallRecordsBaseCode+5, "failed to list of "+loanUserCallRecordsName)

	ErrDeleteByIDsLoanUserCallRecords    = errcode.NewError(loanUserCallRecordsBaseCode+6, "failed to delete by batch ids "+loanUserCallRecordsName)
	ErrGetByConditionLoanUserCallRecords = errcode.NewError(loanUserCallRecordsBaseCode+7, "failed to get "+loanUserCallRecordsName+" details by conditions")
	ErrListByIDsLoanUserCallRecords      = errcode.NewError(loanUserCallRecordsBaseCode+8, "failed to list by batch ids "+loanUserCallRecordsName)
	ErrListByLastIDLoanUserCallRecords   = errcode.NewError(loanUserCallRecordsBaseCode+9, "failed to list by last id "+loanUserCallRecordsName)

	// error codes are globally unique, adding 1 to the previous error code
)
