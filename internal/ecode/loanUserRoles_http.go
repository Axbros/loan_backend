package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanUserRoles business-level http error codes.
// the loanUserRolesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanUserRolesNO = 100
	loanUserRolesName     = "loanUserRoles"
	loanUserRolesBaseCode = errcode.HCode(loanUserRolesNO)

	ErrCreateLoanUserRoles     = errcode.NewError(loanUserRolesBaseCode+1, "failed to create "+loanUserRolesName)
	ErrDeleteByIDLoanUserRoles = errcode.NewError(loanUserRolesBaseCode+2, "failed to delete "+loanUserRolesName)
	ErrUpdateByIDLoanUserRoles = errcode.NewError(loanUserRolesBaseCode+3, "failed to update "+loanUserRolesName)
	ErrGetByIDLoanUserRoles    = errcode.NewError(loanUserRolesBaseCode+4, "failed to get "+loanUserRolesName+" details")
	ErrListLoanUserRoles       = errcode.NewError(loanUserRolesBaseCode+5, "failed to list of "+loanUserRolesName)

	ErrDeleteByIDsLoanUserRoles    = errcode.NewError(loanUserRolesBaseCode+6, "failed to delete by batch ids "+loanUserRolesName)
	ErrGetByConditionLoanUserRoles = errcode.NewError(loanUserRolesBaseCode+7, "failed to get "+loanUserRolesName+" details by conditions")
	ErrListByIDsLoanUserRoles      = errcode.NewError(loanUserRolesBaseCode+8, "failed to list by batch ids "+loanUserRolesName)
	ErrListByLastIDLoanUserRoles   = errcode.NewError(loanUserRolesBaseCode+9, "failed to list by last id "+loanUserRolesName)

	// error codes are globally unique, adding 1 to the previous error code
)
