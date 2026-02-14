package types

import (
	"github.com/shopspring/decimal"
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanDisbursementsRequest request params
type CreateLoanDisbursementsRequest struct {
	BaseinfoID           int        `json:"baseinfoID" binding:""`           // 关联申请单 loan_baseinfo.id
	DisburseAmount       int        `json:"disburseAmount" binding:""`       // 放款金额(单位按你的系统：元/分，建议统一)
	NetAmount            int        `json:"netAmount" binding:""`            // 到账金额(扣除费用后实际到账)
	Status               int        `json:"status" binding:""`               // 放款状态：0待放款 1已放款
	SourceReferrerUserID int64      `json:"sourceReferrerUserID" binding:""` // 用户来源(分享人 loan_users.id，冗余快照，便于查询)
	AuditorUserID        int64      `json:"auditorUserID" binding:""`        // 审核人员(loan_users.id)
	AuditedAt            *time.Time `json:"auditedAt" binding:""`            // 审核通过时间
	PayoutChannelID      int64      `json:"payoutChannelID" binding:""`      // 放款渠道(代付) loan_payment_channels.id
	PayoutOrderNo        string     `json:"payoutOrderNo" binding:""`        // 放款订单号/三方代付单号
	DisbursedAt          *time.Time `json:"disbursedAt" binding:""`          // 放款时间
}

// UpdateLoanDisbursementsByIDRequest request params
type UpdateLoanDisbursementsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键(放款单ID)
	BaseinfoID           int        `json:"baseinfoID" binding:""`           // 关联申请单 loan_baseinfo.id
	DisburseAmount       int        `json:"disburseAmount" binding:""`       // 放款金额(单位按你的系统：元/分，建议统一)
	NetAmount            int        `json:"netAmount" binding:""`            // 到账金额(扣除费用后实际到账)
	Status               int        `json:"status" binding:""`               // 放款状态：0待放款 1已放款
	SourceReferrerUserID int64      `json:"sourceReferrerUserID" binding:""` // 用户来源(分享人 loan_users.id，冗余快照，便于查询)
	AuditorUserID        int64      `json:"auditorUserID" binding:""`        // 审核人员(loan_users.id)
	AuditedAt            *time.Time `json:"auditedAt" binding:""`            // 审核通过时间
	PayoutChannelID      int64      `json:"payoutChannelID" binding:""`      // 放款渠道(代付) loan_payment_channels.id
	PayoutOrderNo        string     `json:"payoutOrderNo" binding:""`        // 放款订单号/三方代付单号
	DisbursedAt          *time.Time `json:"disbursedAt" binding:""`          // 放款时间
}

// LoanDisbursementsObjDetail detail
type LoanDisbursementsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键(放款单ID)
	BaseinfoID           int        `json:"baseinfoID"`           // 关联申请单 loan_baseinfo.id
	DisburseAmount       int        `json:"disburseAmount"`       // 放款金额(单位按你的系统：元/分，建议统一)
	NetAmount            int        `json:"netAmount"`            // 到账金额(扣除费用后实际到账)
	Status               int        `json:"status"`               // 放款状态：0待放款 1已放款
	SourceReferrerUserID int64      `json:"sourceReferrerUserID"` // 用户来源(分享人 loan_users.id，冗余快照，便于查询)
	AuditorUserID        int64      `json:"auditorUserID"`        // 审核人员(loan_users.id)
	AuditedAt            *time.Time `json:"auditedAt"`            // 审核通过时间
	PayoutChannelID      int64      `json:"payoutChannelID"`      // 放款渠道(代付) loan_payment_channels.id
	PayoutOrderNo        string     `json:"payoutOrderNo"`        // 放款订单号/三方代付单号
	DisbursedAt          *time.Time `json:"disbursedAt"`          // 放款时间
	CreatedAt            *time.Time `json:"createdAt"`            // 创建时间(进入待放款时刻)
	UpdatedAt            *time.Time `json:"updatedAt"`            // 更新时间
}

// CreateLoanDisbursementsReply only for api docs
type CreateLoanDisbursementsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

type LoanDisbursementDetailReply struct {
	Disbursement LoanDisbursementsObjDetail `json:"disbursement"`
	Baseinfo     LoanBaseinfoObjDetail      `json:"baseinfo"`
	Audits       LoanAuditsObjDetail        `json:"audits"`
}

// UpdateLoanDisbursementsByIDReply only for api docs
type UpdateLoanDisbursementsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanDisbursementsByIDReply only for api docs
type GetLoanDisbursementsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursements LoanDisbursementsObjDetail `json:"loanDisbursements"`
	} `json:"data"` // return data
}

