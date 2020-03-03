package main

import (
	"io/ioutil"
	"strings"
)

const blacklistFile = "blacklist.txt"

var blacklistMap = map[string]bool{}

func init() {
	data, err := ioutil.ReadFile(blacklistFile)
	if err != nil {
		panic("missing input file: blacklist.txt")
	}

	for _, domain := range strings.Split(string(data), "\n") {
		blacklistMap[domain] = true
	}
}

func domain(email string) string {
	return email[strings.Index(email, "@")+1:]
}

func ValidateBlacklist(email string) (bool, string) {
	if _, ok := blacklistMap[domain(email)]; ok {
		return false, ReasonUntrustedDomain
	}
	return true, ""
}
