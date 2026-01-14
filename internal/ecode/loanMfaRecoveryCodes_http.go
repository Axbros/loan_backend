package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanMfaRecoveryCodes business-level http error codes.
// the loanMfaRecoveryCodesNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanMfaRecoveryCodesNO = 88
	loanMfaRecoveryCodesName     = "loanMfaRecoveryCodes"
	loanMfaRecoveryCodesBaseCode = errcode.HCode(loanMfaRecoveryCodesNO)

	ErrCreateLoanMfaRecoveryCodes     = errcode.NewError(loanMfaRecoveryCodesBaseCode+1, "failed to create "+loanMfaRecoveryCodesName)
	ErrDeleteByIDLoanMfaRecoveryCodes = errcode.NewError(loanMfaRecoveryCodesBaseCode+2, "failed to delete "+loanMfaRecoveryCodesName)
	ErrUpdateByIDLoanMfaRecoveryCodes = errcode.NewError(loanMfaRecoveryCodesBaseCode+3, "failed to update "+loanMfaRecoveryCodesName)
	ErrGetByIDLoanMfaRecoveryCodes    = errcode.NewError(loanMfaRecoveryCodesBaseCode+4, "failed to get "+loanMfaRecoveryCodesName+" details")
	ErrListLoanMfaRecoveryCodes       = errcode.NewError(loanMfaRecoveryCodesBaseCode+5, "failed to list of "+loanMfaRecoveryCodesName)

	ErrDeleteByIDsLoanMfaRecoveryCodes    = errcode.NewError(loanMfaRecoveryCodesBaseCode+6, "failed to delete by batch ids "+loanMfaRecoveryCodesName)
	ErrGetByConditionLoanMfaRecoveryCodes = errcode.NewError(loanMfaRecoveryCodesBaseCode+7, "failed to get "+loanMfaRecoveryCodesName+" details by conditions")
	ErrListByIDsLoanMfaRecoveryCodes      = errcode.NewError(loanMfaRecoveryCodesBaseCode+8, "failed to list by batch ids "+loanMfaRecoveryCodesName)
	ErrListByLastIDLoanMfaRecoveryCodes   = errcode.NewError(loanMfaRecoveryCodesBaseCode+9, "failed to list by last id "+loanMfaRecoveryCodesName)

	// error codes are globally unique, adding 1 to the previous error code
)
