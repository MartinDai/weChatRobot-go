package util

import (
	"net"
	"regexp"
	"strconv"
	"strings"
)

func ValidateAddress(address string) bool {
	parts := strings.Split(address, ":")
	if len(parts) > 2 {
		return false
	}

	host := parts[0]
	port := ""

	if len(parts) == 2 {
		port = parts[1]
	}

	// 校验Host部分
	if !isIP(host) && !isDomainName(host) {
		return false
	}

	// 校验Port部分
	if port != "" && !isPort(port) {
		return false
	}

	return true
}

func isPort(port string) bool {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	if portNum < 0 || portNum > 65535 {
		return false
	}

	return true
}

func isIP(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	}
	return true
}

func isDomainName(domain string) bool {
	domainRegex := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`
	match, _ := regexp.MatchString(domainRegex, domain)
	return match
}
