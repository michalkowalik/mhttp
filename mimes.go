package main

import (
	"fmt"
	"path"
	"strings"
)

var mimes = map[string]string{
	"txt":  "text/html",
	"html": "text/html",
	"htm":  "text/html",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
}

func getContentType(fileExtension string) string {
	mimeType, ok := mimes[fileExtension]
	if !ok {
		return "text/html"
	}

	fmt.Printf("for filename %s mime type %s\n", fileExtension, mimeType)

	return mimeType
}

func getExtension(url string) string {
	fileName := path.Base(url)
	tokens := strings.Split(fileName, ".")
	return strings.ToLower(tokens[len(tokens)-1])
}
