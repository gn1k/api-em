package main

import {
	"os/exec"
	"regexp"
	"encoding/base64"
	"github.com/elliotchance/phpserialize"
	"gopkg.in/yaml.v3"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
}

// Constant variables
const WHMAPI1 = "/usr/sbin/whmapi1"
const CPAPI2 = "/usr/bin/cpapi2"
const UAPI = "/usr/bin/uapi"
const MYSQLCli = "/usr/bin/mysql"
const RSYNCCli = "/usr/bin/rsync"
const CHOWNCli = "/bin/chown"
const BASHCli = "/bin/bash"

// Rsync skeleton
func rsyncSkeleton(target string) ([]byte, error) {
	// Arguments
	args := []string{
		"-ar",
		"--delete",
		Cfg_API.Skeleton,
		target,
	}

	// Run cmd
	cmd := exec.Command(RSYNCCli, args...)
	out, err := cmd.CombinedOutput()
	return out, err	
}

// Chown skeleton
func chownSkeleton(owner, group, target string) ([]byte, error) {
	// Arguments
	args := []string{
		"-R",
		owner + ":" + group,
		target,
	}

	// Run cmd
	cmd := exec.Command(CHOWNCli, args...)
	out, err := cmd.CombinedOutput()
	return out, err	
}

// Map clients
func mapClients(cfgapi ConfigAPI) {
	for _, element := range cfgapi.Clients {
		Clients_Map[element] = true
	}
}

