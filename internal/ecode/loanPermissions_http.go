package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanPermissions business-level http error codes.
// the loanPermissionsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanPermissionsNO = 90
	loanPermissionsName     = "loanPermissions"
	loanPermissionsBaseCode = errcode.HCode(loanPermissionsNO)

	ErrCreateLoanPermissions     = errcode.NewError(loanPermissionsBaseCode+1, "failed to create "+loanPermissionsName)
	ErrDeleteByIDLoanPermissions = errcode.NewError(loanPermissionsBaseCode+2, "failed to delete "+loanPermissionsName)
	ErrUpdateByIDLoanPermissions = errcode.NewError(loanPermissionsBaseCode+3, "failed to update "+loanPermissionsName)
	ErrGetByIDLoanPermissions    = errcode.NewError(loanPermissionsBaseCode+4, "failed to get "+loanPermissionsName+" details")
	ErrListLoanPermissions       = errcode.NewError(loanPermissionsBaseCode+5, "failed to list of "+loanPermissionsName)

	ErrDeleteByIDsLoanPermissions    = errcode.NewError(loanPermissionsBaseCode+6, "failed to delete by batch ids "+loanPermissionsName)
	ErrGetByConditionLoanPermissions = errcode.NewError(loanPermissionsBaseCode+7, "failed to get "+loanPermissionsName+" details by conditions")
	ErrListByIDsLoanPermissions      = errcode.NewError(loanPermissionsBaseCode+8, "failed to list by batch ids "+loanPermissionsName)
	ErrListByLastIDLoanPermissions   = errcode.NewError(loanPermissionsBaseCode+9, "failed to list by last id "+loanPermissionsName)

	// error codes are globally unique, adding 1 to the previous error code
)
