package main

import "strings"

func domain(email string) string {
	return email[strings.LastIndex(email, ".")+1:]
}

func hostname(email string) string {
	return email[strings.Index(email, "@")+1:]
}
