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
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `app_id` varchar(100) DEFAULT NULL COMMENT '应用 ID，同时也是RTC appid',
    `room_id` varchar(100) DEFAULT NULL COMMENT '房间/课堂 ID，同时也是RTC roomid',
    `user_id` varchar(255) DEFAULT NULL COMMENT '用户 ID，同时也是RTC userid',
    `user_name` varchar(255) DEFAULT NULL COMMENT '用户 昵称',
    `user_role` int(11) DEFAULT NULL COMMENT '用户角色， 0：老师，1：学生',
    `user_status` int(11) DEFAULT NULL COMMENT '用户状态，0：在线，1：离线',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '房间/课堂 创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `join_time` bigint(20) DEFAULT NULL COMMENT '最后一次加入时间戳',
    `leave_time` bigint(20) DEFAULT NULL COMMENT '最后一次离开时间戳',
    `is_mic_on` tinyint(4) DEFAULT '0' COMMENT '麦克风状态',
    `is_camera_on` tinyint(4) DEFAULT '0' COMMENT '摄像头状态',
    `is_hands_up` tinyint(4) DEFAULT '0' COMMENT '是否举手',
    `group_speech_join_rtc` tinyint(4) DEFAULT '0' COMMENT '集体发言时是否被选中参与加入RTC房间',
    `rtc_token` varchar(255) DEFAULT NULL COMMENT '进房用的token',
    `conn_id` varchar(255) DEFAULT NULL COMMENT 'socketID',
    `parent_room_id` varchar(100) DEFAULT NULL COMMENT '父房间号',
    `device_id` varchar(128) DEFAULT NULL COMMENT 'device_id',
    `room_type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0:大班课，1:小组课',
    `is_interact` tinyint(4) DEFAULT NULL COMMENT '是否正在视频互动',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_room_id_user_id` (`room_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6104 DEFAULT CHARSET=utf8 COMMENT='房间-用户关联信息'

CREATE TABLE `edu_room_info` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `app_id` varchar(100) DEFAULT NULL COMMENT '应用 ID，同时也是RTC appid',
    `room_id` varchar(100) DEFAULT NULL COMMENT '房间/课堂 ID，同时也是RTC roomid',
    `room_name` varchar(200) DEFAULT NULL COMMENT '课堂名称',
    `room_type` int(11) DEFAULT '0' COMMENT '课堂类型，0：大班课，1：大班小组课',
    `create_user_id` varchar(100) DEFAULT NULL COMMENT '房间/课堂 创建者用户ID，同时也是RTC userid',
    `status` int(11) DEFAULT NULL COMMENT '课堂状态：0：未开始，1：上课中，2：已结束',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '房间/课堂 创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `begin_class_time` bigint(20) DEFAULT NULL COMMENT '上课时间戳',
    `end_class_time` bigint(20) DEFAULT NULL COMMENT '下课时间戳',
    `audio_mute_all` tinyint(4) DEFAULT '1' COMMENT '麦克风状态',
    `video_mute_all` tinyint(4) DEFAULT '1' COMMENT '摄像头状态',
    `enable_group_speech` tinyint(4) DEFAULT '0' COMMENT '集体发言是否打开',
    `enable_interactive` tinyint(4) DEFAULT '0' COMMENT '互动视频是否打开',
    `teacher_name` varchar(100) DEFAULT NULL COMMENT '老师姓名',
    `begin_class_time_real` bigint(20) DEFAULT NULL COMMENT '真实开始上课时间戳',
    `token` varchar(300) DEFAULT NULL COMMENT 'rtc推流token',
    `group_limit` int(11) DEFAULT '0' COMMENT '小组房间最大人数',
    `is_recording` tinyint(4) DEFAULT '0' COMMENT '是否正在录制',
    `group_num` int(8) DEFAULT '0' COMMENT '小组房间数',
    `end_class_time_real` bigint(20) DEFAULT NULL COMMENT '真实结束课堂时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_room_id` (`room_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1629 DEFAULT CHARSET=utf8 COMMENT='房间/课堂信息'

CREATE TABLE `edu_record_info` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `app_id` varchar(100) DEFAULT NULL COMMENT '应用 ID，同时也是RTC appid',
    `room_id` varchar(100) DEFAULT NULL COMMENT '房间/课堂 ID，同时也是RTC roomid',
    `user_id` varchar(100) DEFAULT NULL COMMENT '用户 ID，同时也是RTC userid',
    `room_name` varchar(200) DEFAULT NULL COMMENT '课堂名称',
    `record_status` int(11) DEFAULT NULL COMMENT '录制视频状态：0：录制中，1：录制成功，2：录制失败，3：已删除',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `record_begin_time` bigint(20) DEFAULT NULL COMMENT '开始录制时间戳',
    `record_end_time` bigint(20) DEFAULT NULL COMMENT '结束录制时间戳',
    `vid` varchar(200) DEFAULT '' COMMENT '点播平台 vid',
    `parent_room_id` varchar(100) DEFAULT NULL COMMENT '父房间号',
    `task_id` varchar(100) DEFAULT NULL COMMENT '录制任务id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2634 DEFAULT CHARSET=utf8 COMMENT='房间/课堂录制信息'