package main

import (
	"fmt"
	"net"
	"crypto/tls"
	"bufio"
	"strings"
)

func handle_HELO(connection net.Conn, line string, client *client_) {
	if client.status == "null" {
		client.domain = strings.ReplaceAll(line, "HELO ", "")
		if !is_valid_domain(client.domain) {
			OUTGOING(connection, "504 syntax error\r\n")
		} else {
			payload := fmt.Sprintf("250-mail.siestaq.com greetings %s\r\n", client.domain)
			OUTGOING(connection, payload)
			client.status = "HELO_DONE"
		}
	} else {
		OUTGOING(connection, "503 bad sequence\r\n")
	}
}

func handle_EHLO(connection net.Conn, line string, client *client_) {
	if client.status == "null" {
		client.domain = strings.ReplaceAll(line, "EHLO ", "")
		if !is_valid_domain(client.domain) {
			OUTGOING(connection, "504 syntax error\r\n")
		} else {
			payload := fmt.Sprintf("250-mail.siestaq.com greetings %s\r\n", client.domain)
			OUTGOING(connection, payload)
			OUTGOING(connection, "250-SIZE 10485760\r\n")
			OUTGOING(connection, "250-8BITMIME\r\n")
			OUTGOING(connection, "250-STARTTLS\r\n")
			OUTGOING(connection, "250 ok\r\n")
			client.status = "HELO_DONE"
		}
	} else {
		OUTGOING(connection, "503 bad sequence\r\n")
	}
}

func handle_STARTTLS(connection net.Conn, line string, client *client_) {
	INFO(connection.RemoteAddr(), "loading tls")
	cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/siestaq.com/fullchain.pem", "/etc/letsencrypt/live/siestaq.com/privkey.pem")
	if err != nil {
		WARNING(connection.RemoteAddr(), "error at loading tls: %s", err)
		return
	}
	INFO(connection.RemoteAddr(), "tls loaded")
	OUTGOING(connection, "220 ready\r\n")
	tls_config := &tls.Config{Certificates: []tls.Certificate{cert}}
	tls_connection := tls.Server(connection, tls_config)
	INFO(connection.RemoteAddr(), "start of encryption")
	handle_connection(tls_connection, TLS_TRUE)
	client.status = "quit_after_tls"
}

func handle_MAIL_FROM(connection net.Conn, line string, client *client_) {
	if client.status == "HELO_DONE" {
		client.mail_from = line[11:len(line)-1]
		client.status = "MAIL_FROM_DONE"
		OUTGOING(connection, "250 ok\r\n")
	} else {
		OUTGOING(connection, "503 bad sequence\r\n")
	}
}

func handle_RCPT_TO(connection net.Conn, line string, client *client_) {
	if client.status == "MAIL_FROM_DONE" {
		client.rcpt_to = line[9:len(line)-1]
		client.status = "RCPT_TO_DONE"
		OUTGOING(connection, "250 ok\r\n")
	} else {
		OUTGOING(connection, "503 bad sequence\r\n")
	}
}

func handle_DATA(connection net.Conn, line string, client *client_) {
	if client.status == "RCPT_TO_DONE" {
		client.status = "DATA_PHASE"
		OUTGOING(connection, "354 start mail input; end with <CRLF>.<CRLF>\r\n")
		handle_rest_of_DATA(connection, client)
	} else {
		OUTGOING(connection, "503 bad sequence\r\n")
	}
}

func handle_rest_of_DATA(connection net.Conn, client *client_) {
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			client.status = "DATA_DONE"
			OUTGOING(connection, "250 ok\r\n")
			err := save_mail(client.data) 
			if err != nil {WARNING(connection.RemoteAddr(), "error at writing mail to file: %s", err)} else {INFO(connection.RemoteAddr(), "wrote mail to file")}
			break
		} else {
			if strings.HasPrefix(line, "..") {
				line = line[1:]
			}
			client.data = client.data + line + "\r\n"
		}
	}
}

func handle_RSET(connection net.Conn, line string, client *client_) {
	client.status = "HELO_DONE"
	client.mail_from = ""
	client.rcpt_to = ""
	client.data = ""
	OUTGOING(connection, "250 ok\r\n")
}

func handle_NOOP(connection net.Conn, line string, client *client_) {
	OUTGOING(connection, "250 ok\r\n")
}

func handle_QUIT(connection net.Conn, line string, client *client_) {
	client.status = "quit"
	OUTGOING(connection, "221 mail.siestaq.com have a great day\n\r")
}
