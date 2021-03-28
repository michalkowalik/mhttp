package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Printf("starting mhttp\n")

	// listen
	dstream, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err.Error())
	}
	defer dstream.Close()

	// accept
	for {
		connection, err := dstream.Accept()
		if err != nil {
			panic(err.Error())
		}

		// handle
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	connectionReader := bufio.NewReader(connection)
	requestData := readFromClient(connectionReader)

	fmt.Printf("Request: \n%s\n", string(requestData))
	err := handleResponse(requestData, connection)
	if err != nil {
		log.Printf("ERR: %s\n", err)
	}

}

// by definition, the request ends with an empty line.
func readFromClient(connectionReader *bufio.Reader) []byte {
	var buffer bytes.Buffer

	for {
		line, _, err := connectionReader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}
		buffer.Write(line)
		buffer.Write([]byte("\r\n"))
	}
	return buffer.Bytes()
}

func handleResponse(request []byte, connection net.Conn) error {

	var payload string

	httpResponse, err := ParseRequest(string(request))
	if err != nil {
		log.Printf("ERR: %s\n", err)
		return err
	}

	if httpResponse.Method == "GET" {
		payload = httpResponse.Get()
	}

	// write some response:
	n, err := connection.Write([]byte(payload))
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes of output\n", n)

	return nil
}
