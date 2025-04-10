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
	var version_error []byte
	switch ver {
	case 0, 1, 2, 3, 4:
		version_error = []byte{0, 0}
	default:
		version_error = []byte{0, 35}
	}

	response := make([]byte, 19)

	//copy(response, req[:4])           // message_size param

	copy(response[4:], version_error)            // error_code (represents no error in this case)
	response[6] = 2                              // Number of API keys
	copy(response[7:], req[6:8])                 // api_version
	binary.BigEndian.PutUint16(response[9:], 3)  //             min version
	binary.BigEndian.PutUint16(response[11:], 4) //             max version
	response[13] = 0                             // _tagged_fields
	binary.BigEndian.PutUint32(response[14:], 0) // throttle time
	response[18] = 0

	//copy(response[4:8], req[8:12])    // correlation_id param

	conn.Write(response)
}
