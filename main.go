package main

import (
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.LUTC)
	if len(os.Args) < 2 {log.Fatal("wrong usage")}
	port := os.Args[1]
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {log.Fatalf("error binding port: %s\n", err)}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %s\n", err)
			continue
		}

		go handle_connection(connection)
	}
}
