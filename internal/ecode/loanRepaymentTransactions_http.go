package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanRepaymentTransactions business-level http error codes.
// the loanRepaymentTransactionsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanRepaymentTransactionsNO = 93
	loanRepaymentTransactionsName     = "loanRepaymentTransactions"
	loanRepaymentTransactionsBaseCode = errcode.HCode(loanRepaymentTransactionsNO)

	ErrCreateLoanRepaymentTransactions     = errcode.NewError(loanRepaymentTransactionsBaseCode+1, "failed to create "+loanRepaymentTransactionsName)
	ErrDeleteByIDLoanRepaymentTransactions = errcode.NewError(loanRepaymentTransactionsBaseCode+2, "failed to delete "+loanRepaymentTransactionsName)
	ErrUpdateByIDLoanRepaymentTransactions = errcode.NewError(loanRepaymentTransactionsBaseCode+3, "failed to update "+loanRepaymentTransactionsName)
	ErrGetByIDLoanRepaymentTransactions    = errcode.NewError(loanRepaymentTransactionsBaseCode+4, "failed to get "+loanRepaymentTransactionsName+" details")
	ErrListLoanRepaymentTransactions       = errcode.NewError(loanRepaymentTransactionsBaseCode+5, "failed to list of "+loanRepaymentTransactionsName)

	ErrDeleteByIDsLoanRepaymentTransactions    = errcode.NewError(loanRepaymentTransactionsBaseCode+6, "failed to delete by batch ids "+loanRepaymentTransactionsName)
	ErrGetByConditionLoanRepaymentTransactions = errcode.NewError(loanRepaymentTransactionsBaseCode+7, "failed to get "+loanRepaymentTransactionsName+" details by conditions")
	ErrListByIDsLoanRepaymentTransactions      = errcode.NewError(loanRepaymentTransactionsBaseCode+8, "failed to list by batch ids "+loanRepaymentTransactionsName)
	ErrListByLastIDLoanRepaymentTransactions   = errcode.NewError(loanRepaymentTransactionsBaseCode+9, "failed to list by last id "+loanRepaymentTransactionsName)

	// error codes are globally unique, adding 1 to the previous error code
)
