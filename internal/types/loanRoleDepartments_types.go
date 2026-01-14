package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanRoleDepartmentsRequest request params
type CreateLoanRoleDepartmentsRequest struct {
	RoleID       int64 `json:"roleID" binding:""`
	DepartmentID int64 `json:"departmentID" binding:""`
}

// UpdateLoanRoleDepartmentsByIDRequest request params
type UpdateLoanRoleDepartmentsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// surrogate id
}

// LoanRoleDepartmentsObjDetail detail
type LoanRoleDepartmentsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// surrogate id
	UpdatedAt *time.Time `json:"updatedAt"` // 更新时间
}

// CreateLoanRoleDepartmentsReply only for api docs
type CreateLoanRoleDepartmentsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanRoleDepartmentsByIDReply only for api docs
type UpdateLoanRoleDepartmentsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanRoleDepartmentsByIDReply only for api docs
type GetLoanRoleDepartmentsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoleDepartments LoanRoleDepartmentsObjDetail `json:"loanRoleDepartments"`
	} `json:"data"` // return data
}

// DeleteLoanRoleDepartmentsByIDReply only for api docs
type DeleteLoanRoleDepartmentsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanRoleDepartmentssByIDsReply only for api docs
type DeleteLoanRoleDepartmentssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanRoleDepartmentssRequest request params
type ListLoanRoleDepartmentssRequest struct {
	query.Params
}

// ListLoanRoleDepartmentssReply only for api docs
type ListLoanRoleDepartmentssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoleDepartmentss []LoanRoleDepartmentsObjDetail `json:"loanRoleDepartmentss"`
	} `json:"data"` // return data
}

// DeleteLoanRoleDepartmentssByIDsRequest request params
type DeleteLoanRoleDepartmentssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanRoleDepartmentsByConditionRequest request params
type GetLoanRoleDepartmentsByConditionRequest struct {
	query.Conditions
}

// GetLoanRoleDepartmentsByConditionReply only for api docs
type GetLoanRoleDepartmentsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoleDepartments LoanRoleDepartmentsObjDetail `json:"loanRoleDepartments"`
	} `json:"data"` // return data
}

// ListLoanRoleDepartmentssByIDsRequest request params
type ListLoanRoleDepartmentssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanRoleDepartmentssByIDsReply only for api docs
type ListLoanRoleDepartmentssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRoleDepartmentss []LoanRoleDepartmentsObjDetail `json:"loanRoleDepartmentss"`
	} `json:"data"` // return data
}
