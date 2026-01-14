package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanMfaRecoveryCodesRequest request params
type CreateLoanMfaRecoveryCodesRequest struct {
	UserID   int64      `json:"userID" binding:""`
	CodeHash string     `json:"codeHash" binding:""`
	UsedAt   *time.Time `json:"usedAt" binding:""`
}

// UpdateLoanMfaRecoveryCodesByIDRequest request params
type UpdateLoanMfaRecoveryCodesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID   int64      `json:"userID" binding:""`
	CodeHash string     `json:"codeHash" binding:""`
	UsedAt   *time.Time `json:"usedAt" binding:""`
}

// LoanMfaRecoveryCodesObjDetail detail
type LoanMfaRecoveryCodesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID    int64      `json:"userID"`
	CodeHash  string     `json:"codeHash"`
	UsedAt    *time.Time `json:"usedAt"`
	CreatedAt *time.Time `json:"createdAt"`
}

// CreateLoanMfaRecoveryCodesReply only for api docs
type CreateLoanMfaRecoveryCodesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanMfaRecoveryCodesByIDReply only for api docs
type UpdateLoanMfaRecoveryCodesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanMfaRecoveryCodesByIDReply only for api docs
type GetLoanMfaRecoveryCodesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaRecoveryCodes LoanMfaRecoveryCodesObjDetail `json:"loanMfaRecoveryCodes"`
	} `json:"data"` // return data
}

// DeleteLoanMfaRecoveryCodesByIDReply only for api docs
type DeleteLoanMfaRecoveryCodesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanMfaRecoveryCodessByIDsReply only for api docs
type DeleteLoanMfaRecoveryCodessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanMfaRecoveryCodessRequest request params
type ListLoanMfaRecoveryCodessRequest struct {
	query.Params
}

// ListLoanMfaRecoveryCodessReply only for api docs
type ListLoanMfaRecoveryCodessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaRecoveryCodess []LoanMfaRecoveryCodesObjDetail `json:"loanMfaRecoveryCodess"`
	} `json:"data"` // return data
}

// DeleteLoanMfaRecoveryCodessByIDsRequest request params
type DeleteLoanMfaRecoveryCodessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanMfaRecoveryCodesByConditionRequest request params
type GetLoanMfaRecoveryCodesByConditionRequest struct {
	query.Conditions
}

// GetLoanMfaRecoveryCodesByConditionReply only for api docs
type GetLoanMfaRecoveryCodesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaRecoveryCodes LoanMfaRecoveryCodesObjDetail `json:"loanMfaRecoveryCodes"`
	} `json:"data"` // return data
}

// ListLoanMfaRecoveryCodessByIDsRequest request params
type ListLoanMfaRecoveryCodessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanMfaRecoveryCodessByIDsReply only for api docs
type ListLoanMfaRecoveryCodessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanMfaRecoveryCodess []LoanMfaRecoveryCodesObjDetail `json:"loanMfaRecoveryCodess"`
	} `json:"data"` // return data
}