// Get restictions
func getRestrictionAccount(user string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"Mysql",
		"get_restrictions",
	}

	// Run cmd
	cmd := exec.Command(UAPI, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Get prefix (restriction) database
func getPrefixDatabase(prefix string) string {
	re := regexp.MustCompile(".*prefix: .*")
	prefix_line := re.FindString(prefix)
	prefix_res := strings.Split(prefix_line, "prefix: ")
	return prefix_res[len(prefix_res) - 1]
}

// Create cpanel account in WHM
func createCpanelAccount(user, domain, ctemail, pkgname, owner, password string) ([]byte, error) {
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
		"username=" + URL_encode(user),
		"domain=" + URL_encode(domain),
		"cgi=1",
		"cpmod=paper_lantern",
		"dkim=1",
		"owner=" + URL_encode(owner),
		"pkgname=" + URL_encode(pkgname),
		"reseller=0",
		"password=" + URL_encode(password),
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Create database
func createDatabase(user, db string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"MysqlFE",
		"createdb",
		"db=" + URL_encode(db),
	}

	// Run cmd
	cmd := exec.Command(CPAPI2, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Create dbuser
func createDBUser(user, dbuser, password string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"MysqlFE",
		"createdbuser",
		"dbuser=" + URL_encode(dbuser),
		"password=" + URL_encode(password),
	}

	// Run cmd
	cmd := exec.Command(CPAPI2, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Grant privileges to dbuser
func grantAllPrivileges(user, db, dbuser string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"MysqlFE",
		"setdbuserprivileges",
		"privileges=ALL%20PRIVILEGES",
		"db=" + db,
		"dbuser=" + URL_encode(dbuser),
	}

	// Run cmd
	cmd := exec.Command(CPAPI2, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Create SMTP email account
func createEmailAccount(user, domain, email, pass string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"Email",
		"addpop",
		"domain=" + URL_encode(domain),
		"email=" + URL_encode(email),
		"password=" + URL_encode(pass),
		"quota=unlimited",
	}

	// Run cmd
	cmd := exec.Command(CPAPI2, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Import database
func importDatabase(db, sample string) ([]byte, error) {
	// Arguments
	args := "mysql " + db + " < " + sample

	// Run cmd
	cmd := exec.Command(BASHCli, "-c", args)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Add alias/parked domain
func addAliasDomain(user, domain string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"Park",
		"park",
		"domain=" + URL_encode(domain),
	}

	// Run cmd
	cmd := exec.Command(CPAPI2, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Get exclude domain SSL
func getExcludeDomain(domain, alias string) (string) {
	output := domain + "," +
		"mail." + domain + "," +
		"www." + domain + "," +
		"webmail." + domain + "," +
		"cpanel." + domain + "," +
		"webdisk." + domain + "," +
		"mail." + alias + "," +
		"www." + alias
	return output
}

// Remove scheme http/https
func removeScheme(domain string) (string) {
	output := domain
	if strings.Contains(domain, "https://") {
		output = strings.Replace(output, "https://", "", -1)
	}
	if strings.Contains(domain, "http://") {
		output = strings.Replace(output, "http://", "", -1)
	}
	return output
}

// Do exclude domain SSL
func doExcludeDomain(user, domain, alias string) ([]byte, error) {
	list_exclude := URL_encode(getExcludeDomain(domain, alias))
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"SSL",
		"add_autossl_excluded_domains",
		"domains=" + list_exclude,
	}

	// Run cmd
	cmd := exec.Command(UAPI, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Do autossl check
func doAutoSSLCheck(user string) ([]byte, error) {
	// Arguments
	args := []string{
		"--user=" + URL_encode(user),
		"SSL",
		"start_autossl_check",
	}

	// Run cmd
	cmd := exec.Command(UAPI, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Create db connection
func dbConn(db_user, db_pass, db_name string) (*sql.DB, error) {
	dbDriver := "mysql"
	db, err := sql.Open(dbDriver, db_user + ":" + db_pass + "@/" + db_name)
	return db, err
}

// Update user row
func updateUserRow(db *sql.DB, username, unique_token, pass_hash, emailaddress string) (error) {
	dbName := "email_users"
	insFrom, err := db.Prepare("update " + dbName + " set username=?, unique_token=?, password=?, emailaddress=? where userid=1")
	if err != nil {
		return err
	}
	insFrom.Exec(username, unique_token, pass_hash, emailaddress)
	return nil
}

// Update user pass only
func updateUserPass(db *sql.DB, username, unique_token, pass_hash string) (error) {
	dbName := "email_users"
	insFrom, err := db.Prepare("update " + dbName + " set username=?, unique_token=?, password=? where userid=1")
	if err != nil {
		return err
	}
	insFrom.Exec(username, unique_token, pass_hash)
	return nil
}

// Check cpanel account exist
func checkCPAcct(acc string) bool {
	// Arguments
	args := []string{
		"listaccts",
		"search=" + URL_encode(acc),
		"searchtype=user",
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	// Sub string
	sub := "user: " + acc
	return strings.Contains(string(out), sub)
}

// Check cpanel account suspended
func checkCPAcctSuspended(acc string) bool {
	// Arguments
	args := []string{
		"listsuspended",
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	// Sub string
	sub := "user: " + acc
	return strings.Contains(string(out), sub)
}

// Suspend cpanel account
func suspendCpanelAccount(acc, reason string) ([]byte, error) {
	// Arguments
	args := []string{
		"suspendacct",
		"user=" + URL_encode(acc),
		"reason=" + URL_encode(reason),
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Unsuspend cpanel account
func unsuspendCpanelAccount(acc string) ([]byte, error) {
	// Arguments
	args := []string{
		"unsuspendacct",
		"user=" + URL_encode(acc),
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Terminate/Remove cpanel account
func removeCpanelAccount(acc string) ([]byte, error) {
	// Arguments
	args := []string{
		"removeacct",
		"user=" + URL_encode(acc),
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Change package cpanel account
func changePackageCpanelAccount(acc, pkgname string) ([]byte, error) {
	// Arguments
	args := []string{
		"changepackage",
		"user=" + URL_encode(acc),
		"pkg=" + URL_encode(pkgname),
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}

// Change password dash
func changePasswordDash(cpacc, pass string) error {
	// Read file php config
	cfgstorage_file := fmt.Sprintf(CONFIG_PATH, cpacc)
	content, err := ioutil.ReadFile(cfgstorage_file)
	if err != nil {
		return err
	}
	// Predefine string separator
	define_dbuser_regexp := fmt.Sprintf(DEFINE_CONFIG_PHP_REGEXP, SENDSTUDIO_DATABASE_USER)
	define_dbname_regexp := fmt.Sprintf(DEFINE_CONFIG_PHP_REGEXP, SENDSTUDIO_DATABASE_NAME)
	define_dbpass_regexp := fmt.Sprintf(DEFINE_CONFIG_PHP_REGEXP, SENDSTUDIO_DATABASE_PASS)
	define_dbuser := fmt.Sprintf(DEFINE_CONFIG_PHP, SENDSTUDIO_DATABASE_USER)
	define_dbname := fmt.Sprintf(DEFINE_CONFIG_PHP, SENDSTUDIO_DATABASE_NAME)
	define_dbpass := fmt.Sprintf(DEFINE_CONFIG_PHP, SENDSTUDIO_DATABASE_PASS)

	// Find config
	// Database user
	content_str := string(content)
	re := regexp.MustCompile(".*" + define_dbuser_regexp + ".*")
	cfg_dbuser := re.FindString(content_str)
	res_dbuser := strings.Split(cfg_dbuser, define_dbuser)
	dbuser := removeWeirdCharacter(res_dbuser[len(res_dbuser) - 1])
	re = regexp.MustCompile(".*" + define_dbname_regexp + ".*")
	cfg_dbname := re.FindString(content_str)
	res_dbname := strings.Split(cfg_dbname, define_dbname)
	dbname := removeWeirdCharacter(res_dbname[len(res_dbname) - 1])
	re = regexp.MustCompile(".*" + define_dbpass_regexp + ".*")
	cfg_dbpass := re.FindString(content_str)
	res_dbpass := strings.Split(cfg_dbpass, define_dbpass)
	dbpass := removeWeirdCharacter(res_dbpass[len(res_dbpass) - 1])
	
	// Connect to db
	db, err := dbConn(dbuser, dbpass, dbname)
	if err != nil {
		return err 
	}

	// Get token and password hash
	unique_token := generateUniqueToken(cpacc)
	pass_hash := generatePasswordHash(pass, unique_token)

	// Update row password only
	dbName := "email_users"
	insFrom, err := db.Prepare("update " + dbName + " set unique_token=?, password=? where userid=1")
	if err != nil {
		return err
	}
	insFrom.Exec(unique_token, pass_hash)
	return nil
}

// Check domain exist
func checkDomain(domain string) bool {
	// Arguments
	args := []string{
		"listaccts",
		"search=" + URL_encode(domain),
		"searchtype=domain",
		"searchmethod=exact",
	}

	// Run cmd
	cmd := exec.Command(WHMAPI1, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	// Sub string
	sub := "domain: " + domain
	return strings.Contains(string(out), sub)
}

// Get reason createCpanelAccount success or not
func getReasonCreateCpanelAccount(out string) (string, bool) {
	re := regexp.MustCompile(`.*reason: .*`)
	reason := re.FindString(out)
	success := "Account Creation Ok"
	return reason, strings.Contains(out, success)
}

// Get reason suspendCpanelAccount success or not
func getReasonSuspendCpanelAccount(out string) (string, bool) {
	re := regexp.MustCompile(`.*reason: .*`)
	reason := re.FindString(out)
	success := "reason: OK"
	return reason, strings.Contains(out, success)
}

// Get reason unsuspendCpanelAccount success or not
func getReasonUnsuspendCpanelAccount(out string) (string, bool) {
	re := regexp.MustCompile(`.*reason: .*`)
	reason := re.FindString(out)
	success := "reason: OK"
	return reason, strings.Contains(out, success)
}

// Get reason removeCpanelAccount success or not
func getReasonRemoveCpanelAccount(out string) (string, bool) {
	re := regexp.MustCompile(`.*reason: .*`)
	reason := re.FindString(out)
	success := "account removed"
	return reason, strings.Contains(out, success)
}

// Get reason removeCpanelAccount success or not
func getReasonChangePackageCpanelAccount(out string) (string, bool) {
	re := regexp.MustCompile(`.*reason: .*`)
	reason := re.FindString(out)
	success := "reason: Account Upgrade/Downgrade Complete"
	return reason, strings.Contains(out, success)
}

// Get reason createDatabase success or not
func getReasonCreateDatabase(out string) (string, bool) {
	re := regexp.MustCompile(`.*error: .*`)
	reason := re.FindString(out)
	if reason == "" {
		return reason, true
	}
	return reason, false
}

// Get reason createDBUser success or not
func getReasonCreateDBUser(out string) (string, bool) {
	re := regexp.MustCompile(`.*error: .*`)
	reason := re.FindString(out)
	if reason == "" {
		return reason, true
	}
	return reason, false
}

// Get reason grantAllPrivileges success or not
func getReasonGrantAllPrivileges(out string) (string ,bool) {
	re := regexp.MustCompile(`.*error: .*`)
	reason := re.FindString(out)
	if reason == "" {
		return reason, true
	}
	re = regexp.MustCompile(`.*Cpanel::Exception.*`)
	reason_2 := re.FindString(out)
	if reason_2 != "" {
		return reason, false
	}
	return reason, false
}

// Get reason addAliasDomain success or not
func getReasonAddAliasDomain(out string) (string, bool) {
	re := regexp.MustCompile(`.*reason: .*`)
	reason := re.FindString(out)
	success := "The system successfully parked"
	return reason, strings.Contains(out, success)
}

// Get reason createEmailAccount success or not
func getReasonCreateEmailAccount(out string) (string, bool) {
	re := regexp.MustCompile(`.*error: .*`)
	reason := re.FindString(out)
	if reason == "" {
		return reason, true
	}
	return reason, false
}