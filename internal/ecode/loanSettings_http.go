package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanSettings business-level http error codes.
// the loanSettingsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanSettingsNO       = 72
	loanSettingsName     = "loanSettings"
	loanSettingsBaseCode = errcode.HCode(loanSettingsNO)

	ErrCreateLoanSettings     = errcode.NewError(loanSettingsBaseCode+1, "failed to create "+loanSettingsName)
	ErrDeleteByIDLoanSettings = errcode.NewError(loanSettingsBaseCode+2, "failed to delete "+loanSettingsName)
	ErrUpdateByIDLoanSettings = errcode.NewError(loanSettingsBaseCode+3, "failed to update "+loanSettingsName)
	ErrGetByIDLoanSettings    = errcode.NewError(loanSettingsBaseCode+4, "failed to get "+loanSettingsName+" details")
	ErrListLoanSettings       = errcode.NewError(loanSettingsBaseCode+5, "failed to list of "+loanSettingsName)

	// error codes are globally unique, adding 1 to the previous error code
)
