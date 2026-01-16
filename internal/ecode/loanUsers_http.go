package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanUsers business-level http error codes.
// the loanUsersNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanUsersNO       = 102
	loanUsersName     = "loanUsers"
	loanUsersBaseCode = errcode.HCode(loanUsersNO)

	ErrCreateLoanUsers     = errcode.NewError(loanUsersBaseCode+1, "failed to create "+loanUsersName)
	ErrDeleteByIDLoanUsers = errcode.NewError(loanUsersBaseCode+2, "failed to delete "+loanUsersName)
	ErrUpdateByIDLoanUsers = errcode.NewError(loanUsersBaseCode+3, "failed to update "+loanUsersName)
	ErrGetByIDLoanUsers    = errcode.NewError(loanUsersBaseCode+4, "failed to get "+loanUsersName+" details")
	ErrListLoanUsers       = errcode.NewError(loanUsersBaseCode+5, "failed to list of "+loanUsersName)

	ErrDeleteByIDsLoanUsers    = errcode.NewError(loanUsersBaseCode+6, "failed to delete by batch ids "+loanUsersName)
	ErrGetByConditionLoanUsers = errcode.NewError(loanUsersBaseCode+7, "failed to get "+loanUsersName+" details by conditions")
	ErrListByIDsLoanUsers      = errcode.NewError(loanUsersBaseCode+8, "failed to list by batch ids "+loanUsersName)
	ErrListByLastIDLoanUsers   = errcode.NewError(loanUsersBaseCode+9, "failed to list by last id "+loanUsersName)
	MFAAlreadyEnabled          = errcode.NewError(loanUserRolesBaseCode+14, "mfa already enabled")
	ErrGenerateOTP             = errcode.NewError(loanUsersBaseCode+15, "failed to generate OTP")
	ErrEncrypt                 = errcode.NewError(loanUsersBaseCode+16, "failed to encrypt")
	ErrSecret                  = errcode.NewError(loanUsersBaseCode+17, "failed to decrypt secret")
	ErrValidateSecret          = errcode.NewError(loanUsersBaseCode+18, "failed to validate secret")
	// error codes are globally unique, adding 1 to the previous error code
)
