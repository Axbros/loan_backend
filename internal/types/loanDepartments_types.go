package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanDepartmentsRequest request params
type CreateLoanDepartmentsRequest struct {
	Name     string `json:"name" binding:""`
	ParentID int64  `json:"parentID" binding:""`
	Status   int    `json:"status" binding:""`
}

// UpdateLoanDepartmentsByIDRequest request params
type UpdateLoanDepartmentsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name     string `json:"name" binding:""`
	ParentID int64  `json:"parentID" binding:""`
	Status   int    `json:"status" binding:""`
}

// LoanDepartmentsObjDetail detail
type LoanDepartmentsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Name      string     `json:"name"`
	ParentID  int64      `json:"parentID"`
	Status    int        `json:"status"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// CreateLoanDepartmentsReply only for api docs
type CreateLoanDepartmentsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanDepartmentsByIDReply only for api docs
type UpdateLoanDepartmentsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanDepartmentsByIDReply only for api docs
type GetLoanDepartmentsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartments LoanDepartmentsObjDetail `json:"loanDepartments"`
	} `json:"data"` // return data
}

// DeleteLoanDepartmentsByIDReply only for api docs
type DeleteLoanDepartmentsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanDepartmentssByIDsReply only for api docs
type DeleteLoanDepartmentssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanDepartmentssRequest request params
type ListLoanDepartmentssRequest struct {
	query.Params
}

// ListLoanDepartmentssReply only for api docs
type ListLoanDepartmentssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartmentss []LoanDepartmentsObjDetail `json:"loanDepartmentss"`
	} `json:"data"` // return data
}

// DeleteLoanDepartmentssByIDsRequest request params
type DeleteLoanDepartmentssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanDepartmentsByConditionRequest request params
type GetLoanDepartmentsByConditionRequest struct {
	query.Conditions
}

// GetLoanDepartmentsByConditionReply only for api docs
type GetLoanDepartmentsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartments LoanDepartmentsObjDetail `json:"loanDepartments"`
	} `json:"data"` // return data
}

// ListLoanDepartmentssByIDsRequest request params
type ListLoanDepartmentssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanDepartmentssByIDsReply only for api docs
type ListLoanDepartmentssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDepartmentss []LoanDepartmentsObjDetail `json:"loanDepartmentss"`
	} `json:"data"` // return data
}
