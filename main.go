package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"math/rand"
	"crypto/sha1"
	"crypto/md5"
	"os"
	"fmt"
	//"io"
	"os/exec"
	//"bytes"
	//"log"
	"io/ioutil"
	"bufio"
	"regexp"
	"encoding/base64"
	"github.com/elliotchance/phpserialize"
	"encoding/json"
	"strings"
	"gopkg.in/yaml.v3"
	"database/sql"
	"net/url"
	_ "github.com/go-sql-driver/mysql"
)

//----------------------------------------------------------
// Constant variables
const WHMAPI1 = "/usr/sbin/whmapi1"
const CPAPI2 = "/usr/bin/cpapi2"
const UAPI = "/usr/bin/uapi"
const MYSQLCli = "/usr/bin/mysql"
const RSYNCCli = "/usr/bin/rsync"
const CHOWNCli = "/bin/chown"
const BASHCli = "/bin/bash"

// config.php
const CONFIG_PATH = "/home/%s/public_html/admin/includes/config.php"

const DEFINE_CONFIG_PHP = "define('%s',"

// iem_stash_storage.php
const IEM_STASH_STORAGE_PATH = "/home/%s/public_html/admin/com/storage/iem_stash_storage.php"

// Config.php
const CONFIG_PHP = "PD9waHAKCmRlZmluZSgnU0VORFNUVURJ" +
	"T19EQVRBQkFTRV9UWVBFJywgJ215c3Fs" +
	"Jyk7CmRlZmluZSgnU0VORFNUVURJT19E" +
	"QVRBQkFTRV9VU0VSJywgJyVzJyk7CmRl" +
	"ZmluZSgnU0VORFNUVURJT19EQVRBQkFT" +
	"RV9QQVNTJywgJyVzJyk7CmRlZmluZSgn" +
	"U0VORFNUVURJT19EQVRBQkFTRV9IT1NU" +
	"JywgJ2xvY2FsaG9zdCcpOwpkZWZpbmUo" +
	"J1NFTkRTVFVESU9fREFUQUJBU0VfTkFN" +
	"RScsICclcycpOwpkZWZpbmUoJ1NFTkRT" +
	"VFVESU9fREFUQUJBU0VfVVRGOFBBVENI" +
	"JywgJzEnKTsKZGVmaW5lKCdTRU5EU1RV" +
	"RElPX1RBQkxFUFJFRklYJywgJ2VtYWls" +
	"XycpOwpkZWZpbmUoJ1NFTkRTVFVESU9f" +
	"TElDRU5TRUtFWScsICdUSFZMVERUQlon" +
	"KTsKZGVmaW5lKCdTRU5EU1RVRElPX0FQ" +
	"UExJQ0FUSU9OX1VSTCcsICclcycpOwpk" +
	"ZWZpbmUoJ1NFTkRTVFVESU9fSVNfU0VU" +
	"VVAnLCAxKTsKZGVmaW5lKCdTRU5EU1RV" +
	"RElPX0RFRkFVTFRDSEFSU0VUJywgJ1VU" +
	"Ri04Jyk7"

