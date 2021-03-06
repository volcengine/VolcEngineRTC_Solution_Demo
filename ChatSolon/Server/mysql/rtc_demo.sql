-- MySQL dump 10.13  Distrib 5.7.33, for Linux (x86_64)
--
-- Host: localhost    Database: rtc_demo_db
-- ------------------------------------------------------
-- Server version	5.7.33


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `rtc_demo_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `rtc_demo_db`;

GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'bytedance';
FlUSH PRIVILEGES;


--
-- Table structure for table `user_profile`
--

DROP TABLE IF EXISTS `user_profile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_profile` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary key ID',
    `user_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'user id',
    `user_name` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'user name',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14115 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='user profile information';
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `cs_interact`
--

DROP TABLE IF EXISTS `cs_interact`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cs_interact` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '????????????',
  `interact_id` varchar(100) DEFAULT NULL COMMENT 'interact_ID',
  `owner_room_id` varchar(100) DEFAULT NULL COMMENT 'room_id',
  `owner_user_id` varchar(100) DEFAULT NULL COMMENT 'user_d',
  `rtc_app_id` varchar(100) DEFAULT NULL COMMENT 'app_id',
  `rtc_room_id` varchar(100) DEFAULT NULL COMMENT 'room_id',
  `interact_type` int(11) DEFAULT NULL COMMENT '????????????',
  `status` int(11) DEFAULT NULL COMMENT '????????????',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_interact_id` (`interact_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6408 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='cs????????????';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cs_interact_detail`
--

DROP TABLE IF EXISTS `cs_interact_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cs_interact_detail` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '????????????',
  `interact_id` varchar(100) DEFAULT NULL COMMENT 'interact_ID',
  `from_room_id` varchar(100) DEFAULT NULL COMMENT 'room_id',
  `from_user_id` varchar(100) DEFAULT NULL COMMENT 'user_d',
  `to_room_id` varchar(100) DEFAULT NULL COMMENT 'room_id',
  `to_user_id` varchar(100) DEFAULT NULL COMMENT 'user_d',
  `interact_type` int(11) DEFAULT NULL COMMENT '????????????',
  `status` int(11) DEFAULT NULL COMMENT '????????????',
  `seat_id` int(11) DEFAULT NULL COMMENT 'seat_id',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????????????',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '????????????',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_interact_id` (`interact_id`,`from_room_id`,`from_user_id`,`to_room_id`,`to_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6377 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='cs??????????????????';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cs_room`
--

DROP TABLE IF EXISTS `cs_room`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cs_room` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '????????????',
  `app_id` varchar(100) DEFAULT NULL COMMENT 'app_ID',
  `room_id` varchar(100) DEFAULT NULL COMMENT '?????????room_id',
  `room_name` varchar(200) DEFAULT NULL COMMENT '???????????????',
  `owner_user_id` varchar(100) DEFAULT NULL COMMENT '??????id',
  `owner_user_name` varchar(200) DEFAULT NULL COMMENT '????????????',
  `status` int(11) DEFAULT NULL COMMENT '???????????????',
  `mic` int(11) DEFAULT NULL COMMENT '???????????????',
  `camera` int(11) DEFAULT NULL COMMENT '???????????????',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????????????',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '????????????',
  `ext` varchar(200) DEFAULT NULL COMMENT '????????????',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_room_id` (`room_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6409 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='cs????????????';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cs_room_user`
--

DROP TABLE IF EXISTS `cs_room_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cs_room_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '????????????',
  `app_id` varchar(100) DEFAULT NULL COMMENT 'app_id',
  `room_id` varchar(100) DEFAULT NULL COMMENT 'room_id',
  `user_id` varchar(255) DEFAULT NULL COMMENT 'user_id',
  `user_name` varchar(255) DEFAULT NULL COMMENT '????????????',
  `user_role` int(11) DEFAULT NULL COMMENT '??????????????? 1????????????2?????????',
  `net_status` int(11) DEFAULT NULL COMMENT '??????????????????',
  `interact_status` int(11) DEFAULT NULL COMMENT '??????????????????',
  `mic` tinyint(4) DEFAULT '0' COMMENT '???????????????',
  `camera` tinyint(4) DEFAULT '0' COMMENT '???????????????',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????????????',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '????????????',
  `device_id` varchar(128) DEFAULT NULL COMMENT 'device_id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_room_id_user_id` (`room_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14986 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='cs????????????';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ktv_room`
--


/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-03-01 14:48:03
