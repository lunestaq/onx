package main

import "net"

func main() {
	listener, err := net.Listen("tcp", ":"+CONFIGET(PORT))
	if err != nil {ERROR(nil, "error at binding port: %s", err)}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			WARNING(nil, "error at accepting: %s", err)
			continue
		}

		go handle_connection(connection, TLS_FALSE)
	}
}
