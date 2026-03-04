/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80024
 Source Host           : localhost:3306
 Source Schema         : tophone_loan

 Target Server Type    : MySQL
 Target Server Version : 80024
 File Encoding         : 65001

 Date: 04/03/2026 10:08:58
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
  `audit_type` varchar(255) DEFAULT NULL COMMENT '审核类型(初审1、放款审核2、回款审核3)',
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL COMMENT '审核时间(即审核通过/拒绝时间)',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_audit_baseinfo` (`baseinfo_id`) COMMENT '按申请单查询审核记录',
  KEY `idx_audit_user_time` (`auditor_user_id`,`created_at`) COMMENT '按审核人/时间查询',
  CONSTRAINT `fk_audit_baseinfo` FOREIGN KEY (`baseinfo_id`) REFERENCES `loan_baseinfo` (`id`),
  CONSTRAINT `fk_audit_user` FOREIGN KEY (`auditor_user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='申请审核记录表(审核时间即 created_at)';

-- ----------------------------
-- Records of loan_audits
-- ----------------------------
BEGIN;
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (38, 3, 1, '初审通过', 1, '1', '2026-02-13 22:26:25', '2026-02-13 22:26:25', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (39, 3, 1, '财务审核通过', 1, '2', '2026-02-13 22:26:44', '2026-02-13 22:26:44', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (40, 4, 1, '初审通过', 1, '1', '2026-02-13 22:28:47', '2026-02-13 22:28:47', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (41, 4, 1, '统一', 1, '2', '2026-02-13 22:36:57', '2026-02-13 22:36:57', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (42, 6, 1, '初审通过', 1, '1', '2026-02-13 22:38:26', '2026-02-13 22:38:26', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (43, 1, 1, '初审通过', 1, '1', '2026-02-13 22:41:09', '2026-02-13 22:41:09', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (44, 5, 1, '初审通过', 1, '1', '2026-02-14 14:44:36', '2026-02-14 14:44:36', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (45, 1, 1, '财务审核通过', 1, '2', '2026-02-14 14:45:09', '2026-02-14 14:45:09', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (46, 9, 1, '测试用户 审核通过', 1, '1', '2026-02-14 16:02:16', '2026-02-14 16:02:16', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (47, 9, 1, '财务审核通过 测试', 1, '2', '2026-02-14 16:02:47', '2026-02-14 16:02:47', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (48, 7, 1, '初审通过', 1, '1', '2026-02-24 14:27:15', '2026-02-24 14:27:15', NULL);
INSERT INTO `loan_audits` (`id`, `baseinfo_id`, `audit_result`, `audit_comment`, `auditor_user_id`, `audit_type`, `updated_at`, `created_at`, `deleted_at`) VALUES (49, 7, 1, '优质用户 放款通过', 1, '2', '2026-02-24 14:28:03', '2026-02-24 14:28:03', NULL);
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
  `application_amount` bigint DEFAULT NULL COMMENT '申請金額',
  `audit_status` tinyint DEFAULT '0' COMMENT '審核情況 0待審核 1初审通過 2财务审核通过 -1 審核拒絕',
  `bank_no` varchar(255) DEFAULT NULL COMMENT '銀行卡號',
  `client_ip` varbinary(16) DEFAULT NULL COMMENT '客户端IP地址(IPv4/IPv6)',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `referrer_user_id` bigint DEFAULT NULL COMMENT '邀请人/分享人(loan_users.id)',
  `ref_code` varchar(32) DEFAULT NULL COMMENT '访问时携带的ref(冗余存储便于排查)',
  `loan_days` smallint NOT NULL COMMENT '借款天数(单位：天)',
  `risk_list_status` tinyint NOT NULL DEFAULT '0' COMMENT '名单状态：0正常 1白名单 2黑名单',
  `risk_list_reason` varchar(255) DEFAULT NULL COMMENT '名单原因/来源说明',
  `risk_list_marked_at` datetime DEFAULT NULL COMMENT '名单标记时间',
  PRIMARY KEY (`id`),
  KEY `idx_baseinfo_referrer_user` (`referrer_user_id`) COMMENT '按邀请人查询申请记录',
  KEY `idx_baseinfo_ref_code` (`ref_code`) COMMENT '按ref查询',
  CONSTRAINT `fk_baseinfo_referrer_user` FOREIGN KEY (`referrer_user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=187 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_baseinfo
-- ----------------------------
BEGIN;
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (1, 'Wang', 'Lei', 29, 'M', '16600229988', 'ID_CARD', '110101199401011234', 'oss://idcard/wanglei_front.jpg', 'Android', 'Engineer', 'TechSoft Ltd.', 15000, 0, 0, 0, 100000, 2, '6222020200001234567', 0x3139322E3136382E332E31, '2026-01-14 19:12:30', '2026-02-14 14:45:09', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (2, 'Li', 'Na', 34, 'W', '16600229988', 'ID_CARD', '310101198912123456', NULL, 'iOS', 'Sales', 'TradeCorp', 12000, NULL, NULL, NULL, 80000, 0, '6222020200007654321', 0x3139322E3136382E332E31, '2026-01-14 19:12:36', '2026-01-14 19:12:36', '2026-01-15 13:38:22', 1, 'REFAXBROS01', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (3, 'Li', 'Na', 34, 'W', '16600229988', 'ID_CARD', '310101198912123456', NULL, 'iOS', 'Sales', 'TradeCorp', 12000, NULL, NULL, NULL, 80000, 2, '6222020200007654321', 0x3139322E3136382E332E31, '2026-01-14 19:12:41', '2026-02-13 22:26:44', NULL, 1, 'REFAXBROS01', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (4, 'Chen', 'Yu', 27, 'W', '16600229988', 'ID_CARD', '320101199612124321', NULL, 'Android', 'Designer', 'Creative Studio', 18000, NULL, NULL, NULL, 50000, 2, '6228480402567890123', 0x3139322E3136382E332E31, '2026-01-14 19:12:48', '2026-02-13 22:36:57', NULL, NULL, NULL, 14, 1, '历史还款良好，人工加入白名单', '2026-01-14 19:12:48');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (5, 'Liu', 'Ming', 36, 'M', '16600229988', 'ID_CARD', '510101198801019999', NULL, 'Web', 'Freelancer', NULL, NULL, NULL, NULL, NULL, 30000, 1, NULL, 0x3139322E3136382E332E31, '2026-01-14 19:12:55', '2026-02-14 14:44:36', NULL, NULL, NULL, 7, 2, '命中内部黑名单：多次逾期', '2026-01-14 19:12:55');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (6, 'Test', 'A', 28, 'M', '16600229988', 'ID_CARD', 'TID0001', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 100000, 1, '6222000000000001', 0x0A000001, '2026-01-14 19:39:40', '2026-02-13 22:38:26', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (7, 'Test', 'B', 30, 'W', '16600229988', 'ID_CARD', 'TID0002', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 120000, 2, '6222000000000002', 0x0A000002, '2026-01-14 19:39:40', '2026-02-24 14:28:03', NULL, NULL, NULL, 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (8, 'Test', 'C', 26, 'M', '16600229988', 'ID_CARD', 'TID0003', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 90000, 0, '6222000000000003', 0x0A000003, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (9, 'Test', 'D', 35, 'W', '16600229988', 'ID_CARD', 'TID0004', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 150000, 2, '6222000000000004', 0x0A000004, '2026-01-14 19:39:40', '2026-02-14 16:02:47', NULL, NULL, NULL, 60, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (10, 'Test', 'E', 41, 'M', '16600229988', 'ID_CARD', 'TID0005', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 80000, 0, '6222000000000005', 0x0A000005, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (11, 'Test', 'F', 33, 'W', '16600229988', 'ID_CARD', 'TID0006', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 110000, 0, '6222000000000006', 0x0A000006, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (12, 'Test', 'G', 29, 'M', '16600229988', 'ID_CARD', 'TID0007', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 200000, 0, '6222000000000007', 0x0A000007, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (13, 'Test', 'H', 27, 'W', '16600229988', 'ID_CARD', 'TID0008', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 70000, 0, '6222000000000008', 0x0A000008, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (14, 'Test', 'I', 38, 'M', '16600229988', 'ID_CARD', 'TID0009', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 180000, 0, '6222000000000009', 0x0A000009, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 45, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (15, 'Test', 'J', 24, 'W', '16600229988', 'ID_CARD', 'TID0010', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 60000, 0, '6222000000000010', 0x0A00000A, '2026-01-14 19:39:40', '2026-01-14 19:39:40', NULL, NULL, NULL, 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (66, '王', '磊', 29, 'M', '16600229988', 'ID_CARD', '110101199401011234', 'ID110101199401011234.jpg', 'iOS 17.0', '软件工程师', '北京科技有限公司', 25000, 1, 1, 1, 100000, 0, '6222081001001234567', 0xC0A80165, '2026-01-01 10:00:00', '2026-01-01 10:00:00', NULL, 1, 'REF00001', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (67, '王', '磊', 29, 'M', '16600229988', 'ID_CARD', '110101199401011234', 'ID110101199401011234.jpg', 'iOS 17.0', '软件工程师', '北京科技有限公司', 25000, 1, 1, 1, 100000, 0, '6222081001001234567', 0xC0A80165, '2026-01-01 10:00:00', '2026-01-01 10:00:00', NULL, 1, 'REF00001', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (68, '李', '娜', 34, 'W', '16600229988', 'ID_CARD', '310101198912123456', 'ID310101198912123456.jpg', 'Android 14', '财务经理', '上海金融有限公司', 35000, 1, 1, 1, 80000, 0, '6228480402567890123', 0xC0A80166, '2026-01-02 11:00:00', '2026-01-02 11:00:00', NULL, 2, 'REF00002', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (69, '陈', '宇', 27, 'W', '16600229988', 'ID_CARD', '320101199612124321', 'ID320101199612124321.jpg', 'iOS 16.5', '人力资源专员', '江苏贸易有限公司', 18000, 0, 0, 0, 50000, 0, '6259991234567890123', 0xC0A80167, '2026-01-03 09:30:00', '2026-01-03 09:30:00', NULL, 3, 'REF00003', 14, 1, '优质客户', '2026-01-03 10:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (70, '刘', '明', 36, 'M', '16600229988', 'ID_CARD', '510101198801019999', 'ID510101198801019999.jpg', 'Windows 11', '销售总监', '四川科技有限公司', 45000, 1, 1, 1, 30000, 0, '6226667890123456789', 0xC0A80168, '2026-01-04 14:00:00', '2026-01-04 15:00:00', NULL, 1, 'REF00004', 7, 2, '信用逾期', '2026-01-04 16:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (71, '张', '伟', 41, 'M', '16600229988', 'ID_CARD', '440101198305051234', 'ID440101198305051234.jpg', 'macOS 14', '项目经理', '广东建设有限公司', 50000, 1, 1, 1, 150000, 0, '6229998765432109876', 0xC0A80169, '2026-01-05 08:00:00', '2026-01-05 08:30:00', NULL, 2, 'REF00005', 21, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (72, '赵', '丽', 25, 'W', '16600229988', 'ID_CARD', '210101199808085678', 'ID210101199808085678.jpg', 'Android 13', '行政助理', '辽宁商贸有限公司', 15000, 0, 0, 0, 20000, 0, '6227779876543210987', 0xC0A8016A, '2026-01-06 16:00:00', '2026-01-06 16:00:00', NULL, 3, 'REF00006', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (73, '黄', '浩', 33, 'M', '16600229988', 'ID_CARD', '430101199011118765', 'ID430101199011118765.jpg', 'iOS 17.1', '产品经理', '湖南科技有限公司', 30000, 1, 0, 1, 70000, 0, '6225551234567890123', 0xC0A8016B, '2026-01-07 10:30:00', '2026-01-07 10:30:00', NULL, 1, 'REF00007', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (74, '周', '敏', 28, 'W', '16600229988', 'ID_CARD', '330101199503037654', 'ID330101199503037654.jpg', 'Windows 10', '客服专员', '浙江服务有限公司', 16000, 0, 0, 0, 15000, 0, '6224448765432109876', 0xC0A8016C, '2026-01-08 13:00:00', '2026-01-08 14:00:00', NULL, 2, 'REF00008', 14, 2, '收入不稳定', '2026-01-08 15:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (75, '吴', '强', 38, 'M', '16600229988', 'ID_CARD', '350101198507076543', 'ID350101198507076543.jpg', 'Android 12', '工程师', '福建制造有限公司', 40000, 1, 1, 1, 90000, 0, '6223337654321098765', 0xC0A8016D, '2026-01-09 09:00:00', '2026-01-09 09:15:00', NULL, 3, 'REF00009', 21, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (76, '徐', '芳', 31, 'W', '16600229988', 'ID_CARD', '610101199209094321', 'ID610101199209094321.jpg', 'macOS 13', '设计师', '陕西创意有限公司', 28000, 1, 0, 0, 60000, 0, '6222226543210987654', 0xC0A8016E, '2026-01-10 11:30:00', '2026-01-10 11:30:00', NULL, 1, 'REF00010', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (77, '孙', '杰', 26, 'M', '16600229988', 'ID_CARD', '120101199712123456', 'ID120101199712123456.jpg', 'iOS 16.4', '程序员', '天津科技有限公司', 22000, 0, 0, 1, 40000, 0, '6221115432109876543', 0xC0A8016F, '2026-01-11 15:00:00', '2026-01-11 15:00:00', NULL, 2, 'REF00011', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (78, '马', '丽', 37, 'W', '16600229988', 'ID_CARD', '650101198604048765', 'ID650101198604048765.jpg', 'Windows 11', '教师', '新疆教育有限公司', 20000, 1, 1, 0, 35000, 0, '6220004321098765432', 0xC0A80170, '2026-01-12 08:30:00', '2026-01-12 09:00:00', NULL, 3, 'REF00012', 30, 2, '负债过高', '2026-01-12 10:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (79, '朱', '军', 40, 'M', '16600229988', 'ID_CARD', '500101198310107654', 'ID500101198310107654.jpg', 'Android 14', '医生', '重庆医疗有限公司', 55000, 1, 1, 1, 120000, 0, '6219993210987654321', 0xC0A80171, '2026-01-13 14:30:00', '2026-01-13 14:45:00', NULL, 1, 'REF00013', 21, 1, '优质客户', '2026-01-13 15:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (80, '胡', '欣', 24, 'W', '16600229988', 'ID_CARD', '460101199906066543', 'ID460101199906066543.jpg', 'iOS 17.0', '实习生', '海南旅游有限公司', 8000, 0, 0, 0, 10000, 0, '6218882109876543210', 0xC0A80172, '2026-01-14 10:00:00', '2026-01-14 10:00:00', NULL, 2, 'REF00014', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (81, '林', '涛', 32, 'M', '16600229988', 'ID_CARD', '360101199102025432', 'ID360101199102025432.jpg', 'macOS 14', '摄影师', '江西传媒有限公司', 26000, 0, 0, 1, 50000, 0, '6217771098765432109', 0xC0A80173, '2026-01-15 16:30:00', '2026-01-15 16:30:00', NULL, 3, 'REF00015', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (82, '郭', '燕', 30, 'W', '16600229988', 'ID_CARD', '340101199307074321', 'ID340101199307074321.jpg', 'Windows 10', '护士', '安徽医疗有限公司', 24000, 1, 0, 0, 45000, 0, '6216660987654321098', 0xC0A80174, '2026-01-16 09:30:00', '2026-01-16 10:00:00', NULL, 1, 'REF00016', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (83, '何', '勇', 39, 'M', '16600229988', 'ID_CARD', '540101198408083210', 'ID540101198408083210.jpg', 'Android 13', '建筑工人', '西藏建设有限公司', 32000, 1, 1, 1, 80000, 0, '6215559876543210987', 0xC0A80175, '2026-01-17 13:30:00', '2026-01-17 14:00:00', NULL, 2, 'REF00017', 21, 2, '征信不良', '2026-01-17 15:00:00');
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (84, '高', '静', 27, 'W', '16600229988', 'ID_CARD', '630101199609092109', 'ID630101199609092109.jpg', 'iOS 16.5', '翻译', '青海外贸有限公司', 21000, 0, 0, 0, 25000, 0, '6214448765432109876', 0xC0A80176, '2026-01-18 11:00:00', '2026-01-18 11:00:00', NULL, 3, 'REF00018', 7, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (85, '罗', '刚', 35, 'M', '16600229988', 'ID_CARD', '530101198811111098', 'ID530101198811111098.jpg', 'Windows 11', '厨师', '云南餐饮有限公司', 28000, 1, 0, 1, 60000, 0, '6213337654321098765', 0xC0A80177, '2026-01-19 15:30:00', '2026-01-19 15:30:00', NULL, 1, 'REF00019', 14, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (86, '郑', '艳', 29, 'W', '16600229988', 'ID_CARD', '410101199412120987', 'ID410101199412120987.jpg', 'Android 12', '导购', '河南零售有限公司', 18000, 0, 0, 0, 30000, 0, '6212226543210987654', 0xC0A80178, '2026-01-20 08:00:00', '2026-01-20 08:15:00', NULL, 2, 'REF00020', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (137, '张', '三', 28, 'M', '13800138001', 'ID_CARD', '110101199601011234', '身份证照片链接1', 'iOS', '程序员', '科技有限公司', 15000, 1, 0, 1, 50000, 0, '6222021234567890123', 0xC0A80165, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00001', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (138, '李', '四', 32, 'M', '13900139002', 'ID_CARD', '120101199202021234', '身份证照片链接2', 'Android', '销售', '贸易公司', 8000, 1, 1, 0, 30000, 0, '6228481234567890123', 0xC0A80166, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00002', 15, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (139, '王', '五', 25, 'W', '13700137003', 'ID_CARD', '310101199903031234', '身份证照片链接3', 'iOS', '设计师', '广告公司', 10000, 0, 0, 0, 20000, 0, '6259991234567890123', 0xC0A80167, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00003', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (140, '赵', '六', 40, 'M', '13600136004', 'ID_CARD', '440101198404041234', '身份证照片链接4', 'Android', '经理', '金融公司', 25000, 1, 1, 1, 80000, 0, '6225881234567890123', 0xC0A80168, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00004', 45, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (141, '孙', '七', 29, 'W', '13500135005', 'ID_CARD', '510101199505051234', '身份证照片链接5', 'iOS', '教师', '学校', 8000, 0, 0, 1, 10000, 0, '6226661234567890123', 0xC0A80169, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00005', 10, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (142, '周', '八', 35, 'M', '13400134006', 'ID_CARD', '330101198906061234', '身份证照片链接6', 'Android', '工程师', '制造公司', 12000, 1, 1, 0, 40000, 0, '6229991234567890123', 0xC0A8016A, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00006', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (143, '吴', '九', 27, 'W', '13300133007', 'ID_CARD', '210101199707071234', '身份证照片链接7', 'iOS', '护士', '医院', 7000, 0, 0, 0, 15000, 0, '6227771234567890123', 0xC0A8016B, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00007', 25, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (144, '郑', '十', 38, 'M', '13200132008', 'ID_CARD', '130101198608081234', '身份证照片链接8', 'Android', '司机', '运输公司', 6000, 1, 0, 1, 25000, 0, '6224441234567890123', 0xC0A8016C, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00008', 40, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (145, '冯', '十一', 31, 'W', '13100131009', 'ID_CARD', '610101199309091234', '身份证照片链接9', 'iOS', '会计', '财务公司', 9000, 0, 1, 0, 35000, 0, '6223331234567890123', 0xC0A8016D, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00009', 15, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (146, '陈', '十二', 26, 'M', '13000130010', 'ID_CARD', '500101199810101234', '身份证照片链接10', 'Android', '厨师', '餐饮公司', 7500, 1, 0, 0, 18000, 0, '6221111234567890123', 0xC0A8016E, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00010', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (147, '褚', '十三', 33, 'W', '18800188011', 'ID_CARD', '430101199111111234', '身份证照片链接11', 'iOS', '客服', '电商公司', 6500, 0, 0, 1, 22000, 0, '6220001234567890123', 0xC0A8016F, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00011', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (148, '卫', '十四', 36, 'M', '18900189012', 'ID_CARD', '370101198812121234', '身份证照片链接12', 'Android', '电工', '电力公司', 8500, 1, 1, 1, 45000, 0, '6219991234567890123', 0xC0A80170, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00012', 35, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (149, '蒋', '十五', 24, 'W', '18700187013', 'ID_CARD', '320101200001131234', '身份证照片链接13', 'iOS', '实习生', '互联网公司', 4000, 0, 0, 0, 8000, 0, '6218881234567890123', 0xC0A80171, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00013', 10, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (150, '沈', '十六', 39, 'M', '18600186014', 'ID_CARD', '220101198502141234', '身份证照片链接14', 'Android', '木工', '装修公司', 7000, 1, 0, 1, 28000, 0, '6217771234567890123', 0xC0A80172, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00014', 45, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (151, '韩', '十七', 30, 'W', '18500185015', 'ID_CARD', '140101199403151234', '身份证照片链接15', 'iOS', '文案', '传媒公司', 9500, 0, 1, 0, 32000, 0, '6216661234567890123', 0xC0A80173, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00015', 25, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (152, '杨', '十八', 28, 'M', '18400184016', 'ID_CARD', '650101199604161234', '身份证照片链接16', 'Android', '焊工', '建筑公司', 9000, 1, 1, 0, 38000, 0, '6215551234567890123', 0xC0A80174, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00016', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (153, '朱', '十九', 34, 'W', '18300183017', 'ID_CARD', '540101199005171234', '身份证照片链接17', 'iOS', '主播', '直播公司', 12000, 0, 0, 1, 42000, 0, '6214441234567890123', 0xC0A80175, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00017', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (154, '秦', '二十', 37, 'M', '18200182018', 'ID_CARD', '450101198706181234', '身份证照片链接18', 'Android', '保安', '物业公司', 5000, 1, 0, 0, 12000, 0, '6213331234567890123', 0xC0A80176, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00018', 15, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (155, '尤', '二一', 27, 'W', '18100181019', 'ID_CARD', '410101199707191234', '身份证照片链接19', 'iOS', '导购', '商场', 5500, 0, 1, 1, 16000, 0, '6212221234567890123', 0xC0A80177, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00019', 40, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (156, '许', '二二', 30, 'M', '18000180020', 'ID_CARD', '360101199408201234', '身份证照片链接20', 'Android', '快递员', '物流公司', 6000, 1, 0, 0, 20000, 0, '6211111234567890123', 0xC0A80178, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00020', 25, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (157, '何', '二三', 29, 'W', '17800178021', 'ID_CARD', '350101199509211234', '身份证照片链接21', 'iOS', '翻译', '外贸公司', 11000, 0, 0, 0, 28000, 0, '6210001234567890123', 0xC0A80179, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00021', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (158, '吕', '二四', 35, 'M', '17700177022', 'ID_CARD', '340101198910221234', '身份证照片链接22', 'Android', '理发师', '美发店', 8000, 1, 1, 1, 35000, 0, '6209991234567890123', 0xC0A8017A, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00022', 15, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (159, '施', '二五', 26, 'W', '17600176023', 'ID_CARD', '310101199811231234', '身份证照片链接23', 'iOS', '化妆师', '影楼', 7000, 0, 0, 1, 14000, 0, '6208881234567890123', 0xC0A8017B, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00023', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (160, '张', '二六', 41, 'M', '17500175024', 'ID_CARD', '230101198312241234', '身份证照片链接24', 'Android', '维修工', '售后公司', 7500, 1, 0, 0, 22000, 0, '6207771234567890123', 0xC0A8017C, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00024', 45, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (161, '孔', '二七', 28, 'W', '17400174025', 'ID_CARD', '220101199601251234', '身份证照片链接25', 'iOS', '美甲师', '美甲店', 6000, 0, 1, 0, 10000, 0, '6206661234567890123', 0xC0A8017D, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00025', 10, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (162, '曹', '二八', 32, 'M', '17300173026', 'ID_CARD', '210101199202261234', '身份证照片链接26', 'Android', '健身教练', '健身房', 10000, 1, 1, 1, 40000, 0, '6205551234567890123', 0xC0A8017E, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00026', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (163, '严', '二九', 25, 'W', '17200172027', 'ID_CARD', '190101199903271234', '身份证照片链接27', 'iOS', '保育员', '幼儿园', 5000, 0, 0, 0, 8000, 0, '6204441234567890123', 0xC0A8017F, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00027', 25, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (164, '华', '三十', 38, 'M', '17100171028', 'ID_CARD', '180101198604281234', '身份证照片链接28', 'Android', '外卖员', '餐饮平台', 6500, 1, 0, 1, 18000, 0, '6203331234567890123', 0xC0A80180, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00028', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (165, '金', '三一', 31, 'W', '17000170029', 'ID_CARD', '150101199305291234', '身份证照片链接29', 'iOS', '摄影师', '摄影工作室', 9500, 0, 1, 0, 30000, 0, '6202221234567890123', 0xC0A80181, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00029', 35, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (166, '魏', '三二', 27, 'M', '16600166030', 'ID_CARD', '140101199706301234', '身份证照片链接30', 'Android', '洗车工', '洗车行', 4500, 1, 0, 0, 12000, 0, '6201111234567890123', 0xC0A80182, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00030', 15, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (167, '陶', '三三', 34, 'W', '16500165031', 'ID_CARD', '130101199007311234', '身份证照片链接31', 'iOS', '花艺师', '花店', 6000, 0, 0, 1, 15000, 0, '6200001234567890123', 0xC0A80183, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00031', 40, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (168, '姜', '三四', 36, 'M', '16400164032', 'ID_CARD', '120101198808311234', '身份证照片链接32', 'Android', '水管工', '维修公司', 7000, 1, 1, 0, 25000, 0, '6199991234567890123', 0xC0A80184, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00032', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (169, '戚', '三五', 29, 'W', '16300163033', 'ID_CARD', '110101199509301234', '身份证照片链接33', 'iOS', '图书管理员', '图书馆', 5500, 0, 0, 0, 9000, 0, '6198881234567890123', 0xC0A80185, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00033', 10, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (170, '谢', '三六', 33, 'M', '16200162034', 'ID_CARD', '650101199110311234', '身份证照片链接34', 'Android', '快递站站长', '快递公司', 9000, 1, 0, 1, 32000, 0, '6197771234567890123', 0xC0A80186, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00034', 25, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (171, '邹', '三七', 26, 'W', '16100161035', 'ID_CARD', '640101199811301234', '身份证照片链接35', 'iOS', '前台', '酒店', 5000, 0, 1, 0, 11000, 0, '6196661234567890123', 0xC0A80187, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00035', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (172, '喻', '三八', 40, 'M', '16000160036', 'ID_CARD', '630101198412311234', '身份证照片链接36', 'Android', '木匠', '家具厂', 8000, 1, 1, 1, 45000, 0, '6195551234567890123', 0xC0A80188, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00036', 45, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (173, '柏', '三九', 30, 'W', '15900159037', 'ID_CARD', '620101199401311234', '身份证照片链接37', 'iOS', '调酒师', '酒吧', 7500, 0, 0, 1, 19000, 0, '6194441234567890123', 0xC0A80189, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00037', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (174, '水', '四十', 37, 'M', '15800158038', 'ID_CARD', '610101198702281234', '身份证照片链接38', 'Android', '保安队长', '安保公司', 8500, 1, 0, 0, 24000, 0, '6193331234567890123', 0xC0A8018A, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00038', 15, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (175, '窦', '四一', 28, 'W', '15700157039', 'ID_CARD', '600101199603311234', '身份证照片链接39', 'iOS', '美甲店长', '美甲连锁店', 9000, 0, 1, 0, 26000, 0, '6192221234567890123', 0xC0A8018B, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00039', 25, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (176, '章', '四二', 31, 'M', '15600156040', 'ID_CARD', '590101199304301234', '身份证照片链接40', 'Android', '快递员', '快递公司', 6000, 1, 0, 1, 17000, 0, '6191111234567890123', 0xC0A8018C, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00040', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (177, '云', '四三', 25, 'W', '15500155041', 'ID_CARD', '580101199905311234', '身份证照片链接41', 'iOS', '实习生', '科技公司', 4500, 0, 0, 0, 7000, 0, '6190001234567890123', 0xC0A8018D, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00041', 10, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (178, '苏', '四四', 39, 'M', '15400154042', 'ID_CARD', '570101198506301234', '身份证照片链接42', 'Android', '维修工', '家电维修', 7000, 1, 1, 0, 29000, 0, '6189991234567890123', 0xC0A8018E, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00042', 35, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (179, '潘', '四五', 32, 'W', '15300153043', 'ID_CARD', '560101199207311234', '身份证照片链接43', 'iOS', '会计', '企业', 8500, 0, 0, 1, 21000, 0, '6188881234567890123', 0xC0A8018F, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00043', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (180, '葛', '四六', 35, 'M', '15200152044', 'ID_CARD', '550101198908311234', '身份证照片链接44', 'Android', '工程师', '制造企业', 13000, 1, 1, 1, 50000, 0, '6187771234567890123', 0xC0A80190, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00044', 40, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (181, '奚', '四七', 27, 'W', '15100151045', 'ID_CARD', '540101199709301234', '身份证照片链接45', 'iOS', '护士', '医院', 6500, 0, 1, 0, 13000, 0, '6186661234567890123', 0xC0A80191, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00045', 15, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (182, '范', '四八', 30, 'M', '15000150046', 'ID_CARD', '530101199410311234', '身份证照片链接46', 'Android', '厨师长', '餐厅', 9500, 1, 0, 0, 23000, 0, '6185551234567890123', 0xC0A80192, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00046', 25, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (183, '彭', '四九', 29, 'W', '14900149047', 'ID_CARD', '520101199511301234', '身份证照片链接47', 'iOS', '客服主管', '电商公司', 9000, 0, 0, 1, 27000, 0, '6184441234567890123', 0xC0A80193, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00047', 20, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (184, '郎', '五十', 34, 'M', '14800148048', 'ID_CARD', '510101199012311234', '身份证照片链接48', 'Android', '司机', '网约车公司', 8000, 1, 0, 1, 30000, 0, '6183331234567890123', 0xC0A80194, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00048', 30, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (185, '鲁', '五一', 26, 'W', '14700147049', 'ID_CARD', '500101199801311234', '身份证照片链接49', 'iOS', '设计师', '设计公司', 10000, 0, 1, 0, 33000, 0, '6182221234567890123', 0xC0A80195, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, 1, 'REF00049', 45, 0, NULL, NULL);
INSERT INTO `loan_baseinfo` (`id`, `first_name`, `second_name`, `age`, `gender`, `mobile`, `id_type`, `id_number`, `id_card`, `operator`, `work`, `company`, `salary`, `marital_status`, `has_house`, `has_car`, `application_amount`, `audit_status`, `bank_no`, `client_ip`, `created_at`, `updated_at`, `deleted_at`, `referrer_user_id`, `ref_code`, `loan_days`, `risk_list_status`, `risk_list_reason`, `risk_list_marked_at`) VALUES (186, '韦', '五二', 38, 'M', '14600146050', 'ID_CARD', '490101198602281234', '身份证照片链接50', 'Android', '木工', '装修公司', 7500, 1, 1, 0, 36000, 0, '6181111234567890123', 0xC0A80196, '2026-01-17 14:23:45', '2026-01-17 14:23:45', NULL, NULL, 'REF00050', 15, 0, NULL, NULL);
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
  `admin_user_id` bigint DEFAULT NULL,
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
  `disburse_amount` bigint NOT NULL COMMENT '放款金额(单位按你的系统：元/分，建议统一)',
  `net_amount` bigint NOT NULL COMMENT '到账金额(扣除费用后实际到账)',
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
  CONSTRAINT `fk_disburse_payout_channel` FOREIGN KEY (`payout_channel_id`) REFERENCES `loan_payment_channels` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='放款单/待放款任务表(审核通过后生成，状态待放款->已放款)';

-- ----------------------------
-- Records of loan_disbursements
-- ----------------------------
BEGIN;
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 4, 50000, 30000, 1, 0, 1, '2026-02-13 22:36:57', 2, 'PO20260213223657000', '2026-02-13 22:36:57', '2026-02-13 22:36:57', '2026-02-13 22:36:57', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, 1, 100000, 60000, 1, 0, 1, '2026-02-14 14:45:09', 2, 'PO20260214144508000', '2026-02-14 14:45:09', '2026-02-14 14:45:09', '2026-02-14 14:45:09', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 9, 150000, 120000, 1, 0, 1, '2026-02-14 16:02:47', 1, 'PO20260214160247000', '2026-02-14 16:02:47', '2026-02-14 16:02:47', '2026-02-14 16:02:47', NULL);
INSERT INTO `loan_disbursements` (`id`, `baseinfo_id`, `disburse_amount`, `net_amount`, `status`, `source_referrer_user_id`, `auditor_user_id`, `audited_at`, `payout_channel_id`, `payout_order_no`, `disbursed_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (17, 7, 120000, 96000, 1, 0, 1, '2026-02-24 14:28:03', 1, 'PO20260224142802000', '2026-02-24 14:28:03', '2026-02-24 14:28:03', '2026-02-24 14:28:03', NULL);
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
INSERT INTO `loan_mfa_devices` (`id`, `user_id`, `type`, `name`, `secret_enc`, `is_primary`, `status`, `last_used_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (16, 1, 'TOTP', 'Google Authenticator', 0xA5FA8ADFDBF1B37BC96D880DB015C37B17C350D955D95F9368D5341DE9E1DD22F38B7916614B572CA83EEF39A8088AB59515DF1AE5D26B4644F95440, 1, 1, '2026-02-25 20:14:53', '2026-01-16 14:04:44', '2026-02-25 20:14:53', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_mfa_recovery_codes
-- ----------------------------
DROP TABLE IF EXISTS `loan_mfa_recovery_codes`;
CREATE TABLE `loan_mfa_recovery_codes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `code_hash` varbinary(64) NOT NULL,
  `used_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_rc_user` (`user_id`),
  CONSTRAINT `loan_mfa_recovery_codes_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `loan_users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_mfa_recovery_codes
-- ----------------------------
BEGIN;
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
  `payout_fee_rate` float DEFAULT NULL COMMENT '代付手续费率(如0.003500=0.35%)',
  `payout_fee_fixed` int DEFAULT NULL COMMENT '代付固定手续费(分，若不用可为空)',
  `collect_fee_rate` float DEFAULT NULL COMMENT '代收手续费率',
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
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='统一支付渠道配置表(支持代付放款+代收回款，可禁用/启用，含手续费/限额/结算周期)';

-- ----------------------------
-- Records of loan_payment_channels
-- ----------------------------
BEGIN;
INSERT INTO `loan_payment_channels` (`id`, `code`, `name`, `merchant_no`, `status`, `can_payout`, `can_collect`, `payout_fee_rate`, `payout_fee_fixed`, `collect_fee_rate`, `collect_fee_fixed`, `collect_min_amount`, `collect_max_amount`, `payout_min_amount`, `payout_max_amount`, `settlement_cycle`, `settlement_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'BANK_A', '银行A代付代收', '2020202020202', 1, 1, 1, 20, 0, 0.002, 0, 1000, 20000000, 1000, 50000000, 'T1', '默认T+1结算', '2026-01-14 19:38:20', '2026-01-14 19:38:20', NULL);
INSERT INTO `loan_payment_channels` (`id`, `code`, `name`, `merchant_no`, `status`, `can_payout`, `can_collect`, `payout_fee_rate`, `payout_fee_fixed`, `collect_fee_rate`, `collect_fee_fixed`, `collect_min_amount`, `collect_max_amount`, `payout_min_amount`, `payout_max_amount`, `settlement_cycle`, `settlement_desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 'WALLET_X', '钱包X代付代收', '2020202020202', 1, 1, 1, 40, 100, 0.0025, 50, 1000, 10000000, 1000, 30000000, 'T0', '默认T+0结算', '2026-01-14 19:38:20', '2026-01-14 19:38:20', NULL);
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
  `principal_due` bigint NOT NULL DEFAULT '0' COMMENT '应还本金(建议统一单位：分)',
  `interest_due` bigint NOT NULL DEFAULT '0' COMMENT '应还利息(分)',
  `fee_due` bigint NOT NULL DEFAULT '0' COMMENT '应还费用(分)',
  `penalty_due` int NOT NULL DEFAULT '0' COMMENT '应还罚息(分，逾期产生)',
  `total_due` bigint NOT NULL COMMENT '本期应还总额=本金+利息+费用+罚息(分)',
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
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='还款计划表(支持单期/分期，逾期/已还通过状态体现)';

-- ----------------------------
-- Records of loan_repayment_schedules
-- ----------------------------
BEGIN;
INSERT INTO `loan_repayment_schedules` (`id`, `disbursement_id`, `installment_no`, `due_date`, `principal_due`, `interest_due`, `fee_due`, `penalty_due`, `total_due`, `paid_principal`, `paid_interest`, `paid_fee`, `paid_penalty`, `paid_total`, `status`, `last_paid_at`, `settled_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 14, 1, '2026-02-27', 50000, 0, 50000, 0, 50000, 0, 0, 0, 0, 20000, 0, NULL, NULL, '2026-02-13 22:36:57', '2026-02-24 14:25:04', NULL);
INSERT INTO `loan_repayment_schedules` (`id`, `disbursement_id`, `installment_no`, `due_date`, `principal_due`, `interest_due`, `fee_due`, `penalty_due`, `total_due`, `paid_principal`, `paid_interest`, `paid_fee`, `paid_penalty`, `paid_total`, `status`, `last_paid_at`, `settled_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 15, 1, '2026-02-28', 100000, 0, 100000, 0, 100000, 0, 0, 0, 0, 0, 0, NULL, NULL, '2026-02-14 14:45:09', '2026-02-14 14:45:09', NULL);
INSERT INTO `loan_repayment_schedules` (`id`, `disbursement_id`, `installment_no`, `due_date`, `principal_due`, `interest_due`, `fee_due`, `penalty_due`, `total_due`, `paid_principal`, `paid_interest`, `paid_fee`, `paid_penalty`, `paid_total`, `status`, `last_paid_at`, `settled_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (6, 16, 1, '2026-04-15', 150000, 0, 150000, 0, 150000, 0, 0, 0, 0, 0, 0, NULL, NULL, '2026-02-14 16:02:47', '2026-02-14 16:02:47', NULL);
INSERT INTO `loan_repayment_schedules` (`id`, `disbursement_id`, `installment_no`, `due_date`, `principal_due`, `interest_due`, `fee_due`, `penalty_due`, `total_due`, `paid_principal`, `paid_interest`, `paid_fee`, `paid_penalty`, `paid_total`, `status`, `last_paid_at`, `settled_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (7, 17, 1, '2026-03-26', 120000, 0, 120000, 0, 120000, 0, 0, 0, 0, 20000, 0, NULL, NULL, '2026-02-24 14:28:03', '2026-02-24 14:29:26', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_repayment_transactions
-- ----------------------------
DROP TABLE IF EXISTS `loan_repayment_transactions`;
CREATE TABLE `loan_repayment_transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '回款流水ID',
  `schedule_id` bigint NOT NULL COMMENT '关联期次 loan_repayment_schedules.id(可空：先入账后分配/未分期)',
  `collect_channel_id` bigint DEFAULT NULL COMMENT '回款渠道(代收) loan_payment_channels.id',
  `collect_order_no` varchar(128) DEFAULT NULL COMMENT '回款订单号/三方代收单号(商户单号)',
  `pay_ref` varchar(128) DEFAULT NULL COMMENT '支付渠道流水号/交易号(三方transaction id)',
  `pay_amount` int NOT NULL COMMENT '本次回款金额(分)',
  `pay_method` varchar(32) DEFAULT NULL COMMENT '回款方式(如 BANK_TRANSFER/CARD/WALLET/CASH)',
  `paid_at` datetime DEFAULT NULL COMMENT '回款时间(交易成功时间)',
  `alloc_principal` int NOT NULL DEFAULT '0' COMMENT '本次分配到本金(分)',
  `alloc_interest` int NOT NULL DEFAULT '0' COMMENT '本次分配到利息(分)',
  `alloc_fee` int NOT NULL DEFAULT '0' COMMENT '本次分配到费用(分)',
  `alloc_penalty` int NOT NULL DEFAULT '0' COMMENT '本次分配到罚息(分)',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '流水状态：1成功 0失败 2冲正/撤销',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '备注',
  `created_by` bigint NOT NULL,
  `voucher_file_name` varchar(64) NOT NULL,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间(NULL未删除)',
  PRIMARY KEY (`id`),
  KEY `idx_tx_disburse_time` (`paid_at`) COMMENT '按放款单/时间查回款流水',
  KEY `idx_tx_schedule` (`schedule_id`) COMMENT '按期次查回款流水',
  KEY `idx_tx_collect_channel_time` (`collect_channel_id`,`paid_at`) COMMENT '按回款渠道/时间查流水',
  KEY `idx_tx_collect_order_no` (`collect_order_no`) COMMENT '按回款订单号查询',
  KEY `idx_tx_pay_ref` (`pay_ref`) COMMENT '按交易流水号查询',
  CONSTRAINT `fk_tx_collect_channel` FOREIGN KEY (`collect_channel_id`) REFERENCES `loan_payment_channels` (`id`),
  CONSTRAINT `fk_tx_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `loan_repayment_schedules` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='回款流水表(记录每次实际回款，支持分期/部分还款，含回款渠道与订单号)';

-- ----------------------------
-- Records of loan_repayment_transactions
-- ----------------------------
BEGIN;
INSERT INTO `loan_repayment_transactions` (`id`, `schedule_id`, `collect_channel_id`, `collect_order_no`, `pay_ref`, `pay_amount`, `pay_method`, `paid_at`, `alloc_principal`, `alloc_interest`, `alloc_fee`, `alloc_penalty`, `status`, `remark`, `created_by`, `voucher_file_name`, `created_at`, `updated_at`, `deleted_at`) VALUES (14, 4, 2, 'PI20260224142504000', '', 20000, 'IMPORT', NULL, 20000, 0, 0, 0, 1, '测试回款', 1, '11aa9ba4-dfaa-4497-a61f-0e0d2500e8e8.png', '2026-02-24 14:25:04', '2026-02-24 14:25:04', NULL);
INSERT INTO `loan_repayment_transactions` (`id`, `schedule_id`, `collect_channel_id`, `collect_order_no`, `pay_ref`, `pay_amount`, `pay_method`, `paid_at`, `alloc_principal`, `alloc_interest`, `alloc_fee`, `alloc_penalty`, `status`, `remark`, `created_by`, `voucher_file_name`, `created_at`, `updated_at`, `deleted_at`) VALUES (15, 7, 2, 'PI20260224142925000', '', 20000, 'IMPORT', NULL, 20000, 0, 0, 0, 1, '用户还款200元', 1, '2c1c1edf-8935-4780-afde-7ffad082b13f.png', '2026-02-24 14:29:26', '2026-02-24 14:29:26', NULL);
COMMIT;

-- ----------------------------
-- Table structure for loan_risk_customer
-- ----------------------------
DROP TABLE IF EXISTS `loan_risk_customer`;
CREATE TABLE `loan_risk_customer` (
  `id` int NOT NULL AUTO_INCREMENT,
  `loan_baseinfo_id` int DEFAULT NULL,
  `risk_type` tinyint DEFAULT NULL COMMENT '风险类型 -1 黑名单 1 白名单',
  `risk_reason` varchar(255) DEFAULT NULL COMMENT '风险原因',
  `created_by` int DEFAULT NULL COMMENT 'loan_users_id',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_risk_customer
-- ----------------------------
BEGIN;
INSERT INTO `loan_risk_customer` (`id`, `loan_baseinfo_id`, `risk_type`, `risk_reason`, `created_by`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 1, 1, '征信良好 纳入白名单', 1, '2026-02-06 14:07:07', '2026-02-06 14:07:10', NULL);
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
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of loan_users
-- ----------------------------
BEGIN;
INSERT INTO `loan_users` (`id`, `username`, `password_hash`, `department_id`, `mfa_enabled`, `mfa_required`, `status`, `created_at`, `updated_at`, `deleted_at`, `share_code`) VALUES (1, 'axbros', '$2a$10$Y0Wr7iRD2Z1xxNwV7uDAXuAhG6ECseNjUcjdnpwSdVcE1s3I/duKy', 1, 1, 0, 1, '2026-01-14 17:51:53', '2026-01-14 17:51:53', NULL, 'MvJHGT7XL-cp');
INSERT INTO `loan_users` (`id`, `username`, `password_hash`, `department_id`, `mfa_enabled`, `mfa_required`, `status`, `created_at`, `updated_at`, `deleted_at`, `share_code`) VALUES (2, 'auditor', '$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5lWlP0r8kP7r9v8pJwD3h8m6cK4QK', 1, 0, 0, 1, '2026-01-14 19:36:30', '2026-01-14 19:36:30', NULL, 'REFAUDIT01');
INSERT INTO `loan_users` (`id`, `username`, `password_hash`, `department_id`, `mfa_enabled`, `mfa_required`, `status`, `created_at`, `updated_at`, `deleted_at`, `share_code`) VALUES (3, 'referrer', '$2a$10$7EqJtq98hPqEX7fNZaFWoOhi5lWlP0r8kP7r9v8pJwD3h8m6cK4QK', 1, 0, 0, 1, '2026-01-14 19:36:30', '2026-01-14 19:36:30', NULL, 'REFSHARE01');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
