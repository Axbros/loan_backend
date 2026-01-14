package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanBaseinfoRequest request params
type CreateLoanBaseinfoRequest struct {
	FirstName                     string     `json:"firstName" binding:""`                     // 姓
	SecondName                    string     `json:"secondName" binding:""`                    // 名
	Age                           int        `json:"age" binding:""`                           // 年齡
	Gender                        string     `json:"gender" binding:""`                        // 性別
	IdType                        string     `json:"idType" binding:""`                        // 證件類型
	IdNumber                      string     `json:"idNumber" binding:""`                      // 證件號碼
	IdCard                        string     `json:"idCard" binding:""`                        // 證件
	Operator                      string     `json:"operator" binding:""`                      // 操作系統
	Inviter                       string     `json:"inviter" binding:""`                       // 邀請人
	Work                          string     `json:"work" binding:""`                          // 工作
	Company                       string     `json:"company" binding:""`                       // 公司
	Salary                        int        `json:"salary" binding:""`                        // 薪資
	MaritalStatus                 int        `json:"maritalStatus" binding:""`                 // 婚否
	HasHouse                      int        `json:"hasHouse" binding:""`                      // 是否有房
	PropertyCertificate           string     `json:"propertyCertificate" binding:""`           // 房產證
	HasCar                        int        `json:"hasCar" binding:""`                        // 是否有車
	VehicleRgistrationCertificate string     `json:"vehicleRgistrationCertificate" binding:""` // 行駛證
	ApplicationAmount             int        `json:"applicationAmount" binding:""`             // 申請金額
	AuditStatus                   int        `json:"auditStatus" binding:""`                   // 審核情況 0待審核 1審核通過 -1 審核拒絕
	BankNo                        string     `json:"bankNo" binding:""`                        // 銀行卡號
	ClientIP                      string     `json:"clientIP" binding:""`                      // 客户端IP地址(IPv4/IPv6)
	ReferrerUserID                int64      `json:"referrerUserID" binding:""`                // 邀请人/分享人(loan_users.id)
	RefCode                       string     `json:"refCode" binding:""`                       // 访问时携带的ref(冗余存储便于排查)
	LoanDays                      int        `json:"loanDays" binding:""`                      // 借款天数(单位：天)
	RiskListStatus                int        `json:"riskListStatus" binding:""`                // 名单状态：0正常 1白名单 2黑名单
	RiskListReason                string     `json:"riskListReason" binding:""`                // 名单原因/来源说明
	RiskListMarkedAt              *time.Time `json:"riskListMarkedAt" binding:""`              // 名单标记时间
}

// UpdateLoanBaseinfoByIDRequest request params
type UpdateLoanBaseinfoByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	FirstName                     string     `json:"firstName" binding:""`                     // 姓
	SecondName                    string     `json:"secondName" binding:""`                    // 名
	Age                           int        `json:"age" binding:""`                           // 年齡
	Gender                        string     `json:"gender" binding:""`                        // 性別
	IdType                        string     `json:"idType" binding:""`                        // 證件類型
	IdNumber                      string     `json:"idNumber" binding:""`                      // 證件號碼
	IdCard                        string     `json:"idCard" binding:""`                        // 證件
	Operator                      string     `json:"operator" binding:""`                      // 操作系統
	Inviter                       string     `json:"inviter" binding:""`                       // 邀請人
	Work                          string     `json:"work" binding:""`                          // 工作
	Company                       string     `json:"company" binding:""`                       // 公司
	Salary                        int        `json:"salary" binding:""`                        // 薪資
	MaritalStatus                 int        `json:"maritalStatus" binding:""`                 // 婚否
	HasHouse                      int        `json:"hasHouse" binding:""`                      // 是否有房
	PropertyCertificate           string     `json:"propertyCertificate" binding:""`           // 房產證
	HasCar                        int        `json:"hasCar" binding:""`                        // 是否有車
	VehicleRgistrationCertificate string     `json:"vehicleRgistrationCertificate" binding:""` // 行駛證
	ApplicationAmount             int        `json:"applicationAmount" binding:""`             // 申請金額
	AuditStatus                   int        `json:"auditStatus" binding:""`                   // 審核情況 0待審核 1審核通過 -1 審核拒絕
	BankNo                        string     `json:"bankNo" binding:""`                        // 銀行卡號
	ClientIP                      string     `json:"clientIP" binding:""`                      // 客户端IP地址(IPv4/IPv6)
	ReferrerUserID                int64      `json:"referrerUserID" binding:""`                // 邀请人/分享人(loan_users.id)
	RefCode                       string     `json:"refCode" binding:""`                       // 访问时携带的ref(冗余存储便于排查)
	LoanDays                      int        `json:"loanDays" binding:""`                      // 借款天数(单位：天)
	RiskListStatus                int        `json:"riskListStatus" binding:""`                // 名单状态：0正常 1白名单 2黑名单
	RiskListReason                string     `json:"riskListReason" binding:""`                // 名单原因/来源说明
	RiskListMarkedAt              *time.Time `json:"riskListMarkedAt" binding:""`              // 名单标记时间
}

