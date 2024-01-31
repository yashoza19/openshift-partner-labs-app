-- MySQL dump 10.13  Distrib 8.0.32, for Linux (x86_64)
--
-- Host: a9b2c26c893304a5e8b1ac631f41f96b-190976360.us-east-1.elb.amazonaws.com    Database: openshift_partner_labs_app_dev
-- ------------------------------------------------------
-- Server version	8.0.32

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `avatars`
--

DROP TABLE IF EXISTS `avatars`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `avatars` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `username` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `url` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `labs`
--

DROP TABLE IF EXISTS `labs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `labs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cluster_id` char(36) COLLATE utf8mb4_general_ci NOT NULL,
  `generated_name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `cluster_name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `openshift_version` varchar(16) COLLATE utf8mb4_general_ci NOT NULL,
  `cluster_size` varchar(7) COLLATE utf8mb4_general_ci NOT NULL,
  `company_name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `request_type` varchar(12) COLLATE utf8mb4_general_ci NOT NULL,
  `partner` tinyint(1) NOT NULL,
  `sponsor` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `cloud_provider` varchar(8) COLLATE utf8mb4_general_ci NOT NULL,
  `primary_first` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `primary_last` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `primary_email` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `secondary_first` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `secondary_last` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `secondary_email` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `region` varchar(5) COLLATE utf8mb4_general_ci NOT NULL,
  `always_on` tinyint(1) NOT NULL,
  `project_name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `lease_time` varchar(2) COLLATE utf8mb4_general_ci NOT NULL,
  `description` text COLLATE utf8mb4_general_ci NOT NULL,
  `notes` text COLLATE utf8mb4_general_ci NOT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime NOT NULL,
  `state` varchar(8) COLLATE utf8mb4_general_ci NOT NULL,
  `hold` tinyint(1) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `notes`
--

DROP TABLE IF EXISTS `notes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `labid` int NOT NULL,
  `note` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `notes_labs_id_fk` (`labid`),
  CONSTRAINT `notes_labs_id_fk` FOREIGN KEY (`labid`) REFERENCES `labs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `schema_migration`
--

DROP TABLE IF EXISTS `schema_migration`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migration` (
  `version` varchar(14) COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`version`),
  UNIQUE KEY `schema_migration_version_idx` (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-04-24 13:02:34
