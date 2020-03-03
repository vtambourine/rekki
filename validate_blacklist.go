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

func ValidateBlacklist(email string) (bool, string) {
	if _, ok := blacklistMap[hostname(email)]; ok {
		return false, ReasonUntrustedDomain
	}
	return true, ""
}
