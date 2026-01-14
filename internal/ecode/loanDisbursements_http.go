package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanDisbursements business-level http error codes.
// the loanDisbursementsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanDisbursementsNO = 85
	loanDisbursementsName     = "loanDisbursements"
	loanDisbursementsBaseCode = errcode.HCode(loanDisbursementsNO)

	ErrCreateLoanDisbursements     = errcode.NewError(loanDisbursementsBaseCode+1, "failed to create "+loanDisbursementsName)
	ErrDeleteByIDLoanDisbursements = errcode.NewError(loanDisbursementsBaseCode+2, "failed to delete "+loanDisbursementsName)
	ErrUpdateByIDLoanDisbursements = errcode.NewError(loanDisbursementsBaseCode+3, "failed to update "+loanDisbursementsName)
	ErrGetByIDLoanDisbursements    = errcode.NewError(loanDisbursementsBaseCode+4, "failed to get "+loanDisbursementsName+" details")
	ErrListLoanDisbursements       = errcode.NewError(loanDisbursementsBaseCode+5, "failed to list of "+loanDisbursementsName)

	ErrDeleteByIDsLoanDisbursements    = errcode.NewError(loanDisbursementsBaseCode+6, "failed to delete by batch ids "+loanDisbursementsName)
	ErrGetByConditionLoanDisbursements = errcode.NewError(loanDisbursementsBaseCode+7, "failed to get "+loanDisbursementsName+" details by conditions")
	ErrListByIDsLoanDisbursements      = errcode.NewError(loanDisbursementsBaseCode+8, "failed to list by batch ids "+loanDisbursementsName)
	ErrListByLastIDLoanDisbursements   = errcode.NewError(loanDisbursementsBaseCode+9, "failed to list by last id "+loanDisbursementsName)

	// error codes are globally unique, adding 1 to the previous error code
)
