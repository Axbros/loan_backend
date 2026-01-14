package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// loanUserDeviceApps business-level http error codes.
// the loanUserDeviceAppsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	loanUserDeviceAppsNO = 99
	loanUserDeviceAppsName     = "loanUserDeviceApps"
	loanUserDeviceAppsBaseCode = errcode.HCode(loanUserDeviceAppsNO)

	ErrCreateLoanUserDeviceApps     = errcode.NewError(loanUserDeviceAppsBaseCode+1, "failed to create "+loanUserDeviceAppsName)
	ErrDeleteByIDLoanUserDeviceApps = errcode.NewError(loanUserDeviceAppsBaseCode+2, "failed to delete "+loanUserDeviceAppsName)
	ErrUpdateByIDLoanUserDeviceApps = errcode.NewError(loanUserDeviceAppsBaseCode+3, "failed to update "+loanUserDeviceAppsName)
	ErrGetByIDLoanUserDeviceApps    = errcode.NewError(loanUserDeviceAppsBaseCode+4, "failed to get "+loanUserDeviceAppsName+" details")
	ErrListLoanUserDeviceApps       = errcode.NewError(loanUserDeviceAppsBaseCode+5, "failed to list of "+loanUserDeviceAppsName)

	ErrDeleteByIDsLoanUserDeviceApps    = errcode.NewError(loanUserDeviceAppsBaseCode+6, "failed to delete by batch ids "+loanUserDeviceAppsName)
	ErrGetByConditionLoanUserDeviceApps = errcode.NewError(loanUserDeviceAppsBaseCode+7, "failed to get "+loanUserDeviceAppsName+" details by conditions")
	ErrListByIDsLoanUserDeviceApps      = errcode.NewError(loanUserDeviceAppsBaseCode+8, "failed to list by batch ids "+loanUserDeviceAppsName)
	ErrListByLastIDLoanUserDeviceApps   = errcode.NewError(loanUserDeviceAppsBaseCode+9, "failed to list by last id "+loanUserDeviceAppsName)

	// error codes are globally unique, adding 1 to the previous error code
)
