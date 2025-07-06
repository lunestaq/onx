package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func INFO(ip net.Addr, text string) {
	fmt.Fprintf(os.Stderr, "[%s] (%v) [INFO] %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, text)
}

func WARNING(ip net.Addr, text string, err error) {
	fmt.Fprintf(os.Stderr, "[%s] (%v) [WARNING] %s: %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, text, err)
}

func ERROR(ip net.Addr, text string, err error) {
	fmt.Fprintf(os.Stderr, "[%s] (%v) [ERROR] %s: %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, text, err)
	os.Exit(1)
}

func INCOMING(connection net.Conn, message string) {
	fmt.Fprintf(os.Stderr, "[%s] (%s) INCOMING: %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), connection.RemoteAddr(), message)
}

func OUTGOING(connection net.Conn, message string) {
	fmt.Fprint(connection, message)
	fmt.Fprintf(os.Stderr, "[%s] (%s) OUTGOING: %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), connection.RemoteAddr(), message)
}
