CREATE TABLE IF NOT EXISTS `terminal_connection` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `app_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'app id',
    `conn_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'connection id',
    `addr` varchar(32) NOT NULL DEFAULT '' COMMENT 'frontier ipv4 addr',
    `addr6` varchar(32) NOT NULL DEFAULT '' COMMENT 'frontier ipv6 addr',
    `state`tinyint(4) NOT NULL DEFAULT '-1' COMMENT '0 inactive 1 active',
    `device_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'device id',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    KEY `idx_conn_id` (`conn_id`),
    KEY `idx_device_id` (`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT'terminal connection information';

CREATE TABLE IF NOT EXISTS `conference_user` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `app_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'app id',
    `user_id` varchar(32) BINARY NOT NULL DEFAULT '' COMMENT 'user id',
    `room_id` varchar(32) BINARY NOT NULL DEFAULT '' COMMENT 'room id',
    `device_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'device id',
    `conn_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'connection id',
    `state`tinyint(4) NOT NULL DEFAULT '1' COMMENT '0 inactive 1 active 2 reconnecting',
    `is_host` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'host state',
    `is_mic_on` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'mic state',
    `is_camera_on` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'camera state',
    `is_sharing` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'screen share state',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    KEY `idx_conn_id` (`conn_id`),
    KEY `idx_room_id_user_id` (`room_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT'conference user information';

CREATE TABLE IF NOT EXISTS `conference_room` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `app_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'app id',
    `room_id` varchar(32) BINARY NOT NULL DEFAULT '' COMMENT 'room id',
    `state` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '0 inactive 1 active',
    `record` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'record state',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
	UNIQUE KEY `uniq_room_id` (`room_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT'conference room information';

CREATE TABLE IF NOT EXISTS `conference_video_record` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `room_id` varchar(32) BINARY NOT NULL DEFAULT '' COMMENT 'room id',
    `app_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'app id',
	`vid` varchar(128) NOT NULL DEFAULT '' COMMENT 'video id',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    KEY `idx_room_id` (`room_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT'conference record information';

CREATE TABLE `cs_meeting_user` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `app_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'app id',
    `user_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'user id',
    `user_name` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'user name',
    `user_status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 audience 1 raise_hands 2 on_microphone',
    `room_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'room id',
    `device_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'device id',
    `conn_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'connection id',
    `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0 inactive 1 active 2 reconnecting',
    `is_host` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'host state',
    `is_mic_on` tinyint(1) NOT NULL DEFAULT '1' COMMENT 'mic state',
    `raise_hands_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'raise hands time',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    KEY `idx_conn_id` (`conn_id`),
    KEY `idx_room_id_user_id` (`room_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1243 DEFAULT CHARSET=utf8 COMMENT='chat salon meeting user information'


CREATE TABLE `cs_meeting_room` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `app_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'app id',
    `room_id` varchar(32) NOT NULL DEFAULT '' COMMENT 'room id',
    `room_name` varchar(128) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'room name',
    `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0 inactive 1 active',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_room_id` (`room_id`)
) ENGINE=InnoDB AUTO_INCREMENT=619 DEFAULT CHARSET=utf8 COMMENT='chat salon meeting room information'

