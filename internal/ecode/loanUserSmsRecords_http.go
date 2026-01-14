package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanUserSmsRecords business-level http error codes.
// the loanUserSmsRecordsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanUserSmsRecordsNO = 101
	loanUserSmsRecordsName     = "loanUserSmsRecords"
	loanUserSmsRecordsBaseCode = errcode.HCode(loanUserSmsRecordsNO)

	ErrCreateLoanUserSmsRecords     = errcode.NewError(loanUserSmsRecordsBaseCode+1, "failed to create "+loanUserSmsRecordsName)
	ErrDeleteByIDLoanUserSmsRecords = errcode.NewError(loanUserSmsRecordsBaseCode+2, "failed to delete "+loanUserSmsRecordsName)
	ErrUpdateByIDLoanUserSmsRecords = errcode.NewError(loanUserSmsRecordsBaseCode+3, "failed to update "+loanUserSmsRecordsName)
	ErrGetByIDLoanUserSmsRecords    = errcode.NewError(loanUserSmsRecordsBaseCode+4, "failed to get "+loanUserSmsRecordsName+" details")
	ErrListLoanUserSmsRecords       = errcode.NewError(loanUserSmsRecordsBaseCode+5, "failed to list of "+loanUserSmsRecordsName)

	ErrDeleteByIDsLoanUserSmsRecords    = errcode.NewError(loanUserSmsRecordsBaseCode+6, "failed to delete by batch ids "+loanUserSmsRecordsName)
	ErrGetByConditionLoanUserSmsRecords = errcode.NewError(loanUserSmsRecordsBaseCode+7, "failed to get "+loanUserSmsRecordsName+" details by conditions")
	ErrListByIDsLoanUserSmsRecords      = errcode.NewError(loanUserSmsRecordsBaseCode+8, "failed to list by batch ids "+loanUserSmsRecordsName)
	ErrListByLastIDLoanUserSmsRecords   = errcode.NewError(loanUserSmsRecordsBaseCode+9, "failed to list by last id "+loanUserSmsRecordsName)

	// error codes are globally unique, adding 1 to the previous error code
)
