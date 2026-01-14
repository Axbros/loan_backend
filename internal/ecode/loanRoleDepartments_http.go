package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanRoleDepartments business-level http error codes.
// the loanRoleDepartmentsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanRoleDepartmentsNO = 94
	loanRoleDepartmentsName     = "loanRoleDepartments"
	loanRoleDepartmentsBaseCode = errcode.HCode(loanRoleDepartmentsNO)

	ErrCreateLoanRoleDepartments     = errcode.NewError(loanRoleDepartmentsBaseCode+1, "failed to create "+loanRoleDepartmentsName)
	ErrDeleteByIDLoanRoleDepartments = errcode.NewError(loanRoleDepartmentsBaseCode+2, "failed to delete "+loanRoleDepartmentsName)
	ErrUpdateByIDLoanRoleDepartments = errcode.NewError(loanRoleDepartmentsBaseCode+3, "failed to update "+loanRoleDepartmentsName)
	ErrGetByIDLoanRoleDepartments    = errcode.NewError(loanRoleDepartmentsBaseCode+4, "failed to get "+loanRoleDepartmentsName+" details")
	ErrListLoanRoleDepartments       = errcode.NewError(loanRoleDepartmentsBaseCode+5, "failed to list of "+loanRoleDepartmentsName)

	ErrDeleteByIDsLoanRoleDepartments    = errcode.NewError(loanRoleDepartmentsBaseCode+6, "failed to delete by batch ids "+loanRoleDepartmentsName)
	ErrGetByConditionLoanRoleDepartments = errcode.NewError(loanRoleDepartmentsBaseCode+7, "failed to get "+loanRoleDepartmentsName+" details by conditions")
	ErrListByIDsLoanRoleDepartments      = errcode.NewError(loanRoleDepartmentsBaseCode+8, "failed to list by batch ids "+loanRoleDepartmentsName)
	ErrListByLastIDLoanRoleDepartments   = errcode.NewError(loanRoleDepartmentsBaseCode+9, "failed to list by last id "+loanRoleDepartmentsName)

	// error codes are globally unique, adding 1 to the previous error code
)
