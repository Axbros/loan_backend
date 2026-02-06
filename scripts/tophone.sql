/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 90500 (9.5.0)
 Source Host           : localhost:3306
 Source Schema         : tophone

 Target Server Type    : MySQL
 Target Server Version : 90500 (9.5.0)
 File Encoding         : 65001

 Date: 16/01/2026 18:33:02
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for loan_audits
-- ----------------------------
DROP TABLE IF EXISTS `loan_audits`;
CREATE TABLE `loan_audits` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '审核记录ID',
  `baseinfo_id` int NOT NULL COMMENT '关联申请单 loan_baseinfo.id',
  `audit_result` tinyint NOT NULL COMMENT '审核结果：1通过 -1拒绝',
  `audit_comment` varchar(255) DEFAULT NULL COMMENT '审核备注/原因',
  `auditor_user_id` bigint NOT NULL COMMENT '审核人员(loan_users.id)',
  `audit_type` varchar(255) DEFAULT NULL COMMENT '审核类型(初审、放款审核、回款审核)',
  `created_at` datetime NOT NULL COMMENT '审核时间(即审核通过/拒绝时间)',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_audit_baseinfo` (`baseinfo_id`) COMMENT '按申请单查询审核记录',
  KEY `idx_audit_user_time` (`auditor_user_id`,`created_at`) COMMENT '按审核人/时间查询',
  CONSTRAINT `fk_audit_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`),
  CONSTRAINT `fk_audit_user` FOREIGN KEY (`auditor_user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='申请审核记录表(审核时间即 created_at)';