// Storage.php
const CONFIG_STORAGE = "ewogICAgIkludGVyc3BpcmVFdmVudCI6IHsKICAgICAgICAiSUVNX1NZU1RFTV9TVEFSVFVQX0JFRk9SRSI6IFtdLAogICAgICAgICJJRU1fU1lTVEVNX1NUQVJUVVBfQUZURVIiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vc3lzdGVtbG9nL3N5c3RlbWxvZy5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19zeXN0ZW1sb2ciLAogICAgICAgICAgICAgICAgICAgICAgICAiU2V0UHJ1bmVMb2ciCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NZU1RFTV9TSFVURE9XTl9CRUZPUkUiOiBbXSwKICAgICAgICAiSUVNX1NZU1RFTV9TSFVURE9XTl9BRlRFUiI6IFtdLAogICAgICAgICJJRU1fU1RBVFNBUElfUkVDT1JET1BFTiI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fUFVCTElDX1BBVEglfS9mdW5jdGlvbnMvYXBpL3RyaWdnZXJlbWFpbHMucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJUcmlnZ2VyRW1haWxzX0FQSSIsCiAgICAgICAgICAgICAgICAgICAgICAgICJldmVudEVtYWlsT3BlbiIKICAgICAgICAgICAgICAgICAgICBdCiAgICAgICAgICAgICAgICB9LAogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vZW1haWxldmVudGxvZy9lbWFpbGV2ZW50bG9nLnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX2VtYWlsZXZlbnRsb2ciLAogICAgICAgICAgICAgICAgICAgICAgICAiT3BlblRyYWNrZWQiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NUQVRTQVBJX1JFQ09SRExJTktDTElDSyI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fUFVCTElDX1BBVEglfS9mdW5jdGlvbnMvYXBpL3RyaWdnZXJlbWFpbHMucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJUcmlnZ2VyRW1haWxzX0FQSSIsCiAgICAgICAgICAgICAgICAgICAgICAgICJldmVudExpbmtDbGlja2VkIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0sCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9lbWFpbGV2ZW50bG9nL2VtYWlsZXZlbnRsb2cucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJBZGRvbnNfZW1haWxldmVudGxvZyIsCiAgICAgICAgICAgICAgICAgICAgICAgICJMaW5rQ2xpY2tlZCIKICAgICAgICAgICAgICAgICAgICBdCiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIF0KICAgICAgICB9LAogICAgICAgICJJRU1fTUFSS0VSIjogW10sCiAgICAgICAgIklFTV9TRU5EQVBJX1NFTkRUT1JFQ0lQSUVOVCI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9lbWFpbGV2ZW50bG9nL2VtYWlsZXZlbnRsb2cucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJBZGRvbnNfZW1haWxldmVudGxvZyIsCiAgICAgICAgICAgICAgICAgICAgICAgICJDYW1wYWlnblNlbnQiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX0pPQlNUUklHR0VSRU1BSUxTQVBJX1BST0NFU1NKT0JTRU5EIjogewogICAgICAgICAgICAiMDUwICI6IFsKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L2VtYWlsZXZlbnRsb2cvZW1haWxldmVudGxvZy5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19lbWFpbGV2ZW50bG9nIiwKICAgICAgICAgICAgICAgICAgICAgICAgIlRyaWdnZXJTZW50IgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9KT0JTQVVUT1JFU1BPTkRFUkFQSV9BQ1RJT05KT0IiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vZW1haWxldmVudGxvZy9lbWFpbGV2ZW50bG9nLnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX2VtYWlsZXZlbnRsb2ciLAogICAgICAgICAgICAgICAgICAgICAgICAiQXV0b3Jlc3BvbmRlclNlbnQiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NFTkRTVFVESU9GVU5DVElPTlNfR0VORVJBVEVURVhUTUVOVUxJTktTIjogewogICAgICAgICAgICAiMDUwICI6IFsKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L2RiY2hlY2svZGJjaGVjay5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19kYmNoZWNrIiwKICAgICAgICAgICAgICAgICAgICAgICAgIkdldFRleHRNZW51SXRlbXMiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfSwKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L3N5c3RlbWxvZy9zeXN0ZW1sb2cucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJBZGRvbnNfc3lzdGVtbG9nIiwKICAgICAgICAgICAgICAgICAgICAgICAgIkdldFRleHRNZW51SXRlbXMiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfSwKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L2NoZWNrcGVybWlzc2lvbnMvY2hlY2twZXJtaXNzaW9ucy5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19jaGVja3Blcm1pc3Npb25zIiwKICAgICAgICAgICAgICAgICAgICAgICAgIkdldFRleHRNZW51SXRlbXMiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NFTkRTVFVESU9GVU5DVElPTlNfR0VORVJBVEVNRU5VTElOS1MiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vZHluYW1pY2NvbnRlbnR0YWdzL2R5bmFtaWNjb250ZW50dGFncy5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19keW5hbWljY29udGVudHRhZ3MiLAogICAgICAgICAgICAgICAgICAgICAgICAiU2V0TWVudUl0ZW1zIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0sCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9zcGxpdHRlc3Qvc3BsaXR0ZXN0LnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX3NwbGl0dGVzdCIsCiAgICAgICAgICAgICAgICAgICAgICAgICJTZXRNZW51SXRlbXMiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1VTRVJBUElfR0VUUEVSTUlTU0lPTlRZUEVTIjogewogICAgICAgICAgICAiMDUwICI6IFsKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L2R5bmFtaWNjb250ZW50dGFncy9keW5hbWljY29udGVudHRhZ3MucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJJbnRlcnNwaXJlX0FkZG9ucyIsCiAgICAgICAgICAgICAgICAgICAgICAgICJHZXRBZGRvblBlcm1pc3Npb25zIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0sCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9zcGxpdHRlc3Qvc3BsaXR0ZXN0LnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiSW50ZXJzcGlyZV9BZGRvbnMiLAogICAgICAgICAgICAgICAgICAgICAgICAiR2V0QWRkb25QZXJtaXNzaW9ucyIKICAgICAgICAgICAgICAgICAgICBdCiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIF0KICAgICAgICB9LAogICAgICAgICJJRU1fRENUX0hUTUxFRElUT1JfVElOWU1DRVBMVUdJTiI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9keW5hbWljY29udGVudHRhZ3MvZHluYW1pY2NvbnRlbnR0YWdzLnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX2R5bmFtaWNjb250ZW50dGFncyIsCiAgICAgICAgICAgICAgICAgICAgICAgICJEY3RUaW55TUNFUGx1Z2luSG9vayIKICAgICAgICAgICAgICAgICAgICBdCiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIF0KICAgICAgICB9LAogICAgICAgICJJRU1fRURJVE9SX1RBR19CVVRUT04iOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vZHluYW1pY2NvbnRlbnR0YWdzL2R5bmFtaWNjb250ZW50dGFncy5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19keW5hbWljY29udGVudHRhZ3MiLAogICAgICAgICAgICAgICAgICAgICAgICAiQ3JlYXRlSW5zZXJ0VGFnQnV0dG9uIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9BRERPTl9EWU5BTUlDQ09OVEVOVFRBR1NfR0VUQUxMVEFHUyI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9keW5hbWljY29udGVudHRhZ3MvZHluYW1pY2NvbnRlbnR0YWdzLnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX2R5bmFtaWNjb250ZW50dGFncyIsCiAgICAgICAgICAgICAgICAgICAgICAgICJnZXRBbGxUYWdzIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9BRERPTl9EWU5BTUlDQ09OVEVOVFRBR1NfUkVQTEFDRVRBR0NPTlRFTlQiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vZHluYW1pY2NvbnRlbnR0YWdzL2R5bmFtaWNjb250ZW50dGFncy5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19keW5hbWljY29udGVudHRhZ3MiLAogICAgICAgICAgICAgICAgICAgICAgICAicmVwbGFjZVRhZ0NvbnRlbnQiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NFVFRJTkdTQVBJX0xPQURTRVRUSU5HUyI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9zcGxpdHRlc3Qvc3BsaXR0ZXN0LnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX3NwbGl0dGVzdCIsCiAgICAgICAgICAgICAgICAgICAgICAgICJTZXRTZXR0aW5ncyIKICAgICAgICAgICAgICAgICAgICBdCiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIF0KICAgICAgICB9LAogICAgICAgICJJRU1fQ1JPTl9SVU5BRERPTlMiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vc3BsaXR0ZXN0L2Nyb24ucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiAiU3BsaXR0ZXN0X0Nyb25fR2V0Sm9icyIKICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9TRU5EU1RVRElPRlVOQ1RJT05TX0NMRUFOVVBPTERRVUVVRVMiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vc3BsaXR0ZXN0L3NwbGl0dGVzdC5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19zcGxpdHRlc3QiLAogICAgICAgICAgICAgICAgICAgICAgICAiQ2xlYW51cFBhcnRpYWxTZW5kcyIKICAgICAgICAgICAgICAgICAgICBdCiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIF0KICAgICAgICB9LAogICAgICAgICJJRU1fTkVXU0xFVFRFUlNBUElfREVMRVRFIjogewogICAgICAgICAgICAiMDUwICI6IFsKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L3NwbGl0dGVzdC9zcGxpdHRlc3QucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJBZGRvbnNfc3BsaXR0ZXN0IiwKICAgICAgICAgICAgICAgICAgICAgICAgIkRlbGV0ZU5ld3NsZXR0ZXJzIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9ORVdTTEVUVEVSU19NQU5BR0VORVdTTEVUVEVSUyI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9zcGxpdHRlc3Qvc3BsaXR0ZXN0LnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX3NwbGl0dGVzdCIsCiAgICAgICAgICAgICAgICAgICAgICAgICJNYW5hZ2VOZXdzbGV0dGVycyIKICAgICAgICAgICAgICAgICAgICBdCiAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgIF0KICAgICAgICB9LAogICAgICAgICJJRU1fSk9CU0FQSV9HRVRKT0JMSVNUIjogewogICAgICAgICAgICAiMDUwICI6IFsKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L3NwbGl0dGVzdC9zcGxpdHRlc3QucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJBZGRvbnNfc3BsaXR0ZXN0IiwKICAgICAgICAgICAgICAgICAgICAgICAgIkdlbmVyYXRlSm9iTGlzdFF1ZXJ5IgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9KT0JTQVBJX0dFVEpPQlNUQVRVUyI6IHsKICAgICAgICAgICAgIjA1MCAiOiBbCiAgICAgICAgICAgICAgICB7CiAgICAgICAgICAgICAgICAgICAgImZpbGUiOiAieyVJRU1fQURET05TX1BBVEglfS9zcGxpdHRlc3Qvc3BsaXR0ZXN0LnBocCIsCiAgICAgICAgICAgICAgICAgICAgImZ1bmN0aW9uIjogWwogICAgICAgICAgICAgICAgICAgICAgICAiQWRkb25zX3NwbGl0dGVzdCIsCiAgICAgICAgICAgICAgICAgICAgICAgICJHZXRKb2JTdGF0dXMiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NDSEVEVUxFX1BBVVNFSk9CIjogewogICAgICAgICAgICAiMDUwICI6IFsKICAgICAgICAgICAgICAgIHsKICAgICAgICAgICAgICAgICAgICAiZmlsZSI6ICJ7JUlFTV9BRERPTlNfUEFUSCV9L3NwbGl0dGVzdC9zcGxpdHRlc3QucGhwIiwKICAgICAgICAgICAgICAgICAgICAiZnVuY3Rpb24iOiBbCiAgICAgICAgICAgICAgICAgICAgICAgICJBZGRvbnNfc3BsaXR0ZXN0IiwKICAgICAgICAgICAgICAgICAgICAgICAgIlBhdXNlU2NoZWR1bGUiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NDSEVEVUxFX0VESVRKT0IiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vc3BsaXR0ZXN0L3NwbGl0dGVzdC5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19zcGxpdHRlc3QiLAogICAgICAgICAgICAgICAgICAgICAgICAiRWRpdFNjaGVkdWxlIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9TQ0hFRFVMRV9SRVNVTUVKT0IiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vc3BsaXR0ZXN0L3NwbGl0dGVzdC5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19zcGxpdHRlc3QiLAogICAgICAgICAgICAgICAgICAgICAgICAiUmVzdW1lU2NoZWR1bGUiCiAgICAgICAgICAgICAgICAgICAgXQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICAgICAgfSwKICAgICAgICAiSUVNX1NDSEVEVUxFX0RFTEVURUpPQlMiOiB7CiAgICAgICAgICAgICIwNTAgIjogWwogICAgICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJmaWxlIjogInslSUVNX0FERE9OU19QQVRIJX0vc3BsaXR0ZXN0L3NwbGl0dGVzdC5waHAiLAogICAgICAgICAgICAgICAgICAgICJmdW5jdGlvbiI6IFsKICAgICAgICAgICAgICAgICAgICAgICAgIkFkZG9uc19zcGxpdHRlc3QiLAogICAgICAgICAgICAgICAgICAgICAgICAiRGVsZXRlU2NoZWR1bGVzIgogICAgICAgICAgICAgICAgICAgIF0KICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgXQogICAgICAgIH0sCiAgICAgICAgIklFTV9NQVJLRVJfMjAwOTA3MDEiOiBbXSwKICAgICAgICAiSW50ZXJzcGlyZVRlbXBsYXRlX0JlZm9yZVVuY2FjaGVkVGVtcGxhdGVQYXJzZWQiOiBbXSwKICAgICAgICAiSW50ZXJzcGlyZVRlbXBsYXRlX0FmdGVyVW5jYWNoZWRUZW1wbGF0ZVBhcnNlZCI6IFtdLAogICAgICAgICJJbnRlcnNwaXJlVGVtcGxhdGVfQmVmb3JlVGVtcGxhdGVDYWNoZWQiOiBbXSwKICAgICAgICAiSW50ZXJzcGlyZVRlbXBsYXRlX0FmdGVyVGVtcGxhdGVDYWNoZWQiOiBbXSwKICAgICAgICAiSW50ZXJzcGlyZVRlbXBsYXRlX0JlZm9yZVRlbXBsYXRlSW5jbHVkZWQiOiBbXSwKICAgICAgICAiSW50ZXJzcGlyZVRlbXBsYXRlX0FmdGVyVGVtcGxhdGVJbmNsdWRlZCI6IFtdLAogICAgICAgICJJbnRlcnNwaXJlVGVtcGxhdGVfVGVtcGxhdGVDYXB0dXJlZCI6IFtdCiAgICB9LAogICAgIklFTV9TWVNURU1fU0VUVElOR1MiOiB7CiAgICAgICAgIlNNVFBfU0VSVkVSIjogImVtLnZpbmFob3N0LnZuIiwKICAgICAgICAiU01UUF9VU0VSTkFNRSI6ICIlcyIsCiAgICAgICAgIlNNVFBfUEFTU1dPUkQiOiAiJXMiLAogICAgICAgICJTTVRQX1BPUlQiOiAiMjUiLAogICAgICAgICJCT1VOQ0VfQUREUkVTUyI6ICIiLAogICAgICAgICJCT1VOQ0VfU0VSVkVSIjogIiIsCiAgICAgICAgIkJPVU5DRV9VU0VSTkFNRSI6ICIiLAogICAgICAgICJCT1VOQ0VfUEFTU1dPUkQiOiAiIiwKICAgICAgICAiQk9VTkNFX0lNQVAiOiAiMCIsCiAgICAgICAgIkJPVU5DRV9FWFRSQVNFVFRJTkdTIjogIiIsCiAgICAgICAgIkJPVU5DRV9BR1JFRURFTEVURSI6ICIwIiwKICAgICAgICAiQk9VTkNFX0FHUkVFREVMRVRFQUxMIjogIjAiLAogICAgICAgICJIVE1MRk9PVEVSIjogIiIsCiAgICAgICAgIlRFWFRGT09URVIiOiAiIiwKICAgICAgICAiRk9SQ0VfVU5TVUJMSU5LIjogIjAiLAogICAgICAgICJNQVhIT1VSTFlSQVRFIjogIjAiLAogICAgICAgICJNQVhPVkVSU0laRSI6ICIwIiwKICAgICAgICAiQ1JPTl9FTkFCTEVEIjogIjAiLAogICAgICAgICJERUZBVUxUQ0hBUlNFVCI6ICJVVEYtOCIsCiAgICAgICAgIkVNQUlMX0FERFJFU1MiOiAiJXMiLAogICAgICAgICJJUFRSQUNLSU5HIjogIjEiLAogICAgICAgICJVU0VNVUxUSVBMRVVOU1VCU0NSSUJFIjogIjAiLAogICAgICAgICJDT05UQUNUQ0FOTU9ESUZZRU1BSUwiOiAiMCIsCiAgICAgICAgIk1BWF9JTUFHRVdJRFRIIjogIjcwMCIsCiAgICAgICAgIk1BWF9JTUFHRUhFSUdIVCI6ICI0MDAiLAogICAgICAgICJBTExPV19FTUJFRElNQUdFUyI6ICIxIiwKICAgICAgICAiREVGQVVMVF9FTUJFRElNQUdFUyI6ICIwIiwKICAgICAgICAiQUxMT1dfQVRUQUNITUVOVFMiOiAiMSIsCiAgICAgICAgIkFUVEFDSE1FTlRfU0laRSI6ICIyMDQ4IiwKICAgICAgICAiQ1JPTl9TRU5EIjogIjUiLAogICAgICAgICJDUk9OX0FVVE9SRVNQT05ERVIiOiAiMTAiLAogICAgICAgICJDUk9OX0JPVU5DRSI6ICI2MCIsCiAgICAgICAgIkNST05fVFJJR0dFUkVNQUlMU19TIjogIjAiLAogICAgICAgICJDUk9OX1RSSUdHRVJFTUFJTFNfUCI6ICIxNDQwIiwKICAgICAgICAiQ1JPTl9NQUlOVEVOQU5DRSI6ICIwIiwKICAgICAgICAiRU1BSUxTSVpFX1dBUk5JTkciOiAiNTAwIiwKICAgICAgICAiRU1BSUxTSVpFX01BWElNVU0iOiAiMjA0OCIsCiAgICAgICAgIlNZU1RFTV9NRVNTQUdFIjogIiIsCiAgICAgICAgIlNZU1RFTV9EQVRBQkFTRV9WRVJTSU9OIjogIjUuNi40NSIsCiAgICAgICAgIlNFTkRfVEVTVF9NT0RFIjogIjAiLAogICAgICAgICJSRVNFTkRfTUFYSU1VTSI6ICIzIiwKICAgICAgICAiU0hPV19TTVRQQ09NX09QVElPTiI6ICIwIiwKICAgICAgICAiU0VDVVJJVFlfV1JPTkdfTE9HSU5fV0FJVCI6ICIwIiwKICAgICAgICAiU0VDVVJJVFlfV1JPTkdfTE9HSU5fVEhSRVNIT0xEX0NPVU5UIjogIjAiLAogICAgICAgICJTRUNVUklUWV9XUk9OR19MT0dJTl9USFJFU0hPTERfRFVSQVRJT04iOiAiNjAiLAogICAgICAgICJTRUNVUklUWV9CQU5fRFVSQVRJT04iOiAiNjAiLAogICAgICAgICJDUkVESVRfSU5DTFVERV9BVVRPUkVTUE9OREVSUyI6ICIxIiwKICAgICAgICAiQ1JFRElUX0lOQ0xVREVfVFJJR0dFUlMiOiAiMSIsCiAgICAgICAgIkNSRURJVF9XQVJOSU5HUyI6ICIwIiwKICAgICAgICAiREVGQVVMVF9FTUFJTFNJWkUiOiAiMCIsCiAgICAgICAgIkNST05fU1BMSVRURVNUIjogIjAiCiAgICB9LAogICAgIklFTV9TRVRUSU5HU19XSElURUxBQkVMIjogWwogICAgICAgIHsKICAgICAgICAgICAgIm5hbWUiOiAiQVBQTElDQVRJT05fRkFWSUNPTiIsCiAgICAgICAgICAgICJ2YWx1ZSI6ICJpbWFnZXMvZmF2aWNvbi5pY28iCiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJuYW1lIjogIkFQUExJQ0FUSU9OX0xPR09fSU1BR0UiLAogICAgICAgICAgICAidmFsdWUiOiAiaW1hZ2VzL2xvZ28uanBnIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJMTkdfQWNjb3VudFVwZ3JhZGVNZXNzYWdlIiwKICAgICAgICAgICAgInZhbHVlIjogIllvdSBhcmUgb24gZGF5ICUldHJpYWxfZGF5c19jdXJyZW50JSUgb2YgeW91ciBmcmVlIHRyaWFsLiBZb3UgaGF2ZSAlJXRyaWFsX2RheXNfbGVmdCUlIGRheXMgcmVtYWluaW5nLiIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICAgIm5hbWUiOiAiTE5HX0FwcGxpY2F0aW9uVGl0bGUiLAogICAgICAgICAgICAidmFsdWUiOiAiVmluYWhvc3QiCiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJuYW1lIjogIkxOR19Db3B5cmlnaHQiLAogICAgICAgICAgICAidmFsdWUiOiAiIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJMTkdfRGVmYXVsdF9HbG9iYWxfSFRNTF9Gb290ZXIiLAogICAgICAgICAgICAidmFsdWUiOiAiIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJMTkdfRGVmYXVsdF9HbG9iYWxfVGV4dF9Gb290ZXIiLAogICAgICAgICAgICAidmFsdWUiOiAiIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJMTkdfRnJlZVRyaWFsX0V4cGlyeV9Mb2dpbiIsCiAgICAgICAgICAgICJ2YWx1ZSI6ICJZb3VyIGZyZWUgdHJpYWwgYWNjb3VudCBoYXMgZXhwaXJlZC4gUGxlYXNlIGNvbnRhY3QgeW91ciBhZG1pbiIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICAgIm5hbWUiOiAiU0hPV19JTlRST19WSURFTyIsCiAgICAgICAgICAgICJ2YWx1ZSI6ICIwIgogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJTSE9XX1NNVFBfQ09NX09QVElPTiIsCiAgICAgICAgICAgICJ2YWx1ZSI6ICIxIgogICAgICAgIH0KICAgIF0KfQ=="

