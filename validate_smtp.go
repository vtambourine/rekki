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

	// Connect to mail server
	for _, mx := range mxList {
		conn, err := defaultDialer.Dial("tcp", fmt.Sprintf("%v:%v", mx.Host, defaultPort))
		if t, ok := err.(*net.OpError); ok {
			if t.Timeout() {
				return ReasonTimeout
			}
			return ReasonUnableToConnect
		} else if err != nil {
			return ReasonMailServerError
		}
		c, err = smtp.NewClient(conn, mx.Host)
		if err == nil {
			break
		}
		if c == nil {
			return ReasonMailServerError
		}
	}

	resChan := make(chan string, 1)

	go func() {
		defer c.Close()
		defer c.Quit()

		// Request host name control
		err := c.Hello(hostname(email))
		if err != nil {
			resChan <- ReasonMailServerError
			return
		}

		// Attempt to send an email
		err = c.Mail(defaultFromEmail)
		if err != nil {
			resChan <- ReasonMailServerError
			return
		}

		id, err := c.Text.Cmd("RCPT TO:<%s>", email)
		if err != nil {
			resChan <- ReasonMailServerError
			return
		}

		// Expect 25X successful response
		c.Text.StartResponse(id)
		code, _, err := c.Text.ReadResponse(25)
		c.Text.EndResponse(id)

		// Address doesn't exists
		if code == 550 {
			resChan <- ReasonUnavailableMailbox
			return
		}

		if err != nil {
			resChan <- ReasonMailServerError
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
	// Fetch DNS MX records
	mxList, err := net.LookupMX(hostname(email))
	if err != nil || len(mxList) == 0 {
		return false, ReasonInvalidHostname
	}

	// Check mailbox availability
	res := dialMailbox(email, mxList)
	if res != "" {
		return false, res
	}

	return true, ""
}
