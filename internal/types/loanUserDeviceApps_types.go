package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLoanUserDeviceAppsRequest request params
type CreateLoanUserDeviceAppsRequest struct {
	BaseinfoID  int    `json:"baseinfoID" binding:""`  // 关联 loan_baseinfo.id
	PackageName string `json:"packageName" binding:""` // 应用包名/BundleId(如 com.xxx.app)
	AppName     string `json:"appName" binding:""`     // 应用名称
	VersionName string `json:"versionName" binding:""` // 版本名(如 1.2.3)
	VersionCode int64  `json:"versionCode" binding:""` // 版本号(如 Android versionCode，可选)
}

// UpdateLoanUserDeviceAppsByIDRequest request params
type UpdateLoanUserDeviceAppsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键
	BaseinfoID  int    `json:"baseinfoID" binding:""`  // 关联 loan_baseinfo.id
	PackageName string `json:"packageName" binding:""` // 应用包名/BundleId(如 com.xxx.app)
	AppName     string `json:"appName" binding:""`     // 应用名称
	VersionName string `json:"versionName" binding:""` // 版本名(如 1.2.3)
	VersionCode int64  `json:"versionCode" binding:""` // 版本号(如 Android versionCode，可选)
}

// LoanUserDeviceAppsObjDetail detail
type LoanUserDeviceAppsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键
	BaseinfoID  int        `json:"baseinfoID"`  // 关联 loan_baseinfo.id
	PackageName string     `json:"packageName"` // 应用包名/BundleId(如 com.xxx.app)
	AppName     string     `json:"appName"`     // 应用名称
	VersionName string     `json:"versionName"` // 版本名(如 1.2.3)
	VersionCode int64      `json:"versionCode"` // 版本号(如 Android versionCode，可选)
	CreatedAt   *time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   *time.Time `json:"updatedAt"`   // 更新时间
}

// CreateLoanUserDeviceAppsReply only for api docs
type CreateLoanUserDeviceAppsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLoanUserDeviceAppsByIDReply only for api docs
type UpdateLoanUserDeviceAppsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetLoanUserDeviceAppsByIDReply only for api docs
type GetLoanUserDeviceAppsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserDeviceApps LoanUserDeviceAppsObjDetail `json:"loanUserDeviceApps"`
	} `json:"data"` // return data
}

// DeleteLoanUserDeviceAppsByIDReply only for api docs
type DeleteLoanUserDeviceAppsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// DeleteLoanUserDeviceAppssByIDsReply only for api docs
type DeleteLoanUserDeviceAppssByIDsReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// ListLoanUserDeviceAppssRequest request params
type ListLoanUserDeviceAppssRequest struct {
	query.Params
}

// ListLoanUserDeviceAppssReply only for api docs
type ListLoanUserDeviceAppssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserDeviceAppss []LoanUserDeviceAppsObjDetail `json:"loanUserDeviceAppss"`
	} `json:"data"` // return data
}

// DeleteLoanUserDeviceAppssByIDsRequest request params
type DeleteLoanUserDeviceAppssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLoanUserDeviceAppsByConditionRequest request params
type GetLoanUserDeviceAppsByConditionRequest struct {
	query.Conditions
}

// GetLoanUserDeviceAppsByConditionReply only for api docs
type GetLoanUserDeviceAppsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserDeviceApps LoanUserDeviceAppsObjDetail `json:"loanUserDeviceApps"`
	} `json:"data"` // return data
}

// ListLoanUserDeviceAppssByIDsRequest request params
type ListLoanUserDeviceAppssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLoanUserDeviceAppssByIDsReply only for api docs
type ListLoanUserDeviceAppssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		LoanUserDeviceAppss []LoanUserDeviceAppsObjDetail `json:"loanUserDeviceAppss"`
	} `json:"data"` // return data
}