// Config pattern
const SENDSTUDIO_DATABASE_USER = "SENDSTUDIO_DATABASE_USER"
const SENDSTUDIO_DATABASE_PASS = "SENDSTUDIO_DATABASE_PASS"
const SENDSTUDIO_DATABASE_NAME = "SENDSTUDIO_DATABASE_NAME"

const SMTP_USERNAME = "SMTP_USERNAME"
const SMTP_PASSWORD = "SMTP_PASSWORD"
const EMAIL_ADDRESS = "EMAIL_ADDRESS"

const SMTP_USERNAME_CFG = "\"SMTP_USERNAME\": \"%s\""
const SMTP_PASSWORD_CFG = "\"SMTP_PASSWORD\": \"%s\""
const EMAIL_ADDRESS_CFG = "\"EMAIL_ADDRESS\": \"%s\""

// Default cpanel var
const (
	Def_Ctemail		= "haond@vinahost.vn"
	Def_Pkgname	= "Start"
	Def_Owner		= "vinahost"
)

// Struct config API
type ConfigAPI struct {
	Log struct {
		Dir string					`yaml:"dir"`
		File_Name string				`yaml:"file_name"`
	}								`yaml:"LogAudit"`
	Credential struct {
		User string					`yaml:"user"`
		Password string				`yaml:"password"`
	}								`yaml:"Credential"`
	DB_Sample string					`yaml:"DB_Sample"`
	Skeleton string					`yaml:"Skeleton"`
	Clients []string					`yaml:"ListIPClient"`
}

