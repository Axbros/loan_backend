package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanMfaDevices business-level http error codes.
// the loanMfaDevicesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanMfaDevicesNO = 87
	loanMfaDevicesName     = "loanMfaDevices"
	loanMfaDevicesBaseCode = errcode.HCode(loanMfaDevicesNO)

	ErrCreateLoanMfaDevices     = errcode.NewError(loanMfaDevicesBaseCode+1, "failed to create "+loanMfaDevicesName)
	ErrDeleteByIDLoanMfaDevices = errcode.NewError(loanMfaDevicesBaseCode+2, "failed to delete "+loanMfaDevicesName)
	ErrUpdateByIDLoanMfaDevices = errcode.NewError(loanMfaDevicesBaseCode+3, "failed to update "+loanMfaDevicesName)
	ErrGetByIDLoanMfaDevices    = errcode.NewError(loanMfaDevicesBaseCode+4, "failed to get "+loanMfaDevicesName+" details")
	ErrListLoanMfaDevices       = errcode.NewError(loanMfaDevicesBaseCode+5, "failed to list of "+loanMfaDevicesName)

	ErrDeleteByIDsLoanMfaDevices    = errcode.NewError(loanMfaDevicesBaseCode+6, "failed to delete by batch ids "+loanMfaDevicesName)
	ErrGetByConditionLoanMfaDevices = errcode.NewError(loanMfaDevicesBaseCode+7, "failed to get "+loanMfaDevicesName+" details by conditions")
	ErrListByIDsLoanMfaDevices      = errcode.NewError(loanMfaDevicesBaseCode+8, "failed to list by batch ids "+loanMfaDevicesName)
	ErrListByLastIDLoanMfaDevices   = errcode.NewError(loanMfaDevicesBaseCode+9, "failed to list by last id "+loanMfaDevicesName)

	// error codes are globally unique, adding 1 to the previous error code
)
