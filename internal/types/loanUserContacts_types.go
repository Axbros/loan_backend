package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanUserContactsRequest request params
type CreateLoanUserContactsRequest struct {
	BaseinfoID  int    `json:"baseinfoID" binding:""`  // 关联 loan_baseinfo.id
	ContactName string `json:"contactName" binding:""` // 联系人姓名
	PhoneNumber string `json:"phoneNumber" binding:""` // 联系人手机号/电话
	ContactHash string `json:"contactHash" binding:""` // 联系人去重哈希(如 sha256(name+phone_normalized))
}

// UpdateLoanUserContactsByIDRequest request params
type UpdateLoanUserContactsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键
	BaseinfoID  int    `json:"baseinfoID" binding:""`  // 关联 loan_baseinfo.id
	ContactName string `json:"contactName" binding:""` // 联系人姓名
	PhoneNumber string `json:"phoneNumber" binding:""` // 联系人手机号/电话
	ContactHash string `json:"contactHash" binding:""` // 联系人去重哈希(如 sha256(name+phone_normalized))
}

// LoanUserContactsObjDetail detail
type LoanUserContactsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键
	BaseinfoID  int        `json:"baseinfoID"`  // 关联 loan_baseinfo.id
	ContactName string     `json:"contactName"` // 联系人姓名
	PhoneNumber string     `json:"phoneNumber"` // 联系人手机号/电话
	ContactHash string     `json:"contactHash"` // 联系人去重哈希(如 sha256(name+phone_normalized))
	CreatedAt   *time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   *time.Time `json:"updatedAt"`   // 更新时间
}

// CreateLoanUserContactsReply only for api docs
type CreateLoanUserContactsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanUserContactsByIDReply only for api docs
type UpdateLoanUserContactsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanUserContactsByIDReply only for api docs
type GetLoanUserContactsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserContacts LoanUserContactsObjDetail `json:"loanUserContacts"`
	} `json:"data"` // return data
}

// DeleteLoanUserContactsByIDReply only for api docs
type DeleteLoanUserContactsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanUserContactssByIDsReply only for api docs
type DeleteLoanUserContactssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanUserContactssRequest request params
type ListLoanUserContactssRequest struct {
	query.Params
}

// ListLoanUserContactssReply only for api docs
type ListLoanUserContactssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserContactss []LoanUserContactsObjDetail `json:"loanUserContactss"`
	} `json:"data"` // return data
}

// DeleteLoanUserContactssByIDsRequest request params
type DeleteLoanUserContactssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanUserContactsByConditionRequest request params
type GetLoanUserContactsByConditionRequest struct {
	query.Conditions
}

// GetLoanUserContactsByConditionReply only for api docs
type GetLoanUserContactsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserContacts LoanUserContactsObjDetail `json:"loanUserContacts"`
	} `json:"data"` // return data
}

// ListLoanUserContactssByIDsRequest request params
type ListLoanUserContactssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanUserContactssByIDsReply only for api docs
type ListLoanUserContactssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserContactss []LoanUserContactsObjDetail `json:"loanUserContactss"`
	} `json:"data"` // return data
}