// Struct Config php source
type ConfigInfo struct {
	Action			string			`json:"action"`
	Reason			string			`json:"reason"`
	User			string			`json:"user"`
	Password		string			`json:"password"`
	Domain			string			`json:"domain"`
	Email			string			`json:"email"`
	App_url			string			`json:"app_url"`
	Pkgname		string			`json:"pkgname"`
	map_cfgphp		map[string]string
	map_cfgstorage	map[string]string
}

// Struct Response
type Response struct {
	Success bool					`json:"success"`
	Message string				`json:"message"`
}

//----------------------------------------------------------
// Global variables

// Lock map for deduplicate
var Lock_Map map[string]bool = make(map[string]bool)

// Clients map for check allow
var Clients_Map map[string]bool = make(map[string]bool)

// Config API
var Cfg_API ConfigAPI

// Action map
var Action_Map = map[string]bool {
	"create": true,
	"suspend": true,
	"unsuspend": true,
	"terminate": true,
	"changepackage": true,
	"changepassword": true,
}

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
	u, err := url.Parse(str)
	if err != nil {
		return str
	}
	return u.String()
}

// Remove ', ());
func removeWeirdCharacter(str string) string {
	out := strings.Replace(str, "'", "", -1)
	out = strings.Replace(out, ",", "", -1)
	out = strings.Replace(out, " ", "", -1)
	out = strings.Replace(out, ")", "", -1)
	out = strings.Replace(out, "(", "", -1)
	out = strings.Replace(out, ";", "", -1)
	out = strings.Replace(out, " ", "", -1)
	return out
}

