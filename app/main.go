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
	ver := binary.BigEndian.Uint16(req[6:8])
	fmt.Println(ver)
	var version_error []byte
	switch ver {
	case 0, 1, 2, 3, 4:
		version_error = []byte{0, 0}
	default:
		version_error = []byte{0, 35}
	}

	response := make([]byte, 19)
	//copy(response[:4], req[8:12])     // correlation_id param - 4 bytes
	copy(response[:4], req[:4])             // DONE message_size param - 4 bytes
	copy(response[4:], version_error)       // DONE error_code (represents no error in this case) - 2 bytes
	response[6] = 2                         // Number of API keys - 1 byte
	copy(response[7:], []byte{0, 18})       // API Key api_version - 2 bytes
	copy(response[9:], []byte{0, 3})        // min version - 2 bytes
	copy(response[11:], []byte{0, 4})       // max version - 2 bytes
	response[13] = 0                        // _tagged_fields - 1 byte
	copy(response[14:], []byte{0, 0, 0, 0}) // throttle time - 4 bytes
	response[18] = 0                        // 1 byte

	fmt.Println(response)
	conn.Write(response)
}
