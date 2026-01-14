package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanCollectionLogs business-level http error codes.
// the loanCollectionLogsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanCollectionLogsNO = 82
	loanCollectionLogsName     = "loanCollectionLogs"
	loanCollectionLogsBaseCode = errcode.HCode(loanCollectionLogsNO)

	ErrCreateLoanCollectionLogs     = errcode.NewError(loanCollectionLogsBaseCode+1, "failed to create "+loanCollectionLogsName)
	ErrDeleteByIDLoanCollectionLogs = errcode.NewError(loanCollectionLogsBaseCode+2, "failed to delete "+loanCollectionLogsName)
	ErrUpdateByIDLoanCollectionLogs = errcode.NewError(loanCollectionLogsBaseCode+3, "failed to update "+loanCollectionLogsName)
	ErrGetByIDLoanCollectionLogs    = errcode.NewError(loanCollectionLogsBaseCode+4, "failed to get "+loanCollectionLogsName+" details")
	ErrListLoanCollectionLogs       = errcode.NewError(loanCollectionLogsBaseCode+5, "failed to list of "+loanCollectionLogsName)

	ErrDeleteByIDsLoanCollectionLogs    = errcode.NewError(loanCollectionLogsBaseCode+6, "failed to delete by batch ids "+loanCollectionLogsName)
	ErrGetByConditionLoanCollectionLogs = errcode.NewError(loanCollectionLogsBaseCode+7, "failed to get "+loanCollectionLogsName+" details by conditions")
	ErrListByIDsLoanCollectionLogs      = errcode.NewError(loanCollectionLogsBaseCode+8, "failed to list by batch ids "+loanCollectionLogsName)
	ErrListByLastIDLoanCollectionLogs   = errcode.NewError(loanCollectionLogsBaseCode+9, "failed to list by last id "+loanCollectionLogsName)

	// error codes are globally unique, adding 1 to the previous error code
)