//----------------------------------------------------------
// Function zone

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

// Create account cpanel in WHM
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
		return
	}
	defer f.Close()

	// Write message to file
	t := time.Now()
	if _, err = f.WriteString("\n[" + t.Format("2006.01.02 15:04:05") + "] " + msg); err != nil {
		return
	}
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

// Check cpanel account exist
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
	cfgstorage_file := fmt.Sprintf(IEM_STASH_STORAGE_PATH, cpacc)
	content, err := ioutil.ReadFile(cfgstorage_file)
	if err != nil {
		return err
	}
	// Predefine string separator
	define_dbuser := fmt.Sprintf(DEFINE_CONFIG_PHP, SENDSTUDIO_DATABASE_USER)
	define_dbname := fmt.Sprintf(DEFINE_CONFIG_PHP, SENDSTUDIO_DATABASE_NAME)
	define_dbpass := fmt.Sprintf(DEFINE_CONFIG_PHP, SENDSTUDIO_DATABASE_PASS)

	// Find config
	// Database user
	content_str := string(content)
	re := regexp.MustCompile(define_dbuser)
	cfg_dbuser := re.FindString(content_str)
	res_dbuser := strings.Split(cfg_dbuser, define_dbuser)
	res_dbuser = removeWeirdCharacter(res_dbuser[0])
	re = regexp.MustCompile(define_dbname)
	cfg_dbname := re.FindString(content_str)
	res_dbname := strings.Split(cfg_dbuser, define_dbuser)
	res_dbname = removeWeirdCharacter(res_dbuser[0])
	re = regexp.MustCompile(define_dbpass)
	cfg_dbpass := re.FindString(content_str)
	res_dbpass := strings.Split(cfg_dbuser, define_dbuser)
	res_dbpass = removeWeirdCharacter(res_dbuser[0])
	
	// Connect to db
	db, err := dbConn(cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER],
			cfg.map_cfgphp[SENDSTUDIO_DATABASE_PASS],
			cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME])
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

