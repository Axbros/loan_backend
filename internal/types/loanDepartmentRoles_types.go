package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanDepartmentRolesRequest request params
type CreateLoanDepartmentRolesRequest struct {
	DepartmentID int64 `json:"departmentID" binding:""`
	RoleID       int64 `json:"roleID" binding:""`
}

// UpdateLoanDepartmentRolesByIDRequest request params
type UpdateLoanDepartmentRolesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// surrogate id
}

// LoanDepartmentRolesObjDetail detail
type LoanDepartmentRolesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// surrogate id
	UpdatedAt *time.Time `json:"updatedAt"` // 更新时间
}

// CreateLoanDepartmentRolesReply only for api docs
type CreateLoanDepartmentRolesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanDepartmentRolesByIDReply only for api docs
type UpdateLoanDepartmentRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanDepartmentRolesByIDReply only for api docs
type GetLoanDepartmentRolesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartmentRoles LoanDepartmentRolesObjDetail `json:"loanDepartmentRoles"`
	} `json:"data"` // return data
}

// DeleteLoanDepartmentRolesByIDReply only for api docs
type DeleteLoanDepartmentRolesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanDepartmentRolessByIDsReply only for api docs
type DeleteLoanDepartmentRolessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanDepartmentRolessRequest request params
type ListLoanDepartmentRolessRequest struct {
	query.Params
}

// ListLoanDepartmentRolessReply only for api docs
type ListLoanDepartmentRolessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartmentRoless []LoanDepartmentRolesObjDetail `json:"loanDepartmentRoless"`
	} `json:"data"` // return data
}

// DeleteLoanDepartmentRolessByIDsRequest request params
type DeleteLoanDepartmentRolessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanDepartmentRolesByConditionRequest request params
type GetLoanDepartmentRolesByConditionRequest struct {
	query.Conditions
}

// GetLoanDepartmentRolesByConditionReply only for api docs
type GetLoanDepartmentRolesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartmentRoles LoanDepartmentRolesObjDetail `json:"loanDepartmentRoles"`
	} `json:"data"` // return data
}

// ListLoanDepartmentRolessByIDsRequest request params
type ListLoanDepartmentRolessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanDepartmentRolessByIDsReply only for api docs
type ListLoanDepartmentRolessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartmentRoless []LoanDepartmentRolesObjDetail `json:"loanDepartmentRoless"`
	} `json:"data"` // return data
}
