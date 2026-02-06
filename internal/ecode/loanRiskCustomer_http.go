package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanRiskCustomer business-level http error codes.
// the loanRiskCustomerNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanRiskCustomerNO = 28
	loanRiskCustomerName     = "loanRiskCustomer"
	loanRiskCustomerBaseCode = errcode.HCode(loanRiskCustomerNO)

	ErrCreateLoanRiskCustomer     = errcode.NewError(loanRiskCustomerBaseCode+1, "failed to create "+loanRiskCustomerName)
	ErrDeleteByIDLoanRiskCustomer = errcode.NewError(loanRiskCustomerBaseCode+2, "failed to delete "+loanRiskCustomerName)
	ErrUpdateByIDLoanRiskCustomer = errcode.NewError(loanRiskCustomerBaseCode+3, "failed to update "+loanRiskCustomerName)
	ErrGetByIDLoanRiskCustomer    = errcode.NewError(loanRiskCustomerBaseCode+4, "failed to get "+loanRiskCustomerName+" details")
	ErrListLoanRiskCustomer       = errcode.NewError(loanRiskCustomerBaseCode+5, "failed to list of "+loanRiskCustomerName)

	// error codes are globally unique, adding 1 to the previous error code
)