// Update configure
func updateConfigure(cfg ConfigInfo) error {
	// Phrase 1. Update admin/includes/config.php
	sDec_cfgphp, _ := base64.StdEncoding.DecodeString(CONFIG_PHP)
	out_cfgphp := fmt.Sprintf(string(sDec_cfgphp),
		cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER],
		cfg.map_cfgphp[SENDSTUDIO_DATABASE_PASS],
		cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME],
		cfg.App_url,
	)

	cfgphp_file := fmt.Sprintf(CONFIG_PATH, cfg.User)
	err := ioutil.WriteFile(cfgphp_file, []byte(out_cfgphp), 0644)

	// Phrase 2. Update admin/com/storage/iem_stash_storage.php
	sDec_cfgstorage, _ := base64.StdEncoding.DecodeString(CONFIG_STORAGE)
	prt_SMTP_USERNAME := fmt.Sprintf(SMTP_USERNAME_CFG, cfg.map_cfgstorage[SMTP_USERNAME])
	prt_SMTP_PASSWORD := fmt.Sprintf(SMTP_PASSWORD_CFG, base64.StdEncoding.EncodeToString([]byte(cfg.map_cfgstorage[SMTP_PASSWORD])))
	prt_EMAIL_ADDRESS := fmt.Sprintf(EMAIL_ADDRESS_CFG, cfg.map_cfgstorage[EMAIL_ADDRESS])

	out_cfgstorage := strings.Replace(string(sDec_cfgstorage), SMTP_USERNAME_CFG, prt_SMTP_USERNAME, -1)
	out_cfgstorage = strings.Replace(out_cfgstorage, SMTP_PASSWORD_CFG, prt_SMTP_PASSWORD, -1)
	out_cfgstorage = strings.Replace(out_cfgstorage, EMAIL_ADDRESS_CFG, prt_EMAIL_ADDRESS, -1)

	dat, err := Json_decode(out_cfgstorage)
	if err != nil {
		return err
	}
	out, err := phpserialize.Marshal(dat, nil)

	out_cfg := "<?php /*" + string(out)
	cfgstorage_file := fmt.Sprintf(IEM_STASH_STORAGE_PATH, cfg.User)
	err = ioutil.WriteFile(cfgstorage_file, []byte(out_cfg), 0644)
	return err
}

