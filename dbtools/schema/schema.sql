-- ACCORD PHONEBOOK DATABASE
-- mysql> show grants for 'ec2-user'@'localhost';
-- +-----------------------------------------------------------------------------+
-- | Grants for ec2-user@localhost                                               |
-- +-----------------------------------------------------------------------------+
-- | GRANT USAGE ON *.* TO 'ec2-user'@'localhost'                                |
-- | GRANT ALL PRIVILEGES ON accordtest.* TO 'ec2-user'@'localhost'              |
-- | GRANT ALL PRIVILEGES ON accordtest.accordtest TO 'ec2-user'@'localhost'     |
-- +-----------------------------------------------------------------------------+

DROP DATABASE IF EXISTS accordtest;
CREATE DATABASE accordtest;
USE accordtest;
GRANT ALL PRIVILEGES ON accordtest TO 'ec2-user'@'localhost';
GRANT ALL PRIVILEGES ON accordtest.* TO 'ec2-user'@'localhost';
CREATE TABLE classes (
    ClassCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(25) NOT NULL DEFAULT '',
    Designation CHAR(3) NOT NULL DEFAULT '',
    Description VARCHAR(256) NOT NULL DEFAULT '',
    LastModTime TIMESTAMP,
    LastModBy MEDIUMINT NOT NULL DEFAULT 0,
    PRIMARY KEY (ClassCode)
);

CREATE TABLE companies (
    CoCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    LegalName VARCHAR(50) NOT NULL DEFAULT '',
    CommonName VARCHAR(50) NOT NULL DEFAULT '',
    Address VARCHAR(35) NOT NULL DEFAULT '',
    Address2 VARCHAR(35) NOT NULL DEFAULT '',
    City VARCHAR(25) NOT NULL DEFAULT '',
    State CHAR(25) NOT NULL DEFAULT '',
    PostalCode VARCHAR(10) NOT NULL DEFAULT '',
    Country VARCHAR(25) NOT NULL DEFAULT '',
    Phone VARCHAR(25) NOT NULL DEFAULT '',
    Fax VARCHAR(25) NOT NULL DEFAULT '',
    Email VARCHAR(35) NOT NULL DEFAULT '',
    Designation CHAR(3) NOT NULL NOT NULL DEFAULT '',
    Active SMALLINT NOT NULL DEFAULT 0,
    EmploysPersonnel SMALLINT NOT NULL DEFAULT 0,
    LastModTime TIMESTAMP,
    LastModBy MEDIUMINT NOT NULL DEFAULT 0,
    PRIMARY KEY (CoCode)
);

CREATE TABLE compensation (
    UID MEDIUMINT NOT NULL,
    Type MEDIUMINT NOT NULL
);

CREATE TABLE deductions (
    UID MEDIUMINT NOT NULL,
    Deduction INT NOT NULL
);

CREATE TABLE deductionlist (
    DCode MEDIUMINT NOT NULL,
    Name VARCHAR(25) NOT NULL
);

CREATE TABLE departments (
    DeptCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(25),
    PRIMARY KEY (DeptCode)
);

CREATE TABLE fieldperms (
    RID MEDIUMINT NOT NULL,
    Elem MEDIUMINT NOT NULL,
    Field VARCHAR(25) NOT NULL,
    Perm MEDIUMINT NOT NULL,
    Descr VARCHAR(256)
);

CREATE TABLE jobtitles (
    JobCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    Title VARCHAR(40) NOT NULL DEFAULT '',
    Descr VARCHAR(256) NOT NULL DEFAULT '',
    PRIMARY KEY (JobCode)
);

CREATE TABLE people (
    UID MEDIUMINT NOT NULL AUTO_INCREMENT,
    UserName VARCHAR(20) NOT NULL,
    LastName VARCHAR(25) NOT NULL DEFAULT '',
    MiddleName VARCHAR(25) NOT NULL DEFAULT '',
    FirstName VARCHAR(25) NOT NULL DEFAULT '',
    PreferredName VARCHAR(25) NOT NULL DEFAULT '',
    Salutation VARCHAR(10) NOT NULL DEFAULT '',
    PositionControlNumber VARCHAR(10) NOT NULL DEFAULT '',
    OfficePhone VARCHAR(25) NOT NULL DEFAULT '',
    OfficeFax VARCHAR(25) NOT NULL DEFAULT '',
    CellPhone VARCHAR(25) NOT NULL DEFAULT '',
    PrimaryEmail VARCHAR(35) NOT NULL DEFAULT '',
    SecondaryEmail VARCHAR(35) NOT NULL DEFAULT '',
    BirthMonth TINYINT NOT NULL DEFAULT 0,
    BirthDoM TINYINT NOT NULL DEFAULT 0,
    HomeStreetAddress VARCHAR(35) NOT NULL DEFAULT '',
    HomeStreetAddress2 VARCHAR(25) NOT NULL DEFAULT '',
    HomeCity VARCHAR(25) NOT NULL DEFAULT '',
    HomeState CHAR(2) NOT NULL DEFAULT '',
    HomePostalCode varchar(10) NOT NULL DEFAULT '',
    HomeCountry VARCHAR(25) NOT NULL DEFAULT '',
    JobCode MEDIUMINT NOT NULL DEFAULT 0,
    Hire DATE NOT NULL DEFAULT '2000-01-01 00:00:00',
    Termination DATE NOT NULL DEFAULT '2000-01-01 00:00:00',
    MgrUID MEDIUMINT NOT NULL DEFAULT 0,
    DeptCode MEDIUMINT NOT NULL DEFAULT 0,
    CoCode MEDIUMINT NOT NULL DEFAULT 0,
    ClassCode SMALLINT NOT NULL DEFAULT 0,
    StateOfEmployment VARCHAR(25) NOT NULL DEFAULT '',
    CountryOfEmployment VARCHAR(25) NOT NULL DEFAULT '',
    EmergencyContactName VARCHAR(25) NOT NULL DEFAULT '',
    EmergencyContactPhone VARCHAR(25) NOT NULL DEFAULT '',
    Status SMALLINT NOT NULL DEFAULT 0,
    EligibleForRehire SMALLINT NOT NULL DEFAULT 0,
    HealthInsuranceAccepted SMALLINT NOT NULL DEFAULT 0,
    DentalInsuranceAccepted SMALLINT NOT NULL DEFAULT 0,
    Accepted401K SMALLINT NOT NULL DEFAULT 0,
    LastReview DATE NOT NULL DEFAULT '2000-01-01 00:00:00',
    NextReview DATE NOT NULL DEFAULT '2000-01-01 00:00:00',
    passhash char(128) NOT NULL DEFAULT '',
    RID MEDIUMINT NOT NULL DEFAULT 0,
    LastModTime TIMESTAMP,
    LastModBy MEDIUMINT NOT NULL DEFAULT 0,
    PRIMARY KEY (UID)
);

CREATE TABLE roles (
    RID MEDIUMINT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(25) NOT NULL,
    Descr VARCHAR(256),
    PRIMARY KEY(RID)
);

-- Add the Administrator as the first and only user
INSERT INTO people (UserName,FirstName,LastName) VALUES("administrator","Administrator","Administrator");
