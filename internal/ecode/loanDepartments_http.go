package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanDepartments business-level http error codes.
// the loanDepartmentsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanDepartmentsNO = 84
	loanDepartmentsName     = "loanDepartments"
	loanDepartmentsBaseCode = errcode.HCode(loanDepartmentsNO)

	ErrCreateLoanDepartments     = errcode.NewError(loanDepartmentsBaseCode+1, "failed to create "+loanDepartmentsName)
	ErrDeleteByIDLoanDepartments = errcode.NewError(loanDepartmentsBaseCode+2, "failed to delete "+loanDepartmentsName)
	ErrUpdateByIDLoanDepartments = errcode.NewError(loanDepartmentsBaseCode+3, "failed to update "+loanDepartmentsName)
	ErrGetByIDLoanDepartments    = errcode.NewError(loanDepartmentsBaseCode+4, "failed to get "+loanDepartmentsName+" details")
	ErrListLoanDepartments       = errcode.NewError(loanDepartmentsBaseCode+5, "failed to list of "+loanDepartmentsName)

	ErrDeleteByIDsLoanDepartments    = errcode.NewError(loanDepartmentsBaseCode+6, "failed to delete by batch ids "+loanDepartmentsName)
	ErrGetByConditionLoanDepartments = errcode.NewError(loanDepartmentsBaseCode+7, "failed to get "+loanDepartmentsName+" details by conditions")
	ErrListByIDsLoanDepartments      = errcode.NewError(loanDepartmentsBaseCode+8, "failed to list by batch ids "+loanDepartmentsName)
	ErrListByLastIDLoanDepartments   = errcode.NewError(loanDepartmentsBaseCode+9, "failed to list by last id "+loanDepartmentsName)

	// error codes are globally unique, adding 1 to the previous error code
)
