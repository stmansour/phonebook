package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func readCommandLineArgs() {
	var sd int64
	sd = time.Now().UnixNano()

	dbuPtr := flag.String("B", "ec2-user", "database user name")
	dbnmPtr := flag.String("N", "accordtest", "database name (accordtest, accord)")
	hPtr := flag.String("h", "localhost", "server hostname")
	pPtr := flag.Int("p", 8250, "port on which the server listens")
	dbgPtr := flag.Bool("d", false, "debug mode when true")
	tPtr := flag.Int("t", 0, "test duration time in minutes. 0 means use iterations")
	uPtr := flag.Int("u", 1, "number of users to simulate")
	iPtr := flag.Int("i", 1, "number of iterations, ignored if test duration time is non-zero")
	p := flag.Int64("s", sd, "seed for random numbers. Default is to use a random seed.")
	flag.Parse()
	App.TestIterations = *iPtr // number of iterations (mutually exclusive with TestDuration)
	App.TestUsers = *uPtr      // number of users to test with
	App.TestDuration = *tPtr   // time in minutes
	App.DBName = *dbnmPtr
	App.DBUser = *dbuPtr
	App.Seed = int64(*p)
	App.Host = *hPtr
	App.Port = *pPtr
	App.Debug = *dbgPtr
	rand.Seed(App.Seed)
}

func loadNames() {
	file, err := os.Open("./firstnames.txt")
	errcheck(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		App.FirstNames = append(App.FirstNames, scanner.Text())
	}
	errcheck(scanner.Err())

	file, err = os.Open("./lastnames.txt")
	errcheck(err)
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		App.LastNames = append(App.LastNames, scanner.Text())
	}
	errcheck(scanner.Err())

	file, err = os.Open("./states.txt")
	errcheck(err)
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		App.States = append(App.States, scanner.Text())
	}
	errcheck(scanner.Err())

	file, err = os.Open("./cities.txt")
	errcheck(err)
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		App.Cities = append(App.Cities, scanner.Text())
	}
	errcheck(scanner.Err())

	file, err = os.Open("./streets.txt")
	errcheck(err)
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		App.Streets = append(App.Streets, scanner.Text())
	}
	errcheck(scanner.Err())

	if App.Debug {
		fmt.Printf("FirstNames: %d\n", len(App.FirstNames))
		fmt.Printf("LastNames: %d\n", len(App.LastNames))
		fmt.Printf("Cities: %d\n", len(App.Cities))
		fmt.Printf("States: %d\n", len(App.States))
		fmt.Printf("Streets: %d\n", len(App.Streets))
	}
}

func readAccessRoles() {
	rows, err := App.db.Query("select RID,Name from Roles")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		var r Role
		errcheck(rows.Scan(&r.RID, &r.Name))
		App.Roles = append(App.Roles, r)
	}
	errcheck(rows.Err())
}

func loadCompanies() {
	var code int
	var name string
	App.CoCodeToName = make(map[int]string)
	App.NameToCoCode = make(map[string]int)

	rows, err := App.db.Query("select cocode,CommonName from companies")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		App.CoCodeToName[code] = name
		App.NameToCoCode[name] = code
	}
	errcheck(rows.Err())
}

func loadClasses() {
	var code int
	var name string

	App.NameToClassCode = make(map[string]int)
	App.ClassCodeToName = make(map[int]string)
	rows, err := App.db.Query("select classcode,designation from classes")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		App.NameToClassCode[name] = code
		App.ClassCodeToName[code] = name
	}
	// for k, v := range App.NameToClassCode {
	// 	fmt.Printf("%s %d\n", k, v)
	// }
	errcheck(rows.Err())
}

func loadMaps() {
	var code int
	var name string

	loadCompanies()
	loadClasses()

	App.NameToJobCode = make(map[string]int)
	rows, err := App.db.Query("select jobcode,title from jobtitles")
	errcheck(err)
	defer rows.Close()
	App.JCLo = 99999
	App.JCHi = 0
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		App.NameToJobCode[name] = code
		if code < App.JCLo {
			App.JCLo = code
		}
		if code > App.JCHi {
			App.JCHi = code
		}
	}
	errcheck(rows.Err())

	App.NameToDeptCode = make(map[string]int)
	rows, err = App.db.Query("select deptcode,name from departments order by name")
	errcheck(err)
	defer rows.Close()
	for rows.Next() {
		errcheck(rows.Scan(&code, &name))
		App.NameToDeptCode[name] = code
		if code < App.DeptLo {
			App.DeptLo = code
		}
		if code > App.DeptHi {
			App.DeptHi = code
		}
	}
	errcheck(rows.Err())
	if App.Debug {
		fmt.Printf("DeptLo=%d, DeptHi=%d\n", App.DeptLo, App.DeptHi)
	}

	App.AcceptCodeToName = make(map[int]string)
	for i := ACPTUNKNOWN; i <= ACPTLAST; i++ {
		App.AcceptCodeToName[i] = acceptIntToString(i)
	}

	App.Months = make([]string, len(fmtMonths))
	for i := 0; i < len(fmtMonths); i++ {
		App.Months[i] = fmtMonths[i]
	}

}
