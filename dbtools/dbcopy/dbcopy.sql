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

insert into accordtest.DeductionList select * from accord.DeductionList;
insert into accordtest.FieldPerms select * from accord.FieldPerms;
insert into accordtest.Compensation select * from accord.Compensation;
insert into accordtest.Deductions select * from accord.Deductions;
