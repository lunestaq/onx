package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const PORT               = "port="               // 25
const ROOT_DOMAIN        = "root_domain="        // example.com
const MAIL_DOMAIN        = "mail_domain="        // mail.example.com
const MAILDIR_PATH       = "mail_path="          // ~/Maildir
const enable_disable_tls = 4                     // not yet implemented (tls is always on)
const TLS_FILE_fullchain = "tls_file_fullchain=" // /etc/letsencrypt/live/example.com/fullchain.pem
const TLS_FILE_privkey   = "tls_file_privkey="   // /etc/letsencrypt/live/example.com/privkey.pem
const MAIL_BLACKLIST     = 7                     // not yet implemented

func extract_mail(line string) string {
	return line[strings.Index(line, "<")+1:strings.Index(line, ">")]
}

func is_valid_domain(arg string) bool {
    if len(arg) == 0 || len(arg) > 255 {
        return false
    }

    // Check if argument is an IPv4 address in brackets: [x.x.x.x]
    if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") {
        ip := arg[1 : len(arg)-1]
        parsedIP := net.ParseIP(ip)
        if parsedIP == nil {
            return false
        }
        // Ensure it's IPv4 (not IPv6)
        return parsedIP.To4() != nil
    }

    // Otherwise, validate as domain/subdomain
    var domainRegex = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)(?:\.(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?))*$`)

    return domainRegex.MatchString(arg)
	// thanks chatgpt
}

func save_mail(email_data string) error {
    hostname, err := os.Hostname()
    if err != nil {hostname = "localhost"}
	
	pid := os.Getpid()
    timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
    filename := fmt.Sprintf("%s.%d.%s", timestamp, pid, hostname)

	filepath_to_write := filepath.Join(CONFIGET(MAILDIR_PATH), "new")
	filepath_to_write = filepath.Join(filepath_to_write, filename)
	file, err := os.Create(filepath_to_write)
	if err != nil {return err}
	_, err = file.WriteString(email_data)
	if err != nil {return err}
	err = file.Sync()
	if err != nil {return err}
	err = file.Close()
	return err
}


// config file located at /etc/onx/onx.conf
func CONFIGET(type_of_config string) string {
	options := []string{PORT, ROOT_DOMAIN, MAIL_DOMAIN, MAILDIR_PATH, TLS_FILE_fullchain, TLS_FILE_privkey}

	file, err := os.Open("/etc/onx/onx.conf")
	if err != nil {ERROR(nil, "error at opening config file", err)}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {break}
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\n", "")
		for _, option := range options {
			if strings.HasPrefix(line, option) && type_of_config == option {
				return strings.ReplaceAll(line, option, "")
			}
		}
	}
	ERROR(nil, fmt.Sprintf("there is no config for %s", type_of_config), nil)
	return "null" // this part actually will never be executed but my lsp gives an error so this stays
}
