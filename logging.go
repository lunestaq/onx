package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func INFO(ip net.Addr, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] [INFO] (%v)"+msg+"\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, args)
}

func WARNING(ip net.Addr, msg string, args ...interface{}) {	
	fmt.Fprintf(os.Stderr, "[%s] [INFO] (%v)"+msg+"\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, args)
}

func ERROR(ip net.Addr, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] [INFO] (%v)"+msg+"\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, args)
	os.Exit(1)
}
