package main

func buildPreparedStatements() {
	var err error
	Phonebook.prepstmt.deductList, err = Phonebook.db.Prepare("select deduction from deductions where uid=?")
	errcheck(err)
	Phonebook.prepstmt.getComps, err = Phonebook.db.Prepare("select type from compensation where uid=?")
	errcheck(err)
	Phonebook.prepstmt.myDeductions, err = Phonebook.db.Prepare("select dcode,name from deductionlist")
	errcheck(err)
	Phonebook.prepstmt.adminPersonDetails, err = Phonebook.db.Prepare(
		"select LastName,FirstName,MiddleName,Salutation," + // 4
			"ClassCode,Status,PositionControlNumber," + // 7
			"OfficePhone,OfficeFax,CellPhone,PrimaryEmail," + // 11
			"SecondaryEmail,EligibleForRehire,LastReview,NextReview," + // 15
			"BirthMonth,BirthDOM,HomeStreetAddress,HomeStreetAddress2,HomeCity," + // 20
			"HomeState,HomePostalCode,HomeCountry," + // 23
			"AcceptedHealthInsurance,AcceptedDentalInsurance,Accepted401K," + // 26
			"jobcode,hire,termination," + // 29
			"mgruid,deptcode,cocode,StateOfEmployment," + // 33
			"CountryOfEmployment,PreferredName," + // 35
			"EmergencyContactName,EmergencyContactPhone,RID,username " + // 38
			"from people where uid=?")
	errcheck(err)
	Phonebook.prepstmt.classInfo, err = Phonebook.db.Prepare("select classcode,CoCode,Name,Designation,Description from classes where classcode=?")
	errcheck(err)
	Phonebook.prepstmt.companyInfo, err = Phonebook.db.Prepare("select cocode,LegalName,CommonName,Address,Address2,City,State,PostalCode,Country,Phone,Fax,Email,Designation,Active,EmploysPersonnel from companies where cocode=?")
	errcheck(err)
	Phonebook.prepstmt.GetAllCompanies, err = Phonebook.db.Prepare("select cocode,LegalName,CommonName,Address,Address2,City,State,PostalCode,Country,Phone,Fax,Email,Designation,Active,EmploysPersonnel from companies")
	errcheck(err)
	Phonebook.prepstmt.countersUpdate, err = Phonebook.db.Prepare("update counters set SearchPeople=SearchPeople+?,SearchClasses=SearchClasses+?," +
		"SearchCompanies=SearchCompanies+?,EditPerson=EditPerson+?,ViewPerson=ViewPerson+?,ViewClass=ViewClass+?,ViewCompany=ViewCompany+?," +
		"AdminEditPerson=AdminEditPerson+?,AdminEditClass=AdminEditClass+?,AdminEditCompany=AdminEditCompany+?,DeletePerson=DeletePerson+?," +
		"DeleteClass=DeleteClass+?,DeleteCompany=DeleteCompany+?,SignIn=SignIn+?,Logoff=Logoff+?")
	errcheck(err)
	Phonebook.prepstmt.delClass, err = Phonebook.db.Prepare("DELETE FROM classes WHERE ClassCode=?")
	errcheck(err)
	Phonebook.prepstmt.delCompany, err = Phonebook.db.Prepare("select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode from people where cocode=?")
	errcheck(err)
	Phonebook.prepstmt.delPerson, err = Phonebook.db.Prepare("DELETE FROM people WHERE UID=?")
	errcheck(err)
	Phonebook.prepstmt.delPersonComp, err = Phonebook.db.Prepare("DELETE FROM compensation WHERE UID=?")
	errcheck(err)
	Phonebook.prepstmt.delPersonDeduct, err = Phonebook.db.Prepare("DELETE FROM deductions WHERE UID=?")
	errcheck(err)
	Phonebook.prepstmt.getJobTitle, err = Phonebook.db.Prepare("select title from jobtitles where jobcode=?")
	errcheck(err)
	Phonebook.prepstmt.nameFromUID, err = Phonebook.db.Prepare("select firstname,lastname from people where uid=?")
	errcheck(err)
	Phonebook.prepstmt.deptName, err = Phonebook.db.Prepare("select name from departments where deptcode=?")
	errcheck(err)
	Phonebook.prepstmt.directReports, err = Phonebook.db.Prepare("select uid,lastname,firstname,jobcode,primaryemail,officephone,cellphone from people where mgruid=? AND status>0 order by lastname, firstname")
	errcheck(err)
	Phonebook.prepstmt.personDetail, err = Phonebook.db.Prepare(
		"select lastname,middlename,firstname,preferredname,jobcode,primaryemail," + // 6
			"officephone,cellphone,deptcode,cocode,mgruid,ClassCode," + // 12
			"EmergencyContactName,EmergencyContactPhone," + // 14
			"HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry,OfficeFax " + // 21
			"from people where uid=?")
	errcheck(err)
	Phonebook.prepstmt.adminInsertPerson, err = Phonebook.db.Prepare(
		"INSERT INTO people (Salutation,FirstName,MiddleName,LastName,PreferredName," +
			"EmergencyContactName,EmergencyContactPhone," +
			"PrimaryEmail,SecondaryEmail,OfficePhone,OfficeFax,CellPhone,CoCode,JobCode," +
			"PositionControlNumber,DeptCode," +
			"HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry," +
			"status,EligibleForRehire,Accepted401K,AcceptedDentalInsurance,AcceptedHealthInsurance," +
			"Hire,Termination,ClassCode," +
			"BirthMonth,BirthDOM,mgruid,StateOfEmployment,CountryOfEmployment," +
			"LastReview,NextReview,RID,lastmodby,UserName) " +
			//      1                 10                  20                  30                  40
			"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	errcheck(err)
	Phonebook.prepstmt.adminUpdatePerson, err = Phonebook.db.Prepare(
		"update people set Salutation=?,FirstName=?,MiddleName=?,LastName=?,PreferredName=?," + // 5
			"EmergencyContactName=?,EmergencyContactPhone=?," + // 7
			"PrimaryEmail=?,SecondaryEmail=?,OfficePhone=?,OfficeFax=?,CellPhone=?,CoCode=?,JobCode=?," + // 14
			"PositionControlNumber=?,DeptCode=?," + // 16
			"HomeStreetAddress=?,HomeStreetAddress2=?,HomeCity=?,HomeState=?,HomePostalCode=?,HomeCountry=?," + // 22
			"status=?,EligibleForRehire=?,Accepted401K=?,AcceptedDentalInsurance=?,AcceptedHealthInsurance=?," + // 27
			"Hire=?,Termination=?,ClassCode=?," + // 30
			"BirthMonth=?,BirthDOM=?,mgruid=?,StateOfEmployment=?,CountryOfEmployment=?," + // 35
			"LastReview=?,NextReview=?,lastmodby=?,RID=? " + // 39
			"where people.uid=?")
	errcheck(err)
	Phonebook.prepstmt.adminReadBack, err = Phonebook.db.Prepare("select uid from people where FirstName=? and LastName=? and PrimaryEmail=? and OfficePhone=? and CoCode=? and JobCode=?")
	errcheck(err)
	Phonebook.prepstmt.insertComp, err = Phonebook.db.Prepare("INSERT INTO compensation (uid,type) VALUES(?,?)")
	errcheck(err)
	Phonebook.prepstmt.insertDeduct, err = Phonebook.db.Prepare("INSERT INTO deductions (uid,deduction) VALUES(?,?)")
	errcheck(err)
	Phonebook.prepstmt.insertClass, err = Phonebook.db.Prepare("INSERT INTO classes (CoCode,Name,Designation,Description,lastmodby) VALUES(?,?,?,?,?)")
	errcheck(err)
	Phonebook.prepstmt.classReadBack, err = Phonebook.db.Prepare("select ClassCode from classes where Name=? and Designation=?")
	errcheck(err)
	Phonebook.prepstmt.updateClass, err = Phonebook.db.Prepare("update classes set CoCode=?,Name=?,Designation=?,Description=?,lastmodby=? where ClassCode=?")
	errcheck(err)
	Phonebook.prepstmt.insertCompany, err = Phonebook.db.Prepare("INSERT INTO companies (LegalName,CommonName,Designation," +
		"Email,Phone,Fax,Active,EmploysPersonnel,Address,Address2,City,State,PostalCode,Country,lastmodby) " +
		//      1                 10                  20                  30
		"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	errcheck(err)
	Phonebook.prepstmt.companyReadback, err = Phonebook.db.Prepare("select CoCode from companies where CommonName=? and LegalName=?")
	errcheck(err)
	Phonebook.prepstmt.updateCompany, err = Phonebook.db.Prepare("update companies set LegalName=?,CommonName=?,Designation=?,Email=?,Phone=?,Fax=?,EmploysPersonnel=?,Active=?,Address=?,Address2=?,City=?,State=?,PostalCode=?,Country=?,lastmodby=? where CoCode=?")
	errcheck(err)
	Phonebook.prepstmt.updateMyDetails, err = Phonebook.db.Prepare("update people set PreferredName=?,PrimaryEmail=?,OfficePhone=?,CellPhone=?," +
		"EmergencyContactName=?,EmergencyContactPhone=?," +
		"HomeStreetAddress=?,HomeStreetAddress2=?,HomeCity=?,HomeState=?,HomePostalCode=?,HomeCountry=?,lastmodby=? " +
		"where people.uid=?")
	errcheck(err)
	Phonebook.prepstmt.updatePasswd, err = Phonebook.db.Prepare("update people set passhash=? where uid=?")
	errcheck(err)
	Phonebook.prepstmt.readFieldPerms, err = Phonebook.db.Prepare("select Elem,Field,Perm,Descr from fieldperms where RID=?")
	errcheck(err)
	Phonebook.prepstmt.accessRoles, err = Phonebook.db.Prepare("select RID,Name,Descr from roles")
	errcheck(err)
	Phonebook.prepstmt.getUserCoCode, err = Phonebook.db.Prepare("select cocode from people where uid=?")
	errcheck(err)
	Phonebook.prepstmt.loginInfo, err = Phonebook.db.Prepare("select uid,firstname,preferredname,passhash,rid from people where username=?")
	errcheck(err)
	Phonebook.prepstmt.CompanyClasses, err = Phonebook.db.Prepare("select ClassCode,CoCode,Name,Designation,Description,LastModTime,LastModBy from classes where CoCode=?")
	errcheck(err)
}
