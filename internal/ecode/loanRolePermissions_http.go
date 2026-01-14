package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanRolePermissions business-level http error codes.
// the loanRolePermissionsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanRolePermissionsNO = 95
	loanRolePermissionsName     = "loanRolePermissions"
	loanRolePermissionsBaseCode = errcode.HCode(loanRolePermissionsNO)

	ErrCreateLoanRolePermissions     = errcode.NewError(loanRolePermissionsBaseCode+1, "failed to create "+loanRolePermissionsName)
	ErrDeleteByIDLoanRolePermissions = errcode.NewError(loanRolePermissionsBaseCode+2, "failed to delete "+loanRolePermissionsName)
	ErrUpdateByIDLoanRolePermissions = errcode.NewError(loanRolePermissionsBaseCode+3, "failed to update "+loanRolePermissionsName)
	ErrGetByIDLoanRolePermissions    = errcode.NewError(loanRolePermissionsBaseCode+4, "failed to get "+loanRolePermissionsName+" details")
	ErrListLoanRolePermissions       = errcode.NewError(loanRolePermissionsBaseCode+5, "failed to list of "+loanRolePermissionsName)

	ErrDeleteByIDsLoanRolePermissions    = errcode.NewError(loanRolePermissionsBaseCode+6, "failed to delete by batch ids "+loanRolePermissionsName)
	ErrGetByConditionLoanRolePermissions = errcode.NewError(loanRolePermissionsBaseCode+7, "failed to get "+loanRolePermissionsName+" details by conditions")
	ErrListByIDsLoanRolePermissions      = errcode.NewError(loanRolePermissionsBaseCode+8, "failed to list by batch ids "+loanRolePermissionsName)
	ErrListByLastIDLoanRolePermissions   = errcode.NewError(loanRolePermissionsBaseCode+9, "failed to list by last id "+loanRolePermissionsName)

	// error codes are globally unique, adding 1 to the previous error code
)
