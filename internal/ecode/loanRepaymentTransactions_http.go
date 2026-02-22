package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanRepaymentTransactions business-level http error codes.
// the loanRepaymentTransactionsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanRepaymentTransactionsNO       = 93
	loanRepaymentTransactionsName     = "loanRepaymentTransactions"
	loanRepaymentTransactionsBaseCode = errcode.HCode(loanRepaymentTransactionsNO)

	ErrCreateLoanRepaymentTransactions     = errcode.NewError(loanRepaymentTransactionsBaseCode+1, "failed to create "+loanRepaymentTransactionsName)
	ErrDeleteByIDLoanRepaymentTransactions = errcode.NewError(loanRepaymentTransactionsBaseCode+2, "failed to delete "+loanRepaymentTransactionsName)
	ErrUpdateByIDLoanRepaymentTransactions = errcode.NewError(loanRepaymentTransactionsBaseCode+3, "failed to update "+loanRepaymentTransactionsName)
	ErrGetByIDLoanRepaymentTransactions    = errcode.NewError(loanRepaymentTransactionsBaseCode+4, "failed to get "+loanRepaymentTransactionsName+" details")
	ErrListLoanRepaymentTransactions       = errcode.NewError(loanRepaymentTransactionsBaseCode+5, "failed to list of "+loanRepaymentTransactionsName)

	ErrDeleteByIDsLoanRepaymentTransactions    = errcode.NewError(loanRepaymentTransactionsBaseCode+6, "failed to delete by batch ids "+loanRepaymentTransactionsName)
	ErrGetByConditionLoanRepaymentTransactions = errcode.NewError(loanRepaymentTransactionsBaseCode+7, "failed to get "+loanRepaymentTransactionsName+" details by conditions")
	ErrListByIDsLoanRepaymentTransactions      = errcode.NewError(loanRepaymentTransactionsBaseCode+8, "failed to list by batch ids "+loanRepaymentTransactionsName)
	ErrListByLastIDLoanRepaymentTransactions   = errcode.NewError(loanRepaymentTransactionsBaseCode+9, "failed to list by last id "+loanRepaymentTransactionsName)
	UnsupportedFileType                        = errcode.NewError(loanRepaymentTransactionsBaseCode+10, "unsupported file type")
	ErrCreateFileFolder                        = errcode.NewError(loanRepaymentTransactionsBaseCode+11, "创建文件目录失败！")
	ErrInvalidFilePath                         = errcode.NewError(loanRepaymentTransactionsBaseCode+12, "创建本地文件失败")
	ErrSaveFile                                = errcode.NewError(loanRepaymentTransactionsBaseCode+13, "保存文件失败")
	FileNotFound                               = errcode.NewError(loanRepaymentTransactionsBaseCode+14, "file not found")
	ErrReadFile                                = errcode.NewError(loanRepaymentTransactionsBaseCode+15, "failed to read the file")
	// error codes are globally unique, adding 1 to the previous error code
)
