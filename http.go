package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const WebRoot = "/home/mkowalik/webroot"
const IndexFile = "index.html"

// Http message and methods definition
type Http struct {
	Status      int
	Headers     map[string]string
	RequestData string
	Method      string
	Url         string

	// just for the future reference. it doesn't seem too hard to define.
	keepAlive bool
}

// Parse Request creates new http object
func ParseRequest(request string) (*Http, error) {
	http := new(Http)

	// for simplicity reason, only first line will be parsed:
	lines := strings.Split(request, "\r\n")
	tokens := strings.Split(lines[0], " ")

	if len(tokens) != 3 {
		return nil, errors.New("Bad request")
	}

	http.Method = tokens[0]
	http.Url = tokens[1]

	// for now, we're happily ignore http version
	return http, nil
}

// Handle Get request
// for simplicity reasons, just files from a hardcoded directory for now.
func (h *Http) Get() []byte {

	data, err := h.readFile(h.Url)
	if err != nil {
		h.Status = 404
		return []byte(h.Handle404())
	}

	h.Status = 200
	return append(h.createHeaders(200, len(data)), append(data, []byte("\r\n")...)...)
}

// read a file from url:
func (h *Http) readFile(url string) ([]byte, error) {
	if url == "/" {
		url = "/index.html"
	}

	filePath := path.Join(WebRoot, url)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return []byte("File not found"), err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte("Error reading file"), err
	}
	return data, nil
}

// Handle 404, file not found
func (h *Http) Handle404() []byte {
	return append(h.createHeaders(404, 13), "Not found\r\n\r\n"...)
}

// CGI gateway? -> for the future reference

// GetHello returns generic hello World response
func (h *Http) GetHello() []byte {
	helloString := "Hello World from mhttp/0.0.1\r\n\r\n"
	return append(h.createHeaders(200, len(helloString)), helloString...)
}

// create headers:
func (h *Http) createHeaders(status int, dataLength int) []byte {
	response := fmt.Sprintf("HTTP/1.1 %d", status)

	if status == 200 {
		response += " OK\r\n"
	} else if status == 404 {
		response += " Not Found\r\n"
	} else {
		response += " Internal Server Error\r\n"
	}
	response += "Server: mhttp/0.0.1\r\n"
	response += fmt.Sprintf("Content-Type: %s\r\n", getContentType(getExtension(h.Url)))
	response += fmt.Sprintf("Content-Length: %d\r\n", dataLength)

	return []byte(response + "\r\n")
}
