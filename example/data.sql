CREATE TABLE `company` (
                           `id` bigint unsigned NOT NULL COMMENT '授权方企业本地id',
                           `platform` varchar(10) NOT NULL DEFAULT '' COMMENT '授权方企业的平台方',
                           `corp_id` varchar(100) NOT NULL DEFAULT '' COMMENT '平台授权企业id',
                           `corp_type` tinyint NOT NULL DEFAULT '0' COMMENT '企业类型',
                           `full_corp_name` varchar(100) NOT NULL DEFAULT '' COMMENT '企业全称',
                           `corp_type2` tinyint NOT NULL DEFAULT '0' COMMENT '0 是普通组织\n1 是项目\n2是圈子\n3没有业务表现形式\n4是自建班级群\n10是敏捷组织\n11是培训群敏捷组织',
                           `corp_name` varchar(45) NOT NULL DEFAULT '' COMMENT '企业简称',
                           `industry` varchar(45) NOT NULL DEFAULT '' COMMENT '行业类型',
                           `is_authenticated` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否认证',
                           `license_code` varchar(45) NOT NULL DEFAULT '' COMMENT '序列号',
                           `corp_logo_url` varchar(145) NOT NULL DEFAULT '' COMMENT '企业logo',
                           `invite_url` varchar(145) NOT NULL DEFAULT '' COMMENT '企业邀请链接',
                           `invite_code` varchar(45) NOT NULL DEFAULT '' COMMENT '邀请码，只有自己邀请的企业才会返回邀请码，可用该邀请码统计不同渠道的拉新，否则值为空字符串',
                           `is_ecological_corp` tinyint(1) NOT NULL DEFAULT '0',
                           `auth_level` tinyint NOT NULL DEFAULT '0' COMMENT '企业认证等级：\n\n0：未认证\n1：高级认证\n2：中级认证\n3：初级认证',
                           `auth_channel` varchar(45) NOT NULL DEFAULT '' COMMENT '渠道码',
                           `auth_channel_type` varchar(45) NOT NULL DEFAULT '' COMMENT '渠道类型。为了避免渠道码重复，可与渠道码共同确认渠道。可能为空，非空时当前只有满天星类型，值为STAR_ACTIVITY',
                           `state` decimal(5,2) NOT NULL DEFAULT '0' COMMENT '系统状态：-1：已删除；0：禁用；1：可用',
                           `state_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '系统状态时间',
                           `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `corp_id_idx` (`corp_id`,`platform`) USING BTREE /*!80000 INVISIBLE */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='授权方企业';

CREATE TABLE `biz` (
                           `id` bigint unsigned NOT NULL COMMENT '授权方企业本地id',
                           `platform` varchar(10) NOT NULL DEFAULT '' COMMENT '授权方企业的平台方',
                           `corp_id` varchar(100) NOT NULL DEFAULT '' COMMENT '平台授权企业id',
                           `corp_type` tinyint NOT NULL DEFAULT '0' COMMENT '企业类型',
                           `full_corp_name` varchar(100) NOT NULL DEFAULT '' COMMENT '企业全称',
                           `corp_type2` tinyint NOT NULL DEFAULT '0' COMMENT '0 是普通组织\n1 是项目\n2是圈子\n3没有业务表现形式\n4是自建班级群\n10是敏捷组织\n11是培训群敏捷组织',
                           `corp_name` varchar(45) NOT NULL DEFAULT '' COMMENT '企业简称',
                           `industry` varchar(45) NOT NULL DEFAULT '' COMMENT '行业类型',
                           `is_authenticated` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否认证',
                           `license_code` varchar(45) NOT NULL DEFAULT '' COMMENT '序列号',
                           `corp_logo_url` varchar(145) NOT NULL DEFAULT '' COMMENT '企业logo',
                           `invite_url` varchar(145) NOT NULL DEFAULT '' COMMENT '企业邀请链接',
                           `invite_code` varchar(45) NOT NULL DEFAULT '' COMMENT '邀请码，只有自己邀请的企业才会返回邀请码，可用该邀请码统计不同渠道的拉新，否则值为空字符串',
                           `is_ecological_corp` tinyint(1) NOT NULL DEFAULT '0',
                           `auth_level` tinyint NOT NULL DEFAULT '0' COMMENT '企业认证等级：\n\n0：未认证\n1：高级认证\n2：中级认证\n3：初级认证',
                           `auth_channel` varchar(45) NOT NULL DEFAULT '' COMMENT '渠道码',
                           `auth_channel_type` varchar(45) NOT NULL DEFAULT '' COMMENT '渠道类型。为了避免渠道码重复，可与渠道码共同确认渠道。可能为空，非空时当前只有满天星类型，值为STAR_ACTIVITY',
                           `state` decimal(5,2) NOT NULL DEFAULT '0' COMMENT '系统状态：-1：已删除；0：禁用；1：可用',
                           `state_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '系统状态时间',
                           `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `corp_id_idx` (`corp_id`,`platform`) USING BTREE /*!80000 INVISIBLE */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='授权方企业';