// DeleteLoanDisbursementsByIDReply only for api docs
type DeleteLoanDisbursementsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanDisbursementssByIDsReply only for api docs
type DeleteLoanDisbursementssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanDisbursementssRequest request params
type ListLoanDisbursementssRequest struct {
	query.Params
}

type PageRequest struct {
	Page  int `json:"page" form:"page" binding:"gte=0"`   // 页码（从0开始）
	Limit int `json:"limit" form:"limit" binding:"gte=1"` // 页大小（最小1）
}
type BaseCondition struct {
	Name       string `json:"name" form:"name"`                                                       // 姓名（loan_baseinfo.first_name）
	Age        *int   `json:"age" form:"age" binding:"omitempty,gte=0"`                               // 年龄（可选，非0）
	Gender     string `json:"gender" form:"gender" binding:"omitempty,oneof=M W"`                     // 性别（M/W）
	IDType     string `json:"idType" form:"idType" binding:"omitempty,oneof=ID_CARD PASSPORT DRIVER"` // 证件类型
	IDNo       string `json:"idNo" form:"idNo"`                                                       // 证件号码
	LoanAmount *int64 `json:"loanAmount" form:"loanAmount" binding:"omitempty,gte=0"`                 // 申请金额
	Status     *int   `json:"status" form:"status"`
}
type BaseOverviewRequest struct {
	PageRequest                // 嵌入分页参数（继承 Page/Limit 字段）
	Condition   *BaseCondition `json:"condition" form:"condition"`
}

type ListLoanDisbursementsOverviewResponse struct {
	Total int64                `json:"total"` // 总条数
	List  []*LoanDisbursedList `json:"list"`  // 分页数据列表
}

// ListLoanDisbursementssReply only for api docs
type ListLoanDisbursementssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursementss []LoanDisbursementsObjDetail `json:"loanDisbursementss"`
	} `json:"data"` // return data
}

// DeleteLoanDisbursementssByIDsRequest request params
type DeleteLoanDisbursementssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanDisbursementsByConditionRequest request params
type GetLoanDisbursementsByConditionRequest struct {
	query.Conditions
}

// GetLoanDisbursementsByConditionReply only for api docs
type GetLoanDisbursementsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursements LoanDisbursementsObjDetail `json:"loanDisbursements"`
	} `json:"data"` // return data
}

// ListLoanDisbursementssByIDsRequest request params
type ListLoanDisbursementssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanDisbursementssByIDsReply only for api docs
type ListLoanDisbursementssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanDisbursementss []LoanDisbursementsObjDetail `json:"loanDisbursementss"`
	} `json:"data"` // return data
}

// DisbursementWithChannel 放款记录+支付渠道关联结果
type DisbursementWithChannel struct {
	ID             int64            `json:"id"`              // 放款记录ID
	DisburseAmount *decimal.Decimal `json:"disburse_amount"` // 放款金额
	NetAmount      *decimal.Decimal `json:"net_amount"`      // 净金额
	Status         int              `json:"status"`          // 放款状态
	PayoutOrderNo  string           `json:"payout_order_no"` // 放款订单号
	DisbursedAt    *time.Time       `json:"disbursed_at"`
	ChannelName    string           `json:"channel_name"` // 支付渠道名称（对应 c.name）
}

type LoanDisbursedList struct {
	ID                int64            `json:"id" gorm:"column:id"`                                 // 基础信息ID（b.id）
	FirstName         string           `json:"first_name" gorm:"column:first_name"`                 // 姓名
	Age               int              `json:"age" gorm:"column:age"`                               // 年龄
	Gender            string           `json:"gender" gorm:"column:gender"`                         // 性别（1=男/2=女等）
	IDType            string           `json:"id_type" gorm:"column:id_type"`                       // 证件类型
	IDNumber          string           `json:"id_number" gorm:"column:id_number"`                   // 证件号码
	ApplicationAmount *decimal.Decimal `json:"application_amount" gorm:"column:application_amount"` // 申请金额
	NetAmount         *decimal.Decimal `json:"net_amount" gorm:"column:net_amount"`                 // 净放款金额
	LoanDays          int              `json:"loan_days" gorm:"column:loan_days"`                   // 借款天数
	ChannelName       string           `json:"channel_name" gorm:"column:name"`                     // 支付渠道名称（c.name）
	PayoutOrderNo     string           `json:"payout_order_no" gorm:"column:payout_order_no"`       // 放款订单号
	PayoutFeeRate     *decimal.Decimal `json:"payout_fee_rate" gorm:"column:payout_fee_rate"`       // 渠道费率
}
