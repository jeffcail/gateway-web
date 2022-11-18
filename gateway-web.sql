/*
 Navicat MySQL Data Transfer

 Source Server         : my-mysql
 Source Server Type    : MySQL
 Source Server Version : 50737
 Source Host           : localhost:3306
 Source Schema         : gateway-web

 Target Server Type    : MySQL
 Target Server Version : 50737
 File Encoding         : 65001

 Date: 18/11/2022 18:38:59
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for area
-- ----------------------------
DROP TABLE IF EXISTS `area`;
CREATE TABLE `area` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `area_name` varchar(255) NOT NULL,
  `city_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `update_at` datetime NOT NULL,
  `create_at` datetime NOT NULL,
  `delete_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='area';

-- ----------------------------
-- Records of area
-- ----------------------------
BEGIN;
INSERT INTO `area` VALUES (2, 'area_name', 1, 2, '2019-06-15 00:00:00', '2019-06-15 00:00:00', '2019-06-15 00:00:00');
COMMIT;

-- ----------------------------
-- Table structure for gateway-service_tcp_rule
-- ----------------------------
DROP TABLE IF EXISTS `gateway-service_tcp_rule`;
CREATE TABLE `gateway-service_tcp_rule` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `service_id` bigint(20) NOT NULL COMMENT '服务id',
  `port` int(11) NOT NULL COMMENT '端口',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Tcp网关路由匹配表';

-- ----------------------------
-- Table structure for gateway_admin
-- ----------------------------
DROP TABLE IF EXISTS `gateway_admin`;
CREATE TABLE `gateway_admin` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `user_name` varchar(255) NOT NULL COMMENT '用户名',
  `salt` varchar(50) NOT NULL COMMENT '加密盐值',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `update_at` datetime NOT NULL COMMENT '更新时间',
  `is_delete` tinyint(4) NOT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

-- ----------------------------
-- Records of gateway_admin
-- ----------------------------
BEGIN;
INSERT INTO `gateway_admin` VALUES (1, 'admin', 'admin', '2823d896e9822c0833d41d4904f0c00756d718570fce49b9a379a62c804689d3', '2020-04-10 16:42:05', '2020-04-21 06:35:08', 0);
COMMIT;

-- ----------------------------
-- Table structure for gateway_service_access_control
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_access_control`;
CREATE TABLE `gateway_service_access_control` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `service_id` bigint(20) DEFAULT NULL COMMENT '服务id',
  `open_auth` tinyint(4) DEFAULT NULL COMMENT '是否开启权限 1=开启',
  `black_list` varchar(2000) DEFAULT NULL COMMENT '黑名单ip',
  `white_list` varchar(2000) DEFAULT NULL COMMENT '白名单ip',
  `white_host_name` varchar(2000) DEFAULT NULL COMMENT '白名单主机',
  `clientip_flow_limit` int(11) DEFAULT NULL COMMENT '客户端ip限流',
  `service_flow_limit` int(20) DEFAULT NULL COMMENT '服务端限流',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='网关权限控制表';

-- ----------------------------
-- Table structure for gateway_service_grpc_rule
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_grpc_rule`;
CREATE TABLE `gateway_service_grpc_rule` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `service_id` bigint(20) NOT NULL COMMENT '服务Id',
  `port` int(11) NOT NULL COMMENT '端口',
  `header_transfor` varchar(5000) NOT NULL COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='grpc网关路由匹配表';

-- ----------------------------
-- Table structure for gateway_service_http_rule
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_http_rule`;
CREATE TABLE `gateway_service_http_rule` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `service_id` bigint(20) NOT NULL COMMENT '服务id',
  `rule_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '匹配类型 0=url前缀url_prefix 1=域名domain',
  `rule` varchar(255) NOT NULL COMMENT 'type=domain表示域名，type=url_prefix时表示url前缀',
  `need_https` tinyint(4) NOT NULL DEFAULT '0' COMMENT '支持https 1=支持',
  `need_strip_uri` tinyint(4) NOT NULL DEFAULT '0' COMMENT '启用strip_uri 1=启用',
  `need_websocket` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否支持websocket 1=支持',
  `url-rewrite` varchar(5000) NOT NULL COMMENT 'url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔',
  `header_transfor` varchar(5000) DEFAULT NULL COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='http网关路由匹配表';

-- ----------------------------
-- Table structure for gateway_service_info
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_info`;
CREATE TABLE `gateway_service_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `load_type` varchar(0) NOT NULL COMMENT '服务类型',
  `service_name` varchar(255) NOT NULL COMMENT '服务名称',
  `service_desc` varchar(255) NOT NULL COMMENT '服务描述',
  `create_at` varchar(255) NOT NULL COMMENT '添加时间',
  `update_at` varchar(255) NOT NULL COMMENT '更新时间',
  `is_delete` tinyint(4) NOT NULL COMMENT '是否删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='网关基本信息表';

-- ----------------------------
-- Table structure for gateway_service_load_balance
-- ----------------------------
DROP TABLE IF EXISTS `gateway_service_load_balance`;
CREATE TABLE `gateway_service_load_balance` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `service_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '服务id',
  `check_method` varchar(255) NOT NULL DEFAULT '0' COMMENT '检查方法 0=tcpchk,检测端口是否握手成功',
  `check_timeout` int(10) NOT NULL DEFAULT '0' COMMENT 'check超时时间,单位s',
  `check_interval` int(11) NOT NULL DEFAULT '0' COMMENT '检查间隔, 单位s',
  `round_type` tinyint(255) NOT NULL DEFAULT '2' COMMENT '轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash',
  `ip_list` varchar(2000) NOT NULL COMMENT 'ip列表',
  `weight_list` varchar(2000) NOT NULL COMMENT '权重列表',
  `forbid_list` varchar(2000) NOT NULL COMMENT '禁用ip列表',
  `upstream_connect_timeout` int(11) NOT NULL DEFAULT '0' COMMENT '建立连接超时, 单位s',
  `upstream_header_timeout` int(11) NOT NULL DEFAULT '0' COMMENT '获取header超时, 单位s',
  `upstream_idle_timeout` int(10) NOT NULL DEFAULT '0' COMMENT '链接最大空闲时间, 单位s',
  `upstream_max_idle` int(11) NOT NULL DEFAULT '0' COMMENT '最大空闲链接数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='网关负载表';

SET FOREIGN_KEY_CHECKS = 1;
