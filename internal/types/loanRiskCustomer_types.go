package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanRiskCustomerRequest request params
type CreateLoanRiskCustomerRequest struct {
	LoanBaseinfoID int    `json:"loanBaseinfoID" binding:""`
	RiskType       int    `json:"riskType" binding:""`   // 风险类型 -1 黑名单 1 白名单
	RiskReason     string `json:"riskReason" binding:""` // 风险原因
	CreatedBy      int    `json:"createdBy" binding:""`  // loan_users_id
}

// UpdateLoanRiskCustomerByIDRequest request params
type UpdateLoanRiskCustomerByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	LoanBaseinfoID int    `json:"loanBaseinfoID" binding:""`
	RiskType       int    `json:"riskType" binding:""`   // 风险类型 -1 黑名单 1 白名单
	RiskReason     string `json:"riskReason" binding:""` // 风险原因
	CreatedBy      int    `json:"createdBy" binding:""`  // loan_users_id
}

// LoanRiskCustomerObjDetail detail
type LoanRiskCustomerObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	//LoanBaseinfoID int    `json:"loanBaseinfoID"`
	RiskType   int    `json:"riskType"`   // 风险类型 -1 黑名单 1 白名单
	RiskReason string `json:"riskReason"` // 风险原因

	//CreatedBy int        `json:"createdBy"` // loan_users_id
	CreatedAt *time.Time `json:"createdAt"`
	//UpdatedAt *time.Time `json:"updatedAt"`

	LoanBaseinfo *LoanBaseinfoSimpleObjDetail `json:"loanBaseinfo"`
	OperateUser  *LoanUsersObjSimple          `json:"operateUser"`
}

// CreateLoanRiskCustomerReply only for api docs
type CreateLoanRiskCustomerReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteLoanRiskCustomerByIDReply only for api docs
type DeleteLoanRiskCustomerByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateLoanRiskCustomerByIDReply only for api docs
type UpdateLoanRiskCustomerByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanRiskCustomerByIDReply only for api docs
type GetLoanRiskCustomerByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRiskCustomer LoanRiskCustomerObjDetail `json:"loanRiskCustomer"`
	} `json:"data"` // return data
}

// ListLoanRiskCustomersRequest request params
type ListLoanRiskCustomersRequest struct {
	query.Params
}

// ListLoanRiskCustomersReply only for api docs
type ListLoanRiskCustomersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanRiskCustomers []LoanRiskCustomerObjDetail `json:"loanRiskCustomers"`
	} `json:"data"` // return data
}
