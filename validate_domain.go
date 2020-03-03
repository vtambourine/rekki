package main

import (
	"io/ioutil"
	"strings"
)

// List of top-level domains published by ICANN.
// http://data.iana.org/TLD/tlds-alpha-by-domain.txt
const domainsFile = "domains.txt"

var domainsMap = map[string]bool{}

func init() {
	data, err := ioutil.ReadFile(domainsFile)
	if err != nil {
		panic("missing input file: domains.txt")
	}

	for _, domain := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(domain, "#") {
			domainsMap[domain] = true
		}
	}
}

// TODO: Add support for internationalized top-level domains.
func ValidateDomain(email string) (bool, string) {
	if _, ok := domainsMap[strings.ToUpper(domain(email))]; ok {
		return true, ""
	}
	return false, ReasonInvalidTLD
}
