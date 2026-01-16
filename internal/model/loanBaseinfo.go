package model

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm"
)

type LoanBaseinfo struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	FirstName  string `gorm:"column:first_name;type:varchar(32)" json:"firstName"`   // 姓
	SecondName string `gorm:"column:second_name;type:varchar(32)" json:"secondName"` // 名
	Age        int    `gorm:"column:age;type:int(11)" json:"age"`                    // 年齡
	Gender     string `gorm:"column:gender;type:varchar(4)" json:"gender"`           // 性別
	IdType     string `gorm:"column:id_type;type:varchar(32)" json:"idType"`         // 證件類型
	IdNumber   string `gorm:"column:id_number;type:varchar(32)" json:"idNumber"`     // 證件號碼
	IdCard     string `gorm:"column:id_card;type:varchar(255)" json:"idCard"`        // 證件
	Operator   string `gorm:"column:operator;type:varchar(255)" json:"operator"`     // 操作系統
	Mobile     string `gorm:"column:mobile;type:varchar(32)" json:"mobile"`
	//Inviter                       string     `gorm:"column:inviter;type:varchar(255)" json:"inviter"`                                               // 邀請人
	Work          string `gorm:"column:work;type:varchar(255)" json:"work"`                  // 工作
	Company       string `gorm:"column:company;type:varchar(255)" json:"company"`            // 公司
	Salary        int    `gorm:"column:salary;type:int(11)" json:"salary"`                   // 薪資
	MaritalStatus int    `gorm:"column:marital_status;type:tinyint(4)" json:"maritalStatus"` // 婚否
	HasHouse      int    `gorm:"column:has_house;type:tinyint(4)" json:"hasHouse"`           // 是否有房
	//PropertyCertificate           string     `gorm:"column:property_certificate;type:varchar(255)" json:"propertyCertificate"`                      // 房產證
	HasCar int `gorm:"column:has_car;type:tinyint(4)" json:"hasCar"` // 是否有車
	//VehicleRgistrationCertificate string     `gorm:"column:vehicle_rgistration_certificate;type:varchar(255)" json:"vehicleRgistrationCertificate"` // 行駛證
	ApplicationAmount int        `gorm:"column:application_amount;type:int(11)" json:"applicationAmount"`                  // 申請金額
	AuditStatus       int        `gorm:"column:audit_status;type:tinyint(4);default:0" json:"auditStatus"`                 // 審核情況 0待審核 1審核通過 -1 審核拒絕
	BankNo            string     `gorm:"column:bank_no;type:varchar(255)" json:"bankNo"`                                   // 銀行卡號
	ClientIP          string     `gorm:"column:client_ip;type:varbinary(16)" json:"clientIP"`                              // 客户端IP地址(IPv4/IPv6)
	ReferrerUserID    int64      `gorm:"column:referrer_user_id;type:bigint(20)" json:"referrerUserID"`                    // 邀请人/分享人(loan_users.id)
	RefCode           string     `gorm:"column:ref_code;type:varchar(32)" json:"refCode"`                                  // 访问时携带的ref(冗余存储便于排查)
	LoanDays          int        `gorm:"column:loan_days;type:smallint(6);not null" json:"loanDays"`                       // 借款天数(单位：天)
	RiskListStatus    int        `gorm:"column:risk_list_status;type:tinyint(4);default:0;not null" json:"riskListStatus"` // 名单状态：0正常 1白名单 2黑名单
	RiskListReason    string     `gorm:"column:risk_list_reason;type:varchar(255)" json:"riskListReason"`                  // 名单原因/来源说明
	RiskListMarkedAt  *time.Time `gorm:"column:risk_list_marked_at;type:datetime" json:"riskListMarkedAt"`                 // 名单标记时间
}

// TableName table name
func (m *LoanBaseinfo) TableName() string {
	return "loan_baseinfo"
}

// LoanBaseinfoColumnNames Whitelist for custom query fields to prevent sql injection attacks
var LoanBaseinfoColumnNames = map[string]bool{
	"id":                  true,
	"created_at":          true,
	"updated_at":          true,
	"deleted_at":          true,
	"first_name":          true,
	"second_name":         true,
	"age":                 true,
	"gender":              true,
	"mobile":              true,
	"id_type":             true,
	"id_number":           true,
	"id_card":             true,
	"operator":            true,
	"work":                true,
	"company":             true,
	"salary":              true,
	"marital_status":      true,
	"has_house":           true,
	"has_car":             true,
	"application_amount":  true,
	"audit_status":        true,
	"bank_no":             true,
	"client_ip":           true,
	"referrer_user_id":    true,
	"ref_code":            true,
	"loan_days":           true,
	"risk_list_status":    true,
	"risk_list_reason":    true,
	"risk_list_marked_at": true,
}
