package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

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

	//var correlation_id in32 = 7
	handleRequest(conn)
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	req := make([]byte, 1024)
	conn.Read(req)

	response := make([]byte, 10)
	copy(response, []byte{0, 0, 0, 0}) // message_size param
	/*
		request example:
		so we start copy elements from 8 since it's the first element of correlation_id
		00 00 00 23  // message_size:        35
		00 12        // request_api_key:     18
		00 04        // request_api_version: 4
		6f 7f c6 61  // correlation_id:      1870644833
	*/
	copy(response[4:8], req[8:12])    // [4:] bc we have [0:5] as message_size param, [8:13] correlation_id param in the request
	copy(response[8:], []byte{0, 35}) // error_code 35 (we have only 2 elements to fill)
	/*
		response example:
		00 00 00 00  // message_size:   0 (any value works)
		4f 74 d2 8b  // correlation_id: 1333056139
		00 23        // error_code:     35
	*/
	conn.Write(response)
}
