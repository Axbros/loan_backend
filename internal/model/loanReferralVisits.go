package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

// LoanReferralVisits 邀请链接访问/点击记录表(匿名访问，用于统计点击与转化)
type LoanReferralVisits struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	VisitorID      string     `gorm:"column:visitor_id;type:char(36);not null" json:"visitorID"`            // 访客标识(UUID)，前端cookie生成)
	RefCode        string     `gorm:"column:ref_code;type:varchar(32);not null" json:"refCode"`             // 访问链接携带的ref(share_code)
	ReferrerUserID int64      `gorm:"column:referrer_user_id;type:bigint(20)" json:"referrerUserID"`        // 邀请人(loan_users.id)
	LandingPath    string     `gorm:"column:landing_path;type:varchar(255)" json:"landingPath"`             // 落地页路径
	ClientIP       string     `gorm:"column:client_ip;type:varbinary(16)" json:"clientIP"`                  // 访问IP(IPv4/IPv6)
	UserAgent      string     `gorm:"column:user_agent;type:varchar(255)" json:"userAgent"`                 // 浏览器UA
	FirstSeenAt    *time.Time `gorm:"column:first_seen_at;type:datetime;not null" json:"firstSeenAt"`       // 首次访问时间
	LastSeenAt     *time.Time `gorm:"column:last_seen_at;type:datetime;not null" json:"lastSeenAt"`         // 最近访问时间
	VisitCount     int        `gorm:"column:visit_count;type:int(11);default:1;not null" json:"visitCount"` // 访问次数
}

// LoanReferralVisitsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanReferralVisitsColumnNames = map[string]bool{
	"id":               true,
	"created_at":       true,
	"updated_at":       true,
	"deleted_at":       true,
	"visitor_id":       true,
	"ref_code":         true,
	"referrer_user_id": true,
	"landing_path":     true,
	"client_ip":        true,
	"user_agent":       true,
	"first_seen_at":    true,
	"last_seen_at":     true,
	"visit_count":      true,
}
