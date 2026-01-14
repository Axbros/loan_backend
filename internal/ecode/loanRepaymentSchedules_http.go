package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanRepaymentSchedules business-level http error codes.
// the loanRepaymentSchedulesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanRepaymentSchedulesNO = 92
	loanRepaymentSchedulesName     = "loanRepaymentSchedules"
	loanRepaymentSchedulesBaseCode = errcode.HCode(loanRepaymentSchedulesNO)

	ErrCreateLoanRepaymentSchedules     = errcode.NewError(loanRepaymentSchedulesBaseCode+1, "failed to create "+loanRepaymentSchedulesName)
	ErrDeleteByIDLoanRepaymentSchedules = errcode.NewError(loanRepaymentSchedulesBaseCode+2, "failed to delete "+loanRepaymentSchedulesName)
	ErrUpdateByIDLoanRepaymentSchedules = errcode.NewError(loanRepaymentSchedulesBaseCode+3, "failed to update "+loanRepaymentSchedulesName)
	ErrGetByIDLoanRepaymentSchedules    = errcode.NewError(loanRepaymentSchedulesBaseCode+4, "failed to get "+loanRepaymentSchedulesName+" details")
	ErrListLoanRepaymentSchedules       = errcode.NewError(loanRepaymentSchedulesBaseCode+5, "failed to list of "+loanRepaymentSchedulesName)

	ErrDeleteByIDsLoanRepaymentSchedules    = errcode.NewError(loanRepaymentSchedulesBaseCode+6, "failed to delete by batch ids "+loanRepaymentSchedulesName)
	ErrGetByConditionLoanRepaymentSchedules = errcode.NewError(loanRepaymentSchedulesBaseCode+7, "failed to get "+loanRepaymentSchedulesName+" details by conditions")
	ErrListByIDsLoanRepaymentSchedules      = errcode.NewError(loanRepaymentSchedulesBaseCode+8, "failed to list by batch ids "+loanRepaymentSchedulesName)
	ErrListByLastIDLoanRepaymentSchedules   = errcode.NewError(loanRepaymentSchedulesBaseCode+9, "failed to list by last id "+loanRepaymentSchedulesName)

	// error codes are globally unique, adding 1 to the previous error code
)
