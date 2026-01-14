package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanPaymentChannels business-level http error codes.
// the loanPaymentChannelsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanPaymentChannelsNO = 89
	loanPaymentChannelsName     = "loanPaymentChannels"
	loanPaymentChannelsBaseCode = errcode.HCode(loanPaymentChannelsNO)

	ErrCreateLoanPaymentChannels     = errcode.NewError(loanPaymentChannelsBaseCode+1, "failed to create "+loanPaymentChannelsName)
	ErrDeleteByIDLoanPaymentChannels = errcode.NewError(loanPaymentChannelsBaseCode+2, "failed to delete "+loanPaymentChannelsName)
	ErrUpdateByIDLoanPaymentChannels = errcode.NewError(loanPaymentChannelsBaseCode+3, "failed to update "+loanPaymentChannelsName)
	ErrGetByIDLoanPaymentChannels    = errcode.NewError(loanPaymentChannelsBaseCode+4, "failed to get "+loanPaymentChannelsName+" details")
	ErrListLoanPaymentChannels       = errcode.NewError(loanPaymentChannelsBaseCode+5, "failed to list of "+loanPaymentChannelsName)

	ErrDeleteByIDsLoanPaymentChannels    = errcode.NewError(loanPaymentChannelsBaseCode+6, "failed to delete by batch ids "+loanPaymentChannelsName)
	ErrGetByConditionLoanPaymentChannels = errcode.NewError(loanPaymentChannelsBaseCode+7, "failed to get "+loanPaymentChannelsName+" details by conditions")
	ErrListByIDsLoanPaymentChannels      = errcode.NewError(loanPaymentChannelsBaseCode+8, "failed to list by batch ids "+loanPaymentChannelsName)
	ErrListByLastIDLoanPaymentChannels   = errcode.NewError(loanPaymentChannelsBaseCode+9, "failed to list by last id "+loanPaymentChannelsName)

	// error codes are globally unique, adding 1 to the previous error code
)