CREATE TABLE `edu_user_room_info` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '????????????',
    `app_id` varchar(100) DEFAULT NULL COMMENT '?????? ID???????????????RTC appid',
    `room_id` varchar(100) DEFAULT NULL COMMENT '??????/?????? ID???????????????RTC roomid',
    `user_id` varchar(255) DEFAULT NULL COMMENT '?????? ID???????????????RTC userid',
    `user_name` varchar(255) DEFAULT NULL COMMENT '?????? ??????',
    `user_role` int(11) DEFAULT NULL COMMENT '??????????????? 0????????????1?????????',
    `user_status` int(11) DEFAULT NULL COMMENT '???????????????0????????????1?????????',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '??????/?????? ????????????',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '????????????',
    `join_time` bigint(20) DEFAULT NULL COMMENT '???????????????????????????',
    `leave_time` bigint(20) DEFAULT NULL COMMENT '???????????????????????????',
    `is_mic_on` tinyint(4) DEFAULT '0' COMMENT '???????????????',
    `is_camera_on` tinyint(4) DEFAULT '0' COMMENT '???????????????',
    `is_hands_up` tinyint(4) DEFAULT '0' COMMENT '????????????',
    `group_speech_join_rtc` tinyint(4) DEFAULT '0' COMMENT '??????????????????????????????????????????RTC??????',
    `rtc_token` varchar(255) DEFAULT NULL COMMENT '????????????token',
    `conn_id` varchar(255) DEFAULT NULL COMMENT 'socketID',
    `parent_room_id` varchar(100) DEFAULT NULL COMMENT '????????????',
    `device_id` varchar(128) DEFAULT NULL COMMENT 'device_id',
    `room_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0:????????????1:?????????',
    `is_interact` tinyint(4) DEFAULT NULL COMMENT '????????????????????????',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_room_id_user_id` (`room_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6104 DEFAULT CHARSET=utf8 COMMENT='??????-??????????????????'

CREATE TABLE `edu_room_info` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '????????????',
    `app_id` varchar(100) DEFAULT NULL COMMENT '?????? ID???????????????RTC appid',
    `room_id` varchar(100) DEFAULT NULL COMMENT '??????/?????? ID???????????????RTC roomid',
    `room_name` varchar(200) DEFAULT NULL COMMENT '????????????',
    `room_type` int(11) DEFAULT '0' COMMENT '???????????????0???????????????1??????????????????',
    `create_user_id` varchar(100) DEFAULT NULL COMMENT '??????/?????? ???????????????ID???????????????RTC userid',
    `status` int(11) DEFAULT NULL COMMENT '???????????????0???????????????1???????????????2????????????',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '??????/?????? ????????????',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '????????????',
    `begin_class_time` bigint(20) DEFAULT NULL COMMENT '???????????????',
    `end_class_time` bigint(20) DEFAULT NULL COMMENT '???????????????',
    `audio_mute_all` tinyint(4) DEFAULT '1' COMMENT '???????????????',
    `video_mute_all` tinyint(4) DEFAULT '1' COMMENT '???????????????',
    `enable_group_speech` tinyint(4) DEFAULT '0' COMMENT '????????????????????????',
    `enable_interactive` tinyint(4) DEFAULT '0' COMMENT '????????????????????????',
    `teacher_name` varchar(100) DEFAULT NULL COMMENT '????????????',
    `begin_class_time_real` bigint(20) DEFAULT NULL COMMENT '???????????????????????????',
    `token` varchar(300) DEFAULT NULL COMMENT 'rtc??????token',
    `group_limit` int(11) DEFAULT '0' COMMENT '????????????????????????',
    `is_recording` tinyint(4) DEFAULT '0' COMMENT '??????????????????',
    `group_num` int(8) DEFAULT '0' COMMENT '???????????????',
    `end_class_time_real` bigint(20) DEFAULT NULL COMMENT '????????????????????????',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_room_id` (`room_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1629 DEFAULT CHARSET=utf8 COMMENT='??????/????????????'

CREATE TABLE `edu_record_info` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '????????????',
    `app_id` varchar(100) DEFAULT NULL COMMENT '?????? ID???????????????RTC appid',
    `room_id` varchar(100) DEFAULT NULL COMMENT '??????/?????? ID???????????????RTC roomid',
    `user_id` varchar(100) DEFAULT NULL COMMENT '?????? ID???????????????RTC userid',
    `room_name` varchar(200) DEFAULT NULL COMMENT '????????????',
    `record_status` int(11) DEFAULT NULL COMMENT '?????????????????????0???????????????1??????????????????2??????????????????3????????????',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????????????',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '????????????',
    `record_begin_time` bigint(20) DEFAULT NULL COMMENT '?????????????????????',
    `record_end_time` bigint(20) DEFAULT NULL COMMENT '?????????????????????',
    `vid` varchar(200) DEFAULT '' COMMENT '???????????? vid',
    `parent_room_id` varchar(100) DEFAULT NULL COMMENT '????????????',
    `task_id` varchar(100) DEFAULT NULL COMMENT '????????????id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2634 DEFAULT CHARSET=utf8 COMMENT='??????/??????????????????'