package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanBaseinfo business-level http error codes.
// the loanBaseinfoNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanBaseinfoNO = 80
	loanBaseinfoName     = "loanBaseinfo"
	loanBaseinfoBaseCode = errcode.HCode(loanBaseinfoNO)

	ErrCreateLoanBaseinfo     = errcode.NewError(loanBaseinfoBaseCode+1, "failed to create "+loanBaseinfoName)
	ErrDeleteByIDLoanBaseinfo = errcode.NewError(loanBaseinfoBaseCode+2, "failed to delete "+loanBaseinfoName)
	ErrUpdateByIDLoanBaseinfo = errcode.NewError(loanBaseinfoBaseCode+3, "failed to update "+loanBaseinfoName)
	ErrGetByIDLoanBaseinfo    = errcode.NewError(loanBaseinfoBaseCode+4, "failed to get "+loanBaseinfoName+" details")
	ErrListLoanBaseinfo       = errcode.NewError(loanBaseinfoBaseCode+5, "failed to list of "+loanBaseinfoName)

	ErrDeleteByIDsLoanBaseinfo    = errcode.NewError(loanBaseinfoBaseCode+6, "failed to delete by batch ids "+loanBaseinfoName)
	ErrGetByConditionLoanBaseinfo = errcode.NewError(loanBaseinfoBaseCode+7, "failed to get "+loanBaseinfoName+" details by conditions")
	ErrListByIDsLoanBaseinfo      = errcode.NewError(loanBaseinfoBaseCode+8, "failed to list by batch ids "+loanBaseinfoName)
	ErrListByLastIDLoanBaseinfo   = errcode.NewError(loanBaseinfoBaseCode+9, "failed to list by last id "+loanBaseinfoName)

	// error codes are globally unique, adding 1 to the previous error code
)
