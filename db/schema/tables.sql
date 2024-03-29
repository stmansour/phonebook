CREATE TABLE classes (
    ClassCode MEDIUMINT NOT NULL AUTO_INCREMENT,
    CoCode MEDIUMINT NOT NULL DEFAULT 0,
    Name VARCHAR(50) NOT NULL DEFAULT '',
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
    Email VARCHAR(50) NOT NULL DEFAULT '',
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

CREATE TABLE counters (
    SearchPeople BIGINT NOT NULL DEFAULT 0,
    SearchClasses BIGINT NOT NULL DEFAULT 0,
    SearchCompanies BIGINT NOT NULL DEFAULT 0,
    EditPerson BIGINT NOT NULL DEFAULT 0,
    ViewPerson BIGINT NOT NULL DEFAULT 0,
    ViewClass BIGINT NOT NULL DEFAULT 0,
    ViewCompany BIGINT NOT NULL DEFAULT 0,
    AdminEditPerson BIGINT NOT NULL DEFAULT 0,
    AdminEditClass BIGINT NOT NULL DEFAULT 0,
    AdminEditCompany BIGINT NOT NULL DEFAULT 0,
    DeletePerson BIGINT NOT NULL DEFAULT 0,
    DeleteClass BIGINT NOT NULL DEFAULT 0,
    DeleteCompany BIGINT NOT NULL DEFAULT 0,
    SignIn BIGINT NOT NULL DEFAULT 0,
    Logoff BIGINT NOT NULL DEFAULT 0,
    LastModTime TIMESTAMP
);

INSERT INTO counters (SearchPeople) VALUES(0);

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
    UserName VARCHAR(20) NOT NULL DEFAULT '',
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
    Status SMALLINT NOT NULL DEFAULT 0,                       -- 0 = inactive, 1 = active
    EligibleForRehire SMALLINT NOT NULL DEFAULT 0,
    AcceptedHealthInsurance SMALLINT NOT NULL DEFAULT 0,
    AcceptedDentalInsurance SMALLINT NOT NULL DEFAULT 0,
    Accepted401K SMALLINT NOT NULL DEFAULT 0,
    LastReview DATE NOT NULL DEFAULT '2000-01-01 00:00:00',
    NextReview DATE NOT NULL DEFAULT '2000-01-01 00:00:00',
    passhash char(128) NOT NULL DEFAULT '',
    RID MEDIUMINT NOT NULL DEFAULT 0,
    ImagePath VARCHAR(200) NOT NULL DEFAULT '',
    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY (UID)
);

CREATE TABLE roles (
    RID MEDIUMINT NOT NULL AUTO_INCREMENT,
    Name VARCHAR(25) NOT NULL,
    Descr VARCHAR(256),
    PRIMARY KEY(RID)
);

CREATE TABLE sessions (
    UID BIGINT NOT NULL,
    UserName VARCHAR(40) NOT NULL DEFAULT '',
    Cookie VARCHAR(40) NOT NULL DEFAULT '',
    DtExpire DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
    UserAgent VARCHAR(256) NOT NULL DEFAULT '',
    IP VARCHAR(40) NOT NULL DEFAULT ''
);

-- Add the Administrator as the first and only user
-- INSERT INTO people (UserName,FirstName,LastName) VALUES("administrator","Administrator","Administrator");

-- this table is needed by WithCredentials
CREATE TABLE license (
    LID BIGINT NOT NULL AUTO_INCREMENT,
    UID BIGINT NOT NULL DEFAULT 0,
    State VARCHAR(25) NOT NULL DEFAULT '',
    LicenseNo VARCHAR(128) NOT NULL DEFAULT '',
    FLAGS BIGINT NOT NULL DEFAULT 0,
    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY (LID)
);
