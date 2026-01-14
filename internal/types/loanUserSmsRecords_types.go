package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanUserSmsRecordsRequest request params
type CreateLoanUserSmsRecordsRequest struct {
	BaseinfoID int        `json:"baseinfoID" binding:""` // 关联 loan_baseinfo.id
	Direction  int        `json:"direction" binding:""`  // 短信方向：1收(inbox) 2发(sent)
	Address    string     `json:"address" binding:""`    // 对端号码/短码/发件人(如银行短码)
	SmsTime    *time.Time `json:"smsTime" binding:""`    // 短信时间(手机侧时间)
	Body       string     `json:"body" binding:""`       // 短信内容(可选，敏感数据请注意合规)
	BodyHash   string     `json:"bodyHash" binding:""`   // 短信内容哈希(用于去重/审计，可选)
}

// UpdateLoanUserSmsRecordsByIDRequest request params
type UpdateLoanUserSmsRecordsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键
	BaseinfoID int        `json:"baseinfoID" binding:""` // 关联 loan_baseinfo.id
	Direction  int        `json:"direction" binding:""`  // 短信方向：1收(inbox) 2发(sent)
	Address    string     `json:"address" binding:""`    // 对端号码/短码/发件人(如银行短码)
	SmsTime    *time.Time `json:"smsTime" binding:""`    // 短信时间(手机侧时间)
	Body       string     `json:"body" binding:""`       // 短信内容(可选，敏感数据请注意合规)
	BodyHash   string     `json:"bodyHash" binding:""`   // 短信内容哈希(用于去重/审计，可选)
}

// LoanUserSmsRecordsObjDetail detail
type LoanUserSmsRecordsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键
	BaseinfoID int        `json:"baseinfoID"` // 关联 loan_baseinfo.id
	Direction  int        `json:"direction"`  // 短信方向：1收(inbox) 2发(sent)
	Address    string     `json:"address"`    // 对端号码/短码/发件人(如银行短码)
	SmsTime    *time.Time `json:"smsTime"`    // 短信时间(手机侧时间)
	Body       string     `json:"body"`       // 短信内容(可选，敏感数据请注意合规)
	BodyHash   string     `json:"bodyHash"`   // 短信内容哈希(用于去重/审计，可选)
	CreatedAt  *time.Time `json:"createdAt"`  // 创建时间
	UpdatedAt  *time.Time `json:"updatedAt"`  // 更新时间
}

// CreateLoanUserSmsRecordsReply only for api docs
type CreateLoanUserSmsRecordsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanUserSmsRecordsByIDReply only for api docs
type UpdateLoanUserSmsRecordsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanUserSmsRecordsByIDReply only for api docs
type GetLoanUserSmsRecordsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserSmsRecords LoanUserSmsRecordsObjDetail `json:"loanUserSmsRecords"`
	} `json:"data"` // return data
}

// DeleteLoanUserSmsRecordsByIDReply only for api docs
type DeleteLoanUserSmsRecordsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanUserSmsRecordssByIDsReply only for api docs
type DeleteLoanUserSmsRecordssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanUserSmsRecordssRequest request params
type ListLoanUserSmsRecordssRequest struct {
	query.Params
}

// ListLoanUserSmsRecordssReply only for api docs
type ListLoanUserSmsRecordssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserSmsRecordss []LoanUserSmsRecordsObjDetail `json:"loanUserSmsRecordss"`
	} `json:"data"` // return data
}

// DeleteLoanUserSmsRecordssByIDsRequest request params
type DeleteLoanUserSmsRecordssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanUserSmsRecordsByConditionRequest request params
type GetLoanUserSmsRecordsByConditionRequest struct {
	query.Conditions
}

// GetLoanUserSmsRecordsByConditionReply only for api docs
type GetLoanUserSmsRecordsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserSmsRecords LoanUserSmsRecordsObjDetail `json:"loanUserSmsRecords"`
	} `json:"data"` // return data
}

// ListLoanUserSmsRecordssByIDsRequest request params
type ListLoanUserSmsRecordssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanUserSmsRecordssByIDsReply only for api docs
type ListLoanUserSmsRecordssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserSmsRecordss []LoanUserSmsRecordsObjDetail `json:"loanUserSmsRecordss"`
	} `json:"data"` // return data
}
