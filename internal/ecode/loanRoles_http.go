package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanRoles business-level http error codes.
// the loanRolesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanRolesNO = 96
	loanRolesName     = "loanRoles"
	loanRolesBaseCode = errcode.HCode(loanRolesNO)

	ErrCreateLoanRoles     = errcode.NewError(loanRolesBaseCode+1, "failed to create "+loanRolesName)
	ErrDeleteByIDLoanRoles = errcode.NewError(loanRolesBaseCode+2, "failed to delete "+loanRolesName)
	ErrUpdateByIDLoanRoles = errcode.NewError(loanRolesBaseCode+3, "failed to update "+loanRolesName)
	ErrGetByIDLoanRoles    = errcode.NewError(loanRolesBaseCode+4, "failed to get "+loanRolesName+" details")
	ErrListLoanRoles       = errcode.NewError(loanRolesBaseCode+5, "failed to list of "+loanRolesName)

	ErrDeleteByIDsLoanRoles    = errcode.NewError(loanRolesBaseCode+6, "failed to delete by batch ids "+loanRolesName)
	ErrGetByConditionLoanRoles = errcode.NewError(loanRolesBaseCode+7, "failed to get "+loanRolesName+" details by conditions")
	ErrListByIDsLoanRoles      = errcode.NewError(loanRolesBaseCode+8, "failed to list by batch ids "+loanRolesName)
	ErrListByLastIDLoanRoles   = errcode.NewError(loanRolesBaseCode+9, "failed to list by last id "+loanRolesName)

	// error codes are globally unique, adding 1 to the previous error code
)
