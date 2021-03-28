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
// for simplicity reasons, just text files from a hardcoded directory for now. nothing fancy-schmancy
func (h *Http) Get() string {

	data, err := h.readFile(h.Url)
	if err != nil {
		h.Status = 404
		return h.Handle404()
	}

	h.Status = 200
	return h.createHeaders(200) +
		data + "\r\n"
}

// read a file from url:
func (h *Http) readFile(url string) (string, error) {
	if url == "/" {
		url = "/index.html"
	}

	filePath := path.Join(WebRoot, url)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "File not found", err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "Error reading file", err
	}
	return string(data), nil
}

// Handle 404, file not found
func (h *Http) Handle404() string {
	return h.createHeaders(404) +
		"Not found" +
		"\r\n\r\n"
}

// CGI gateway? -> for the future reference

// Hello world response:

// GetHello returns generic hello World response
func (h *Http) GetHello() string {
	return h.createHeaders(200) +
		"Hello World from mhttp/0.0.1\r\n" +
		"\r\n"
}

// create headers:
func (h *Http) createHeaders(status int) string {
	response := fmt.Sprintf("HTTP/1.1 %d", status)

	if status == 200 {
		response += " OK\r\n"
	} else if status == 404 {
		response += " Not Found\r\n"
	} else {
		response += " Internal Server Error\r\n"
	}
	response += "Server: mhttp/0.0.1\r\n"
	response += "Content-Type: text/html\n\n"

	return response + "\r\n"
}
