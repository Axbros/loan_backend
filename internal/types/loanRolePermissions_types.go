package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanRolePermissionsRequest request params
type CreateLoanRolePermissionsRequest struct {
	RoleID       int64 `json:"roleID" binding:""`
	PermissionID int64 `json:"permissionID" binding:""`
}

// UpdateLoanRolePermissionsByIDRequest request params
type UpdateLoanRolePermissionsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// surrogate id
}

// LoanRolePermissionsObjDetail detail
type LoanRolePermissionsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// surrogate id
	UpdatedAt *time.Time `json:"updatedAt"` // 更新时间
}

// CreateLoanRolePermissionsReply only for api docs
type CreateLoanRolePermissionsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanRolePermissionsByIDReply only for api docs
type UpdateLoanRolePermissionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanRolePermissionsByIDReply only for api docs
type GetLoanRolePermissionsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRolePermissions LoanRolePermissionsObjDetail `json:"loanRolePermissions"`
	} `json:"data"` // return data
}

// DeleteLoanRolePermissionsByIDReply only for api docs
type DeleteLoanRolePermissionsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanRolePermissionssByIDsReply only for api docs
type DeleteLoanRolePermissionssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanRolePermissionssRequest request params
type ListLoanRolePermissionssRequest struct {
	query.Params
}

// ListLoanRolePermissionssReply only for api docs
type ListLoanRolePermissionssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRolePermissionss []LoanRolePermissionsObjDetail `json:"loanRolePermissionss"`
	} `json:"data"` // return data
}

// DeleteLoanRolePermissionssByIDsRequest request params
type DeleteLoanRolePermissionssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanRolePermissionsByConditionRequest request params
type GetLoanRolePermissionsByConditionRequest struct {
	query.Conditions
}

// GetLoanRolePermissionsByConditionReply only for api docs
type GetLoanRolePermissionsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRolePermissions LoanRolePermissionsObjDetail `json:"loanRolePermissions"`
	} `json:"data"` // return data
}

// ListLoanRolePermissionssByIDsRequest request params
type ListLoanRolePermissionssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanRolePermissionssByIDsReply only for api docs
type ListLoanRolePermissionssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRolePermissionss []LoanRolePermissionsObjDetail `json:"loanRolePermissionss"`
	} `json:"data"` // return data
}
