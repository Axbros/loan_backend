package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanMfaDevicesRequest request params
type CreateLoanMfaDevicesRequest struct {
	UserID     int64      `json:"userID" binding:""`
	Type       string     `json:"type" binding:""`
	Name       string     `json:"name" binding:""`
	SecretEnc  string     `json:"secretEnc" binding:""`
	IsPrimary  int        `json:"isPrimary" binding:""`
	Status     int        `json:"status" binding:""`
	LastUsedAt *time.Time `json:"lastUsedAt" binding:""`
}

// UpdateLoanMfaDevicesByIDRequest request params
type UpdateLoanMfaDevicesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID     int64      `json:"userID" binding:""`
	Type       string     `json:"type" binding:""`
	Name       string     `json:"name" binding:""`
	SecretEnc  string     `json:"secretEnc" binding:""`
	IsPrimary  int        `json:"isPrimary" binding:""`
	Status     int        `json:"status" binding:""`
	LastUsedAt *time.Time `json:"lastUsedAt" binding:""`
}

// LoanMfaDevicesObjDetail detail
type LoanMfaDevicesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID     int64      `json:"userID"`
	Type       string     `json:"type"`
	Name       string     `json:"name"`
	SecretEnc  string     `json:"secretEnc"`
	IsPrimary  int        `json:"isPrimary"`
	Status     int        `json:"status"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

// CreateLoanMfaDevicesReply only for api docs
type CreateLoanMfaDevicesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanMfaDevicesByIDReply only for api docs
type UpdateLoanMfaDevicesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanMfaDevicesByIDReply only for api docs
type GetLoanMfaDevicesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaDevices LoanMfaDevicesObjDetail `json:"loanMfaDevices"`
	} `json:"data"` // return data
}

// DeleteLoanMfaDevicesByIDReply only for api docs
type DeleteLoanMfaDevicesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanMfaDevicessByIDsReply only for api docs
type DeleteLoanMfaDevicessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanMfaDevicessRequest request params
type ListLoanMfaDevicessRequest struct {
	query.Params
}

// ListLoanMfaDevicessReply only for api docs
type ListLoanMfaDevicessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaDevicess []LoanMfaDevicesObjDetail `json:"loanMfaDevicess"`
	} `json:"data"` // return data
}

// DeleteLoanMfaDevicessByIDsRequest request params
type DeleteLoanMfaDevicessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanMfaDevicesByConditionRequest request params
type GetLoanMfaDevicesByConditionRequest struct {
	query.Conditions
}

// GetLoanMfaDevicesByConditionReply only for api docs
type GetLoanMfaDevicesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaDevices LoanMfaDevicesObjDetail `json:"loanMfaDevices"`
	} `json:"data"` // return data
}

// ListLoanMfaDevicessByIDsRequest request params
type ListLoanMfaDevicessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanMfaDevicessByIDsReply only for api docs
type ListLoanMfaDevicessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaDevicess []LoanMfaDevicesObjDetail `json:"loanMfaDevicess"`
	} `json:"data"` // return data
}
