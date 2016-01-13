package main

func buildPreparedStatements() {
	var err error
	App.prepstmt.deductList, err = App.db.Prepare("select deduction from deductions where uid=?")
	errcheck(err)
	App.prepstmt.getComps, err = App.db.Prepare("select type from compensation where uid=?")
	errcheck(err)
	App.prepstmt.myDeductions, err = App.db.Prepare("select dcode,name from deductionlist")
	errcheck(err)
	App.prepstmt.adminPersonDetails, err = App.db.Prepare(
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
			"EmergencyContactName,EmergencyContactPhone,RID " + // 38
			"from people where uid=?")
	errcheck(err)
	App.prepstmt.classInfo, err = App.db.Prepare("select classcode,Name,Designation,Description from classes where classcode=?")
	errcheck(err)
	App.prepstmt.companyInfo, err = App.db.Prepare("select cocode,LegalName,CommonName,Address,Address2,City,State,PostalCode,Country,Phone,Fax,Email,Designation,Active,EmploysPersonnel from companies where cocode=?")
	errcheck(err)
	App.prepstmt.countersUpdate, err = App.db.Prepare("update counters set SearchPeople=?,SearchClasses=?,SearchCompanies=?,EditPerson=?,ViewPerson=?,ViewClass=?,ViewCompany=?,AdminEditPerson=?,AdminEditClass=?,AdminEditCompany=?,DeletePerson=?,DeleteClass=?,DeleteCompany=?,SignIn=?,Logoff=?")
	errcheck(err)
	App.prepstmt.delClass, err = App.db.Prepare("DELETE FROM classes WHERE ClassCode=?")
	errcheck(err)
	App.prepstmt.delCompany, err = App.db.Prepare("select uid,lastname,firstname,preferredname,jobcode,primaryemail,officephone,cellphone,deptcode from people where cocode=?")
	errcheck(err)
	App.prepstmt.delPerson, err = App.db.Prepare("DELETE FROM people WHERE UID=?")
	errcheck(err)
	App.prepstmt.delPersonComp, err = App.db.Prepare("DELETE FROM compensation WHERE UID=?")
	errcheck(err)
	App.prepstmt.delPersonDeduct, err = App.db.Prepare("DELETE FROM deductions WHERE UID=?")
	errcheck(err)
	App.prepstmt.getJobTitle, err = App.db.Prepare("select title from jobtitles where jobcode=?")
	errcheck(err)
	App.prepstmt.nameFromUID, err = App.db.Prepare("select firstname,lastname from people where uid=?")
	errcheck(err)
	App.prepstmt.deptName, err = App.db.Prepare("select name from departments where deptcode=?")
	errcheck(err)
	App.prepstmt.directReports, err = App.db.Prepare("select uid,lastname,firstname,jobcode,primaryemail,officephone,cellphone from people where mgruid=? AND status>0 order by lastname, firstname")
	errcheck(err)
	App.prepstmt.personDetail, err = App.db.Prepare(
		"select lastname,middlename,firstname,preferredname,jobcode,primaryemail," + // 6
			"officephone,cellphone,deptcode,cocode,mgruid,ClassCode," + // 12
			"EmergencyContactName,EmergencyContactPhone," + // 14
			"HomeStreetAddress,HomeStreetAddress2,HomeCity,HomeState,HomePostalCode,HomeCountry " + // 20
			"from people where uid=?")
	errcheck(err)
	App.prepstmt.adminInsertPerson, err = App.db.Prepare(
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
	App.prepstmt.adminUpdatePerson, err = App.db.Prepare(
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
	App.prepstmt.adminReadBack, err = App.db.Prepare("select uid from people where FirstName=? and LastName=? and PrimaryEmail=? and OfficePhone=? and CoCode=? and JobCode=?")
	errcheck(err)
	App.prepstmt.insertComp, err = App.db.Prepare("INSERT INTO compensation (uid,type) VALUES(?,?)")
	errcheck(err)
	App.prepstmt.insertDeduct, err = App.db.Prepare("INSERT INTO deductions (uid,deduction) VALUES(?,?)")
	errcheck(err)
	App.prepstmt.insertClass, err = App.db.Prepare("INSERT INTO classes (Name,Designation,Description,lastmodby) VALUES(?,?,?,?)")
	errcheck(err)
	App.prepstmt.classReadBack, err = App.db.Prepare("select ClassCode from classes where Name=? and Designation=?")
	errcheck(err)
	App.prepstmt.updateClass, err = App.db.Prepare("update classes set Name=?,Designation=?,Description=?,lastmodby=? where ClassCode=?")
	errcheck(err)
	App.prepstmt.insertCompany, err = App.db.Prepare("INSERT INTO companies (LegalName,CommonName,Designation," +
		"Email,Phone,Fax,Active,EmploysPersonnel,Address,Address2,City,State,PostalCode,Country,lastmodby) " +
		//      1                 10                  20                  30
		"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	errcheck(err)
	App.prepstmt.companyReadback, err = App.db.Prepare("select CoCode from companies where CommonName=? and LegalName=?")
	errcheck(err)
	App.prepstmt.updateCompany, err = App.db.Prepare("update companies set LegalName=?,CommonName=?,Designation=?,Email=?,Phone=?,Fax=?,EmploysPersonnel=?,Active=?,Address=?,Address2=?,City=?,State=?,PostalCode=?,Country=?,lastmodby=? where CoCode=?")
	errcheck(err)
	App.prepstmt.updateMyDetails, err = App.db.Prepare("update people set PreferredName=?,PrimaryEmail=?,OfficePhone=?,CellPhone=?," +
		"EmergencyContactName=?,EmergencyContactPhone=?," +
		"HomeStreetAddress=?,HomeStreetAddress2=?,HomeCity=?,HomeState=?,HomePostalCode=?,HomeCountry=?,lastmodby=? " +
		"where people.uid=?")
	errcheck(err)
	App.prepstmt.updatePasswd, err = App.db.Prepare("update people set passhash=? where uid=?")
	errcheck(err)
	App.prepstmt.readFieldPerms, err = App.db.Prepare("select Elem,Field,Perm,Descr from fieldperms where RID=?")
	errcheck(err)
	App.prepstmt.accessRoles, err = App.db.Prepare("select RID,Name,Descr from roles")
	errcheck(err)
	App.prepstmt.getUserCoCode, err = App.db.Prepare("select cocode from people where uid=?")
	errcheck(err)
	App.prepstmt.loginInfo, err = App.db.Prepare("select uid,firstname,preferredname,passhash,rid from people where username=?")
	errcheck(err)
}
