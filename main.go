package main

import (
	//"github.com/gin-gonic/gin"
	//"net/http"
	"time"
	"math/rand"
	"os"
	"fmt"
	//"io"
	"os/exec"
	//"bytes"
	//"log"
	"io/ioutil"
	"bufio"
)

//----------------------------------------------------------
// Constant variables
const WHMCli = "/usr/sbin/whmapi1"
const CPCli = "/usr/bin/cpapi2"
const MYSQLCli = "/usr/bin/mysql"

const (
	Def_Ctemail	= "haond@vinahost.vn"
	Def_Pkgname	= "Start"
	Def_Owner	= "vinahost"
)

// Struct
type Lock struct {
	Task string
}

// Global variables
var lock_array []Lock

//----------------------------------------------------------
// Tool zone

// Random string with fixed length
func StringRand(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

// Read file
func readFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	return content, err
}

// Read file lines
func readFileLines(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return lines, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	file.Close()
	return lines, err
}

//----------------------------------------------------------
// Function zone

// Create account cpanel in WHM
func createAccountWHM(user, domain, ctemail, pkgname, owner, password string) ([]byte, error) {
	if ctemail == "" {
		ctemail = Def_Ctemail
	}
	if pkgname == "" {
		pkgname = Def_Pkgname
	}
	if owner == "" {
		owner = Def_Owner
	}

        // Arguments
        args := []string{
		"createacct",
                "username=" + user,
                "domain=" + domain,
                "cgi=1",
                "contactemail=" + ctemail,
                "cpmod=paper_lantern",
                "dkim=1",
                "owner=" + owner,
                "pkgname=" + pkgname,
                "reseller=0",
		"password=" + password,
	}

	// Run cmd
	cmd := exec.Command(WHMCli, args...)
        out, err := cmd.CombinedOutput()
	return out, err
}

// Create database
func createDatabase(user, db string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + user,
		"MysqlFE",
		"createdb",
		"db=" + db,
	}

	// Run cmd
	cmd := exec.Command(CPCli, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Create dbuser
func createDBUser(user, dbuser, password string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + user,
		"MysqlFE",
		"createdbuser",
		"dbuser=" + dbuser,
		"password=" + password,
	}

	// Run cmd
	cmd := exec.Command(CPCli, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Grant privileges to dbuser
func grantAllPrivileges(user, db, dbuser string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + user,
		"MysqlFE",
		"setdbuserprivileges",
		"privileges=ALL%20PRIVILEGES",
		"db=" + db,
		"dbuser=" + dbuser,
	}

	// Run cmd
	cmd := exec.Command(CPCli, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Import database
func importDatabase(db, sample string) ([]byte, error) {
	// Arguments
	args := []string{
		"db",
		"<",
		"sample",
	}

	// Run cmd
	cmd := exec.Command(CPCli, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Update configure
func updateConfigure(working_dir string) {
	content, err := readFile("/etc/profile")
	if err != nil {
		fmt.Println("Got some error.")
		return
	}

	
	//
	fmt.Println(string(content))
}

// Post handler
func postHandler() {
	//fmt.Printf("combined out:\n%s\n", string(out))
}

//-----------------------------------------------------
// Main function
func main() {
	//router := gin.Default()

	//s := &http.Server{
	//	Addr:           ":8080",
	//	Handler:        router,
	//	ReadTimeout:    10 * time.Second,
	//	WriteTimeout:   10 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//router.GET("/", func(c *gin.Context) {
	//	var msg struct {
	//		Name    string `json:"user"`
	//		Message string
	//		Number  int
	//	}
	//	msg.Name = "Lena"
	//	msg.Message = "hey"
	//	msg.Number = 123
	//	
	//	c.JSON(http.StatusOK, msg)
	//})

	//s.ListenAndServe()

	//out, err := createAccountWHM("haond1", "domain.com", "alo@gmail.com", "Start", "vinahost")
	//if err != nil {
        //        log.Fatalf("cmd.Run() failed with %s\n", err)
        //} else {
        //        fmt.Printf("combined out:\n%s\n", string(out))
        //}

	//out, err := createDBUser("haond", "haond_user1", "vinahost@888")
	//fmt.Println(string(out))
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//out, err := grantAllPrivileges("haond", "haond_db", "haond_user1")
	//fmt.Println(string(out))


	lines, err := readFileLines("/etc/profile")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, eachline := range lines {
		fmt.Println(eachline)
	}
}
