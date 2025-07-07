package main

import (
	"regexp"
	"strings"
	"net"
	"os"
	"fmt"
	"strconv"
	"path/filepath"
	"time"
)

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

func save_mail(emailData string) error {
    hostname, err := os.Hostname()
    if err != nil {hostname = "localhost"}
	
	pid := os.Getpid()
    timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
    filename := fmt.Sprintf("%s.%d.%s", timestamp, pid, hostname)

	file := filepath.Join("/home/siesta/Maildir/new", filename)
    return os.WriteFile(file, []byte(emailData), 0666)
}

