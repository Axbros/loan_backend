package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanAudits business-level http error codes.
// the loanAuditsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanAuditsNO       = 78
	loanAuditsName     = "loanAudits"
	loanAuditsBaseCode = errcode.HCode(loanAuditsNO)

	ErrCreateLoanAudits     = errcode.NewError(loanAuditsBaseCode+1, "failed to create "+loanAuditsName)
	ErrDeleteByIDLoanAudits = errcode.NewError(loanAuditsBaseCode+2, "failed to delete "+loanAuditsName)
	ErrUpdateByIDLoanAudits = errcode.NewError(loanAuditsBaseCode+3, "failed to update "+loanAuditsName)
	ErrGetByIDLoanAudits    = errcode.NewError(loanAuditsBaseCode+4, "failed to get "+loanAuditsName+" details")
	ErrListLoanAudits       = errcode.NewError(loanAuditsBaseCode+5, "failed to list of "+loanAuditsName)

	ErrDeleteByIDsLoanAudits    = errcode.NewError(loanAuditsBaseCode+6, "failed to delete by batch ids "+loanAuditsName)
	ErrGetByConditionLoanAudits = errcode.NewError(loanAuditsBaseCode+7, "failed to get "+loanAuditsName+" details by conditions")
	ErrListByIDsLoanAudits      = errcode.NewError(loanAuditsBaseCode+8, "failed to list by batch ids "+loanAuditsName)
	ErrListByLastIDLoanAudits   = errcode.NewError(loanAuditsBaseCode+9, "failed to list by last id "+loanAuditsName)

	// error codes are globally unique, adding 1 to the previous error code
)
