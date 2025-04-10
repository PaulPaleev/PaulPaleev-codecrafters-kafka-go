package main

import (
	"encoding/binary"
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
	api_version := req[6:8]
	ver := binary.BigEndian.Uint16(api_version)
	var version_error []byte
	switch ver {
	case 0, 1, 2, 3, 4:
		version_error = []byte{0, 0}
	default:
		version_error = []byte{0, 35}
	}

	response := make([]byte, 13)
	copy(response, req[0:4])          // message_size param
	copy(response, api_version)       // api_version
	copy(response[4:8], req[8:12])    // correlation_id param
	copy(response[8:], version_error) // error_code (represents no error in this case)
	conn.Write(response)
	ret_key := []byte{2, 0, 18, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0}
	conn.Write(ret_key)
}
