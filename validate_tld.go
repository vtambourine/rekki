package main

import (
	"io/ioutil"
	"strings"
)

// List of top-level domains published by ICANN.
// http://data.iana.org/TLD/tlds-alpha-by-domain.txt
const tldsFile = "tlds.txt"

var tldsMap = map[string]bool{}

func init() {
	data, err := ioutil.ReadFile(tldsFile)
	if err != nil {
		panic("missing input file: tlds.txt")
	}

	for _, domain := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(domain, "#") {
			tldsMap[domain] = true
		}
	}
}

func tld(email string) string {
	return email[strings.LastIndex(email, ".")+1:]
}

// TODO: Add support for internationalized top-level domains.
func ValidateTLD(email string) (bool, string) {
	if _, ok := tldsMap[strings.ToUpper(tld(email))]; ok {
		return true, ""
	}
	return false, ReasonInvalidTLD
}
