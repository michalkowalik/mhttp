package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

var listenAddress *string = flag.String("l", "localhost:8080", "local address")

func main() {
	fmt.Printf("starting mhttp\n")

	addr, _ := net.ResolveTCPAddr("tcp", *listenAddress)
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	defer listener.Close()

	// accept
	for {
		var connection *net.TCPConn
		connection, err := listener.AcceptTCP()

		if err != nil {
			panic(err.Error())
		}
		connection.SetKeepAlive(true)
		connection.SetKeepAlivePeriod(time.Second * 30)

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

	var payload []byte

	httpResponse, err := ParseRequest(string(request))
	if err != nil {
		log.Printf("ERR: %s\n", err)
		return err
	}

	if httpResponse.Method == "GET" {
		payload = httpResponse.Get()
	}

	// write some response:
	n, err := connection.Write(payload)
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes of output\n", n)

	return nil
}
