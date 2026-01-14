package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanBaseinfoFiles business-level http error codes.
// the loanBaseinfoFilesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanBaseinfoFilesNO = 79
	loanBaseinfoFilesName     = "loanBaseinfoFiles"
	loanBaseinfoFilesBaseCode = errcode.HCode(loanBaseinfoFilesNO)

	ErrCreateLoanBaseinfoFiles     = errcode.NewError(loanBaseinfoFilesBaseCode+1, "failed to create "+loanBaseinfoFilesName)
	ErrDeleteByIDLoanBaseinfoFiles = errcode.NewError(loanBaseinfoFilesBaseCode+2, "failed to delete "+loanBaseinfoFilesName)
	ErrUpdateByIDLoanBaseinfoFiles = errcode.NewError(loanBaseinfoFilesBaseCode+3, "failed to update "+loanBaseinfoFilesName)
	ErrGetByIDLoanBaseinfoFiles    = errcode.NewError(loanBaseinfoFilesBaseCode+4, "failed to get "+loanBaseinfoFilesName+" details")
	ErrListLoanBaseinfoFiles       = errcode.NewError(loanBaseinfoFilesBaseCode+5, "failed to list of "+loanBaseinfoFilesName)

	ErrDeleteByIDsLoanBaseinfoFiles    = errcode.NewError(loanBaseinfoFilesBaseCode+6, "failed to delete by batch ids "+loanBaseinfoFilesName)
	ErrGetByConditionLoanBaseinfoFiles = errcode.NewError(loanBaseinfoFilesBaseCode+7, "failed to get "+loanBaseinfoFilesName+" details by conditions")
	ErrListByIDsLoanBaseinfoFiles      = errcode.NewError(loanBaseinfoFilesBaseCode+8, "failed to list by batch ids "+loanBaseinfoFilesName)
	ErrListByLastIDLoanBaseinfoFiles   = errcode.NewError(loanBaseinfoFilesBaseCode+9, "failed to list by last id "+loanBaseinfoFilesName)

	// error codes are globally unique, adding 1 to the previous error code
)
