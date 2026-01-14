package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

// LoanUserDeviceApps 设备软件列表采集表(匿名表单采集，与loan_baseinfo关联)
type LoanUserDeviceApps struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	BaseinfoID  int    `gorm:"column:baseinfo_id;type:int(11);not null" json:"baseinfoID"`        // 关联 loan_baseinfo.id
	PackageName string `gorm:"column:package_name;type:varchar(255);not null" json:"packageName"` // 应用包名/BundleId(如 com.xxx.app)
	AppName     string `gorm:"column:app_name;type:varchar(255)" json:"appName"`                  // 应用名称
	VersionName string `gorm:"column:version_name;type:varchar(64)" json:"versionName"`           // 版本名(如 1.2.3)
	VersionCode int64  `gorm:"column:version_code;type:bigint(20)" json:"versionCode"`            // 版本号(如 Android versionCode，可选)
}

// LoanUserDeviceAppsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanUserDeviceAppsColumnNames = map[string]bool{
	"id":           true,
	"created_at":   true,
	"updated_at":   true,
	"deleted_at":   true,
	"baseinfo_id":  true,
	"package_name": true,
	"app_name":     true,
	"version_name": true,
	"version_code": true,
}
