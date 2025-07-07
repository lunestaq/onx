package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println(CONFIGET(PORT), CONFIGET(ROOT_DOMAIN), CONFIGET(MAIL_DOMAIN), CONFIGET(MAIL_PATH), CONFIGET(TLS_FILE_fullchain), CONFIGET(TLS_FILE_privkey))
	if len(os.Args) < 2 {INFO(nil, "wrong usage"); os.Exit(1)}
	listener, err := net.Listen("tcp", ":"+os.Args[1])
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
