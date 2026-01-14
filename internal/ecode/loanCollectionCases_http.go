package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanCollectionCases business-level http error codes.
// the loanCollectionCasesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanCollectionCasesNO = 81
	loanCollectionCasesName     = "loanCollectionCases"
	loanCollectionCasesBaseCode = errcode.HCode(loanCollectionCasesNO)

	ErrCreateLoanCollectionCases     = errcode.NewError(loanCollectionCasesBaseCode+1, "failed to create "+loanCollectionCasesName)
	ErrDeleteByIDLoanCollectionCases = errcode.NewError(loanCollectionCasesBaseCode+2, "failed to delete "+loanCollectionCasesName)
	ErrUpdateByIDLoanCollectionCases = errcode.NewError(loanCollectionCasesBaseCode+3, "failed to update "+loanCollectionCasesName)
	ErrGetByIDLoanCollectionCases    = errcode.NewError(loanCollectionCasesBaseCode+4, "failed to get "+loanCollectionCasesName+" details")
	ErrListLoanCollectionCases       = errcode.NewError(loanCollectionCasesBaseCode+5, "failed to list of "+loanCollectionCasesName)

	ErrDeleteByIDsLoanCollectionCases    = errcode.NewError(loanCollectionCasesBaseCode+6, "failed to delete by batch ids "+loanCollectionCasesName)
	ErrGetByConditionLoanCollectionCases = errcode.NewError(loanCollectionCasesBaseCode+7, "failed to get "+loanCollectionCasesName+" details by conditions")
	ErrListByIDsLoanCollectionCases      = errcode.NewError(loanCollectionCasesBaseCode+8, "failed to list by batch ids "+loanCollectionCasesName)
	ErrListByLastIDLoanCollectionCases   = errcode.NewError(loanCollectionCasesBaseCode+9, "failed to list by last id "+loanCollectionCasesName)

	// error codes are globally unique, adding 1 to the previous error code
)
