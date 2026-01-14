package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanReferralVisits business-level http error codes.
// the loanReferralVisitsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanReferralVisitsNO = 91
	loanReferralVisitsName     = "loanReferralVisits"
	loanReferralVisitsBaseCode = errcode.HCode(loanReferralVisitsNO)

	ErrCreateLoanReferralVisits     = errcode.NewError(loanReferralVisitsBaseCode+1, "failed to create "+loanReferralVisitsName)
	ErrDeleteByIDLoanReferralVisits = errcode.NewError(loanReferralVisitsBaseCode+2, "failed to delete "+loanReferralVisitsName)
	ErrUpdateByIDLoanReferralVisits = errcode.NewError(loanReferralVisitsBaseCode+3, "failed to update "+loanReferralVisitsName)
	ErrGetByIDLoanReferralVisits    = errcode.NewError(loanReferralVisitsBaseCode+4, "failed to get "+loanReferralVisitsName+" details")
	ErrListLoanReferralVisits       = errcode.NewError(loanReferralVisitsBaseCode+5, "failed to list of "+loanReferralVisitsName)

	ErrDeleteByIDsLoanReferralVisits    = errcode.NewError(loanReferralVisitsBaseCode+6, "failed to delete by batch ids "+loanReferralVisitsName)
	ErrGetByConditionLoanReferralVisits = errcode.NewError(loanReferralVisitsBaseCode+7, "failed to get "+loanReferralVisitsName+" details by conditions")
	ErrListByIDsLoanReferralVisits      = errcode.NewError(loanReferralVisitsBaseCode+8, "failed to list by batch ids "+loanReferralVisitsName)
	ErrListByLastIDLoanReferralVisits   = errcode.NewError(loanReferralVisitsBaseCode+9, "failed to list by last id "+loanReferralVisitsName)

	// error codes are globally unique, adding 1 to the previous error code
)