// Post handler
func postHandler(c *gin.Context) {
	// Prepare configure
	var cfg ConfigInfo
	var response Response

	// Get IP client, if blocked
	if Clients_Map[c.ClientIP()] == false {
		response.Success = false
		response.Message = "Your IP: " + c.ClientIP() + " has blocked"
		writeAuditLog(response.Message)
		c.JSON(http.StatusForbidden, response)
		return
	}

	err := c.BindJSON(&cfg)
	// Check bind json fail
	if err != nil {
		response.Success = false
		response.Message = "Your json data parse failed."
		writeAuditLog(response.Message)
		c.JSON(http.StatusForbidden, response)
		return
	}

	// Check post action support
	if Action_Map[cfg.Action] != true {
		response.Success = false
		response.Message = "Post action: " + cfg.Action + " not allow."
		writeAuditLog(response.Message)
		c.JSON(http.StatusOK, response)
		return
	}

	// Check lock map
	if Lock_Map[cfg.User] {
		response.Success = false
		response.Message = "Error user: " + cfg.User + " on handle."
		writeAuditLog(response.Message)
		c.JSON(http.StatusOK, response)
		return
	}

	// Lock
	Lock_Map[cfg.User] = true
	// Defer unlock
	defer func() {delete(Lock_Map, cfg.User)}()

	//cfg.User = "haond1"
	//cfg.Password = "xxxxx"
	//cfg.Domain = "haond1.com"
	//cfg.Email = "userhaond1.com"
	//cfg.App_url = "https://12347.em.vinahost.vn"

	// Action create
	if cfg.Action == "create" {
		// Generate configure
		cfg.map_cfgphp = make(map[string]string)
		cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER] = cfg.User + "_db"
		cfg.map_cfgphp[SENDSTUDIO_DATABASE_PASS] = StringRand(16)
		cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] = cfg.User + "_db"

		cfg.map_cfgstorage = make(map[string]string)
		cfg.map_cfgstorage[SMTP_USERNAME] = cfg.Email
		cfg.map_cfgstorage[SMTP_PASSWORD] = StringRand(16)
		cfg.map_cfgstorage[EMAIL_ADDRESS] = cfg.Email

		// Create cpanel account
		out, err := createCpanelAccount(cfg.User, cfg.Domain, cfg.Email, cfg.Pkgname, "vinahost", StringRand(16))
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error create account cpanel: " + cfg.User + ", " + err.Error()
			} else {
				response.Message = "Error create account cpanel: " + cfg.User + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check create cpanel account false
		out, check := getReasonCreateCpanelAccount(out)
		if check == false {
			response.Success = false
			response.Message = "Error create account cpanel: " + cfg.User + ", " + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success create account cpanel: " + cfg.User
			writeAuditLog(response.Message)
		}

		// Rsync skeleton
		target := "/home/" + cfg.User + "/public_html/"
		out, err = rsyncSkeleton(target)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error rsync skeleton: " + target + ", " + err.Error()
			} else {
				response.Message = "Error rsync skeleton: " + target + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success rsync skeleton: " + target
			writeAuditLog(response.Message)
		}

		// Chown skeleton
		out, err = chownSkeleton(cfg.User, cfg.User, target)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error chown skeleton: " + err.Error()
			} else {
				response.Message = "Error chown skeleton: " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success chown skeleton: "
			writeAuditLog(response.Message)
		}

		// Update configure php
		err = updateConfigure(cfg)
		if err != nil {
			response.Success = false
			response.Message = "Error update configure php source: " + err.Error()
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success update configure php source: "
			writeAuditLog(response.Message)
		}

		// Create database
		out, err = createDatabase(cfg.User, cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME])
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error create database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " + err.Error()
			} else {
				response.Message = "Error create database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check create database false
		out, check = getReasonCreateDatabase(out)
		if check == false {
			response.Success = false
			response.Message = "Error create database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success create database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME]
			writeAuditLog(response.Message)
		}

		// Create db user
		out, err = createDBUser(cfg.User, cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER], cfg.map_cfgphp[SENDSTUDIO_DATABASE_PASS])
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error create dbuser: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER] + ", " + err.Error()
			} else {
				response.Message = "Error create dbuser: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER] + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check create db user false
		out, check = getReasonCreateDBUser(out)
		if check == false {
			response.Success = false
			response.Message = "Error create dbuser: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER] + ", " + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success create dbuser: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER]
			writeAuditLog(response.Message)
		}

		// Grant all privileges
		out, err = grantAllPrivileges(cfg.User, cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME], cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER])
		if err != nil {
			response.Success = false
			response.Message = "Error grant all privileges: " +
							cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER] +
							"to database " +
							cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " +
							err.Error() + "\n" + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check grant all privileges false
		out, check = getReasonGrantAllPrivileges(out)
		if check == false {
			response.Success = false
			response.Message = "Error grant all privileges: user " +
							cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER] +
							"to database " +
							cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " +
							string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success grant all privileges: user " +
							cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER] +
							"to database " +
							cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME]
			writeAuditLog(response.Message)
		}

		// Import database
		out, err = importDatabase(cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME], Cfg_API.DB_Sample)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error import database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " + err.Error()
			} else {
				response.Message = "Error import database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Message = "Success import database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME]
			writeAuditLog(response.Message)
		}
		
		// Create alias
		alias := removeScheme(cfg.App_url)
		out, err = addAliasDomain(cfg.User, alias)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error create alias: " + alias + ", " + err.Error()
			} else {
				response.Message = "Error create alias: " + alias + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}

		// Run autoSSL - do exclude domain
		out, err = doExcludeDomain(cfg.User, cfg.Domain, removeScheme(cfg.App_url))
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error do exclude autoSSL: " + cfg.User + ", " + err.Error()
			} else {
				response.Message = "Error do exclude autoSSL: " + cfg.User + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Run autoSSL - do autoSSL check
		out, err = doAutoSSLCheck(cfg.User)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error do run autoSSL: " + err.Error()
			} else {
				response.Message = "Error do run autoSSL: " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}

		// Create email account
		stringSlice := strings.Split(cfg.Email, "@")
		email_user := stringSlice[0]
		out, err = createEmailAccount(cfg.User, cfg.Domain, email_user, cfg.Password)

		// Update email.users in database
		// Create db connection
		db, err := dbConn(cfg.map_cfgphp[SENDSTUDIO_DATABASE_USER],
			cfg.map_cfgphp[SENDSTUDIO_DATABASE_PASS],
			cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME])
		if err != nil {
			response.Success = false
			response.Message = "Error connect to database: " + cfg.map_cfgphp[SENDSTUDIO_DATABASE_NAME] + ", " + err.Error()
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}

		// Get token and password hash
		unique_token := generateUniqueToken(cfg.User)
		pass_hash := generatePasswordHash(cfg.Password, unique_token)
		// Run update row
		err = updateUserRow(db, cfg.User, unique_token, pass_hash, cfg.Email)
		if err != nil {
			response.Success = false
			response.Message = "Error update token/password: '" + cfg.Password + "' for user " + cfg.User + ", " + err.Error()
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		db.Close()

		// Announce success
		response.Success = true
		response.Message = "Install success: " + cfg.User + "/" + cfg.Domain
		writeAuditLog(response.Message)
		c.JSON(http.StatusOK, response)
		return
	}

	// Action suspend
	if cfg.Action == "suspend" {
		// Suspend cpanel account
		out, err := suspendCpanelAccount(cfg.User, cfg.Reason)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error suspend account cpanel: " + cfg.User + ", " + err.Error()
			} else {
				response.Message = "Error suspend account cpanel: " + cfg.User + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check suspend account fail
		out, check := getReasonSuspendCpanelAccount(out)
		if check == false {
			response.Success = false
			response.Message = "Error suspend account cpanel: " + cfg.User + ", " + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Success = true
			response.Message = "Success suspend account cpanel: " + cfg.User
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// Action unsuspend
	if cfg.Action == "unsuspend" {
		// Unsuspend cpanel account
		out, err := unsuspendCpanelAccount(cfg.User)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error unsuspend account cpanel: " + cfg.User + ", " + err.Error()
			} else {
				response.Message = "Error unsuspend account cpanel: " + cfg.User + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check unsuspend account fail
		out, check := getReasonUnsuspendCpanelAccount(out)
		if check == false {
			response.Success = false
			response.Message = "Error unsuspend account cpanel: " + cfg.User + ", " + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Success = true
			response.Message = "Success unsuspend account cpanel: " + cfg.User
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// Action terminate/remove
	if cfg.Action == "terminate" {
		// Terminate/remove account
		out, err := removeCpanelAccount(cfg.User)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error remove account cpanel: " + cfg.User + ", " + err.Error()
			} else {
				response.Message = "Error remove account cpanel: " + cfg.User + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check remove account fail
		out, check := getReasonRemoveCpanelAccount(out)
		if check == false {
			response.Success = false
			response.Message = "Error remove account cpanel: " + cfg.User + ", " + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Success = true
			response.Message = "Success remove account cpanel: " + cfg.User
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// Action change package
	if cfg.Action == "changepackage" {
		// Change package
		out, err := changePackageCpanelAccount(cfg.User)
		if err != nil {
			response.Success = false
			if string(out) == "" {
				response.Message = "Error change package account cpanel: " + cfg.User + ", " + err.Error()
			} else {
				response.Message = "Error change package account cpanel: " + cfg.User + ", " + err.Error() + "\n" + string(out)
			}
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
		// Check change package fail
		out, check := getReasonChangePackageCpanelAccount(out)
		if check == false {
			response.Success = false
			response.Message = "Error change package account cpanel: " + cfg.User + ", " + string(out)
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Success = true
			response.Message = "Success change package account cpanel: " + cfg.User
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// Action change password dash
	if cfg.Action == "changepackage" {
		// Change password dash
		err := changePasswordDash(cfg.User, cfg.Password)
		if err != nil {
			response.Success = false
			response.Message = "Error change password dash user: " + cfg.User + " - " + cfg.Password + ", " + err.Error()
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		} else {
			response.Success = true
			response.Message = "Success change password dash user: " + cfg.User + " - " + cfg.Password
			writeAuditLog(response.Message)
			c.JSON(http.StatusOK, response)
			return
		}
	}
}

//-----------------------------------------------------

// Main function
func main() {
	// Read configure file
	Cfg_API = readConfigAPI("/root/go/src/github.com/gn1k/api-em/config.yaml")
	mapClients(Cfg_API)
	Cfg_API.Skeleton = addSlash(Cfg_API.Skeleton)
	Cfg_API.Log.Dir = addSlash(Cfg_API.Log.Dir)

	router := gin.Default()

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		Cfg_API.Credential.User: Cfg_API.Credential.Password,
	}))

	authorized.POST("/", postHandler)
	router.RunTLS(":12345", "/etc/api-em/vinahost.crt", "/etc/api-em/vinahost.key")
}
