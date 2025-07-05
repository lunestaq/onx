package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client_ struct {
	status	  string
	domain    string
	mail_from string
	rcpt_to   string
	data      string
}

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
		fmt.Fprint(connection, "250 ok\r\n")
	} else {
		fmt.Fprint(connection, "503\r\n")
	}
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

func handle_QUIT(connection net.Conn, line string, client *client_) {
	client.status = "quit"
	log.Printf("[%s] is quitting", connection.RemoteAddr())
	fmt.Fprint(connection, "221 mail.siestaq.com have a great day\n\r")
}

func handle_smtp(connection net.Conn) {
	defer connection.Close()
	log.Printf("new connection: %s\n", connection.RemoteAddr())
	client := client_{status: "null", data: ""}
	funcmap := map[string]func(net.Conn, string, *client_) {
		"HELO": handle_HELO,
		"EHLO": handle_EHLO,
		"MAIL FROM:": handle_MAIL_FROM,
		"RCPT TO:": handle_RCPT_TO,
		"DATA": handle_DATA,
		"QUIT": handle_QUIT,
	}	
	
	fmt.Fprintf(connection, "220 mail.siestaq.com ready\r\n")
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("recieved [%s]: %s", connection.RemoteAddr(), line)

		for prefix, function := range funcmap {
			if strings.HasPrefix(strings.ToUpper(line), prefix) {
				function(connection, line, &client)
				if client.status == "quit" {
					fmt.Print("\n\n\n\n\n\n")
					fmt.Printf("status: %s\ndomain: %s\nmail from: %s\nrcpt to: %s\ndata: %s\n\n",client.status,client.domain,client.mail_from,client.rcpt_to,client.data)
					return
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %s\n", err)
	}
}
