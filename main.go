package main

import (
	"log"
	"net"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)
	log.Println("ONXYIA SERVER STARTING")
	listener, err := net.Listen("tcp", ":2525")
	if err != nil {log.Fatalf("error on binding port: %v", err)}
	defer listener.Close()
	log.Println("started listening")

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("error on accepting connection: %v\n", err)
			continue
		}
		log.Printf("got new connection: %v\n", connection.RemoteAddr())
		handle_connection(connection)
	}
}
