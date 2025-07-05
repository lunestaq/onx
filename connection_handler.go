package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type session__ struct {
	client_ip     string
	client_domain string
	mail_from     string
	rcpt_to       []string
	data	      []string
	state         string
}


func handle_connection(connection net.Conn) {
	defer connection.Close()

	session := &session__{
		client_ip: connection.RemoteAddr().String(),
		state: "CONNECTED",
	}

	scanner := bufio.NewScanner(connection)
	fmt.Fprint(connection, "220 siestaq.com ONXYIA Ready\r\n")
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("C [%s]: %s\n", session.client_ip, line)

		if session.state == "DATA" {
			if line == "." {
				session.state = "DATA_DONE"
				// handle_data
				fmt.Fprint(connection, "250 OK\r\n")
			} else {
				session.data = append(session.data, line)
			}
			continue
		}

		command := strings.ToUpper(line[:4])
		switch command {
		case "EHLO":
			session.client_domain = strings.ReplaceAll(line, "EHLO ", "")
			fmt.Fprintf(connection, "250-siestaq.com Hello %s\r\n", session.client_domain)
			fmt.Fprint(connection, "250-SIZE 10485760\r\n")
			fmt.Fprint(connection, "250-8BITMIME\r\n")
			fmt.Fprint(connection, "250 HELP\r\n")
			session.state = "EHLO_DONE"

		case "HELO":
			session.client_domain = strings.ReplaceAll(line, "HELO ", "")
			fmt.Fprintf(connection, "250 siestaq.com Hello %s\r\n", session.client_domain)
			session.state = "HELO_DONE"

		case "MAIL":
			session.mail_from = strings.ReplaceAll(line, "MAIL FROM:", "")
			fmt.Fprint(connection, "250 OK\r\n")
			session.state = "MAIL_FROM_DONE"

		case "RCPT":
		session.rcpt_to = append(session.rcpt_to, strings.ReplaceAll(line, "RCPT TO:", ""))
		fmt.Fprint(connection, "250 OK\r\n")
		session.state = "RCPT_TO_DONE"

		case "DATA":
		if session.state != "RCPT_TO_DONE" {
			fmt.Fprint(connection, "503 Bad sequence of commands\r\n")
		} else {
			fmt.Fprint(connection, "354 End data with <CR><LF>.<CR><LF>\r\n")
			session.state = "DATA"
		}
		
		case "RSET":
			session.mail_from = ""
			session.rcpt_to = nil
			session.data = nil
			session.state = "CONNECTED"
			fmt.Fprint(connection, "250 OK\r\n")

		case "NOOP":
			fmt.Fprint(connection, "250 OK\r\n")
		
		case "QUIT":
			fmt.Fprint(connection, "221 Bye\r\n")
			return
		
		case "HELP":
			fmt.Fprint(connection, "214-Commands supported:\r\n")
			fmt.Fprint(connection, "214 HELO EHLO MAIL RCPT DATA RSET NOOP QUIT HELP\r\n")

		default:
			fmt.Fprint(connection, "502 Command not implemented\r\n")
		}
	}
}