// LoanBaseinfoObjDetail detail
type LoanBaseinfoObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	FirstName                     string              `json:"firstName"`                     // 姓
	SecondName                    string              `json:"secondName"`                    // 名
	Age                           int                 `json:"age"`                           // 年齡
	Gender                        string              `json:"gender"`                        // 性別
	IdType                        string              `json:"idType"`                        // 證件類型
	IdNumber                      string              `json:"idNumber"`                      // 證件號碼
	IdCard                        string              `json:"idCard"`                        // 證件
	Operator                      string              `json:"operator"`                      // 操作系統
	Inviter                       string              `json:"inviter"`                       // 邀請人
	Work                          string              `json:"work"`                          // 工作
	Company                       string              `json:"company"`                       // 公司
	Salary                        int                 `json:"salary"`                        // 薪資
	MaritalStatus                 int                 `json:"maritalStatus"`                 // 婚否
	HasHouse                      int                 `json:"hasHouse"`                      // 是否有房
	PropertyCertificate           string              `json:"propertyCertificate"`           // 房產證
	HasCar                        int                 `json:"hasCar"`                        // 是否有車
	VehicleRgistrationCertificate string              `json:"vehicleRgistrationCertificate"` // 行駛證
	ApplicationAmount             int                 `json:"applicationAmount"`             // 申請金額
	AuditStatus                   int                 `json:"auditStatus"`                   // 審核情況 0待審核 1審核通過 -1 審核拒絕
	BankNo                        string              `json:"bankNo"`                        // 銀行卡號
	ClientIP                      string              `json:"clientIP"`                      // 客户端IP地址(IPv4/IPv6)
	CreatedAt                     *time.Time          `json:"createdAt"`
	UpdatedAt                     *time.Time          `json:"updatedAt"`
	ReferrerUserID                int64               `json:"referrerUserID"`   // 邀请人/分享人(loan_users.id)
	RefCode                       string              `json:"refCode"`          // 访问时携带的ref(冗余存储便于排查)
	LoanDays                      int                 `json:"loanDays"`         // 借款天数(单位：天)
	RiskListStatus                int                 `json:"riskListStatus"`   // 名单状态：0正常 1白名单 2黑名单
	RiskListReason                string              `json:"riskListReason"`   // 名单原因/来源说明
	RiskListMarkedAt              *time.Time          `json:"riskListMarkedAt"` // 名单标记时间
	Files                         map[string][]string `json:"files"`
}

type LoanBaseinfoSimpleObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	FirstName         string `json:"firstName"`         // 姓
	SecondName        string `json:"secondName"`        // 名
	Age               int    `json:"age"`               // 年齡
	Gender            string `json:"gender"`            // 性別
	IdType            string `json:"idType"`            // 證件類型
	IdNumber          string `json:"idNumber"`          // 證件號碼
	ApplicationAmount int    `json:"applicationAmount"` // 申請金額
	AuditStatus       int    `json:"auditStatus"`       // 審核情況 0待審核 1審核通過 -1 審核拒絕
	ReferrerUserID    int64  `json:"referrerUserID"`    // 邀请人/分享人(loan_users.id)
	LoanDays          int    `json:"loanDays"`          // 借款天数(单位：天)
}

// CreateLoanBaseinfoReply only for api docs
type CreateLoanBaseinfoReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanBaseinfoByIDReply only for api docs
type UpdateLoanBaseinfoByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanBaseinfoByIDReply only for api docs
type GetLoanBaseinfoByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfo LoanBaseinfoObjDetail `json:"loanBaseinfo"`
	} `json:"data"` // return data
}

// DeleteLoanBaseinfoByIDReply only for api docs
type DeleteLoanBaseinfoByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanBaseinfosByIDsReply only for api docs
type DeleteLoanBaseinfosByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanBaseinfosRequest request params
type ListLoanBaseinfosRequest struct {
	query.Params
}

// ListLoanBaseinfosReply only for api docs
type ListLoanBaseinfosReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfos []LoanBaseinfoObjDetail `json:"loanBaseinfos"`
	} `json:"data"` // return data
}

// DeleteLoanBaseinfosByIDsRequest request params
type DeleteLoanBaseinfosByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanBaseinfoByConditionRequest request params
type GetLoanBaseinfoByConditionRequest struct {
	query.Conditions
}

// GetLoanBaseinfoByConditionReply only for api docs
type GetLoanBaseinfoByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfo LoanBaseinfoObjDetail `json:"loanBaseinfo"`
	} `json:"data"` // return data
}

// ListLoanBaseinfosByIDsRequest request params
type ListLoanBaseinfosByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanBaseinfosByIDsReply only for api docs
type ListLoanBaseinfosByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfos []LoanBaseinfoObjDetail `json:"loanBaseinfos"`
	} `json:"data"` // return data
}
