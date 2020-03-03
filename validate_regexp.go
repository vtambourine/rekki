package main

import "regexp"

// TODO: Refine email regular expression.
// Extend naïve implementation beyond latin alphanumeric characters and standard TLDs.
var emailRegexp = regexp.MustCompile("^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}$")

func ValidateRegexp(email string) (bool, string) {
	if emailRegexp.Match([]byte(email)) {
		return true, ""
	}
	return false, ReasonRegexpMismatch
}
