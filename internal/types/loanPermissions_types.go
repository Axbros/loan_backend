package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanPermissionsRequest request params
type CreateLoanPermissionsRequest struct {
	Code     string `json:"code" binding:""`
	Name     string `json:"name" binding:""`
	Type     string `json:"type" binding:""`
	Resource string `json:"resource" binding:""`
}

// UpdateLoanPermissionsByIDRequest request params
type UpdateLoanPermissionsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Code     string `json:"code" binding:""`
	Name     string `json:"name" binding:""`
	Type     string `json:"type" binding:""`
	Resource string `json:"resource" binding:""`
}

// LoanPermissionsObjDetail detail
type LoanPermissionsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Resource  string     `json:"resource"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// CreateLoanPermissionsReply only for api docs
type CreateLoanPermissionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanPermissionsByIDReply only for api docs
type UpdateLoanPermissionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanPermissionsByIDReply only for api docs
type GetLoanPermissionsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPermissions LoanPermissionsObjDetail `json:"loanPermissions"`
	} `json:"data"` // return data
}

// DeleteLoanPermissionsByIDReply only for api docs
type DeleteLoanPermissionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanPermissionssByIDsReply only for api docs
type DeleteLoanPermissionssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanPermissionssRequest request params
type ListLoanPermissionssRequest struct {
	query.Params
}

// ListLoanPermissionssReply only for api docs
type ListLoanPermissionssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPermissionss []LoanPermissionsObjDetail `json:"loanPermissionss"`
	} `json:"data"` // return data
}

// DeleteLoanPermissionssByIDsRequest request params
type DeleteLoanPermissionssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanPermissionsByConditionRequest request params
type GetLoanPermissionsByConditionRequest struct {
	query.Conditions
}

// GetLoanPermissionsByConditionReply only for api docs
type GetLoanPermissionsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPermissions LoanPermissionsObjDetail `json:"loanPermissions"`
	} `json:"data"` // return data
}

// ListLoanPermissionssByIDsRequest request params
type ListLoanPermissionssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanPermissionssByIDsReply only for api docs
type ListLoanPermissionssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPermissionss []LoanPermissionsObjDetail `json:"loanPermissionss"`
	} `json:"data"` // return data
}
