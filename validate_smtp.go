package main

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
)

var defaultDialer = net.Dialer{}

// TODO: Extract to configuration.
var defaultPort = 25
var defaultFromEmail = "noreply@rekki.com"

func dialMailbox(email string, mxList []*net.MX) (err error) {
	var c *smtp.Client
	for _, mx := range mxList {
		conn, err := defaultDialer.Dial("tcp", fmt.Sprintf("%v:%v", mx.Host, defaultPort))
		if t, ok := err.(*net.OpError); ok {
			if t.Timeout() {
				return errors.New(ReasonTimeout)
			}
			return errors.New(ReasonUnableToConnect)
		} else if err != nil {
			return errors.New(ReasonMailserverError)
		}
		c, err = smtp.NewClient(conn, mx.Host)
		if err == nil {
			break
		}
		if c == nil {
			return errors.New(ReasonMailserverError)
		}
	}

	resChan := make(chan error, 1)

	go func() {
		defer c.Close()
		defer c.Quit()

		err := c.Hello(domain(email))
		if err != nil {
			resChan <- errors.New(ReasonMailserverError)
			return
		}

		err = c.Mail(defaultFromEmail)
		if err != nil {
			resChan <- errors.New(ReasonMailserverError)
			return
		}

		id, err := c.Text.Cmd("RCPT TO:<%s>", email)
		if err != nil {
			resChan <- errors.New(ReasonMailserverError)
			return
		}

		c.Text.StartResponse(id)
		code, _, err := c.Text.ReadResponse(25)
		c.Text.EndResponse(id)

		if code == 550 {
			resChan <- errors.New(ReasonUnavailableMailbox)
			return
		}

		if err != nil {
			resChan <- errors.New(ReasonMailserverError)
			return
		}

		resChan <- nil
	}()

	select {
	case q := <-resChan:
		return q
	}
}

func ValidateSMTP(email string) (bool, string) {
	mxList, err := net.LookupMX(domain(email))
	if err != nil || len(mxList) == 0 {
		//log.Println("First stop")
		return false, ReasonInvalidHostname
	}
	//log.Println(mxList)
	err = dialMailbox(email, mxList)
	if err != nil {
		//log.Println("Something went wrong", err)
		return false, err.Error()
	}

	return true, ""
}
