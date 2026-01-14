package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanReferralVisitsRequest request params
type CreateLoanReferralVisitsRequest struct {
	VisitorID      string     `json:"visitorID" binding:""`      // 访客标识(UUID)，前端cookie生成)
	RefCode        string     `json:"refCode" binding:""`        // 访问链接携带的ref(share_code)
	ReferrerUserID int64      `json:"referrerUserID" binding:""` // 邀请人(loan_users.id)
	LandingPath    string     `json:"landingPath" binding:""`    // 落地页路径
	ClientIP       string     `json:"clientIP" binding:""`       // 访问IP(IPv4/IPv6)
	UserAgent      string     `json:"userAgent" binding:""`      // 浏览器UA
	FirstSeenAt    *time.Time `json:"firstSeenAt" binding:""`    // 首次访问时间
	LastSeenAt     *time.Time `json:"lastSeenAt" binding:""`     // 最近访问时间
	VisitCount     int        `json:"visitCount" binding:""`     // 访问次数
}

// UpdateLoanReferralVisitsByIDRequest request params
type UpdateLoanReferralVisitsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键
	VisitorID      string     `json:"visitorID" binding:""`      // 访客标识(UUID)，前端cookie生成)
	RefCode        string     `json:"refCode" binding:""`        // 访问链接携带的ref(share_code)
	ReferrerUserID int64      `json:"referrerUserID" binding:""` // 邀请人(loan_users.id)
	LandingPath    string     `json:"landingPath" binding:""`    // 落地页路径
	ClientIP       string     `json:"clientIP" binding:""`       // 访问IP(IPv4/IPv6)
	UserAgent      string     `json:"userAgent" binding:""`      // 浏览器UA
	FirstSeenAt    *time.Time `json:"firstSeenAt" binding:""`    // 首次访问时间
	LastSeenAt     *time.Time `json:"lastSeenAt" binding:""`     // 最近访问时间
	VisitCount     int        `json:"visitCount" binding:""`     // 访问次数
}

// LoanReferralVisitsObjDetail detail
type LoanReferralVisitsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键
	VisitorID      string     `json:"visitorID"`      // 访客标识(UUID)，前端cookie生成)
	RefCode        string     `json:"refCode"`        // 访问链接携带的ref(share_code)
	ReferrerUserID int64      `json:"referrerUserID"` // 邀请人(loan_users.id)
	LandingPath    string     `json:"landingPath"`    // 落地页路径
	ClientIP       string     `json:"clientIP"`       // 访问IP(IPv4/IPv6)
	UserAgent      string     `json:"userAgent"`      // 浏览器UA
	FirstSeenAt    *time.Time `json:"firstSeenAt"`    // 首次访问时间
	LastSeenAt     *time.Time `json:"lastSeenAt"`     // 最近访问时间
	VisitCount     int        `json:"visitCount"`     // 访问次数
	CreatedAt      *time.Time `json:"createdAt"`      // 创建时间
	UpdatedAt      *time.Time `json:"updatedAt"`      // 更新时间
}

// CreateLoanReferralVisitsReply only for api docs
type CreateLoanReferralVisitsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanReferralVisitsByIDReply only for api docs
type UpdateLoanReferralVisitsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanReferralVisitsByIDReply only for api docs
type GetLoanReferralVisitsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanReferralVisits LoanReferralVisitsObjDetail `json:"loanReferralVisits"`
	} `json:"data"` // return data
}

// DeleteLoanReferralVisitsByIDReply only for api docs
type DeleteLoanReferralVisitsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanReferralVisitssByIDsReply only for api docs
type DeleteLoanReferralVisitssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanReferralVisitssRequest request params
type ListLoanReferralVisitssRequest struct {
	query.Params
}

// ListLoanReferralVisitssReply only for api docs
type ListLoanReferralVisitssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanReferralVisitss []LoanReferralVisitsObjDetail `json:"loanReferralVisitss"`
	} `json:"data"` // return data
}

// DeleteLoanReferralVisitssByIDsRequest request params
type DeleteLoanReferralVisitssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanReferralVisitsByConditionRequest request params
type GetLoanReferralVisitsByConditionRequest struct {
	query.Conditions
}

// GetLoanReferralVisitsByConditionReply only for api docs
type GetLoanReferralVisitsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanReferralVisits LoanReferralVisitsObjDetail `json:"loanReferralVisits"`
	} `json:"data"` // return data
}

// ListLoanReferralVisitssByIDsRequest request params
type ListLoanReferralVisitssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanReferralVisitssByIDsReply only for api docs
type ListLoanReferralVisitssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanReferralVisitss []LoanReferralVisitsObjDetail `json:"loanReferralVisitss"`
	} `json:"data"` // return data
}
