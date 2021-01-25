-- MySQL dump 10.13  Distrib 5.7.22, for osx10.12 (x86_64)
--
-- Host: localhost    Database: accord
-- ------------------------------------------------------
-- Server version	5.7.22

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

--
-- Table structure for table `RentableUseType`
--

DROP TABLE IF EXISTS `RentableUseType`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `RentableUseType` (
  `UTID` bigint(20) NOT NULL AUTO_INCREMENT,
  `RID` bigint(20) NOT NULL DEFAULT '0',
  `BID` bigint(20) NOT NULL DEFAULT '0',
  `UseType` smallint(6) NOT NULL DEFAULT '0',
  `DtStart` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `DtStop` datetime NOT NULL DEFAULT '1970-01-01 00:00:00',
  `Comment` varchar(2048) NOT NULL DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` bigint(20) NOT NULL DEFAULT '0',
  `CreateTS` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`UTID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `RentableUseType`
--

LOCK TABLES `RentableUseType` WRITE;
/*!40000 ALTER TABLE `RentableUseType` DISABLE KEYS */;
/*!40000 ALTER TABLE `RentableUseType` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `classes`
--

DROP TABLE IF EXISTS `classes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `classes` (
  `ClassCode` mediumint(9) NOT NULL AUTO_INCREMENT,
  `CoCode` mediumint(9) NOT NULL DEFAULT '0',
  `Name` varchar(50) NOT NULL DEFAULT '',
  `Designation` char(3) NOT NULL DEFAULT '',
  `Description` varchar(256) NOT NULL DEFAULT '',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` mediumint(9) NOT NULL DEFAULT '0',
  PRIMARY KEY (`ClassCode`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `classes`
--

LOCK TABLES `classes` WRITE;
/*!40000 ALTER TABLE `classes` DISABLE KEYS */;
INSERT INTO `classes` VALUES (1,0,'Formula Gray','FOG','','2018-12-03 19:58:25',0),(2,0,'Sexsi Senorita','SES','','2018-12-03 19:58:25',0),(3,0,'Strong bod','STB','','2018-12-03 19:58:25',0),(4,0,'Access Asia','ACA','','2018-12-03 19:58:25',0),(5,0,'Intelacard','INT','','2018-12-03 19:58:25',0),(6,0,'Integra Wealth','INW','','2018-12-03 19:58:25',0),(7,0,'Destiny Realty Solutions','DRS','','2018-12-03 19:58:25',0);
/*!40000 ALTER TABLE `classes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `companies`
--

DROP TABLE IF EXISTS `companies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `companies` (
  `CoCode` mediumint(9) NOT NULL AUTO_INCREMENT,
  `LegalName` varchar(50) NOT NULL DEFAULT '',
  `CommonName` varchar(50) NOT NULL DEFAULT '',
  `Address` varchar(35) NOT NULL DEFAULT '',
  `Address2` varchar(35) NOT NULL DEFAULT '',
  `City` varchar(25) NOT NULL DEFAULT '',
  `State` char(25) NOT NULL DEFAULT '',
  `PostalCode` varchar(10) NOT NULL DEFAULT '',
  `Country` varchar(25) NOT NULL DEFAULT '',
  `Phone` varchar(25) NOT NULL DEFAULT '',
  `Fax` varchar(25) NOT NULL DEFAULT '',
  `Email` varchar(50) NOT NULL DEFAULT '',
  `Designation` char(3) NOT NULL DEFAULT '',
  `Active` smallint(6) NOT NULL DEFAULT '0',
  `EmploysPersonnel` smallint(6) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` mediumint(9) NOT NULL DEFAULT '0',
  PRIMARY KEY (`CoCode`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `companies`
--

LOCK TABLES `companies` WRITE;
/*!40000 ALTER TABLE `companies` DISABLE KEYS */;
INSERT INTO `companies` VALUES (1,'Teradyne Inc','Teradyne Inc','10325 South Carolina','','Fort Smith','NM','11091','USA','(129) 296-1934','(853) 401-4817','TeradyneInc6795@abiz.com','TEI',1,0,'2018-12-03 19:58:25',0),(2,'Jefferson-Pilot Co.','Jefferson-Pilot Co.','89324 Lehua','','Coral Springs','OK','31073','USA','(564) 727-2003','(750) 370-4740','Jefferson-Pilot2702@gmail.com','JEC',0,0,'2018-12-03 19:58:25',0),(3,'The Neiman Marcus Group I','The Neiman Marcus Group I','93101 Smith','','Bellevue','ME','65750','USA','(632) 707-0548','(965) 426-2338','T5065@hotmail.com','TNM',1,0,'2018-12-03 19:58:25',0),(4,'MPS Group Inc.','MPS Group Inc.','81466 Wood','','Newburgh','MD','04271','USA','(590) 501-5552','(708) 293-2820','MPSGroupInc.311@yahoo.com','MGI',1,0,'2018-12-03 19:58:25',0),(5,'UnumProvident Corporation','UnumProvident Corporation','90927 Delaware','','Albany','AL','99491','USA','(500) 606-1045','(769) 702-6163','UnumProvidentC2624@bdiddy.com','UNC',1,1,'2018-12-03 19:58:25',0),(6,'United Defense Industries','United Defense Industries','28756 North','','Alexandria','NV','91403','USA','(969) 592-9079','(596) 842-9399','UnitedDefense4779@bdiddy.com','UDI',0,0,'2018-12-03 19:58:25',0),(7,'IT Group Inc.','IT Group Inc.','96133 Pleasant','','Houston','HI','55442','USA','(443) 401-7972','(337) 661-0665','I8835@belcore.com','IGI',1,0,'2018-12-03 19:58:25',0);
/*!40000 ALTER TABLE `companies` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `compensation`
--

DROP TABLE IF EXISTS `compensation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `compensation` (
  `UID` mediumint(9) NOT NULL,
  `Type` mediumint(9) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `compensation`
--

LOCK TABLES `compensation` WRITE;
/*!40000 ALTER TABLE `compensation` DISABLE KEYS */;
/*!40000 ALTER TABLE `compensation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `counters`
--

DROP TABLE IF EXISTS `counters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `counters` (
  `SearchPeople` bigint(20) NOT NULL DEFAULT '0',
  `SearchClasses` bigint(20) NOT NULL DEFAULT '0',
  `SearchCompanies` bigint(20) NOT NULL DEFAULT '0',
  `EditPerson` bigint(20) NOT NULL DEFAULT '0',
  `ViewPerson` bigint(20) NOT NULL DEFAULT '0',
  `ViewClass` bigint(20) NOT NULL DEFAULT '0',
  `ViewCompany` bigint(20) NOT NULL DEFAULT '0',
  `AdminEditPerson` bigint(20) NOT NULL DEFAULT '0',
  `AdminEditClass` bigint(20) NOT NULL DEFAULT '0',
  `AdminEditCompany` bigint(20) NOT NULL DEFAULT '0',
  `DeletePerson` bigint(20) NOT NULL DEFAULT '0',
  `DeleteClass` bigint(20) NOT NULL DEFAULT '0',
  `DeleteCompany` bigint(20) NOT NULL DEFAULT '0',
  `SignIn` bigint(20) NOT NULL DEFAULT '0',
  `Logoff` bigint(20) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `counters`
--

LOCK TABLES `counters` WRITE;
/*!40000 ALTER TABLE `counters` DISABLE KEYS */;
INSERT INTO `counters` VALUES (0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,'2018-12-03 19:35:17');
/*!40000 ALTER TABLE `counters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deductionlist`
--

DROP TABLE IF EXISTS `deductionlist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deductionlist` (
  `DCode` mediumint(9) NOT NULL,
  `Name` varchar(25) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deductionlist`
--

LOCK TABLES `deductionlist` WRITE;
/*!40000 ALTER TABLE `deductionlist` DISABLE KEYS */;
INSERT INTO `deductionlist` VALUES (0,'Unknown'),(1,'401K'),(2,'401K Loan'),(3,'Child Support'),(4,'Dental'),(5,'FSA'),(6,'GARN'),(7,'Group Life'),(8,'Housing'),(9,'Medical'),(10,'Miscded'),(11,'Taxes');
/*!40000 ALTER TABLE `deductionlist` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deductions`
--

DROP TABLE IF EXISTS `deductions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deductions` (
  `UID` mediumint(9) NOT NULL,
  `Deduction` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deductions`
--

LOCK TABLES `deductions` WRITE;
/*!40000 ALTER TABLE `deductions` DISABLE KEYS */;
/*!40000 ALTER TABLE `deductions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `departments`
--

DROP TABLE IF EXISTS `departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `departments` (
  `DeptCode` mediumint(9) NOT NULL AUTO_INCREMENT,
  `Name` varchar(25) DEFAULT NULL,
  PRIMARY KEY (`DeptCode`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `departments`
--

LOCK TABLES `departments` WRITE;
/*!40000 ALTER TABLE `departments` DISABLE KEYS */;
INSERT INTO `departments` VALUES (1,'Accounting'),(2,'Administrative'),(3,'Capital Improvements'),(4,'Courtesy Patrol'),(5,'Customer Service'),(6,'Fitness Center'),(7,'Food & Beverage'),(8,'Guest Services'),(9,'Housekeeping'),(10,'Landscaping'),(11,'Maintenance'),(12,'Product Development'),(13,'Product Sales'),(14,'Serviced Apt Sales'),(15,'Trad Apt Sales'),(16,'Unknown'),(17,'Warehouse');
/*!40000 ALTER TABLE `departments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fieldperms`
--

DROP TABLE IF EXISTS `fieldperms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fieldperms` (
  `RID` mediumint(9) NOT NULL,
  `Elem` mediumint(9) NOT NULL,
  `Field` varchar(25) NOT NULL,
  `Perm` mediumint(9) NOT NULL,
  `Descr` varchar(256) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fieldperms`
--

LOCK TABLES `fieldperms` WRITE;
/*!40000 ALTER TABLE `fieldperms` DISABLE KEYS */;
INSERT INTO `fieldperms` VALUES (1,1,'Status',23,'Indicates whether the person is an active employee.'),(1,1,'EligibleForRehire',23,'Indicates whether a past employee can be rehired.'),(1,1,'UID',19,'A unique identifier associated with the employee. Once created, it never changes.'),(1,1,'Salutation',31,'\'Mr.\', \'Mrs.\', \'Ms.\', etc.'),(1,1,'FirstName',31,'The person\'s common name.'),(1,1,'MiddleName',31,'The person\'s middle name.'),(1,1,'LastName',31,'The person\'s surname or last name.'),(1,1,'PreferredName',95,'Less formal name but more commonly used, for example \'Mike\' rather than \'Michael\'.'),(1,1,'PrimaryEmail',95,'The primary email address to use for this person.'),(1,1,'OfficePhone',95,'This person\'s office telephone number.'),(1,1,'CellPhone',95,'This person\'s cellphone number.'),(1,1,'EmergencyContactName',95,'Name of someone to contact in the event of an emergency.'),(1,1,'EmergencyContactPhone',95,'Phone number for the emergency contact.'),(1,1,'HomeStreetAddress',95,'def'),(1,1,'HomeStreetAddress2',95,'def'),(1,1,'HomeCity',95,'def'),(1,1,'HomeState',95,'def'),(1,1,'HomePostalCode',95,'def'),(1,1,'HomeCountry',95,'def'),(1,1,'PrimaryEmail',95,'def'),(1,1,'SecondaryEmail',95,'def'),(1,1,'OfficePhone',95,'def'),(1,1,'OfficeFax',95,'def'),(1,1,'CellPhone',95,'def'),(1,1,'BirthDOM',31,'def'),(1,1,'BirthMonth',31,'def'),(1,1,'CoCode',31,'The company code associated with this user.'),(1,1,'JobCode',31,'def'),(1,1,'ClassCode',31,'def'),(1,1,'DeptCode',31,'def'),(1,1,'PositionControlNumber',31,'def'),(1,1,'MgrUID',31,'def'),(1,1,'Accepted401K',31,'def'),(1,1,'AcceptedDentalInsurance',31,'def'),(1,1,'AcceptedHealthInsurance',31,'def'),(1,1,'Hire',31,'def'),(1,1,'Termination',31,'def'),(1,1,'LastReview',31,'def'),(1,1,'NextReview',31,'def'),(1,1,'StateOfEmployment',31,'def'),(1,1,'CountryOfEmployment',31,'def'),(1,1,'Comps',31,'def'),(1,1,'Deductions',31,'def'),(1,1,'MyDeductions',31,'def'),(1,1,'RID',31,'def'),(1,1,'Role',31,'Permissions role'),(1,1,'ElemEntity',31,'Permissions to delete the entity'),(1,2,'CoCode',31,'def'),(1,2,'LegalName',31,'def'),(1,2,'CommonName',31,'def'),(1,2,'Address',31,'def'),(1,2,'Address2',31,'def'),(1,2,'City',31,'def'),(1,2,'State',31,'def'),(1,2,'PostalCode',31,'def'),(1,2,'Country',31,'def'),(1,2,'Phone',31,'def'),(1,2,'Fax',31,'def'),(1,2,'Email',31,'def'),(1,2,'Designation',31,'def'),(1,2,'Active',31,'def'),(1,2,'EmploysPersonnel',31,'def'),(1,2,'ElemEntity',31,'def'),(1,3,'ClassCode',31,'def'),(1,3,'CoCode',31,'The parent company for this business unit'),(1,3,'Name',31,'def'),(1,3,'Designation',31,'def'),(1,3,'Description',31,'def'),(1,3,'ElemEntity',31,'def'),(1,4,'Shutdown',256,'Permission to shutdown the service'),(1,4,'Restart',256,'Permission to restart the service'),(2,1,'Status',23,'Indicates whether the person is an active employee.'),(2,1,'EligibleForRehire',23,'Indicates whether a past employee can be rehired.'),(2,1,'UID',19,'A unique identifier associated with the employee. Once created, it never changes.'),(2,1,'Salutation',95,'\'Mr.\', \'Mrs.\', \'Ms.\', etc.'),(2,1,'FirstName',95,'The person\'s common name.'),(2,1,'MiddleName',95,'The person\'s middle name.'),(2,1,'LastName',95,'The person\'s surname or last name.'),(2,1,'PreferredName',95,'Less formal name but more commonly used, for example \'Mike\' rather than \'Michael\'.'),(2,1,'PrimaryEmail',95,'The primary email address to use for this person.'),(2,1,'OfficePhone',95,'This person\'s office telephone number.'),(2,1,'CellPhone',95,'This person\'s cellphone number.'),(2,1,'EmergencyContactName',95,'Name of someone to contact in the event of an emergency.'),(2,1,'EmergencyContactPhone',95,'Phone number for the emergency contact.'),(2,1,'HomeStreetAddress',95,'def'),(2,1,'HomeStreetAddress2',95,'def'),(2,1,'HomeCity',95,'def'),(2,1,'HomeState',95,'def'),(2,1,'HomePostalCode',95,'def'),(2,1,'HomeCountry',95,'def'),(2,1,'PrimaryEmail',95,'def'),(2,1,'SecondaryEmail',95,'def'),(2,1,'OfficePhone',95,'def'),(2,1,'OfficeFax',95,'def'),(2,1,'CellPhone',95,'def'),(2,1,'BirthDOM',31,'def'),(2,1,'BirthMonth',31,'def'),(2,1,'CoCode',31,'The company code associated with this user.'),(2,1,'JobCode',31,'def'),(2,1,'DeptCode',31,'def'),(2,1,'ClassCode',31,'def'),(2,1,'PositionControlNumber',31,'def'),(2,1,'MgrUID',31,'def'),(2,1,'Accepted401K',31,'def'),(2,1,'AcceptedDentalInsurance',31,'def'),(2,1,'AcceptedHealthInsurance',31,'def'),(2,1,'Hire',31,'def'),(2,1,'Termination',31,'def'),(2,1,'LastReview',31,'def'),(2,1,'NextReview',31,'def'),(2,1,'StateOfEmployment',31,'def'),(2,1,'CountryOfEmployment',31,'def'),(2,1,'Comps',31,'def'),(2,1,'Deductions',31,'def'),(2,1,'MyDeductions',31,'def'),(2,1,'Role',1,'Permissions Role'),(2,1,'RID',17,'def'),(2,1,'ElemEntity',0,'Permissions to delete the entity'),(2,2,'CoCode',17,'def'),(2,2,'LegalName',17,'def'),(2,2,'CommonName',17,'def'),(2,2,'Address',17,'def'),(2,2,'Address2',17,'def'),(2,2,'City',17,'def'),(2,2,'State',17,'def'),(2,2,'PostalCode',17,'def'),(2,2,'Country',17,'def'),(2,2,'Phone',17,'def'),(2,2,'Fax',17,'def'),(2,2,'Email',17,'def'),(2,2,'Designation',17,'def'),(2,2,'Active',17,'def'),(2,2,'EmploysPersonnel',17,'def'),(2,2,'ElemEntity',0,'def'),(2,3,'ClassCode',17,'def'),(2,3,'CoCode',17,'def'),(2,3,'Name',17,'The parent company for this business unit'),(2,3,'Designation',17,'def'),(2,3,'Description',17,'def'),(2,3,'ElemEntity',0,'def'),(2,4,'Shutdown',0,'Permission to shutdown the service'),(2,4,'Restart',0,'Permission to restart the service'),(3,1,'Status',17,'Indicates whether the person is an active employee.'),(3,1,'EligibleForRehire',23,'Indicates whether a past employee can be rehired.'),(3,1,'UID',19,'A unique identifier associated with the employee. Once created, it never changes.'),(3,1,'Salutation',17,'\'Mr.\', \'Mrs.\', \'Ms.\', etc.'),(3,1,'FirstName',17,'The person\'s common name.'),(3,1,'MiddleName',17,'The person\'s middle name.'),(3,1,'LastName',17,'The person\'s surname or last name.'),(3,1,'PreferredName',81,'Less formal name but more commonly used, for example \'Mike\' rather than \'Michael\'.'),(3,1,'PrimaryEmail',81,'The primary email address to use for this person.'),(3,1,'OfficePhone',81,'This person\'s office telephone number.'),(3,1,'CellPhone',81,'This person\'s cellphone number.'),(3,1,'EmergencyContactName',112,'Name of someone to contact in the event of an emergency.'),(3,1,'EmergencyContactPhone',112,'Phone number for the emergency contact.'),(3,1,'HomeStreetAddress',112,'def'),(3,1,'HomeStreetAddress2',112,'def'),(3,1,'HomeCity',112,'def'),(3,1,'HomeState',112,'def'),(3,1,'HomePostalCode',112,'def'),(3,1,'HomeCountry',112,'def'),(3,1,'PrimaryEmail',81,'def'),(3,1,'SecondaryEmail',81,'def'),(3,1,'OfficePhone',81,'def'),(3,1,'OfficeFax',81,'def'),(3,1,'CellPhone',81,'def'),(3,1,'BirthDOM',48,'def'),(3,1,'BirthMonth',48,'def'),(3,1,'CoCode',17,'The company code associated with this user.'),(3,1,'JobCode',17,'def'),(3,1,'DeptCode',17,'def'),(3,1,'ClassCode',17,'def'),(3,1,'MgrUID',17,'def'),(3,1,'Accepted401K',17,'def'),(3,1,'AcceptedDentalInsurance',17,'def'),(3,1,'AcceptedHealthInsurance',17,'def'),(3,1,'PositionControlNumber',17,'def'),(3,1,'Hire',48,'def'),(3,1,'Termination',17,'def'),(3,1,'LastReview',0,'def'),(3,1,'NextReview',0,'def'),(3,1,'StateOfEmployment',17,'def'),(3,1,'CountryOfEmployment',17,'def'),(3,1,'Comps',17,'def'),(3,1,'Deductions',17,'def'),(3,1,'MyDeductions',17,'def'),(3,1,'RID',0,'def'),(3,1,'Role',0,'Permissions Role'),(3,1,'ElemEntity',0,'Permissions to delete the entity'),(3,2,'CoCode',31,'def'),(3,2,'LegalName',31,'def'),(3,2,'CommonName',31,'def'),(3,2,'Address',31,'def'),(3,2,'Address2',31,'def'),(3,2,'City',31,'def'),(3,2,'State',31,'def'),(3,2,'PostalCode',31,'def'),(3,2,'Country',31,'def'),(3,2,'Phone',31,'def'),(3,2,'Fax',31,'def'),(3,2,'Email',31,'def'),(3,2,'Designation',31,'def'),(3,2,'Active',31,'def'),(3,2,'EmploysPersonnel',31,'def'),(3,2,'ElemEntity',0,'def'),(3,3,'ClassCode',31,'def'),(3,3,'CoCode',31,'The parent company for this business unit'),(3,3,'Name',31,'def'),(3,3,'Designation',31,'def'),(3,3,'Description',31,'def'),(3,3,'ElemEntity',0,'def'),(3,4,'Shutdown',0,'Permission to shutdown the service'),(3,4,'Restart',0,'Permission to restart the service'),(4,1,'Status',1,'Indicates whether the person is an active employee.'),(4,1,'EligibleForRehire',1,'Indicates whether a past employee can be rehired.'),(4,1,'UID',17,'A unique identifier associated with the employee. Once created, it never changes.'),(4,1,'Salutation',1,'\'Mr.\', \'Mrs.\', \'Ms.\', etc.'),(4,1,'FirstName',1,'The person\'s common name.'),(4,1,'MiddleName',1,'The person\'s middle name.'),(4,1,'LastName',1,'The person\'s surname or last name.'),(4,1,'PreferredName',65,'Less formal name but more commonly used, for example \'Mike\' rather than \'Michael\'.'),(4,1,'PrimaryEmail',65,'The primary email address to use for this person.'),(4,1,'OfficePhone',65,'This person\'s office telephone number.'),(4,1,'CellPhone',65,'This person\'s cellphone number.'),(4,1,'EmergencyContactName',193,'Name of someone to contact in the event of an emergency.'),(4,1,'EmergencyContactPhone',193,'Phone number for the emergency contact.'),(4,1,'HomeStreetAddress',193,'def'),(4,1,'HomeStreetAddress2',193,'def'),(4,1,'HomeCity',193,'def'),(4,1,'HomeState',193,'def'),(4,1,'HomePostalCode',193,'def'),(4,1,'HomeCountry',81,'def'),(4,1,'PrimaryEmail',81,'def'),(4,1,'SecondaryEmail',81,'def'),(4,1,'OfficePhone',81,'def'),(4,1,'OfficeFax',81,'def'),(4,1,'CellPhone',81,'def'),(4,1,'BirthDOM',160,'def'),(4,1,'BirthMonth',160,'def'),(4,1,'CoCode',160,'The company code associated with this user.'),(4,1,'JobCode',160,'def'),(4,1,'DeptCode',160,'def'),(4,1,'ClassCode',17,'def'),(4,1,'PositionControlNumber',160,'def'),(4,1,'MgrUID',17,'def'),(4,1,'Accepted401K',160,'def'),(4,1,'AcceptedDentalInsurance',160,'def'),(4,1,'AcceptedHealthInsurance',160,'def'),(4,1,'Hire',160,'def'),(4,1,'Termination',32,'def'),(4,1,'LastReview',32,'def'),(4,1,'NextReview',32,'def'),(4,1,'StateOfEmployment',160,'def'),(4,1,'CountryOfEmployment',160,'def'),(4,1,'Comps',160,'Compensation type(s) for this person.'),(4,1,'Deductions',160,'The deductions for this person.'),(4,1,'MyDeductions',160,'The deductions for this person.'),(4,1,'RID',17,'def'),(4,1,'Role',0,'Permissions Rol'),(4,1,'ElemEntity',0,'Permissions to delete the entity'),(4,2,'CoCode',1,'def'),(4,2,'LegalName',1,'def'),(4,2,'CommonName',1,'def'),(4,2,'Address',1,'def'),(4,2,'Address2',1,'def'),(4,2,'City',1,'def'),(4,2,'State',1,'def'),(4,2,'PostalCode',1,'def'),(4,2,'Country',1,'def'),(4,2,'Phone',1,'def'),(4,2,'Fax',1,'def'),(4,2,'Email',1,'def'),(4,2,'Designation',1,'def'),(4,2,'Active',1,'def'),(4,2,'EmploysPersonnel',1,'def'),(4,3,'ClassCode',1,'def'),(4,3,'CoCode',1,'The parent company for this business unit'),(4,3,'Name',1,'def'),(4,3,'Designation',1,'def'),(4,3,'Description',1,'def'),(4,3,'ElemEntity',0,'def'),(4,4,'Shutdown',0,'Permission to shutdown the service'),(4,4,'Restart',0,'Permission to restart the service'),(5,1,'Status',23,'Indicates whether the person is an active employee.'),(5,1,'EligibleForRehire',1,'Indicates whether a past employee can be rehired.'),(5,1,'UID',3,'A unique identifier associated with the employee. Once created, it never changes.'),(5,1,'Salutation',4,'\'Mr.\', \'Mrs.\', \'Ms.\', etc.'),(5,1,'FirstName',8,'The person\'s common name.'),(5,1,'MiddleName',16,'The person\'s middle name.'),(5,1,'LastName',0,'The person\'s surname or last name.'),(5,1,'PreferredName',17,'Less formal name but more commonly used, for example \'Mike\' rather than \'Michael\'.'),(5,1,'PrimaryEmail',1,'The primary email address to use for this person.'),(5,1,'OfficePhone',0,'This person\'s office telephone number.'),(5,1,'CellPhone',7,'This person\'s cellphone number.'),(5,1,'EmergencyContactName',0,'Name of someone to contact in the event of an emergency.'),(5,1,'EmergencyContactPhone',95,'Phone number for the emergency contact.'),(5,1,'HomeStreetAddress',95,'def'),(5,1,'HomeStreetAddress2',1,'def'),(5,1,'HomeCity',0,'def'),(5,1,'HomeState',95,'def'),(5,1,'HomePostalCode',0,'def'),(5,1,'HomeCountry',95,'def'),(5,1,'PrimaryEmail',95,'def'),(5,1,'SecondaryEmail',0,'def'),(5,1,'OfficePhone',95,'def'),(5,1,'OfficeFax',0,'def'),(5,1,'CellPhone',95,'def'),(5,1,'BirthDOM',0,'def'),(5,1,'BirthMonth',31,'def'),(5,1,'CoCode',0,'The company code associated with this user.'),(5,1,'JobCode',31,'def'),(5,1,'ClassCode',0,'def'),(5,1,'DeptCode',31,'def'),(5,1,'PositionControlNumber',0,'def'),(5,1,'MgrUID',31,'def'),(5,1,'Accepted401K',0,'def'),(5,1,'AcceptedDentalInsurance',31,'def'),(5,1,'AcceptedHealthInsurance',0,'def'),(5,1,'Hire',31,'def'),(5,1,'Termination',0,'def'),(5,1,'LastReview',31,'def'),(5,1,'NextReview',0,'def'),(5,1,'StateOfEmployment',31,'def'),(5,1,'CountryOfEmployment',0,'def'),(5,1,'Comps',31,'def'),(5,1,'Deductions',17,'def'),(5,1,'MyDeductions',17,'def'),(5,1,'RID',17,'def'),(5,1,'Role',0,'Permissions Rol'),(5,1,'ElemEntity',0,'Permissions to delete the entity'),(5,2,'CoCode',31,'def'),(5,2,'LegalName',0,'def'),(5,2,'CommonName',31,'def'),(5,2,'Address',31,'def'),(5,2,'Address2',0,'def'),(5,2,'City',31,'def'),(5,2,'State',31,'def'),(5,2,'PostalCode',31,'def'),(5,2,'Country',0,'def'),(5,2,'Phone',31,'def'),(5,2,'Fax',0,'def'),(5,2,'Email',31,'def'),(5,2,'Designation',31,'def'),(5,2,'Active',0,'def'),(5,2,'EmploysPersonnel',31,'def'),(5,2,'ElemEntity',31,'def'),(5,3,'ClassCode',31,'def'),(5,3,'CoCode',31,'The parent company for this business unit'),(5,3,'Name',31,'def'),(5,3,'Designation',31,'def'),(5,3,'Description',0,'def'),(5,3,'ElemEntity',31,'def'),(5,4,'Shutdown',256,'Permission to shutdown the service'),(5,4,'Restart',256,'Permission to restart the service'),(6,1,'Status',23,'Indicates whether the person is an active employee.'),(6,1,'EligibleForRehire',23,'Indicates whether a past employee can be rehired.'),(6,1,'UID',19,'A unique identifier associated with the employee. Once created, it never changes.'),(6,1,'Salutation',95,'\'Mr.\', \'Mrs.\', \'Ms.\', etc.'),(6,1,'FirstName',95,'The person\'s common name.'),(6,1,'MiddleName',95,'The person\'s middle name.'),(6,1,'LastName',95,'The person\'s surname or last name.'),(6,1,'PreferredName',95,'Less formal name but more commonly used, for example \'Mike\' rather than \'Michael\'.'),(6,1,'PrimaryEmail',95,'The primary email address to use for this person.'),(6,1,'OfficePhone',95,'This person\'s office telephone number.'),(6,1,'CellPhone',95,'This person\'s cellphone number.'),(6,1,'EmergencyContactName',95,'Name of someone to contact in the event of an emergency.'),(6,1,'EmergencyContactPhone',95,'Phone number for the emergency contact.'),(6,1,'HomeStreetAddress',95,'def'),(6,1,'HomeStreetAddress2',95,'def'),(6,1,'HomeCity',95,'def'),(6,1,'HomeState',95,'def'),(6,1,'HomePostalCode',95,'def'),(6,1,'HomeCountry',95,'def'),(6,1,'PrimaryEmail',95,'def'),(6,1,'SecondaryEmail',95,'def'),(6,1,'OfficePhone',95,'def'),(6,1,'OfficeFax',95,'def'),(6,1,'CellPhone',95,'def'),(6,1,'BirthDOM',31,'def'),(6,1,'BirthMonth',31,'def'),(6,1,'CoCode',31,'The company code associated with this user.'),(6,1,'JobCode',31,'def'),(6,1,'DeptCode',31,'def'),(6,1,'ClassCode',31,'def'),(6,1,'PositionControlNumber',31,'def'),(6,1,'MgrUID',31,'def'),(6,1,'Accepted401K',31,'def'),(6,1,'AcceptedDentalInsurance',31,'def'),(6,1,'AcceptedHealthInsurance',31,'def'),(6,1,'Hire',31,'def'),(6,1,'Termination',31,'def'),(6,1,'LastReview',31,'def'),(6,1,'NextReview',31,'def'),(6,1,'StateOfEmployment',31,'def'),(6,1,'CountryOfEmployment',31,'def'),(6,1,'Comps',31,'def'),(6,1,'Deductions',31,'def'),(6,1,'MyDeductions',31,'def'),(6,1,'Role',1,'Permissions Rol'),(6,1,'RID',17,'def'),(6,1,'ElemEntity',0,'Permissions to delete the entity'),(6,2,'CoCode',31,'def'),(6,2,'LegalName',31,'def'),(6,2,'CommonName',31,'def'),(6,2,'Address',31,'def'),(6,2,'Address2',31,'def'),(6,2,'City',31,'def'),(6,2,'State',31,'def'),(6,2,'PostalCode',31,'def'),(6,2,'Country',31,'def'),(6,2,'Phone',31,'def'),(6,2,'Fax',31,'def'),(6,2,'Email',31,'def'),(6,2,'Designation',31,'def'),(6,2,'Active',31,'def'),(6,2,'EmploysPersonnel',31,'def'),(6,2,'ElemEntity',0,'def'),(6,3,'ClassCode',31,'def'),(6,3,'CoCode',31,'The parent company for this business unit'),(6,3,'Name',31,'def'),(6,3,'Designation',31,'def'),(6,3,'Description',31,'def'),(6,3,'ElemEntity',0,'def'),(6,4,'Shutdown',0,'Permission to shutdown the service'),(6,4,'Restart',0,'Permission to restart the service'),(7,1,'Status',23,'Indicates whether the person is an active employee.'),(7,1,'EligibleForRehire',23,'Indicates whether a past employee can be rehired.'),(7,1,'UID',19,'A unique identifier associated with the employee. Once created, it never changes.'),(7,1,'Salutation',95,'\'Mr.\', \'Mrs.\', \'Ms.\', etc.'),(7,1,'FirstName',95,'The person\'s common name.'),(7,1,'MiddleName',95,'The person\'s middle name.'),(7,1,'LastName',95,'The person\'s surname or last name.'),(7,1,'PreferredName',95,'Less formal name but more commonly used, for example \'Mike\' rather than \'Michael\'.'),(7,1,'PrimaryEmail',95,'The primary email address to use for this person.'),(7,1,'OfficePhone',95,'This person\'s office telephone number.'),(7,1,'CellPhone',95,'This person\'s cellphone number.'),(7,1,'EmergencyContactName',95,'Name of someone to contact in the event of an emergency.'),(7,1,'EmergencyContactPhone',95,'Phone number for the emergency contact.'),(7,1,'HomeStreetAddress',95,'def'),(7,1,'HomeStreetAddress2',95,'def'),(7,1,'HomeCity',95,'def'),(7,1,'HomeState',95,'def'),(7,1,'HomePostalCode',95,'def'),(7,1,'HomeCountry',95,'def'),(7,1,'PrimaryEmail',95,'def'),(7,1,'SecondaryEmail',95,'def'),(7,1,'OfficePhone',95,'def'),(7,1,'OfficeFax',95,'def'),(7,1,'CellPhone',95,'def'),(7,1,'BirthDOM',31,'def'),(7,1,'BirthMonth',31,'def'),(7,1,'CoCode',31,'The company code associated with this user.'),(7,1,'JobCode',31,'def'),(7,1,'DeptCode',31,'def'),(7,1,'ClassCode',31,'def'),(7,1,'PositionControlNumber',31,'def'),(7,1,'MgrUID',31,'def'),(7,1,'Accepted401K',31,'def'),(7,1,'AcceptedDentalInsurance',31,'def'),(7,1,'AcceptedHealthInsurance',31,'def'),(7,1,'Hire',31,'def'),(7,1,'Termination',31,'def'),(7,1,'LastReview',31,'def'),(7,1,'NextReview',31,'def'),(7,1,'StateOfEmployment',31,'def'),(7,1,'CountryOfEmployment',31,'def'),(7,1,'Comps',31,'def'),(7,1,'Deductions',31,'def'),(7,1,'MyDeductions',31,'def'),(7,1,'Role',1,'Permissions Rol'),(7,1,'RID',17,'def'),(7,1,'ElemEntity',31,'Permissions to create/delete the entity'),(7,2,'CoCode',31,'def'),(7,2,'LegalName',31,'def'),(7,2,'CommonName',31,'def'),(7,2,'Address',31,'def'),(7,2,'Address2',31,'def'),(7,2,'City',31,'def'),(7,2,'State',31,'def'),(7,2,'PostalCode',31,'def'),(7,2,'Country',31,'def'),(7,2,'Phone',31,'def'),(7,2,'Fax',31,'def'),(7,2,'Email',31,'def'),(7,2,'Designation',31,'def'),(7,2,'Active',31,'def'),(7,2,'EmploysPersonnel',31,'def'),(7,2,'ElemEntity',31,'def'),(7,3,'ClassCode',31,'def'),(7,3,'CoCode',31,'The parent company of this business unit.'),(7,3,'Name',31,'def'),(7,3,'Designation',31,'def'),(7,3,'Description',31,'def'),(7,3,'ElemEntity',31,'def'),(7,4,'Shutdown',0,'Permission to shutdown the service'),(7,4,'Restart',0,'Permission to restart the service');
/*!40000 ALTER TABLE `fieldperms` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jobtitles`
--

DROP TABLE IF EXISTS `jobtitles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `jobtitles` (
  `JobCode` mediumint(9) NOT NULL AUTO_INCREMENT,
  `Title` varchar(40) NOT NULL DEFAULT '',
  `Descr` varchar(256) NOT NULL DEFAULT '',
  PRIMARY KEY (`JobCode`)
) ENGINE=InnoDB AUTO_INCREMENT=87 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jobtitles`
--

LOCK TABLES `jobtitles` WRITE;
/*!40000 ALTER TABLE `jobtitles` DISABLE KEYS */;
INSERT INTO `jobtitles` VALUES (1,'Accounting Assistant',''),(2,'Accounting Associate',''),(3,'Accounting Manager',''),(4,'Administrative Assistant',''),(5,'Assistant Manager',''),(6,'Associate Developer',''),(7,'Chief Executive Officer',''),(8,'Chief Financial Officer',''),(9,'Chief Operating Officer',''),(10,'Chief Technology Officer',''),(11,'General Manager',''),(12,'HR & Payroll Manager',''),(13,'Intern',''),(14,'Night Auditor',''),(15,'Office Manager',''),(16,'Procurement Specialist',''),(17,'Special Projects Associate',''),(18,'Director of Procurements',''),(19,'Call Center Associate',''),(20,'Call Center Manager',''),(21,'Courtesy Patrol Driver',''),(22,'Courtesy Patrol Manager',''),(23,'Courtesy Patrol Officer',''),(24,'Courtesy Patrol Supervisor',''),(25,'Designer',''),(26,'Creative Director',''),(27,'Development Coordinator',''),(28,'Director of Fragrence',''),(29,'Director of Sales',''),(30,'Principal',''),(31,'Studio Manager',''),(32,'Visual Arts Director',''),(33,'Fitness Center Attendant',''),(34,'Fitness Center Manager',''),(35,'Food and Beverage Manager',''),(36,'Executive Chef',''),(37,'Food & Bev Associate',''),(38,'Bar Manager',''),(39,'Bartender',''),(40,'Host',''),(41,'Waitstaff',''),(42,'Guest Services Associate',''),(43,'Concierge Manager',''),(44,'Concierge',''),(45,'Guest Services Manager',''),(46,'Housekeeping Manager',''),(47,'Housekeeping Supervisor',''),(48,'Common Area Housekeeper',''),(49,'Laundry Associate',''),(50,'Laundry Attendant',''),(51,'Serviced Apt Housekeeper',''),(52,'Svc Apt Housekeeping Associate',''),(53,'Traditional Apt Housekeeper',''),(54,'Grounds Associate',''),(55,'Grounds Supervisor',''),(56,'Maintenance Associate',''),(57,'Maintenance Manager',''),(58,'Maintenance Supervisor',''),(59,'Makeready Associate',''),(60,'Cap Improvement Associate',''),(61,'Cap Improvement Supervisor',''),(62,'Customer Service Associate',''),(63,'Customer Service Manager',''),(64,'Product Sales Associate',''),(65,'Product Sales Manager',''),(66,'National Accounts Manager',''),(67,'Leasing Associate',''),(68,'Leasing Manager',''),(69,'Packer',''),(70,'Repair Room Staff',''),(71,'Seasonal Associate',''),(72,'Checker',''),(73,'Warehouse Associate',''),(74,'Warehouse Manager',''),(75,'Warehouse Supervisor',''),(76,'Retail Clerk',''),(77,'Store Manager',''),(78,'Asst Store Manager',''),(79,'Traditional Apt Housekeeping Assoc',''),(80,'Makeready Technician',''),(81,'Serviced Apt Sales Associate',''),(82,'Unknown',''),(83,'Construction Manager',''),(84,'Dishwasher',''),(85,'NY Retail Boutique Manager',''),(86,'Marketing & Sales Manager','');
/*!40000 ALTER TABLE `jobtitles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `people`
--

DROP TABLE IF EXISTS `people`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `people` (
  `UID` mediumint(9) NOT NULL AUTO_INCREMENT,
  `UserName` varchar(20) NOT NULL DEFAULT '',
  `LastName` varchar(25) NOT NULL DEFAULT '',
  `MiddleName` varchar(25) NOT NULL DEFAULT '',
  `FirstName` varchar(25) NOT NULL DEFAULT '',
  `PreferredName` varchar(25) NOT NULL DEFAULT '',
  `Salutation` varchar(10) NOT NULL DEFAULT '',
  `PositionControlNumber` varchar(10) NOT NULL DEFAULT '',
  `OfficePhone` varchar(25) NOT NULL DEFAULT '',
  `OfficeFax` varchar(25) NOT NULL DEFAULT '',
  `CellPhone` varchar(25) NOT NULL DEFAULT '',
  `PrimaryEmail` varchar(35) NOT NULL DEFAULT '',
  `SecondaryEmail` varchar(35) NOT NULL DEFAULT '',
  `BirthMonth` tinyint(4) NOT NULL DEFAULT '0',
  `BirthDoM` tinyint(4) NOT NULL DEFAULT '0',
  `HomeStreetAddress` varchar(35) NOT NULL DEFAULT '',
  `HomeStreetAddress2` varchar(25) NOT NULL DEFAULT '',
  `HomeCity` varchar(25) NOT NULL DEFAULT '',
  `HomeState` char(2) NOT NULL DEFAULT '',
  `HomePostalCode` varchar(10) NOT NULL DEFAULT '',
  `HomeCountry` varchar(25) NOT NULL DEFAULT '',
  `JobCode` mediumint(9) NOT NULL DEFAULT '0',
  `Hire` date NOT NULL DEFAULT '2000-01-01',
  `Termination` date NOT NULL DEFAULT '2000-01-01',
  `MgrUID` mediumint(9) NOT NULL DEFAULT '0',
  `DeptCode` mediumint(9) NOT NULL DEFAULT '0',
  `CoCode` mediumint(9) NOT NULL DEFAULT '0',
  `ClassCode` smallint(6) NOT NULL DEFAULT '0',
  `StateOfEmployment` varchar(25) NOT NULL DEFAULT '',
  `CountryOfEmployment` varchar(25) NOT NULL DEFAULT '',
  `EmergencyContactName` varchar(25) NOT NULL DEFAULT '',
  `EmergencyContactPhone` varchar(25) NOT NULL DEFAULT '',
  `Status` smallint(6) NOT NULL DEFAULT '0',
  `EligibleForRehire` smallint(6) NOT NULL DEFAULT '0',
  `AcceptedHealthInsurance` smallint(6) NOT NULL DEFAULT '0',
  `AcceptedDentalInsurance` smallint(6) NOT NULL DEFAULT '0',
  `Accepted401K` smallint(6) NOT NULL DEFAULT '0',
  `LastReview` date NOT NULL DEFAULT '2000-01-01',
  `NextReview` date NOT NULL DEFAULT '2000-01-01',
  `passhash` char(128) NOT NULL DEFAULT '',
  `RID` mediumint(9) NOT NULL DEFAULT '0',
  `LastModTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `LastModBy` mediumint(9) NOT NULL DEFAULT '0',
  `CreateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreateBy` bigint(20) NOT NULL DEFAULT '0',
  `ImagePath` varchar(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`UID`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `people`
--

LOCK TABLES `people` WRITE;
/*!40000 ALTER TABLE `people` DISABLE KEYS */;
INSERT INTO `people` VALUES (1,'bthorton','Thorton','','Billy','','','','','','','','',0,0,'','','','','','',0,'2000-01-01','2000-01-01',0,0,0,0,'','','','',0,0,0,0,0,'2000-01-01','2000-01-01','5070ea6fea9d36140f6328e1c811e31b35a01c83fd3950c99a3cc8e026a8bfc990df1726eb253627645f647c3b8988ad2d27fa54e7785bf9959b21682b0805fa',4,'2018-12-03 19:52:14',0,'2021-01-24 23:14:26',0,''),(2,'eschneider','schneider','','elinore','angelena','','','(782) 718-9354','(902) 634-6255','(106) 720-6835','elinoreschneider34@abiz.com','eschneider7319@hotmail.com',0,0,'89985 Pinon','','Sunnyvale','DC','06192','USA',29,'2000-01-01','2000-01-01',0,12,5,6,'','','siobhandavenport','(186) 686-0936',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(3,'kwalsh','walsh','','kena','maria','','','(809) 906-3161','(975) 891-0425','(127) 249-3437','kwalsh7710@gmail.com','kwalsh8718@yahoo.com',0,0,'75339 Smith','','Santa Cruz','MI','89063','USA',12,'2000-01-01','2000-01-01',0,5,5,7,'','','elbabooker','(112) 746-8170',1,0,0,0,0,'2000-01-01','2000-01-01','5070ea6fea9d36140f6328e1c811e31b35a01c83fd3950c99a3cc8e026a8bfc990df1726eb253627645f647c3b8988ad2d27fa54e7785bf9959b21682b0805fa',1,'2019-04-15 02:32:00',0,'2021-01-24 23:14:26',0,''),(4,'ahart','hart','','almeta','zenia','','','(117) 719-9285','(636) 765-7066','(375) 816-5647','almetahart192@gmail.com','ahart5017@bdiddy.com',0,0,'9689 Cherry','','Cedar Rapids','NV','12886','USA',58,'2000-01-01','2000-01-01',0,3,5,2,'','','ruperthernandez','(103) 213-2744',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(5,'abarlow','barlow','','arron','tracee','','','(832) 801-6683','(386) 611-5518','(342) 510-3245','arronb8746@gmail.com','abarlow1344@gmail.com',0,0,'95256 A','','Frederick','MT','47729','USA',6,'2000-01-01','2000-01-01',0,8,5,2,'','','faithhyde','(338) 151-7190',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(6,'gnorton','norton','','giuseppe','ivonne','','','(221) 385-8940','(428) 939-2480','(328) 520-8337','giuseppen9343@gmail.com','gnorton2770@gmail.com',0,0,'66706 Apache','','Apple Valley','OR','91412','USA',34,'2000-01-01','2000-01-01',0,9,5,1,'','','griseldaalexander','(779) 964-4651',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(7,'ltate','tate','','lamont','josef','','','(101) 892-1643','(306) 461-7723','(229) 515-9340','ltate1162@aol.com','ltate3721@comcast.net',0,0,'65372 Shore','','Winter Haven','MA','43831','USA',56,'2000-01-01','2000-01-01',0,1,5,1,'','','maymebrock','(730) 231-8377',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(8,'csanford','sanford','','carlena','brad','','','(126) 861-5069','(250) 693-2347','(839) 540-6768','carlenasanford911@aol.com','csanford2232@yahoo.com',0,0,'67781 7th','','Kennewick','AZ','77795','USA',31,'2000-01-01','2000-01-01',0,7,5,4,'','','rockycobb','(513) 226-7218',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(9,'cengland','england','','claris','caron','','','(711) 787-5278','(912) 565-7081','(571) 671-1706','cengland9244@hotmail.com','clarise338@bdiddy.com',0,0,'15022 New Hampshire','','Savannah','WI','03925','USA',9,'2000-01-01','2000-01-01',0,5,5,5,'','','cirabaxter','(358) 901-3173',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(10,'rgreer','greer','','rozanne','mitsuko','','','(968) 323-1891','(415) 458-7937','(149) 737-9376','rozanneg5933@yahoo.com','rgreer5735@hotmail.com',0,0,'46646 Mesquite','','Chandler','VT','71404','USA',19,'2000-01-01','2000-01-01',0,5,5,7,'','','latricewalton','(776) 699-8667',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(11,'amckay','mckay','','aleisha','barbera','','','(687) 936-2893','(136) 838-3304','(962) 607-2470','amckay9026@hotmail.com','aleisham2403@comcast.net',0,0,'46822 South Dakota','','Roseville','AR','81598','USA',3,'2000-01-01','2000-01-01',0,2,5,2,'','','laynemorris','(834) 734-8807',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(12,'lphelps','phelps','','lavonna','jacqueline','','','(946) 729-0558','(171) 822-8341','(938) 147-5236','lavonnap7705@aol.com','lavonnaphelps466@abiz.com',0,0,'33877 Juniper','','Dayton','MD','46260','USA',33,'2000-01-01','2000-01-01',0,3,5,7,'','','krishnafreeman','(798) 704-1934',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(13,'erandolph','randolph','','ermelinda','josef','','','(414) 645-9599','(756) 566-6942','(324) 719-4963','ermelindarandolph104@comcast.net','ermelindarandolph248@yahoo.com',0,0,'98257 Pecan','','San Antonio','AL','00458','USA',4,'2000-01-01','2000-01-01',0,1,5,4,'','','alejandrajohns','(514) 141-9766',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(14,'vhobbs','hobbs','','vickey','jamee','','','(929) 893-4263','(119) 166-2440','(589) 155-8625','vickeyh8599@bdiddy.com','vhobbs9759@bdiddy.com',0,0,'53397 Redwood','','Berkeley','ME','74058','USA',12,'2000-01-01','2000-01-01',0,8,5,1,'','','nathanielmaynard','(408) 798-5449',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(15,'jhaney','haney','','joella','quentin','','','(200) 649-8670','(817) 422-9400','(406) 312-6788','joellahaney436@gmail.com','jhaney6268@aol.com',0,0,'88480 Church','','Shreveport','CO','78467','USA',7,'2000-01-01','2000-01-01',0,2,5,5,'','','vikkimartin','(603) 214-9176',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(16,'jswanson','swanson','','janessa','angelica','','','(446) 461-9100','(525) 826-7240','(611) 390-9728','janessas138@aol.com','jswanson498@hotmail.com',0,0,'55784 New Hampshire','','Warren','ME','92866','USA',34,'2000-01-01','2000-01-01',0,3,5,7,'','','samellaclark','(905) 936-5701',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(17,'jbrock','brock','','jeannie','margart','','','(643) 970-0056','(533) 462-5856','(117) 327-7022','jbrock8700@yahoo.com','jbrock3391@aol.com',0,0,'7677 Lee','','Eugene','MO','99527','USA',5,'2000-01-01','2000-01-01',0,2,5,6,'','','jonahdeleon','(231) 285-8037',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(18,'sburt','burt','','shana','joanna','','','(443) 593-8265','(882) 839-8320','(951) 777-6485','shanab8696@comcast.net','shanaburt110@hotmail.com',0,0,'37762 Second','','Fort Wayne','SD','81619','USA',4,'2000-01-01','2000-01-01',0,2,5,7,'','','ottomacias','(788) 916-0281',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(19,'wmeyer','meyer','','winifred','brittanie','','','(892) 310-8147','(507) 755-9811','(500) 613-3812','winifredmeyer301@yahoo.com','wmeyer784@comcast.net',0,0,'6680 Cypress','','Lowell','TX','90674','USA',16,'2000-01-01','2000-01-01',0,1,5,6,'','','taunyamathis','(423) 671-6217',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(20,'rbarry','barry','','refugio','terese','','','(352) 881-6127','(794) 312-9469','(116) 487-1616','refugiob958@hotmail.com','rbarry1418@hotmail.com',0,0,'80364 South Dakota','','Layton','GA','22488','USA',43,'2000-01-01','2000-01-01',0,1,5,2,'','','robbyncameron','(357) 274-2318',1,0,0,0,0,'2000-01-01','2000-01-01','76498229db1024a919bbbebc3ce6feb702ed3cb6f115926ac5a677b43f387ba3c169dce72c75793f0e5301fe616a3645f71e011f4b21eb0dee6ec849c6d9d42d',1,'2018-12-03 19:58:25',0,'2021-01-24 23:14:26',0,''),(21,'bthorton1','Thorton','','Billy','','','','','','','','',0,0,'','','','','','',0,'2000-01-01','2000-01-01',0,0,0,0,'','','','',0,0,0,0,0,'2000-01-01','2000-01-01','5070ea6fea9d36140f6328e1c811e31b35a01c83fd3950c99a3cc8e026a8bfc990df1726eb253627645f647c3b8988ad2d27fa54e7785bf9959b21682b0805fa',4,'2018-12-03 19:59:08',0,'2021-01-24 23:14:26',0,'');
/*!40000 ALTER TABLE `people` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `roles` (
  `RID` mediumint(9) NOT NULL AUTO_INCREMENT,
  `Name` varchar(25) DEFAULT NULL,
  `Descr` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`RID`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'Administrator','This role has permission to do everything'),(2,'Human Resources','This role has full permissions on people, read and print permissions for Companies and Classes.'),(3,'Finance','This role has full permissions on Companies and Classes, read and print permissions on People.'),(4,'Viewer','This role has read-only permissions on everything. Viewers can modify their own information.'),(5,'Tester','This role is for testing'),(6,'OfficeAdministrator','This role is both HR and Finance.'),(7,'OfficeInfoAdministrator','This role is like Office Administrator but also enables delete.');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sessions`
--

DROP TABLE IF EXISTS `sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sessions` (
  `UID` bigint(20) NOT NULL,
  `UserName` varchar(40) NOT NULL DEFAULT '',
  `Cookie` varchar(40) NOT NULL DEFAULT '',
  `DtExpire` datetime NOT NULL DEFAULT '2000-01-01 00:00:00',
  `UserAgent` varchar(256) NOT NULL DEFAULT '',
  `IP` varchar(40) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sessions`
--

LOCK TABLES `sessions` WRITE;
/*!40000 ALTER TABLE `sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `sessions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-01-24 15:14:45
