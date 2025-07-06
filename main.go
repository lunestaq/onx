package main

import (
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {ERROR(nil, "wrong usage")}
	listener, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {ERROR(nil, err)}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			WARNING(nil, err)
			continue
		}

		go handle_connection(connection, TLS_FALSE)
	}
}
