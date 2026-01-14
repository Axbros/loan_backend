package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanUserContacts business-level http error codes.
// the loanUserContactsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanUserContactsNO = 98
	loanUserContactsName     = "loanUserContacts"
	loanUserContactsBaseCode = errcode.HCode(loanUserContactsNO)

	ErrCreateLoanUserContacts     = errcode.NewError(loanUserContactsBaseCode+1, "failed to create "+loanUserContactsName)
	ErrDeleteByIDLoanUserContacts = errcode.NewError(loanUserContactsBaseCode+2, "failed to delete "+loanUserContactsName)
	ErrUpdateByIDLoanUserContacts = errcode.NewError(loanUserContactsBaseCode+3, "failed to update "+loanUserContactsName)
	ErrGetByIDLoanUserContacts    = errcode.NewError(loanUserContactsBaseCode+4, "failed to get "+loanUserContactsName+" details")
	ErrListLoanUserContacts       = errcode.NewError(loanUserContactsBaseCode+5, "failed to list of "+loanUserContactsName)

	ErrDeleteByIDsLoanUserContacts    = errcode.NewError(loanUserContactsBaseCode+6, "failed to delete by batch ids "+loanUserContactsName)
	ErrGetByConditionLoanUserContacts = errcode.NewError(loanUserContactsBaseCode+7, "failed to get "+loanUserContactsName+" details by conditions")
	ErrListByIDsLoanUserContacts      = errcode.NewError(loanUserContactsBaseCode+8, "failed to list by batch ids "+loanUserContactsName)
	ErrListByLastIDLoanUserContacts   = errcode.NewError(loanUserContactsBaseCode+9, "failed to list by last id "+loanUserContactsName)

	// error codes are globally unique, adding 1 to the previous error code
)
