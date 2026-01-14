package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

// LoanBaseinfoFiles 基础信息附件表(匿名用户上传，按type区分证件/材料，存OSS地址)
type LoanBaseinfoFiles struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	BaseinfoID int    `gorm:"column:baseinfo_id;type:int(11);not null" json:"baseinfoID"` // 关联 loan_baseinfo.id
	Type       string `gorm:"column:type;type:varchar(64);not null" json:"type"`          // 文件类型(如 ID_CARD_FRONT / ID_CARD_BACK / TAX_CERT 等)
	OssURL     string `gorm:"column:oss_url;type:varchar(1024);not null" json:"ossURL"`   // OSS访问地址(或CDN地址)
	OssKey     string `gorm:"column:oss_key;type:varchar(512)" json:"ossKey"`             // OSS对象Key(内部定位/删除用，可选)
	FileName   string `gorm:"column:file_name;type:varchar(255)" json:"fileName"`         // 原始文件名
	MimeType   string `gorm:"column:mime_type;type:varchar(64)" json:"mimeType"`          // 文件MIME类型
	SizeBytes  int64  `gorm:"column:size_bytes;type:bigint(20)" json:"sizeBytes"`         // 文件大小(字节)
	Sha256     string `gorm:"column:sha256;type:char(64)" json:"sha256"`                  // 文件哈希(sha256，用于去重/校验)
}

// LoanBaseinfoFilesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanBaseinfoFilesColumnNames = map[string]bool{
	"id":          true,
	"created_at":  true,
	"updated_at":  true,
	"deleted_at":  true,
	"baseinfo_id": true,
	"type":        true,
	"oss_url":     true,
	"oss_key":     true,
	"file_name":   true,
	"mime_type":   true,
	"size_bytes":  true,
	"sha256":      true,
}
