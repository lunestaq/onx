package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func INFO(ip net.Addr, message interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] [INFO] (%s): %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, message)
}

func WARNING(ip net.Addr, message interface{}) {	
	fmt.Fprintf(os.Stderr, "[%s] [WARNING] (%s): %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, message)
}

func ERROR(ip net.Addr, message interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] [ERROR] (%v): %s\n", time.Now().UTC().Format("02-01-2006 15-04-05.000"), ip, message)
	os.Exit(1)
}
