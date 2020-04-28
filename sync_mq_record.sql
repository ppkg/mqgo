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

 Date: 28/04/2020 15:37:43
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sync_mq_record
-- ----------------------------
DROP TABLE IF EXISTS `sync_mq_record`;
CREATE TABLE `sync_mq_record`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sync_mq_info_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '消息ID',
  `queue` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '队列名称',
  `create_date` datetime(0) NULL DEFAULT current_timestamp(0),
  `exchange` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '交换机名称',
  `route_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由名称',
  `response` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '响应结果',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
