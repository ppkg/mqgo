/*
 Navicat Premium Data Transfer

 Source Server         : 10.1.1.245
 Source Server Type    : MySQL
 Source Server Version : 100213
 Source Host           : 10.1.1.245:3306
 Source Schema         : datacenter

 Target Server Type    : MySQL
 Target Server Version : 100213
 File Encoding         : 65001

 Date: 28/04/2020 15:37:33
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sync_mq_info
-- ----------------------------
DROP TABLE IF EXISTS `sync_mq_info`;
CREATE TABLE `sync_mq_info`  (
  `create_date` datetime(0) NULL DEFAULT current_timestamp(0),
  `id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `queue` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '队列名称',
  `third_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '第三方ID，会员ID或者宠物ID或者其它',
  `request` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '请求',
  `response` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '响应结果',
  `type` tinyint(4) NULL DEFAULT NULL COMMENT '数据类型，1用户数据，2宠物数据，3核销数据',
  `is_sync` bit(1) NULL DEFAULT b'0' COMMENT '是否已经同步',
  `retry_count` int(11) NULL DEFAULT 0 COMMENT '重试次数',
  `retry_datetime` datetime(0) NULL DEFAULT current_timestamp(0) COMMENT '重试时间',
  `platform_id` int(11) NULL DEFAULT NULL COMMENT ' 平台ID',
  `channel_id` int(11) NULL DEFAULT NULL COMMENT '渠道ID',
  `sync_date` datetime(0) NULL DEFAULT NULL COMMENT '同步时间',
  `exchange` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '交换机名称',
  `route_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由key',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `INDEX`(`is_sync`, `queue`, `type`) USING BTREE COMMENT '查询索引'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
