INSERT accordtest.people (LastName, FirstName, MiddleName, Salutation, PositionControlNumber, OfficePhone, OfficeFax, CellPhone, PrimaryEmail,
	SecondaryEmail, birthMonth, birthDoM, HomeStreetAddress, HomeStreetAddress2, HomeCity, HomeState, HomePostalCode, HomeCountry, jobcode, hire, 
	termination, mgruid, deptcode, cocode, StateOfEmployment, CountryOfEmployment, PreferredName, EmergencyContactName, EmergencyContactPhone, 
	status, EligibleForRehire, AcceptedHealthInsurance, AcceptedDentalInsurance, Accepted401K, LastReview, NextReview, username, passhash, 
	classcode, lastmodtime, lastmodby, RID) 
SELECT LastName, FirstName, MiddleName, Salutation, PositionControlNumber, OfficePhone, OfficeFax, CellPhone, PrimaryEmail, SecondaryEmail, 
birthMonth, birthDoM, HomeStreetAddress, HomeStreetAddress2, HomeCity, HomeState, HomePostalCode, HomeCountry, jobcode, hire, termination, 
mgruid, deptcode, cocode, StateOfEmployment, CountryOfEmployment, PreferredName, EmergencyContactName, EmergencyContactPhone, status, 
EligibleForRehire, AcceptedHealthInsurance, AcceptedDentalInsurance, Accepted401K, LastReview, NextReview, username, passhash, classcode, 
lastmodtime, lastmodby, RID from accord.people;


INSERT accordtest.companies (LegalName, CommonName, Address, Address2, City, State, PostalCode, Country, Phone, Fax, Email, Designation, Active, 
	EmploysPersonnel, lastmodby, lastmodtime)
SELECT LegalName, CommonName, Address, Address2, City, State, PostalCode, Country, Phone, Fax, Email, Designation, Active, EmploysPersonnel, 
lastmodby, lastmodtime from accord.companies;


INSERT accordtest.classes (Name, Designation, Description, lastmodby, lastmodtime)
SELECT Name, Designation, Description, lastmodby, lastmodtime from accord.classes;

insert into accordtest.deductionlist select * from accord.deductionlist;
insert into accordtest.fieldperms select * from accord.fieldperms;
insert into accordtest.compensation select * from accord.compensation;
insert into accordtest.deductions select * from accord.deductions;
