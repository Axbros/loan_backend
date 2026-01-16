package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=1,max=72"`
}

// CreateLoanUsersRequest request params
type CreateLoanUsersRequest struct {
	Username     string `json:"username" binding:""`
	PasswordHash string `json:"passwordHash" binding:""`
	DepartmentID int64  `json:"departmentID" binding:""`
	MfaEnabled   int    `json:"mfaEnabled" binding:""`
	MfaRequired  int    `json:"mfaRequired" binding:""`
	Status       int    `json:"status" binding:""`
	ShareCode    string `json:"shareCode" binding:""` // 分享邀请码(用于生成分享链接，建议唯一)
}

// UpdateLoanUsersByIDRequest request params
type UpdateLoanUsersByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Username     string `json:"username" binding:""`
	PasswordHash string `json:"passwordHash" binding:""`
	DepartmentID int64  `json:"departmentID" binding:""`
	MfaEnabled   int    `json:"mfaEnabled" binding:""`
	MfaRequired  int    `json:"mfaRequired" binding:""`
	Status       int    `json:"status" binding:""`
	ShareCode    string `json:"shareCode" binding:""` // 分享邀请码(用于生成分享链接，建议唯一)
}

// LoanUsersObjDetail detail
type LoanUsersObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Username     string     `json:"username"`
	PasswordHash string     `json:"passwordHash"`
	DepartmentID int64      `json:"departmentID"`
	MfaEnabled   int        `json:"mfaEnabled"`
	MfaRequired  int        `json:"mfaRequired"`
	Status       int        `json:"status"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	ShareCode    string     `json:"shareCode"` // 分享邀请码(用于生成分享链接，建议唯一)
}

// CreateLoanUsersReply only for api docs
type CreateLoanUsersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanUsersByIDReply only for api docs
type UpdateLoanUsersByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanUsersByIDReply only for api docs
type GetLoanUsersByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUsers LoanUsersObjDetail `json:"loanUsers"`
	} `json:"data"` // return data
}

// DeleteLoanUsersByIDReply only for api docs
type DeleteLoanUsersByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanUserssByIDsReply only for api docs
type DeleteLoanUserssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanUserssRequest request params
type ListLoanUserssRequest struct {
	query.Params
}

// ListLoanUserssReply only for api docs
type ListLoanUserssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserss []LoanUsersObjDetail `json:"loanUserss"`
	} `json:"data"` // return data
}

// DeleteLoanUserssByIDsRequest request params
type DeleteLoanUserssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanUsersByConditionRequest request params
type GetLoanUsersByConditionRequest struct {
	query.Conditions
}

// GetLoanUsersByConditionReply only for api docs
type GetLoanUsersByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUsers LoanUsersObjDetail `json:"loanUsers"`
	} `json:"data"` // return data
}

// ListLoanUserssByIDsRequest request params
type ListLoanUserssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanUserssByIDsReply only for api docs
type ListLoanUserssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserss []LoanUsersObjDetail `json:"loanUserss"`
	} `json:"data"` // return data
}

type BindMFARequest struct {
	OTP string `json:"otp" binding:""`
}
