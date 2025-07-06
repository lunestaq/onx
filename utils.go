package main

import (
	"regexp"
	"strings"
	"net"
)

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