-- ----------------------------
-- Records of loan_audits
-- ----------------------------
BEGIN;
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `created_at`, `deleted_at`) VALUES (1, 1, 1, '初审通过，资料齐全', 2, NULL, '2026-01-12 19:52:29', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `created_at`, `deleted_at`) VALUES (2, 1, 1, '复审确认收入稳定，风险可控', 2, NULL, '2026-01-13 19:52:37', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_baseinfo
-- ----------------------------
DROP TABLE IF EXISTS `loan_baseinfo`;
CREATE TABLE `loan_baseinfo` (
  `id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(32) DEFAULT NULL COMMENT '姓',
  `second_name` varchar(32) DEFAULT NULL COMMENT '名',
  `age` int DEFAULT NULL COMMENT '年齡',
  `gender` varchar(4) DEFAULT NULL COMMENT '性別',
  `mobile` varchar(18) NOT NULL COMMENT '手机号码',
  `id_type` varchar(32) DEFAULT NULL COMMENT '證件類型',
  `id_number` varchar(32) DEFAULT NULL COMMENT '證件號碼',
  `id_card` varchar(255) DEFAULT NULL COMMENT '證件',
  `operator` varchar(255) DEFAULT NULL COMMENT '操作系統',
  `work` varchar(255) DEFAULT NULL COMMENT '工作',
  `company` varchar(255) DEFAULT NULL COMMENT '公司',
  `salary` int DEFAULT NULL COMMENT '薪資',
  `marital_status` tinyint DEFAULT NULL COMMENT '婚否',
  `has_house` tinyint DEFAULT NULL COMMENT '是否有房',
  `has_car` tinyint DEFAULT NULL COMMENT '是否有車',
  `application_amount` int DEFAULT NULL COMMENT '申請金額',
  `audit_status` tinyint DEFAULT '0' COMMENT '審核情況 0待審核 1審核通過 -1 審核拒絕',
  `bank_no` varchar(255) DEFAULT NULL COMMENT '銀行卡號',
  `client_ip` varbinary(16) DEFAULT NULL COMMENT '客户端IP地址(IPv4/IPv6)',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `referrer_user_id` bigint DEFAULT NULL COMMENT '邀请人/分享人(loan_users.id)',
  `ref_code` varchar(32) DEFAULT NULL COMMENT '访问时携带的ref(冗余存储便于排查)',
  `loan_days` smallint NOT NULL COMMENT '借款天数(单位：天)',
  PRIMARY KEY (`id`),
  KEY `idx_baseinfo_referrer_user` (`referrer_user_id`) COMMENT '按邀请人查询申请记录',
  KEY `idx_baseinfo_ref_code` (`ref_code`) COMMENT '按ref查询',
  CONSTRAINT `fk_baseinfo_referrer_user` FOREIGN KEY (`referrer_user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=87 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_baseinfo
-- ----------------------------
BEGIN;
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (1, 'Wang', 'Lei', 29, 'M', '16600229988', 'ID_CARD', '110101199401011234', 'oss://idcard/wanglei_front.jpg', 'Android', 'Engineer', 'TechSoft Ltd.', 15000, 0, 0, 0, 100000, 0, '6222020200001234567', 0x3139322E3136382E332E31, '2026-01-14 19:12:30', '2026-01-14 19:12:30', '2026-01-15 13:20:21', NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (2, 'Li', 'Na', 34, 'F', '16600229988', 'ID_CARD', '310101198912123456', NULL, 'iOS', 'Sales', 'TradeCorp', 12000, NULL, NULL, NULL, 80000, 0, '6222020200007654321', 0x3139322E3136382E332E31, '2026-01-14 19:12:36', '2026-01-14 19:12:36', '2026-01-15 13:38:22', 1, 'REFAXBROS01', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (3, 'Li', 'Na', 34, 'F', '16600229988', 'ID_CARD', '310101198912123456', NULL, 'iOS', 'Sales', 'TradeCorp', 12000, NULL, NULL, NULL, 80000, 0, '6222020200007654321', 0x3139322E3136382E332E31, '2026-01-14 19:12:41', '2026-01-14 19:12:41', NULL, 1, 'REFAXBROS01', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (4, 'Chen', 'Yu', 27, 'F', '16600229988', 'ID_CARD', '320101199612124321', NULL, 'Android', 'Designer', 'Creative Studio', 18000, NULL, NULL, NULL, 50000, 1, '6228480402567890123', 0x3139322E3136382E332E31, '2026-01-14 19:12:48', '2026-01-14 19:12:48', NULL, NULL, NULL, 14, 1, '历史还款良好，人工加入白名单', '2026-01-14 19:12:48');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (5, 'Liu', 'Ming', 36, 'M', '16600229988', 'ID_CARD', '510101198801019999', NULL, 'Web', 'Freelancer', NULL, NULL, NULL, NULL, NULL, 30000, -1, NULL, 0x3139322E3136382E332E31, '2026-01-14 19:12:55', '2026-01-14 19:12:55', NULL, NULL, NULL, 7, 2, '命中内部黑名单：多次逾期', '2026-01-14 19:12:55');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (6, 'Test', 'A', 28, 'M', '16600229988', 'ID_CARD', 'TID0001', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 100000, 1, '6222000000000001', 0x0A000001, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (7, 'Test', 'B', 30, 'F', '16600229988', 'ID_CARD', 'TID0002', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 120000, 1, '6222000000000002', 0x0A000002, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (8, 'Test', 'C', 26, 'M', '16600229988', 'ID_CARD', 'TID0003', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 90000, 1, '6222000000000003', 0x0A000003, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (9, 'Test', 'D', 35, 'F', '16600229988', 'ID_CARD', 'TID0004', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 150000, 1, '6222000000000004', 0x0A000004, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 60, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (10, 'Test', 'E', 41, 'M', '16600229988', 'ID_CARD', 'TID0005', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 80000, 1, '6222000000000005', 0x0A000005, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (11, 'Test', 'F', 33, 'F', '16600229988', 'ID_CARD', 'TID0006', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 110000, 1, '6222000000000006', 0x0A000006, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (12, 'Test', 'G', 29, 'M', '16600229988', 'ID_CARD', 'TID0007', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 200000, 1, '6222000000000007', 0x0A000007, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (13, 'Test', 'H', 27, 'F', '16600229988', 'ID_CARD', 'TID0008', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 70000, 1, '6222000000000008', 0x0A000008, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (14, 'Test', 'I', 38, 'M', '16600229988', 'ID_CARD', 'TID0009', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 180000, 1, '6222000000000009', 0x0A000009, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 45, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (15, 'Test', 'J', 24, 'F', '16600229988', 'ID_CARD', 'TID0010', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 60000, 1, '6222000000000010', 0x0A00000A, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (66, '王', '磊', 29, 'M', '16600229988', 'ID_CARD', '110101199401011234', 'ID110101199401011234.jpg', 'iOS 17.0', '软件工程师', '北京科技有限公司', 25000, 1, 1, 1, 100000, 0, '6222081001001234567', 0xC0A80165, '2026-01-01 10:00:00', '2026-01-01 10:00:00', NULL, 1, 'REF00001', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (67, '王', '磊', 29, 'M', '16600229988', 'ID_CARD', '110101199401011234', 'ID110101199401011234.jpg', 'iOS 17.0', '软件工程师', '北京科技有限公司', 25000, 1, 1, 1, 100000, 0, '6222081001001234567', 0xC0A80165, '2026-01-01 10:00:00', '2026-01-01 10:00:00', NULL, 1, 'REF00001', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (68, '李', '娜', 34, 'F', '16600229988', 'ID_CARD', '310101198912123456', 'ID310101198912123456.jpg', 'Android 14', '财务经理', '上海金融有限公司', 35000, 1, 1, 1, 80000, 1, '6228480402567890123', 0xC0A80166, '2026-01-02 11:00:00', '2026-01-02 11:00:00', NULL, 2, 'REF00002', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (69, '陈', '宇', 27, 'F', '16600229988', 'ID_CARD', '320101199612124321', 'ID320101199612124321.jpg', 'iOS 16.5', '人力资源专员', '江苏贸易有限公司', 18000, 0, 0, 0, 50000, 0, '6259991234567890123', 0xC0A80167, '2026-01-03 09:30:00', '2026-01-03 09:30:00', NULL, 3, 'REF00003', 14, 1, '优质客户', '2026-01-03 10:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (70, '刘', '明', 36, 'M', '16600229988', 'ID_CARD', '510101198801019999', 'ID510101198801019999.jpg', 'Windows 11', '销售总监', '四川科技有限公司', 45000, 1, 1, 1, 30000, -1, '6226667890123456789', 0xC0A80168, '2026-01-04 14:00:00', '2026-01-04 15:00:00', NULL, 1, 'REF00004', 7, 2, '信用逾期', '2026-01-04 16:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (71, '张', '伟', 41, 'M', '16600229988', 'ID_CARD', '440101198305051234', 'ID440101198305051234.jpg', 'macOS 14', '项目经理', '广东建设有限公司', 50000, 1, 1, 1, 150000, 1, '6229998765432109876', 0xC0A80169, '2026-01-05 08:00:00', '2026-01-05 08:30:00', NULL, 2, 'REF00005', 21, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (72, '赵', '丽', 25, 'F', '16600229988', 'ID_CARD', '210101199808085678', 'ID210101199808085678.jpg', 'Android 13', '行政助理', '辽宁商贸有限公司', 15000, 0, 0, 0, 20000, 0, '6227779876543210987', 0xC0A8016A, '2026-01-06 16:00:00', '2026-01-06 16:00:00', NULL, 3, 'REF00006', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (73, '黄', '浩', 33, 'M', '16600229988', 'ID_CARD', '430101199011118765', 'ID430101199011118765.jpg', 'iOS 17.1', '产品经理', '湖南科技有限公司', 30000, 1, 0, 1, 70000, 0, '6225551234567890123', 0xC0A8016B, '2026-01-07 10:30:00', '2026-01-07 10:30:00', NULL, 1, 'REF00007', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (74, '周', '敏', 28, 'F', '16600229988', 'ID_CARD', '330101199503037654', 'ID330101199503037654.jpg', 'Windows 10', '客服专员', '浙江服务有限公司', 16000, 0, 0, 0, 15000, -1, '6224448765432109876', 0xC0A8016C, '2026-01-08 13:00:00', '2026-01-08 14:00:00', NULL, 2, 'REF00008', 14, 2, '收入不稳定', '2026-01-08 15:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (75, '吴', '强', 38, 'M', '16600229988', 'ID_CARD', '350101198507076543', 'ID350101198507076543.jpg', 'Android 12', '工程师', '福建制造有限公司', 40000, 1, 1, 1, 90000, 1, '6223337654321098765', 0xC0A8016D, '2026-01-09 09:00:00', '2026-01-09 09:15:00', NULL, 3, 'REF00009', 21, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (76, '徐', '芳', 31, 'F', '16600229988', 'ID_CARD', '610101199209094321', 'ID610101199209094321.jpg', 'macOS 13', '设计师', '陕西创意有限公司', 28000, 1, 0, 0, 60000, 0, '6222226543210987654', 0xC0A8016E, '2026-01-10 11:30:00', '2026-01-10 11:30:00', NULL, 1, 'REF00010', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (77, '孙', '杰', 26, 'M', '16600229988', 'ID_CARD', '120101199712123456', 'ID120101199712123456.jpg', 'iOS 16.4', '程序员', '天津科技有限公司', 22000, 0, 0, 1, 40000, 0, '6221115432109876543', 0xC0A8016F, '2026-01-11 15:00:00', '2026-01-11 15:00:00', NULL, 2, 'REF00011', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (78, '马', '丽', 37, 'F', '16600229988', 'ID_CARD', '650101198604048765', 'ID650101198604048765.jpg', 'Windows 11', '教师', '新疆教育有限公司', 20000, 1, 1, 0, 35000, -1, '6220004321098765432', 0xC0A80170, '2026-01-12 08:30:00', '2026-01-12 09:00:00', NULL, 3, 'REF00012', 30, 2, '负债过高', '2026-01-12 10:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (79, '朱', '军', 40, 'M', '16600229988', 'ID_CARD', '500101198310107654', 'ID500101198310107654.jpg', 'Android 14', '医生', '重庆医疗有限公司', 55000, 1, 1, 1, 120000, 1, '6219993210987654321', 0xC0A80171, '2026-01-13 14:30:00', '2026-01-13 14:45:00', NULL, 1, 'REF00013', 21, 1, '优质客户', '2026-01-13 15:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (80, '胡', '欣', 24, 'F', '16600229988', 'ID_CARD', '460101199906066543', 'ID460101199906066543.jpg', 'iOS 17.0', '实习生', '海南旅游有限公司', 8000, 0, 0, 0, 10000, 0, '6218882109876543210', 0xC0A80172, '2026-01-14 10:00:00', '2026-01-14 10:00:00', NULL, 2, 'REF00014', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (81, '林', '涛', 32, 'M', '16600229988', 'ID_CARD', '360101199102025432', 'ID360101199102025432.jpg', 'macOS 14', '摄影师', '江西传媒有限公司', 26000, 0, 0, 1, 50000, 0, '6217771098765432109', 0xC0A80173, '2026-01-15 16:30:00', '2026-01-15 16:30:00', NULL, 3, 'REF00015', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (82, '郭', '燕', 30, 'F', '16600229988', 'ID_CARD', '340101199307074321', 'ID340101199307074321.jpg', 'Windows 10', '护士', '安徽医疗有限公司', 24000, 1, 0, 0, 45000, 1, '6216660987654321098', 0xC0A80174, '2026-01-16 09:30:00', '2026-01-16 10:00:00', NULL, 1, 'REF00016', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (83, '何', '勇', 39, 'M', '16600229988', 'ID_CARD', '540101198408083210', 'ID540101198408083210.jpg', 'Android 13', '建筑工人', '西藏建设有限公司', 32000, 1, 1, 1, 80000, -1, '6215559876543210987', 0xC0A80175, '2026-01-17 13:30:00', '2026-01-17 14:00:00', NULL, 2, 'REF00017', 21, 2, '征信不良', '2026-01-17 15:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (84, '高', '静', 27, 'F', '16600229988', 'ID_CARD', '630101199609092109', 'ID630101199609092109.jpg', 'iOS 16.5', '翻译', '青海外贸有限公司', 21000, 0, 0, 0, 25000, 0, '6214448765432109876', 0xC0A80176, '2026-01-18 11:00:00', '2026-01-18 11:00:00', NULL, 3, 'REF00018', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (85, '罗', '刚', 35, 'M', '16600229988', 'ID_CARD', '530101198811111098', 'ID530101198811111098.jpg', 'Windows 11', '厨师', '云南餐饮有限公司', 28000, 1, 0, 1, 60000, 0, '6213337654321098765', 0xC0A80177, '2026-01-19 15:30:00', '2026-01-19 15:30:00', NULL, 1, 'REF00019', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (86, '郑', '艳', 29, 'F', '16600229988', 'ID_CARD', '410101199412120987', 'ID410101199412120987.jpg', 'Android 12', '导购', '河南零售有限公司', 18000, 0, 0, 0, 30000, 1, '6212226543210987654', 0xC0A80178, '2026-01-20 08:00:00', '2026-01-20 08:15:00', NULL, 2, 'REF00020', 30, 0, NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_baseinfo_files
-- ----------------------------
DROP TABLE IF EXISTS `loan_baseinfo_files`;
CREATE TABLE `loan_baseinfo_files` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `baseinfo_id` int NOT NULL COMMENT '关联 loan_baseinfo.id',
  `type` varchar(64) NOT NULL COMMENT '文件类型(如 ID_CARD_FRONT / ID_CARD_BACK / TAX_CERT 等)',
  `oss_url` varchar(1024) NOT NULL COMMENT 'OSS访问地址(或CDN地址)',
  `oss_key` varchar(512) DEFAULT NULL COMMENT 'OSS对象Key(内部定位/删除用，可选)',
  `file_name` varchar(255) DEFAULT NULL COMMENT '原始文件名',
  `mime_type` varchar(64) DEFAULT NULL COMMENT '文件MIME类型',
  `size_bytes` bigint DEFAULT NULL COMMENT '文件大小(字节)',
  `sha256` char(64) DEFAULT NULL COMMENT '文件哈希(sha256，用于去重/校验)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_baseinfo` (`baseinfo_id`) COMMENT '按基础信息查询附件',
  KEY `idx_type` (`type`) COMMENT '按类型查询附件',
  CONSTRAINT `fk_baseinfo_files_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='基础信息附件表(匿名用户上传，按type区分证件/材料，存OSS地址)';

-- ----------------------------
-- Records of loan_baseinfo_files
-- ----------------------------
BEGIN;
INSERT INTO `loan_baseinfo_files` (`id`, `baseinfo_id`, `type`, `oss_url`, `oss_key`, `file_name`, `mime_type`, `size_bytes`, `sha256`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 1, 'ID_CARD_FRONT', 'https://www.tophone.cc/images/1.png', 'loan/idcard/1_front.jpg', 'id_front.jpg', 'image/jpeg', 245678, 'a3b1c4d5e6f7890a3b1c4d5e6f7890a3b1c4d5e6f7890a3b1c4d5e6f7890', '2026-01-14 19:24:36', NULL, NULL);
INSERT INTO `loan_baseinfo_files` (`id`, `baseinfo_id`, `type`, `oss_url`, `oss_key`, `file_name`, `mime_type`, `size_bytes`, `sha256`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 1, 'ID_CARD_BACK', 'https://www.tophone.cc/images/2.png', 'loan/idcard/1_back.jpg', 'id_back.jpg', 'image/jpeg', 231456, 'b4c2d5e6f7890a3b1c4d5e6f7890a3b1c4d5e6f7890a3b1c4d5e6f7890a', '2026-01-14 19:24:36', NULL, NULL);
INSERT INTO `loan_baseinfo_files` (`id`, `baseinfo_id`, `type`, `oss_url`, `oss_key`, `file_name`, `mime_type`, `size_bytes`, `sha256`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 1, 'TAX_CERT', 'https://www.tophone.cc/images/3.png', 'loan/tax/1_2023.pdf', 'tax_2023.pdf', 'application/pdf', 1048576, NULL, '2026-01-14 19:24:43', NULL, NULL);
INSERT INTO `loan_baseinfo_files` (`id`, `baseinfo_id`, `type`, `oss_url`, `oss_key`, `file_name`, `mime_type`, `size_bytes`, `sha256`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 2, 'AVATAR', 'https://www.tophone.cc/images/4.png', NULL, 'bank_card.jpg', 'image/jpeg', NULL, NULL, '2026-01-14 19:24:49', NULL, NULL);
INSERT INTO `loan_baseinfo_files` (`id`, `baseinfo_id`, `type`, `oss_url`, `oss_key`, `file_name`, `mime_type`, `size_bytes`, `sha256`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 3, 'AVATAR', 'https://www.tophone.cc/images/5.png', NULL, NULL, NULL, NULL, NULL, '2026-01-14 19:24:56', NULL, NULL);
INSERT INTO `loan_baseinfo_files` (`id`, `baseinfo_id`, `type`, `oss_url`, `oss_key`, `file_name`, `mime_type`, `size_bytes`, `sha256`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 3, 'OTHER', 'https://www.tophone.cc/images/6.png', NULL, NULL, NULL, NULL, NULL, '2026-01-14 19:24:56', NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_collection_cases
-- ----------------------------
DROP TABLE IF EXISTS `loan_collection_cases`;
CREATE TABLE `loan_collection_cases` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '催收任务ID',
  `disbursement_id` bigint NOT NULL COMMENT '关联放款单 loan_disbursements.id',
  `schedule_id` bigint DEFAULT NULL COMMENT '关联逾期期次 loan_repayment_schedules.id(按期催收可用，整单催收可为空)',
  `collector_user_id` bigint NOT NULL COMMENT '催收人员 loan_users.id',
  `assigned_by_user_id` bigint NOT NULL COMMENT '分配人(管理员) loan_users.id',
  `assigned_at` datetime NOT NULL COMMENT '分配时间',
  `priority` tinyint NOT NULL DEFAULT '2' COMMENT '优先级：1高 2中 3低',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '任务状态：0待处理 1跟进中 2已完成 3已取消',
  `due_amount` int DEFAULT NULL COMMENT '逾期应还金额快照(分，可选，用于列表展示)',
  `overdue_days` int DEFAULT NULL COMMENT '逾期天数快照(可选，用于列表展示)',
  `completed_at` datetime DEFAULT NULL COMMENT '完成时间(点击完成时)',
  `completed_note` varchar(255) DEFAULT NULL COMMENT '完成备注(例如用户承诺X天内还款)',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_case_schedule_active` (`schedule_id`,`status`) COMMENT '避免重复任务(建议结合业务使用)',
  KEY `idx_case_collector_status` (`collector_user_id`,`status`) COMMENT '催收人员的任务列表(按状态)',
  KEY `idx_case_disbursement` (`disbursement_id`) COMMENT '按放款单查询任务',
  KEY `idx_case_schedule` (`schedule_id`) COMMENT '按期次查询任务',
  KEY `idx_case_assigned_at` (`assigned_at`) COMMENT '按分配时间查询',
  KEY `fk_case_assigned_by` (`assigned_by_user_id`),
  CONSTRAINT `fk_case_assigned_by` FOREIGN KEY (`assigned_by_user_id`) REFERENCES `loan_users` (`id`),
  CONSTRAINT `fk_case_collector` FOREIGN KEY (`collector_user_id`) REFERENCES `loan_users` (`id`),
  CONSTRAINT `fk_case_disbursement` FOREIGN KEY (`disbursement_id`) REFERENCES `loan_disbursements` (`id`),
  CONSTRAINT `fk_case_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `loan_repayment_schedules` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='催收任务表(管理员批量分配逾期任务给催收人员，催收人员完成并备注)';

-- ----------------------------
-- Records of loan_collection_cases
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_collection_logs
-- ----------------------------
DROP TABLE IF EXISTS `loan_collection_logs`;
CREATE TABLE `loan_collection_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '跟进记录ID',
  `case_id` bigint NOT NULL COMMENT '关联催收任务 loan_collection_cases.id',
  `collector_user_id` bigint NOT NULL COMMENT '催收人员 loan_users.id',
  `action_type` varchar(32) DEFAULT NULL COMMENT '动作类型(如 CALL/SMS/VISIT/OTHER，可选)',
  `content` varchar(500) NOT NULL COMMENT '跟进内容/备注(例如用户承诺3天内还款)',
  `next_follow_up_at` datetime DEFAULT NULL COMMENT '下次跟进时间(可选)',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_log_case_time` (`case_id`,`created_at`) COMMENT '按任务查询跟进记录',
  KEY `idx_log_collector_time` (`collector_user_id`,`created_at`) COMMENT '按催收人员查询跟进记录',
  CONSTRAINT `fk_log_case` FOREIGN KEY (`case_id`) REFERENCES `loan_collection_cases` (`id`),
  CONSTRAINT `fk_log_collector` FOREIGN KEY (`collector_user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='催收跟进记录表(一条任务可多次记录沟通内容/承诺/计划)';

-- ----------------------------
-- Records of loan_collection_logs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_department_roles
-- ----------------------------
DROP TABLE IF EXISTS `loan_department_roles`;
CREATE TABLE `loan_department_roles` (
  `department_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate id',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`department_id`,`role_id`),
  UNIQUE KEY `uk_id` (`id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `loan_department_roles_ibfk_1` FOREIGN KEY (`department_id`) REFERENCES `loan_departments` (`id`),
  CONSTRAINT `loan_department_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `loan_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_department_roles
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_departments
-- ----------------------------
DROP TABLE IF EXISTS `loan_departments`;
CREATE TABLE `loan_departments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `parent_id` bigint DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_departments
-- ----------------------------
BEGIN;
INSERT INTO `loan_departments` (`id`, `name`, `parent_id`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '管理员', NULL, 1, '2026-01-16 11:40:28', '2026-01-16 11:40:33', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_disbursements
-- ----------------------------
DROP TABLE IF EXISTS `loan_disbursements`;
CREATE TABLE `loan_disbursements` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键(放款单ID)',
  `baseinfo_id` int NOT NULL COMMENT '关联申请单 loan_baseinfo.id',
  `disburse_amount` int NOT NULL COMMENT '放款金额(单位按你的系统：元/分，建议统一)',
  `net_amount` int NOT NULL COMMENT '到账金额(扣除费用后实际到账)',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '放款状态：0待放款 1已放款',
  `source_referrer_user_id` bigint DEFAULT NULL COMMENT '用户来源(分享人 loan_users.id，冗余快照，便于查询)',
  `auditor_user_id` bigint DEFAULT NULL COMMENT '审核人员(loan_users.id)',
  `audited_at` datetime DEFAULT NULL COMMENT '审核通过时间',
  `payout_channel_id` bigint DEFAULT NULL COMMENT '放款渠道(代付) loan_payment_channels.id',
  `payout_order_no` varchar(128) DEFAULT NULL COMMENT '放款订单号/三方代付单号',
  `disbursed_at` datetime DEFAULT NULL COMMENT '放款时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间(进入待放款时刻)',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_baseinfo_disbursement` (`baseinfo_id`) COMMENT '一个申请单只生成一个放款单(如允许多次放款可去掉)',
  KEY `idx_status` (`status`) COMMENT '按放款状态筛选',
  KEY `idx_auditor_time` (`auditor_user_id`,`audited_at`) COMMENT '按审核人/审核时间筛选',
  KEY `idx_channel_time` (`payout_channel_id`,`disbursed_at`) COMMENT '按渠道/放款时间筛选',
  KEY `idx_source_referrer` (`source_referrer_user_id`) COMMENT '按来源分享人筛选',
  KEY `idx_disburse_payout_channel_time` (`payout_channel_id`,`disbursed_at`) COMMENT '按放款渠道/放款时间查询',
  KEY `idx_disburse_payout_order_no` (`payout_order_no`) COMMENT '按放款订单号查询',
  CONSTRAINT `fk_disburse_auditor` FOREIGN KEY (`auditor_user_id`) REFERENCES `loan_users` (`id`),
  CONSTRAINT `fk_disburse_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`),
  CONSTRAINT `fk_disburse_payout_channel` FOREIGN KEY (`payout_channel_id`) REFERENCES `loan_payment_channels` (`id`),
  CONSTRAINT `fk_disburse_referrer` FOREIGN KEY (`source_referrer_user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='放款单/待放款任务表(审核通过后生成，状态待放款->已放款)';

-- ----------------------------
-- Records of loan_disbursements
-- ----------------------------
BEGIN;
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (41, 1, 100000, 98000, 0, 3, 2, '2026-01-14 19:39:55', 1, NULL, NULL, '2026-01-14 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (42, 2, 200000, 196000, 1, 3, 2, '2026-01-12 19:39:55', 1, 'PO202401010001', '2026-01-13 19:39:55', '2026-01-12 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (43, 3, 150000, 147000, 1, NULL, 2, '2026-01-11 19:39:55', 2, 'PO202401010002', '2026-01-12 19:39:55', '2026-01-11 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (44, 4, 80000, 78000, 0, NULL, 2, '2026-01-14 19:39:55', 1, NULL, NULL, '2026-01-14 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (45, 5, 50000, 49000, 1, 3, 1, '2026-01-09 19:39:55', 2, 'PO202401010003', '2026-01-10 19:39:55', '2026-01-09 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (46, 6, 120000, 118000, 1, 3, 1, '2026-01-08 19:39:55', 1, 'PO202401010004', '2026-01-09 19:39:55', '2026-01-08 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (47, 7, 300000, 294000, 0, NULL, 2, '2026-01-14 19:39:55', 2, NULL, NULL, '2026-01-14 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (48, 8, 180000, 176000, 1, 3, 1, '2026-01-04 19:39:55', 1, 'PO202401010005', '2026-01-05 19:39:55', '2026-01-04 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (49, 9, 90000, 88200, 1, NULL, 2, '2026-01-13 19:39:55', 2, 'PO202401010006', '2026-01-14 07:39:55', '2026-01-13 19:39:55', '2026-01-14 19:39:55', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (50, 10, 110000, 107800, 0, 3, 2, '2026-01-14 19:39:55', 1, NULL, NULL, '2026-01-14 19:39:55', '2026-01-14 19:39:55', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_login_audit
-- ----------------------------
DROP TABLE IF EXISTS `loan_login_audit`;
CREATE TABLE `loan_login_audit` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `login_type` varchar(16) NOT NULL,
  `ip` varchar(64) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `success` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_login_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_login_audit
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_mfa_devices
-- ----------------------------
DROP TABLE IF EXISTS `loan_mfa_devices`;
CREATE TABLE `loan_mfa_devices` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `type` varchar(16) NOT NULL,
  `name` varchar(64) NOT NULL,
  `secret_enc` varbinary(255) DEFAULT NULL,
  `is_primary` tinyint NOT NULL DEFAULT '1',
  `status` tinyint NOT NULL DEFAULT '0',
  `last_used_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_mfa_user` (`user_id`),
  CONSTRAINT `loan_mfa_devices_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_mfa_devices
-- ----------------------------
BEGIN;
INSERT INTO `loan_mfa_devices` (`id`, `user_id`, `type`, `name`, `secret_enc`, `is_primary`, `status`, `last_used_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 1, 'TOTP', 'Google Authenticator', 0xA5FA8ADFDBF1B37BC96D880DB015C37B17C350D955D95F9368D5341DE9E1DD22F38B7916614B572CA83EEF39A8088AB59515DF1AE5D26B4644F95440, 1, 1, '2026-01-16 14:22:36', '2026-01-16 14:04:44', '2026-01-16 14:22:36', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_payment_channels
-- ----------------------------
DROP TABLE IF EXISTS `loan_payment_channels`;
CREATE TABLE `loan_payment_channels` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` varchar(64) NOT NULL COMMENT '渠道编码(唯一，如 BANK_A、WALLET_X)',
  `name` varchar(128) NOT NULL COMMENT '渠道名称',
  `merchant_no` varchar(128) NOT NULL COMMENT '商户号/商户ID(该渠道分配给平台的商户标识)',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '渠道状态：1启用 0禁用',
  `can_payout` tinyint NOT NULL DEFAULT '1' COMMENT '是否支持代付/放款：1是 0否',
  `can_collect` tinyint NOT NULL DEFAULT '1' COMMENT '是否支持代收/回款：1是 0否',
  `payout_fee_rate` decimal(10,6) DEFAULT NULL COMMENT '代付手续费率(如0.003500=0.35%)',
  `payout_fee_fixed` int DEFAULT NULL COMMENT '代付固定手续费(分，若不用可为空)',
  `collect_fee_rate` decimal(10,6) DEFAULT NULL COMMENT '代收手续费率',
  `collect_fee_fixed` int DEFAULT NULL COMMENT '代收固定手续费(分)',
  `collect_min_amount` int DEFAULT NULL COMMENT '最小代收金额(分)',
  `collect_max_amount` int DEFAULT NULL COMMENT '最大代收金额(分)',
  `payout_min_amount` int DEFAULT NULL COMMENT '最小代付金额(分)',
  `payout_max_amount` int DEFAULT NULL COMMENT '最大代付金额(分)',
  `settlement_cycle` varchar(32) DEFAULT NULL COMMENT '结算周期(如 T0/T1/D1/W1/M1，可按你们渠道定义)',
  `settlement_desc` varchar(255) DEFAULT NULL COMMENT '结算说明/备注',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_payment_channel_code` (`code`) COMMENT '渠道编码唯一',
  KEY `idx_payment_channel_status` (`status`) COMMENT '按状态筛选'
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='统一支付渠道配置表(支持代付放款+代收回款，可禁用/启用，含手续费/限额/结算周期)';

-- ----------------------------
-- Records of loan_payment_channels
-- ----------------------------
BEGIN;
INSERT INTO `loan_payment_channels` (`id`, `code`, `name`, `merchant_no`, `status`, `can_payout`, `can_collect`, `payout_fee_rate`, `payout_fee_fixed`, `collect_fee_rate`, `collect_fee_fixed`, `collect_min_amount`, `collect_max_amount`, `payout_min_amount`, `payout_max_amount`, `settlement_cycle`, `settlement_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'BANK_A', '银行A代付代收', '2020202020202', 1, 1, 1, 0.003500, 0, 0.002000, 0, 1000, 20000000, 1000, 50000000, 'T1', '默认T+1结算', '2026-01-14 19:38:20', '2026-01-14 19:38:20', NULL);
INSERT INTO `loan_payment_channels` (`id`, `code`, `name`, `merchant_no`, `status`, `can_payout`, `can_collect`, `payout_fee_rate`, `payout_fee_fixed`, `collect_fee_rate`, `collect_fee_fixed`, `collect_min_amount`, `collect_max_amount`, `payout_min_amount`, `payout_max_amount`, `settlement_cycle`, `settlement_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 'WALLET_X', '钱包X代付代收', '2020202020202', 1, 1, 1, 0.004000, 100, 0.002500, 50, 1000, 10000000, 1000, 30000000, 'T0', '默认T+0结算', '2026-01-14 19:38:20', '2026-01-14 19:38:20', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_permissions
-- ----------------------------
DROP TABLE IF EXISTS `loan_permissions`;
CREATE TABLE `loan_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(128) NOT NULL,
  `name` varchar(128) NOT NULL,
  `type` varchar(16) DEFAULT NULL,
  `resource` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_permissions
-- ----------------------------
BEGIN;
INSERT INTO `loan_permissions` (`id`, `code`, `name`, `type`, `resource`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'customer:view', '客户查看', NULL, NULL, '2026-01-14 18:21:50', '2026-01-14 18:21:50', NULL);
INSERT INTO `loan_permissions` (`id`, `code`, `name`, `type`, `resource`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 'customer:edit', '客户编辑', NULL, NULL, '2026-01-14 18:21:50', '2026-01-14 18:21:50', NULL);
INSERT INTO `loan_permissions` (`id`, `code`, `name`, `type`, `resource`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 'loan:approve', '贷款审核', NULL, NULL, '2026-01-14 18:21:50', '2026-01-14 18:21:50', NULL);
INSERT INTO `loan_permissions` (`id`, `code`, `name`, `type`, `resource`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 'loan:disburse', '贷款放款', NULL, NULL, '2026-01-14 18:21:50', '2026-01-14 18:21:50', NULL);
INSERT INTO `loan_permissions` (`id`, `code`, `name`, `type`, `resource`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 'repay:view', '还款查看', NULL, NULL, '2026-01-14 18:21:50', '2026-01-14 18:21:50', NULL);
INSERT INTO `loan_permissions` (`id`, `code`, `name`, `type`, `resource`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 'repay:collect', '还款催收', NULL, NULL, '2026-01-14 18:21:50', '2026-01-14 18:21:50', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_referral_visits
-- ----------------------------
DROP TABLE IF EXISTS `loan_referral_visits`;
CREATE TABLE `loan_referral_visits` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `visitor_id` char(36) NOT NULL COMMENT '访客标识(UUID)，前端cookie生成)',
  `ref_code` varchar(32) NOT NULL COMMENT '访问链接携带的ref(share_code)',
  `referrer_user_id` bigint DEFAULT NULL COMMENT '邀请人(loan_users.id)',
  `landing_path` varchar(255) DEFAULT NULL COMMENT '落地页路径',
  `client_ip` varbinary(16) DEFAULT NULL COMMENT '访问IP(IPv4/IPv6)',
  `user_agent` varchar(255) DEFAULT NULL COMMENT '浏览器UA',
  `first_seen_at` datetime NOT NULL COMMENT '首次访问时间',
  `last_seen_at` datetime NOT NULL COMMENT '最近访问时间',
  `visit_count` int NOT NULL DEFAULT '1' COMMENT '访问次数',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_visit_visitor_ref` (`visitor_id`,`ref_code`) COMMENT '同访客同邀请码去重',
  KEY `idx_visit_referrer_user` (`referrer_user_id`) COMMENT '按邀请人查访问',
  KEY `idx_visit_ref_code` (`ref_code`) COMMENT '按邀请码查访问',
  CONSTRAINT `fk_visit_referrer_user` FOREIGN KEY (`referrer_user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='邀请链接访问/点击记录表(匿名访问，用于统计点击与转化)';

-- ----------------------------
-- Records of loan_referral_visits
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_repayment_schedules
-- ----------------------------
DROP TABLE IF EXISTS `loan_repayment_schedules`;
CREATE TABLE `loan_repayment_schedules` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '还款计划ID(期次记录)',
  `disbursement_id` bigint NOT NULL COMMENT '关联放款单 loan_disbursements.id',
  `installment_no` int NOT NULL COMMENT '期次(从1开始)',
  `due_date` date NOT NULL COMMENT '应还日期',
  `principal_due` int NOT NULL DEFAULT '0' COMMENT '应还本金(建议统一单位：分)',
  `interest_due` int NOT NULL DEFAULT '0' COMMENT '应还利息(分)',
  `fee_due` int NOT NULL DEFAULT '0' COMMENT '应还费用(分)',
  `penalty_due` int NOT NULL DEFAULT '0' COMMENT '应还罚息(分，逾期产生)',
  `total_due` int NOT NULL COMMENT '本期应还总额=本金+利息+费用+罚息(分)',
  `paid_principal` int NOT NULL DEFAULT '0' COMMENT '已还本金(分)',
  `paid_interest` int NOT NULL DEFAULT '0' COMMENT '已还利息(分)',
  `paid_fee` int NOT NULL DEFAULT '0' COMMENT '已还费用(分)',
  `paid_penalty` int NOT NULL DEFAULT '0' COMMENT '已还罚息(分)',
  `paid_total` int NOT NULL DEFAULT '0' COMMENT '已还总额(分)',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '期次状态：0未还清 1已还清 2逾期',
  `last_paid_at` datetime DEFAULT NULL COMMENT '最近一次还款时间',
  `settled_at` datetime DEFAULT NULL COMMENT '结清时间(本期还清时)',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_disburse_installment` (`disbursement_id`,`installment_no`) COMMENT '同一放款单期次唯一',
  KEY `idx_schedule_due` (`due_date`,`status`) COMMENT '按到期日/状态查催收列表',
  KEY `idx_schedule_disburse` (`disbursement_id`) COMMENT '按放款单查计划',
  CONSTRAINT `fk_schedule_disbursement` FOREIGN KEY (`disbursement_id`) REFERENCES `loan_disbursements` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='还款计划表(支持单期/分期，逾期/已还通过状态体现)';

-- ----------------------------
-- Records of loan_repayment_schedules
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_repayment_transactions
-- ----------------------------
DROP TABLE IF EXISTS `loan_repayment_transactions`;
CREATE TABLE `loan_repayment_transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '回款流水ID',
  `disbursement_id` bigint NOT NULL COMMENT '关联放款单 loan_disbursements.id',
  `schedule_id` bigint DEFAULT NULL COMMENT '关联期次 loan_repayment_schedules.id(可空：先入账后分配/未分期)',
  `collect_channel_id` bigint DEFAULT NULL COMMENT '回款渠道(代收) loan_payment_channels.id',
  `collect_order_no` varchar(128) DEFAULT NULL COMMENT '回款订单号/三方代收单号(商户单号)',
  `pay_ref` varchar(128) DEFAULT NULL COMMENT '支付渠道流水号/交易号(三方transaction id)',
  `pay_amount` int NOT NULL COMMENT '本次回款金额(分)',
  `pay_method` varchar(32) DEFAULT NULL COMMENT '回款方式(如 BANK_TRANSFER/CARD/WALLET/CASH)',
  `paid_at` datetime NOT NULL COMMENT '回款时间(交易成功时间)',
  `alloc_principal` int NOT NULL DEFAULT '0' COMMENT '本次分配到本金(分)',
  `alloc_interest` int NOT NULL DEFAULT '0' COMMENT '本次分配到利息(分)',
  `alloc_fee` int NOT NULL DEFAULT '0' COMMENT '本次分配到费用(分)',
  `alloc_penalty` int NOT NULL DEFAULT '0' COMMENT '本次分配到罚息(分)',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '流水状态：1成功 0失败 2冲正/撤销',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_tx_disburse_time` (`disbursement_id`,`paid_at`) COMMENT '按放款单/时间查回款流水',
  KEY `idx_tx_schedule` (`schedule_id`) COMMENT '按期次查回款流水',
  KEY `idx_tx_collect_channel_time` (`collect_channel_id`,`paid_at`) COMMENT '按回款渠道/时间查流水',
  KEY `idx_tx_collect_order_no` (`collect_order_no`) COMMENT '按回款订单号查询',
  KEY `idx_tx_pay_ref` (`pay_ref`) COMMENT '按交易流水号查询',
  CONSTRAINT `fk_tx_collect_channel` FOREIGN KEY (`collect_channel_id`) REFERENCES `loan_payment_channels` (`id`),
  CONSTRAINT `fk_tx_disbursement` FOREIGN KEY (`disbursement_id`) REFERENCES `loan_disbursements` (`id`),
  CONSTRAINT `fk_tx_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `loan_repayment_schedules` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='回款流水表(记录每次实际回款，支持分期/部分还款，含回款渠道与订单号)';

-- ----------------------------
-- Records of loan_repayment_transactions
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_role_departments
-- ----------------------------
DROP TABLE IF EXISTS `loan_role_departments`;
CREATE TABLE `loan_role_departments` (
  `role_id` bigint NOT NULL,
  `department_id` bigint NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate id',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`role_id`,`department_id`),
  UNIQUE KEY `uk_id` (`id`),
  KEY `department_id` (`department_id`),
  CONSTRAINT `loan_role_departments_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `loan_roles` (`id`),
  CONSTRAINT `loan_role_departments_ibfk_2` FOREIGN KEY (`department_id`) REFERENCES `loan_departments` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_role_departments
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `loan_role_permissions`;
CREATE TABLE `loan_role_permissions` (
  `role_id` bigint NOT NULL,
  `permission_id` bigint NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate id',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`role_id`,`permission_id`),
  UNIQUE KEY `uk_id` (`id`),
  KEY `permission_id` (`permission_id`),
  CONSTRAINT `loan_role_permissions_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `loan_roles` (`id`),
  CONSTRAINT `loan_role_permissions_ibfk_2` FOREIGN KEY (`permission_id`) REFERENCES `loan_permissions` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_role_permissions
-- ----------------------------
BEGIN;
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (1, 1, '2026-01-14 18:22:02', NULL, 1, '2026-01-14 18:22:02');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (1, 2, '2026-01-14 18:22:02', NULL, 2, '2026-01-14 18:22:02');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (1, 3, '2026-01-14 18:22:02', NULL, 3, '2026-01-14 18:22:02');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (1, 4, '2026-01-14 18:22:02', NULL, 4, '2026-01-14 18:22:02');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (1, 5, '2026-01-14 18:22:02', NULL, 5, '2026-01-14 18:22:02');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (1, 6, '2026-01-14 18:22:02', NULL, 6, '2026-01-14 18:22:02');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (2, 5, '2026-01-14 18:22:13', NULL, 7, '2026-01-14 18:22:13');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (2, 6, '2026-01-14 18:22:13', NULL, 8, '2026-01-14 18:22:13');
INSERT INTO `loan_role_permissions` (`role_id`, `permission_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (3, 3, '2026-01-14 18:22:20', NULL, 9, '2026-01-14 18:22:20');
COMMIT;

-- ----------------------------
-- Table structure for loan_roles
-- ----------------------------
DROP TABLE IF EXISTS `loan_roles`;
CREATE TABLE `loan_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(64) NOT NULL,
  `name` varchar(128) NOT NULL,
  `data_scope` varchar(32) NOT NULL DEFAULT 'DEPT',
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_roles
-- ----------------------------
BEGIN;
INSERT INTO `loan_roles` (`id`, `code`, `name`, `data_scope`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'admin', '系统管理员', 'ALL', 1, '2026-01-14 18:21:42', '2026-01-14 18:21:42', NULL);
INSERT INTO `loan_roles` (`id`, `code`, `name`, `data_scope`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 'collector', '催收人员', 'DEPT', 1, '2026-01-14 18:21:42', '2026-01-14 18:21:42', NULL);
INSERT INTO `loan_roles` (`id`, `code`, `name`, `data_scope`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 'auditor', '审核人员', 'DEPT', 1, '2026-01-14 18:21:42', '2026-01-14 18:21:42', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_user_call_records
-- ----------------------------
DROP TABLE IF EXISTS `loan_user_call_records`;
CREATE TABLE `loan_user_call_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `baseinfo_id` int NOT NULL COMMENT '关联 loan_baseinfo.id',
  `call_type` tinyint NOT NULL COMMENT '通话类型：1呼入 2呼出 3未接 4拒接(按采集端定义)',
  `phone_number` varchar(32) DEFAULT NULL COMMENT '对端号码/电话',
  `phone_normalized` varchar(32) DEFAULT NULL COMMENT '标准化号码(去空格/国家码等，可选)',
  `call_time` datetime DEFAULT NULL COMMENT '通话开始时间(手机侧时间)',
  `duration_seconds` int NOT NULL DEFAULT '0' COMMENT '通话时长(秒，未接/拒接一般为0)',
  `call_hash` char(64) DEFAULT NULL COMMENT '记录去重哈希(如 sha256(type+phone+call_time+duration)，可选)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_calls_baseinfo` (`baseinfo_id`) COMMENT '按申请单查询通话记录',
  KEY `idx_calls_time` (`call_time`) COMMENT '按通话时间查询',
  KEY `idx_calls_phone` (`phone_normalized`) COMMENT '按标准化号码查询(可用于风控)',
  KEY `idx_calls_hash` (`call_hash`) COMMENT '按去重哈希查询',
  CONSTRAINT `fk_calls_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='通话记录采集表(匿名表单采集，与loan_baseinfo关联)';

-- ----------------------------
-- Records of loan_user_call_records
-- ----------------------------
BEGIN;
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 3, 2, '+86 13800138000', '8613800138000', '2026-01-10 09:15:20', 125, 'c1570026172aa42d89074e9a09fff26ae2c7777e6875bce6e683df2ba93fb9ef', '2026-01-10 09:18:00', '2026-01-10 09:18:00', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 3, 1, '13900139000', '8613900139000', '2026-01-10 10:20:15', 88, '9e5a04336451617baac52cbd0ecfc2790912575ae7c9b3f5d0bdc396a8aa3b7c', '2026-01-10 10:22:00', '2026-01-10 10:22:00', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 3, 3, '13700137000', '8613700137000', '2026-01-10 11:05:30', 0, 'bab35c29d52684955f9b224a9b624c73276ce9d4c33690013c2e40479cf50c59', '2026-01-10 11:06:00', '2026-01-10 11:06:00', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 3, 4, '13600136000', '8613600136000', '2026-01-10 14:30:45', 0, '017eb062a7e9829786d75610283af220cfb54edf157ba50dc51f5d278a063f56', '2026-01-10 14:31:00', '2026-01-10 14:31:00', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 3, 2, '13500135000', '8613500135000', '2026-01-10 15:10:00', 350, 'b4641a6f4c245ed7180660280c3ed13f43e7ee6c58493e365aec81b2434c01ed', '2026-01-10 15:16:00', '2026-01-10 15:16:00', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 3, 1, '13400134000', '8613400134000', '2026-01-10 16:20:00', 15, '43f3c1c5337abb3b55220ff4ba201223bb2abf31086f4d16ff4ce38140369d5c', '2026-01-10 16:20:20', '2026-01-10 16:20:20', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 3, 3, '18800188000', '8618800188000', '2026-01-10 17:00:10', 0, '528253ad4e59bc8fa1b5b11c86e8bda46e23cf5f751edc12a2aa12b2c7486ac1', '2026-01-10 17:00:20', '2026-01-10 17:00:20', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 3, 4, '18900189000', '8618900189000', '2026-01-10 18:45:00', 0, '71663ca5421460222171c19baa8e38c867791cfb352400c9a773d0886969d478', '2026-01-10 18:45:10', '2026-01-10 18:45:10', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 3, 2, '19900199000', '8619900199000', '2026-01-10 19:10:30', 180, '41501e645800a8bec75d11dcb3a0f58dd359d1c0b621d14bcb4ba4960541387f', '2026-01-10 19:13:30', '2026-01-10 19:13:30', NULL);
INSERT INTO `loan_user_call_records` (`id`, `baseinfo_id`, `call_type`, `phone_number`, `phone_normalized`, `call_time`, `duration_seconds`, `call_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 3, 3, '17700177000', '8617700177000', '2026-01-10 20:05:15', 0, '4007ea6e66c088386210f6c2fd4d904e1251ce37597f1d0f41485c6200f09c41', '2026-01-10 20:05:25', '2026-01-10 20:05:25', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_user_contacts
-- ----------------------------
DROP TABLE IF EXISTS `loan_user_contacts`;
CREATE TABLE `loan_user_contacts` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `baseinfo_id` int NOT NULL COMMENT '关联 loan_baseinfo.id',
  `contact_name` varchar(128) DEFAULT NULL COMMENT '联系人姓名',
  `phone_number` varchar(32) DEFAULT NULL COMMENT '联系人手机号/电话',
  `contact_hash` char(64) DEFAULT NULL COMMENT '联系人去重哈希(如 sha256(name+phone_normalized))',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_contacts_baseinfo` (`baseinfo_id`) COMMENT '按申请单查询通讯录',
  KEY `idx_contacts_hash` (`contact_hash`) COMMENT '按去重哈希查询',
  CONSTRAINT `fk_contacts_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户通讯录采集表(匿名表单采集，与loan_baseinfo关联)';

-- ----------------------------
-- Records of loan_user_contacts
-- ----------------------------
BEGIN;
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 3, '父亲', '13800138001', '3629b5d3f51a7da15e01f8e211e1d3d75ff15a70df0bfeec61f322eee968cbf8', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 3, '母亲', '13800138002', '105cfa7ea4c8472d8e612c333b0695930a55876c9458189419084533203cb9ab', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 3, '老婆', '13800138003', 'd5803c78bd01d0a316bea57e0985559d4ed66de5de05af0d9d8866336de852db', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 3, '张三-同事', '13800138004', '35713607ab96b311171a5d0ce99f26a0ed2e09582977e994a733da75c784a97a', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 3, '李四-朋友', '13800138005', 'ab55450af9c8317739671671572470bef5bb04b89fd450fb47ca4a9ad2ace2e8', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 3, '王总-公司', '13800138006', 'fb9be39780b9388b854c42c8796f250a2fd66c76fc15c891b6f7d1f0670fa885', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 3, '人事部-李姐', '13800138007', 'f583ade13e0c3a571d4f249a3e3c7c4c14b3fcd357f1f05546bfb6fb08102d0b', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 3, '小区物业', '13800138008', '096dacf797c263b82eb4cf9066427777ec64e21b50b7dc227782919d4e16080e', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 3, '顺丰快递员', '13800138009', 'a348ce033639acc93546b27c632ce656804d1eab260521c821c9f1fbc11ec22d', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
INSERT INTO `loan_user_contacts` (`id`, `baseinfo_id`, `contact_name`, `phone_number`, `contact_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 3, '表哥-王磊', '13800138010', '84ba333f132b95e6c2af09074bb1e69a036037ecc6fc91f2ca2a4f0a47e5ca23', '2026-01-10 00:00:00', '2026-01-10 00:00:00', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_user_device_apps
-- ----------------------------
DROP TABLE IF EXISTS `loan_user_device_apps`;
CREATE TABLE `loan_user_device_apps` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `baseinfo_id` int NOT NULL COMMENT '关联 loan_baseinfo.id',
  `package_name` varchar(255) NOT NULL COMMENT '应用包名/BundleId(如 com.xxx.app)',
  `app_name` varchar(255) DEFAULT NULL COMMENT '应用名称',
  `version_name` varchar(64) DEFAULT NULL COMMENT '版本名(如 1.2.3)',
  `version_code` bigint DEFAULT NULL COMMENT '版本号(如 Android versionCode，可选)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_apps_baseinfo` (`baseinfo_id`) COMMENT '按申请单查询应用列表',
  KEY `idx_apps_package` (`package_name`) COMMENT '按包名查询(可用于风控)',
  CONSTRAINT `fk_apps_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备软件列表采集表(匿名表单采集，与loan_baseinfo关联)';

-- ----------------------------
-- Records of loan_user_device_apps
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for loan_user_roles
-- ----------------------------
DROP TABLE IF EXISTS `loan_user_roles`;
CREATE TABLE `loan_user_roles` (
  `user_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'surrogate id',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`user_id`,`role_id`),
  UNIQUE KEY `uk_id` (`id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `loan_user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `loan_users` (`id`),
  CONSTRAINT `loan_user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `loan_roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_user_roles
-- ----------------------------
BEGIN;
INSERT INTO `loan_user_roles` (`user_id`, `role_id`, `created_at`, `deleted_at`, `id`, `updated_at`) VALUES (1, 1, '2026-01-14 18:22:26', NULL, 1, '2026-01-14 18:22:26');
COMMIT;

-- ----------------------------
-- Table structure for loan_user_sms_records
-- ----------------------------
DROP TABLE IF EXISTS `loan_user_sms_records`;
CREATE TABLE `loan_user_sms_records` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `baseinfo_id` int NOT NULL COMMENT '关联 loan_baseinfo.id',
  `direction` tinyint NOT NULL COMMENT '短信方向：1收(inbox) 2发(sent)',
  `address` varchar(64) DEFAULT NULL COMMENT '对端号码/短码/发件人(如银行短码)',
  `sms_time` datetime DEFAULT NULL COMMENT '短信时间(手机侧时间)',
  `body` text COMMENT '短信内容(可选，敏感数据请注意合规)',
  `body_hash` char(64) DEFAULT NULL COMMENT '短信内容哈希(用于去重/审计，可选)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_sms_baseinfo` (`baseinfo_id`) COMMENT '按申请单查询短信',
  KEY `idx_sms_time` (`sms_time`) COMMENT '按短信时间查询',
  CONSTRAINT `fk_sms_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='短信记录采集表(匿名表单采集，与loan_baseinfo关联)';

-- ----------------------------
-- Records of loan_user_sms_records
-- ----------------------------
BEGIN;
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 3, 1, '95588', '2026-01-10 08:30:15', '【工商银行】您尾号1234的储蓄卡1月10日08:29入账工资15000元，余额28560.78元。', '6ddeb8f1a94ed1bc276b197bf805be79f96da2a13ae545a395e27cdc2c0fd703', '2026-01-10 08:31:00', '2026-01-10 08:31:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 3, 2, '13800138000', '2026-01-10 09:45:20', '今天下午有空吗？想约你喝杯咖啡', '53e7cc4a352afacfb3abb27996d700f0ac9e6b7cbd3c55086083f00385de544b', '2026-01-10 09:46:00', '2026-01-10 09:46:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 3, 1, '10086', '2026-01-10 10:15:05', '【中国移动】您本月已使用流量8.5GB，剩余2.5GB，请注意流量使用。', 'ef852cd4e4e674b451e0f7f5349134d0cb9ac6fbadae0b0b072e0cf35d9a6cfb', '2026-01-10 10:16:00', '2026-01-10 10:16:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 3, 2, '13900139000', '2026-01-10 11:20:30', '王总，项目方案我已经发您邮箱了，麻烦抽空看下，有问题随时沟通。', '1d21206d887a274e2520716118f3e41f2f78ba1eef9171ca6b27de81737ef0ea', '2026-01-10 11:21:00', '2026-01-10 11:21:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 3, 1, '95311', '2026-01-10 13:05:10', '【顺丰速运】您的快递(单号SF1234567890)已到达XX小区快递柜，取件码123456，有效期24小时。', '3c6b82193e8c4aa7a99af2633a986dada27812079c22ce0c0f18bf02828befd6', '2026-01-10 13:06:00', '2026-01-10 13:06:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 3, 2, '18800188000', '2026-01-10 14:10:45', '提醒一下，本月房贷请于15日前存入尾号5678的银行卡，金额8500元。', '151aa0aa0c172ce0e51f9fa57889b8bfbacd8460778b164712e6104b3290ab83', '2026-01-10 14:11:00', '2026-01-10 14:11:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 3, 1, '1069000000', '2026-01-10 15:30:20', '【XX金融】您的验证码是876543，5分钟内有效，请勿泄露给他人。', 'f56a59ca53cc7268c41fc71e13aa3eb1a1923194e95579f5092398cd6b538b71', '2026-01-10 15:31:00', '2026-01-10 15:31:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (8, 3, 2, '17700177000', '2026-01-10 16:40:15', '李经理，我明天上午需要请假半天去办理证件，工作已交接给小张，望批准。', 'e499d522777105d950a278d3f83f85b711112d45432b861b9866ce0cd7261e87', '2026-01-10 16:41:00', '2026-01-10 16:41:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, 3, 1, '4008888888', '2026-01-10 17:25:30', '【XX信用卡】您本期账单金额5890.50元，还款日1月20日，最低还款589.05元。', '5dfa9f0f34f6f026fe7d3b9d38f0e1181f09feacbfb9ec80a22dbb4169cd298c', '2026-01-10 17:26:00', '2026-01-10 17:26:00', NULL);
INSERT INTO `loan_user_sms_records` (`id`, `baseinfo_id`, `direction`, `address`, `sms_time`, `body`, `body_hash`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, 3, 2, '19900199000', '2026-01-10 18:50:05', '爸妈，我今晚加班晚点回家，不用等我吃饭了。', 'e28d9ceff670b4aea0d9141ed4b0883002340532d4e13e8d503183ea2e26545d', '2026-01-10 18:51:00', '2026-01-10 18:51:00', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_users
-- ----------------------------
DROP TABLE IF EXISTS `loan_users`;
CREATE TABLE `loan_users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `department_id` bigint NOT NULL,
  `mfa_enabled` tinyint NOT NULL DEFAULT '0',
  `mfa_required` tinyint NOT NULL DEFAULT '0',
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `share_code` varchar(32) DEFAULT NULL COMMENT '分享邀请码(用于生成分享链接，建议唯一)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `uk_users_share_code` (`share_code`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_users
-- ----------------------------
BEGIN;
INSERT INTO `loan_users` (`id`, `username`, `password_hash`, `department_id`, `mfa_enabled`, `mfa_required`, `status`, `created_at`, `updated_at`, `deleted_at`, `share_code`) VALUES (1, 'axbros', '$2a$10$Y0Wr7iRD2Z1xxNwV7uDAXuAhG6ECseNjUcjdnpwSdVcE1s3I/duKy', 1, 1, 0, 1, '2026-01-14 17:51:53', '2026-01-14 17:51:53', NULL, 'MvJHGT7XL-cp');
INSERT INTO `loan_users` (`id`, `username`, `password_hash`, `department_id`, `mfa_enabled`, `mfa_required`, `status`, `created_at`, `updated_at`, `deleted_at`, `share_code`) VALUES (2, 'auditor', '$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5lWlP0r8kP7r9v8pJwD3h8m6cK4QK', 1, 0, 0, 1, '2026-01-14 19:36:30', '2026-01-14 19:36:30', NULL, 'REFAUDIT01');
INSERT INTO `loan_users` (`id`, `username`, `password_hash`, `department_id`, `mfa_enabled`, `mfa_required`, `status`, `created_at`, `updated_at`, `deleted_at`, `share_code`) VALUES (3, 'referrer', '$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5lWlP0r8kP7r9v8pJwD3h8m6cK4QK', 1, 0, 0, 1, '2026-01-14 19:36:30', '2026-01-14 19:36:30', NULL, 'REFSHARE01');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
