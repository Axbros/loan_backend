package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanUserCallRecordsRequest request params
type CreateLoanUserCallRecordsRequest struct {
	BaseinfoID      int        `json:"baseinfoID" binding:""`      // 关联 loan_baseinfo.id
	CallType        int        `json:"callType" binding:""`        // 通话类型：1呼入 2呼出 3未接 4拒接(按采集端定义)
	PhoneNumber     string     `json:"phoneNumber" binding:""`     // 对端号码/电话
	PhoneNormalized string     `json:"phoneNormalized" binding:""` // 标准化号码(去空格/国家码等，可选)
	CallTime        *time.Time `json:"callTime" binding:""`        // 通话开始时间(手机侧时间)
	DurationSeconds int        `json:"durationSeconds" binding:""` // 通话时长(秒，未接/拒接一般为0)
	CallHash        string     `json:"callHash" binding:""`        // 记录去重哈希(如 sha256(type+phone+call_time+duration)，可选)
}

// UpdateLoanUserCallRecordsByIDRequest request params
type UpdateLoanUserCallRecordsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键
	BaseinfoID      int        `json:"baseinfoID" binding:""`      // 关联 loan_baseinfo.id
	CallType        int        `json:"callType" binding:""`        // 通话类型：1呼入 2呼出 3未接 4拒接(按采集端定义)
	PhoneNumber     string     `json:"phoneNumber" binding:""`     // 对端号码/电话
	PhoneNormalized string     `json:"phoneNormalized" binding:""` // 标准化号码(去空格/国家码等，可选)
	CallTime        *time.Time `json:"callTime" binding:""`        // 通话开始时间(手机侧时间)
	DurationSeconds int        `json:"durationSeconds" binding:""` // 通话时长(秒，未接/拒接一般为0)
	CallHash        string     `json:"callHash" binding:""`        // 记录去重哈希(如 sha256(type+phone+call_time+duration)，可选)
}

// LoanUserCallRecordsObjDetail detail
type LoanUserCallRecordsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键
	BaseinfoID      int        `json:"baseinfoID"`      // 关联 loan_baseinfo.id
	CallType        int        `json:"callType"`        // 通话类型：1呼入 2呼出 3未接 4拒接(按采集端定义)
	PhoneNumber     string     `json:"phoneNumber"`     // 对端号码/电话
	PhoneNormalized string     `json:"phoneNormalized"` // 标准化号码(去空格/国家码等，可选)
	CallTime        *time.Time `json:"callTime"`        // 通话开始时间(手机侧时间)
	DurationSeconds int        `json:"durationSeconds"` // 通话时长(秒，未接/拒接一般为0)
	CallHash        string     `json:"callHash"`        // 记录去重哈希(如 sha256(type+phone+call_time+duration)，可选)
	CreatedAt       *time.Time `json:"createdAt"`       // 创建时间
	UpdatedAt       *time.Time `json:"updatedAt"`       // 更新时间
}

// CreateLoanUserCallRecordsReply only for api docs
type CreateLoanUserCallRecordsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanUserCallRecordsByIDReply only for api docs
type UpdateLoanUserCallRecordsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanUserCallRecordsByIDReply only for api docs
type GetLoanUserCallRecordsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserCallRecords LoanUserCallRecordsObjDetail `json:"loanUserCallRecords"`
	} `json:"data"` // return data
}

// DeleteLoanUserCallRecordsByIDReply only for api docs
type DeleteLoanUserCallRecordsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanUserCallRecordssByIDsReply only for api docs
type DeleteLoanUserCallRecordssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanUserCallRecordssRequest request params
type ListLoanUserCallRecordssRequest struct {
	query.Params
}

// ListLoanUserCallRecordssReply only for api docs
type ListLoanUserCallRecordssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserCallRecordss []LoanUserCallRecordsObjDetail `json:"loanUserCallRecordss"`
	} `json:"data"` // return data
}

// DeleteLoanUserCallRecordssByIDsRequest request params
type DeleteLoanUserCallRecordssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanUserCallRecordsByConditionRequest request params
type GetLoanUserCallRecordsByConditionRequest struct {
	query.Conditions
}

// GetLoanUserCallRecordsByConditionReply only for api docs
type GetLoanUserCallRecordsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserCallRecords LoanUserCallRecordsObjDetail `json:"loanUserCallRecords"`
	} `json:"data"` // return data
}

// ListLoanUserCallRecordssByIDsRequest request params
type ListLoanUserCallRecordssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanUserCallRecordssByIDsReply only for api docs
type ListLoanUserCallRecordssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserCallRecordss []LoanUserCallRecordsObjDetail `json:"loanUserCallRecordss"`
	} `json:"data"` // return data
}
