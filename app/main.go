package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleRequest(conn)
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	req := make([]byte, 1024)
	conn.Read(req)

	response := make([]byte, 19)
	copy(response, req[0:4])       // message_size param
	copy(response, req[6:8])       // api_version
	copy(response[4:8], req[8:12]) // correlation_id param
	copy(response[8:], []byte{})   // error_code (represents no error in this case)
	conn.Write(response)
}
