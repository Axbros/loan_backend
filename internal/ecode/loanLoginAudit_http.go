package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanLoginAudit business-level http error codes.
// the loanLoginAuditNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanLoginAuditNO = 86
	loanLoginAuditName     = "loanLoginAudit"
	loanLoginAuditBaseCode = errcode.HCode(loanLoginAuditNO)

	ErrCreateLoanLoginAudit     = errcode.NewError(loanLoginAuditBaseCode+1, "failed to create "+loanLoginAuditName)
	ErrDeleteByIDLoanLoginAudit = errcode.NewError(loanLoginAuditBaseCode+2, "failed to delete "+loanLoginAuditName)
	ErrUpdateByIDLoanLoginAudit = errcode.NewError(loanLoginAuditBaseCode+3, "failed to update "+loanLoginAuditName)
	ErrGetByIDLoanLoginAudit    = errcode.NewError(loanLoginAuditBaseCode+4, "failed to get "+loanLoginAuditName+" details")
	ErrListLoanLoginAudit       = errcode.NewError(loanLoginAuditBaseCode+5, "failed to list of "+loanLoginAuditName)

	ErrDeleteByIDsLoanLoginAudit    = errcode.NewError(loanLoginAuditBaseCode+6, "failed to delete by batch ids "+loanLoginAuditName)
	ErrGetByConditionLoanLoginAudit = errcode.NewError(loanLoginAuditBaseCode+7, "failed to get "+loanLoginAuditName+" details by conditions")
	ErrListByIDsLoanLoginAudit      = errcode.NewError(loanLoginAuditBaseCode+8, "failed to list by batch ids "+loanLoginAuditName)
	ErrListByLastIDLoanLoginAudit   = errcode.NewError(loanLoginAuditBaseCode+9, "failed to list by last id "+loanLoginAuditName)

	// error codes are globally unique, adding 1 to the previous error code
)
