package main

import (
	"fmt"
	"net"
	"net/smtp"
)

var defaultDialer = net.Dialer{}

// TODO: Extract to configuration.
var defaultPort = 25
var defaultFromEmail = "noreply@rekki.com"

func dialMailbox(email string, mxList []*net.MX) (result string) {
	var c *smtp.Client
	for _, mx := range mxList {
		conn, err := defaultDialer.Dial("tcp", fmt.Sprintf("%v:%v", mx.Host, defaultPort))
		if t, ok := err.(*net.OpError); ok {
			if t.Timeout() {
				return ReasonTimeout
			}
			return ReasonUnableToConnect
		} else if err != nil {
			return ReasonMailserverError
		}
		c, err = smtp.NewClient(conn, mx.Host)
		if err == nil {
			break
		}
		if c == nil {
			return ReasonMailserverError
		}
	}

	resChan := make(chan string, 1)

	go func() {
		defer c.Close()
		defer c.Quit()

		err := c.Hello(hostname(email))
		if err != nil {
			resChan <- ReasonMailserverError
			return
		}

		err = c.Mail(defaultFromEmail)
		if err != nil {
			resChan <- ReasonMailserverError
			return
		}

		id, err := c.Text.Cmd("RCPT TO:<%s>", email)
		if err != nil {
			resChan <- ReasonMailserverError
			return
		}

		c.Text.StartResponse(id)
		code, _, err := c.Text.ReadResponse(25)
		c.Text.EndResponse(id)

		if code == 550 {
			resChan <- ReasonUnavailableMailbox
			return
		}

		if err != nil {
			resChan <- ReasonMailserverError
			return
		}

		resChan <- ""
	}()

	select {
	case q := <-resChan:
		return q
	}
}

func ValidateSMTP(email string) (bool, string) {
	mxList, err := net.LookupMX(hostname(email))
	if err != nil || len(mxList) == 0 {
		return false, ReasonInvalidHostname
	}
	res := dialMailbox(email, mxList)
	if res != "" {
		return false, res
	}

	return true, ""
}
