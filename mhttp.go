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

// by definition, the request ends with an empty line.
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
	response := "HTTP/1.0 200 OK\r\n"
	response += "Server: mhttp/0.0.1\r\n"
	response += "Content-Type: text/html\r\n"
	response += "\r\n Hello World \r\n\r\n"

	// write some response:
	n, err := connection.Write([]byte(response))
	if err != nil {
		return err
	}
	fmt.Printf("Written %d byts of output\n", n)

	return nil
}
