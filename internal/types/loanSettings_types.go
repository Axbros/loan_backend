package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanSettingsRequest request params
type CreateLoanSettingsRequest struct {
	Name   string `json:"name" binding:""`
	Remark string `json:"remark" binding:""`
	Value  string `json:"value" binding:""`
}

// UpdateLoanSettingsByIDRequest request params
type UpdateLoanSettingsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name   string `json:"name" binding:""`
	Remark string `json:"remark" binding:""`
	Value  string `json:"value" binding:""`
}

// LoanSettingsObjDetail detail
type LoanSettingsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Name      string     `json:"name"`
	Value     string     `json:"value"`
	Remark    string     `json:"remark"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// CreateLoanSettingsReply only for api docs
type CreateLoanSettingsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteLoanSettingsByIDReply only for api docs
type DeleteLoanSettingsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateLoanSettingsByIDReply only for api docs
type UpdateLoanSettingsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanSettingsByIDReply only for api docs
type GetLoanSettingsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanSettings LoanSettingsObjDetail `json:"loanSettings"`
	} `json:"data"` // return data
}

// ListLoanSettingssRequest request params
type ListLoanSettingssRequest struct {
	query.Params
}

// ListLoanSettingssReply only for api docs
type ListLoanSettingssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanSettingss []LoanSettingsObjDetail `json:"loanSettingss"`
	} `json:"data"` // return data
}
