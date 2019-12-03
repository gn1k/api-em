package main

import (
	"math/rand"
	"io/ioutil"
	"os"
	"encoding/json"
	"strings"
	"crypto/sha1"
	"crypto/md5"
	"time"
	"bufio"
	"fmt"
	"github.com/elliotchance/phpserialize"
	"gopkg.in/yaml.v3"
)

//----------------------------------------------------------
// Tool zone

// Random string with fixed length
func StringRand(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	//specials := "~=+%^*()[]{}/!@#$?|"
	// Not contain: ();' for php code and mysql syntax
	specials := "~=+^*[]{}/!@#$?|_"
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

// Generate unique token
func generateUniqueToken(username string) string {
	token := StringRand(32)
	secs := time.Now().Unix()
	token = fmt.Sprintf("%x", sha1.Sum([]byte(token))) + fmt.Sprintf("%x", md5.Sum([]byte(string(secs)))) + fmt.Sprintf("%x", md5.Sum([]byte(username)))
	return token
}

// Password
func generatePasswordHash(password, token string) string {
	passwordHash := fmt.Sprintf("%x", md5.Sum([]byte(token))) + fmt.Sprintf("%x", md5.Sum([]byte(password)))
	passwordHash = fmt.Sprintf("%x", md5.Sum([]byte(passwordHash)))
	return passwordHash
}

// Json decode to map
func Json_decode(data string) (interface{}, error) {
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

// Url encode
func URL_encode(str string) string {
	enc := strings.ReplaceAll(str, "!", "%21")
	enc = strings.ReplaceAll(enc, "#", "%23")
	enc = strings.ReplaceAll(enc, "$", "%24")
	enc = strings.ReplaceAll(enc, "*", "%2A")
	enc = strings.ReplaceAll(enc, "+", "%2B")
	enc = strings.ReplaceAll(enc, "/", "%2F")
	enc = strings.ReplaceAll(enc, "=", "%3D")
	enc = strings.ReplaceAll(enc, "?", "%3F")
	enc = strings.ReplaceAll(enc, "@", "%40")
	enc = strings.ReplaceAll(enc, "[", "%5B")
	enc = strings.ReplaceAll(enc, "]", "%5D")
	enc = strings.ReplaceAll(enc, "^", "%5E")
	enc = strings.ReplaceAll(enc, "{", "%7B")
	enc = strings.ReplaceAll(enc, "|", "%7C")
	enc = strings.ReplaceAll(enc, "}", "%7D")
	return enc
}

// Remove ', ();
func removeWeirdCharacter(str string) string {
	out := strings.Replace(str, "', '", "", 1)
	out = strings.Replace(out, "');", "", 1)
	return out
}

// Add slash
func addSlash(str string) string {
	if len(str) == 0 {
		return "/"
	}
	if str[len(str)-1:] != "/" {
		str += "/"
	}
	return str
}

// Write audit log
func writeAuditLog(msg string) {
	// Log dir path
	log_path := Cfg_API.Log.Dir
	_, err := os.Stat(log_path)
	if os.IsNotExist(err) {
		os.MkdirAll(log_path, os.ModePerm)
	}

	// File path
	log_path += Cfg_API.Log.File_Name
	// Open append
	f, err := os.OpenFile(log_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	// Write message to file
	t := time.Now()
	if _, err = f.WriteString("\n[" + t.Format("2006.01.02 15:04:05") + "] " + msg); err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Read config file
func readConfigAPI(file string) (ConfigAPI, error) {
	var cfgapi ConfigAPI
	f, err := os.Open(file)
	if err != nil {
		return cfgapi, err
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfgapi)
	return cfgapi, err
}