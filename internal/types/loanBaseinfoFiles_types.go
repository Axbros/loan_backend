package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanBaseinfoFilesRequest request params
type CreateLoanBaseinfoFilesRequest struct {
	BaseinfoID int    `json:"baseinfoID" binding:""` // 关联 loan_baseinfo.id
	Type       string `json:"type" binding:""`       // 文件类型(如 ID_CARD_FRONT / ID_CARD_BACK / TAX_CERT 等)
	OssURL     string `json:"ossURL" binding:""`     // OSS访问地址(或CDN地址)
	OssKey     string `json:"ossKey" binding:""`     // OSS对象Key(内部定位/删除用，可选)
	FileName   string `json:"fileName" binding:""`   // 原始文件名
	MimeType   string `json:"mimeType" binding:""`   // 文件MIME类型
	SizeBytes  int64  `json:"sizeBytes" binding:""`  // 文件大小(字节)
	Sha256     string `json:"sha256" binding:""`     // 文件哈希(sha256，用于去重/校验)
}

// UpdateLoanBaseinfoFilesByIDRequest request params
type UpdateLoanBaseinfoFilesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键
	BaseinfoID int    `json:"baseinfoID" binding:""` // 关联 loan_baseinfo.id
	Type       string `json:"type" binding:""`       // 文件类型(如 ID_CARD_FRONT / ID_CARD_BACK / TAX_CERT 等)
	OssURL     string `json:"ossURL" binding:""`     // OSS访问地址(或CDN地址)
	OssKey     string `json:"ossKey" binding:""`     // OSS对象Key(内部定位/删除用，可选)
	FileName   string `json:"fileName" binding:""`   // 原始文件名
	MimeType   string `json:"mimeType" binding:""`   // 文件MIME类型
	SizeBytes  int64  `json:"sizeBytes" binding:""`  // 文件大小(字节)
	Sha256     string `json:"sha256" binding:""`     // 文件哈希(sha256，用于去重/校验)
}

// LoanBaseinfoFilesObjDetail detail
type LoanBaseinfoFilesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键
	BaseinfoID int        `json:"baseinfoID"` // 关联 loan_baseinfo.id
	Type       string     `json:"type"`       // 文件类型(如 ID_CARD_FRONT / ID_CARD_BACK / TAX_CERT 等)
	OssURL     string     `json:"ossURL"`     // OSS访问地址(或CDN地址)
	OssKey     string     `json:"ossKey"`     // OSS对象Key(内部定位/删除用，可选)
	FileName   string     `json:"fileName"`   // 原始文件名
	MimeType   string     `json:"mimeType"`   // 文件MIME类型
	SizeBytes  int64      `json:"sizeBytes"`  // 文件大小(字节)
	Sha256     string     `json:"sha256"`     // 文件哈希(sha256，用于去重/校验)
	CreatedAt  *time.Time `json:"createdAt"`  // 创建时间
	UpdatedAt  *time.Time `json:"updatedAt"`  // 更新时间
}

// CreateLoanBaseinfoFilesReply only for api docs
type CreateLoanBaseinfoFilesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanBaseinfoFilesByIDReply only for api docs
type UpdateLoanBaseinfoFilesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanBaseinfoFilesByIDReply only for api docs
type GetLoanBaseinfoFilesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfoFiles LoanBaseinfoFilesObjDetail `json:"loanBaseinfoFiles"`
	} `json:"data"` // return data
}

// DeleteLoanBaseinfoFilesByIDReply only for api docs
type DeleteLoanBaseinfoFilesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanBaseinfoFilessByIDsReply only for api docs
type DeleteLoanBaseinfoFilessByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanBaseinfoFilessRequest request params
type ListLoanBaseinfoFilessRequest struct {
	query.Params
}

// ListLoanBaseinfoFilessReply only for api docs
type ListLoanBaseinfoFilessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfoFiless []LoanBaseinfoFilesObjDetail `json:"loanBaseinfoFiless"`
	} `json:"data"` // return data
}

// DeleteLoanBaseinfoFilessByIDsRequest request params
type DeleteLoanBaseinfoFilessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanBaseinfoFilesByConditionRequest request params
type GetLoanBaseinfoFilesByConditionRequest struct {
	query.Conditions
}

// GetLoanBaseinfoFilesByConditionReply only for api docs
type GetLoanBaseinfoFilesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfoFiles LoanBaseinfoFilesObjDetail `json:"loanBaseinfoFiles"`
	} `json:"data"` // return data
}

// ListLoanBaseinfoFilessByIDsRequest request params
type ListLoanBaseinfoFilessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanBaseinfoFilessByIDsReply only for api docs
type ListLoanBaseinfoFilessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanBaseinfoFiless []LoanBaseinfoFilesObjDetail `json:"loanBaseinfoFiless"`
	} `json:"data"` // return data
}
