DROP DATABASE IF EXISTS accordtest;
CREATE DATABASE accordtest;
USE accordtest;
GRANT ALL PRIVILEGES ON Accord TO 'ec2-user'@'localhost';

CREATE TABLE classes (
    ClassCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(25),
    Designation CHAR(3) NOT NULL,
    Description VARCHAR(256),
    LastModTime TIMESTAMP,
    LastModBy MEDIUMINT NOT NULL,
    PRIMARY KEY (ClassCode)
);

CREATE TABLE companies (
    CoCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    LegalName VARCHAR(50) NOT NULL,
    CommonName VARCHAR(50) NOT NULL,
    Address VARCHAR(35),
    Address2 VARCHAR(35),
    City VARCHAR(25),
    State CHAR(25),
    PostalCode VARCHAR(10),
    Country VARCHAR(25),
    Phone VARCHAR(25),
    Fax VARCHAR(25),
    Email VARCHAR(35),
    Designation CHAR(3) NOT NULL,
    Active SMALLINT NOT NULL,
    EmploysPersonnel SMALLINT NOT NULL,
    LastModTime TIMESTAMP,
    LastModBy MEDIUMINT NOT NULL,
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
    Name VARCHAR(25)
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
    PRIMARY KEY (RID)
);

CREATE TABLE jobtitles (
    JobCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    Title VARCHAR(40),
    DeptCode MEDIUMINT,
    Department VARCHAR(25)
);

CREATE TABLE people (
    UID MEDIUMINT NOT NULL AUTO_INCREMENT,
    UserName VARCHAR(20),
    LastName VARCHAR(25),
    MiddleName VARCHAR(25),
    FirstName VARCHAR(25),
    PreferredName VARCHAR(25),
    Salutation VARCHAR(10),
    PositionControlNumber VARCHAR(10),
    OfficePhone VARCHAR(25),
    OfficeFax VARCHAR(25),
    CellPhone VARCHAR(25),
    PrimaryEmail VARCHAR(35),
    SecondaryEmail VARCHAR(35),
    BirthMonth TINYINT,
    BirthDoM TINYINT,
    HomeStreetAddress VARCHAR(35),
    HomeStreetAddress2 VARCHAR(25),
    HomeCity VARCHAR(25),
    HomeState CHAR(2),
    HomePostalCode varchar(10),
    HomeCountry VARCHAR(25),
    JobCode MEDIUMINT,
    Hire DATE,
    Termination DATE,
    MgrUID MEDIUMINT,
    DeptCode MEDIUMINT,
    CoCode MEDIUMINT,
    ClassCode SMALLINT,
    StateOfEmployment VARCHAR(25),
    CountryOfEmployment VARCHAR(25),
    EmergencyContactName VARCHAR(25),
    EmergencyContactPhone VARCHAR(25),
    Status SMALLINT,
    EligibleForRehire SMALLINT,
    HealthInsuranceAccepted SMALLINT,
    DentalInsuranceAccepted SMALLINT,
    Accepted401K SMALLINT,
    LastReview DATE,
    NextReview DATE,
    passhash char(128),
    RID MEDIUMINT,
    LastModTime TIMESTAMP,
    LastModBy MEDIUMINT NOT NULL,
    PRIMARY KEY (UID)
);

CREATE TABLE roles (
    RID MEDIUMINT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(25) NOT NULL,
    Descr VARCHAR(256),
    PRIMARY KEY(RID)
);
