package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanPaymentChannelsRequest request params
type CreateLoanPaymentChannelsRequest struct {
	Code             string `json:"code" binding:""`             // 渠道编码(唯一，如 BANK_A、WALLET_X)
	Name             string `json:"name" binding:""`             // 渠道名称
	MerchantNo       string `json:"merchantNo" binding:""`       // 商户号/商户ID(该渠道分配给平台的商户标识)
	Status           int    `json:"status" binding:""`           // 渠道状态：1启用 0禁用
	CanPayout        int    `json:"canPayout" binding:""`        // 是否支持代付/放款：1是 0否
	CanCollect       int    `json:"canCollect" binding:""`       // 是否支持代收/回款：1是 0否
	PayoutFeeRate    string `json:"payoutFeeRate" binding:""`    // 代付手续费率(如0.003500=0.35%)
	PayoutFeeFixed   int    `json:"payoutFeeFixed" binding:""`   // 代付固定手续费(分，若不用可为空)
	CollectFeeRate   string `json:"collectFeeRate" binding:""`   // 代收手续费率
	CollectFeeFixed  int    `json:"collectFeeFixed" binding:""`  // 代收固定手续费(分)
	CollectMinAmount int    `json:"collectMinAmount" binding:""` // 最小代收金额(分)
	CollectMaxAmount int    `json:"collectMaxAmount" binding:""` // 最大代收金额(分)
	PayoutMinAmount  int    `json:"payoutMinAmount" binding:""`  // 最小代付金额(分)
	PayoutMaxAmount  int    `json:"payoutMaxAmount" binding:""`  // 最大代付金额(分)
	SettlementCycle  string `json:"settlementCycle" binding:""`  // 结算周期(如 T0/T1/D1/W1/M1，可按你们渠道定义)
	SettlementDesc   string `json:"settlementDesc" binding:""`   // 结算说明/备注
}

// UpdateLoanPaymentChannelsByIDRequest request params
type UpdateLoanPaymentChannelsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键
	Code             string `json:"code" binding:""`             // 渠道编码(唯一，如 BANK_A、WALLET_X)
	Name             string `json:"name" binding:""`             // 渠道名称
	MerchantNo       string `json:"merchantNo" binding:""`       // 商户号/商户ID(该渠道分配给平台的商户标识)
	Status           int    `json:"status" binding:""`           // 渠道状态：1启用 0禁用
	CanPayout        int    `json:"canPayout" binding:""`        // 是否支持代付/放款：1是 0否
	CanCollect       int    `json:"canCollect" binding:""`       // 是否支持代收/回款：1是 0否
	PayoutFeeRate    string `json:"payoutFeeRate" binding:""`    // 代付手续费率(如0.003500=0.35%)
	PayoutFeeFixed   int    `json:"payoutFeeFixed" binding:""`   // 代付固定手续费(分，若不用可为空)
	CollectFeeRate   string `json:"collectFeeRate" binding:""`   // 代收手续费率
	CollectFeeFixed  int    `json:"collectFeeFixed" binding:""`  // 代收固定手续费(分)
	CollectMinAmount int    `json:"collectMinAmount" binding:""` // 最小代收金额(分)
	CollectMaxAmount int    `json:"collectMaxAmount" binding:""` // 最大代收金额(分)
	PayoutMinAmount  int    `json:"payoutMinAmount" binding:""`  // 最小代付金额(分)
	PayoutMaxAmount  int    `json:"payoutMaxAmount" binding:""`  // 最大代付金额(分)
	SettlementCycle  string `json:"settlementCycle" binding:""`  // 结算周期(如 T0/T1/D1/W1/M1，可按你们渠道定义)
	SettlementDesc   string `json:"settlementDesc" binding:""`   // 结算说明/备注
}

// LoanPaymentChannelsObjDetail detail
type LoanPaymentChannelsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键
	Code             string     `json:"code"`             // 渠道编码(唯一，如 BANK_A、WALLET_X)
	Name             string     `json:"name"`             // 渠道名称
	MerchantNo       string     `json:"merchantNo"`       // 商户号/商户ID(该渠道分配给平台的商户标识)
	Status           int        `json:"status"`           // 渠道状态：1启用 0禁用
	CanPayout        int        `json:"canPayout"`        // 是否支持代付/放款：1是 0否
	CanCollect       int        `json:"canCollect"`       // 是否支持代收/回款：1是 0否
	PayoutFeeRate    string     `json:"payoutFeeRate"`    // 代付手续费率(如0.003500=0.35%)
	PayoutFeeFixed   int        `json:"payoutFeeFixed"`   // 代付固定手续费(分，若不用可为空)
	CollectFeeRate   string     `json:"collectFeeRate"`   // 代收手续费率
	CollectFeeFixed  int        `json:"collectFeeFixed"`  // 代收固定手续费(分)
	CollectMinAmount int        `json:"collectMinAmount"` // 最小代收金额(分)
	CollectMaxAmount int        `json:"collectMaxAmount"` // 最大代收金额(分)
	PayoutMinAmount  int        `json:"payoutMinAmount"`  // 最小代付金额(分)
	PayoutMaxAmount  int        `json:"payoutMaxAmount"`  // 最大代付金额(分)
	SettlementCycle  string     `json:"settlementCycle"`  // 结算周期(如 T0/T1/D1/W1/M1，可按你们渠道定义)
	SettlementDesc   string     `json:"settlementDesc"`   // 结算说明/备注
	CreatedAt        *time.Time `json:"createdAt"`        // 创建时间
	UpdatedAt        *time.Time `json:"updatedAt"`        // 更新时间
}

// CreateLoanPaymentChannelsReply only for api docs
type CreateLoanPaymentChannelsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanPaymentChannelsByIDReply only for api docs
type UpdateLoanPaymentChannelsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanPaymentChannelsByIDReply only for api docs
type GetLoanPaymentChannelsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPaymentChannels LoanPaymentChannelsObjDetail `json:"loanPaymentChannels"`
	} `json:"data"` // return data
}

// DeleteLoanPaymentChannelsByIDReply only for api docs
type DeleteLoanPaymentChannelsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanPaymentChannelssByIDsReply only for api docs
type DeleteLoanPaymentChannelssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanPaymentChannelssRequest request params
type ListLoanPaymentChannelssRequest struct {
	query.Params
}

// ListLoanPaymentChannelssReply only for api docs
type ListLoanPaymentChannelssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPaymentChannelss []LoanPaymentChannelsObjDetail `json:"loanPaymentChannelss"`
	} `json:"data"` // return data
}

// DeleteLoanPaymentChannelssByIDsRequest request params
type DeleteLoanPaymentChannelssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanPaymentChannelsByConditionRequest request params
type GetLoanPaymentChannelsByConditionRequest struct {
	query.Conditions
}

// GetLoanPaymentChannelsByConditionReply only for api docs
type GetLoanPaymentChannelsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPaymentChannels LoanPaymentChannelsObjDetail `json:"loanPaymentChannels"`
	} `json:"data"` // return data
}

// ListLoanPaymentChannelssByIDsRequest request params
type ListLoanPaymentChannelssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanPaymentChannelssByIDsReply only for api docs
type ListLoanPaymentChannelssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanPaymentChannelss []LoanPaymentChannelsObjDetail `json:"loanPaymentChannelss"`
	} `json:"data"` // return data
}
