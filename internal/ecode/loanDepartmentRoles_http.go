package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanDepartmentRoles business-level http error codes.
// the loanDepartmentRolesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanDepartmentRolesNO = 83
	loanDepartmentRolesName     = "loanDepartmentRoles"
	loanDepartmentRolesBaseCode = errcode.HCode(loanDepartmentRolesNO)

	ErrCreateLoanDepartmentRoles     = errcode.NewError(loanDepartmentRolesBaseCode+1, "failed to create "+loanDepartmentRolesName)
	ErrDeleteByIDLoanDepartmentRoles = errcode.NewError(loanDepartmentRolesBaseCode+2, "failed to delete "+loanDepartmentRolesName)
	ErrUpdateByIDLoanDepartmentRoles = errcode.NewError(loanDepartmentRolesBaseCode+3, "failed to update "+loanDepartmentRolesName)
	ErrGetByIDLoanDepartmentRoles    = errcode.NewError(loanDepartmentRolesBaseCode+4, "failed to get "+loanDepartmentRolesName+" details")
	ErrListLoanDepartmentRoles       = errcode.NewError(loanDepartmentRolesBaseCode+5, "failed to list of "+loanDepartmentRolesName)

	ErrDeleteByIDsLoanDepartmentRoles    = errcode.NewError(loanDepartmentRolesBaseCode+6, "failed to delete by batch ids "+loanDepartmentRolesName)
	ErrGetByConditionLoanDepartmentRoles = errcode.NewError(loanDepartmentRolesBaseCode+7, "failed to get "+loanDepartmentRolesName+" details by conditions")
	ErrListByIDsLoanDepartmentRoles      = errcode.NewError(loanDepartmentRolesBaseCode+8, "failed to list by batch ids "+loanDepartmentRolesName)
	ErrListByLastIDLoanDepartmentRoles   = errcode.NewError(loanDepartmentRolesBaseCode+9, "failed to list by last id "+loanDepartmentRolesName)

	// error codes are globally unique, adding 1 to the previous error code
)
