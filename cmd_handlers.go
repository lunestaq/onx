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
		client.domain = line[5:]
		client.status = "HELO_DONE"
		fmt.Fprintf(connection, "250-mail.siestaq.com Hello %s\r\n", client.domain)
	} else {
		fmt.Fprint(connection, "503\r\n")
	}
}

func handle_EHLO(connection net.Conn, line string, client *client_) {
	if client.status == "null" {
		client.domain = line[5:]
		client.status = "HELO_DONE"
		fmt.Fprintf(connection, "250-mail.siestaq.com Hello %s\r\n", client.domain)
		fmt.Fprint(connection, "250-SIZE 10485760\r\n")
		fmt.Fprint(connection, "250-8BITMEME\r\n")
		fmt.Fprint(connection, "250-STARTTLS\r\n")
		fmt.Fprint(connection, "250 ok\r\n")
	} else {
		fmt.Fprint(connection, "503\r\n")
	}
}

func handle_STARTTLS(connection net.Conn, line string, client *client_) {
	cert, err := tls.LoadX509KeyPair("/etc/letsencrypt/live/siestaq.com/fullchain.pem", "/etc/letsencrypt/live/siestaq.com/privkey.pem")
	if err != nil {
		WARNING(connection.RemoteAddr(), err)
		return
	}
	//var tls_config *tls.Config
	//var tls_connection *tls.Conn
	fmt.Fprint(connection, "220 ready\r\n")
	tls_config := &tls.Config{Certificates: []tls.Certificate{cert}}
	tls_connection := tls.Server(connection, tls_config)
	handle_connection(tls_connection, TLS_TRUE)
	client.status = "quit_after_tls"
}

func handle_MAIL_FROM(connection net.Conn, line string, client *client_) {
	if client.status == "HELO_DONE" {
		client.mail_from = line[11:len(line)-1]
		client.status = "MAIL_FROM_DONE"
		fmt.Fprint(connection, "250 ok\r\n")
	} else {
		fmt.Fprint(connection, "503\r\n")
	}
}

func handle_RCPT_TO(connection net.Conn, line string, client *client_) {
	if client.status == "MAIL_FROM_DONE" {
		client.rcpt_to = line[9:len(line)-1]
		client.status = "RCPT_TO_DONE"
		fmt.Fprint(connection, "250 ok\r\n")
	} else {
		fmt.Fprint(connection, "503\r\n")
	}
}

func handle_DATA(connection net.Conn, line string, client *client_) {
	if client.status == "RCPT_TO_DONE" {
		client.status = "DATA_PHASE"
		fmt.Fprint(connection, "354 Start mail input; end with <CRLF>.<CRLF>\r\n")
		handle_rest_of_DATA(connection, client)
	} else {
		fmt.Fprint(connection, "503\r\n")
	}
}

func handle_rest_of_DATA(connection net.Conn, client *client_) {
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			client.status = "DATA_DONE"
			fmt.Fprint(connection, "250 ok\r\n")
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
	fmt.Fprint(connection, "250 ok\r\n")
}

func handle_NOOP(connection net.Conn, line string, client *client_) {
	fmt.Fprint(connection, "250 ok\r\n")
}

func handle_QUIT(connection net.Conn, line string, client *client_) {
	client.status = "quit"
	fmt.Fprint(connection, "221 mail.siestaq.com have a great day\n\r")
}
