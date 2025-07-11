package main

import (
	"bufio"
	"net"
	"strings"
)

const TLS_TRUE, TLS_FALSE = true, false
type client_ struct {
	status	  string
	domain    string
	mail_from string
	rcpt_to   string
	data      string
} 

func handle_connection(connection net.Conn, is_tls bool) {
	defer connection.Close()
	client := client_{status: "null", data: ""}
	funcmap := map[string]func(net.Conn, string, *client_) {
		"HELO ": handle_HELO,
		"EHLO ": handle_EHLO,
		"STARTTLS": handle_STARTTLS,
		"MAIL FROM:": handle_MAIL_FROM,
		"RCPT TO:": handle_RCPT_TO,
		"DATA": handle_DATA,
		"RSET": handle_RSET,
		"NOOP": handle_NOOP,
		"QUIT": handle_QUIT,
	}

	if is_tls == false {
		INFO(connection.RemoteAddr(), "connected")
		OUTGOING(connection, "220 "+CONFIGET(MAIL_DOMAIN)+" ready\r\n")
	}

	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		line := scanner.Text()
		INCOMING(connection, line)

		for prefix, function := range funcmap {
			if strings.HasPrefix(strings.ToUpper(line), prefix) {
				function(connection, line, &client)
				if client.status == "quit" {
					if client.data != "" {
						err := save_mail(client.data) 
						if err != nil {WARNING(connection.RemoteAddr(), "error at writing mail to file: %s", err)} else {INFO(connection.RemoteAddr(), "saved mail")}
					}
					return
				}
				if client.status == "quit_after_tls" {return}
				goto done
			}	
		}
		OUTGOING(connection, "500 unrecognized\r\n")
		done:
	}
	if err := scanner.Err(); err != nil {
		WARNING(connection.RemoteAddr(), "error at receiving: %s", err)
	}
}
