package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	"net"
	"regexp"
)

var validateName = validation.StringMatch(
	regexp.MustCompile(`^[A-Za-z0-9\p{Han}-_]{1,63}$`),
	"expected value to be 1 - 63 characters and only support chinese, english, numbers, '-', '_'",
)

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateIpAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	res := net.ParseIP(value)

	if res == nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ip address, got error parsing: %s", k, value))
	}

	return
}

func validateSubnetType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Reserve" && value != "Normal" && value != "Physical" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid subnet type, got error parsing: %s", k, value))
	}
	return
}

var instancePasswordUpperPattern = regexp.MustCompile(`[A-Z]`)
var instancePasswordLowerPattern = regexp.MustCompile(`[a-z]`)
var instancePasswordNumPattern = regexp.MustCompile(`[0-9]`)
var instancePasswordSpecialPattern = regexp.MustCompile(`[` + "`" + `()~!@#$%^&*-+=_|{}\[\]:;'<>,.?/]`)
var instancePasswordPattern = regexp.MustCompile(`^[A-Za-z0-9` + "`" + `()~!@#$%^&*-+=_|{}\[\]:;'<>,.?/]{8,30}$`)

func validateInstancePassword(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !instancePasswordPattern.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q is invalid, should have between 8-30 characters and any characters must be legal, got %q", k, value))
	}

	categoryCount := 0
	if instancePasswordUpperPattern.MatchString(value) {
		categoryCount++
	}

	if instancePasswordLowerPattern.MatchString(value) {
		categoryCount++
	}

	if instancePasswordNumPattern.MatchString(value) {
		categoryCount++
	}

	if instancePasswordSpecialPattern.MatchString(value) {
		categoryCount++
	}

	if categoryCount < 3 {
		errors = append(errors, fmt.Errorf("%q is invalid, should have least 3 items of capital letters, lower case letters, numbers and special characters, got %q", k, value))
	}

	return
}

var ipPattern = regexp.MustCompile(`^(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)$`)

func validateIp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !ipPattern.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q is invalid, got %q", k, value))
	}
	return
}
