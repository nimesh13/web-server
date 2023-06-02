package main

import (
	"fmt"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"
)

// Get the mime type for the file requested
func GetContentType(filename string) string {
	return mime.TypeByExtension(filepath.Ext(filename))
}

// Parses the request header to provide method, requestURI and protocol
func ParseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	words := strings.Split(line, " ")
	if len(words) == 0 {
		return "", "", "", false
	}
	if len(words) < 2 || len(words) > 3 {
		return "", "", "", false
	}
	if len(words) == 2 {
		if words[0] != "GET" {
			return "", "", "", false
		}
		return words[0], words[1], DEFAULT_HTTP_PROTOCOL, true
	}
	method, requestURI, proto = words[0], words[1], words[2]
	if !strings.HasPrefix(proto, "HTTP/") {
		return "", "", "", false
	}
	base_version_number := strings.Split(proto, "/")[1]
	version_number := strings.Split(base_version_number, ".")
	if len(version_number) != 2 {
		return "", "", "", false
	}
	if version_number[0] == "" || version_number[1] == "" {
		return "", "", "", false
	}
	return method, requestURI, proto, true
}

// StatusText to return the text for the HTTP codes
func StatusMessage(code int) (message string, explanation string) {
	switch code {
	case StatusOK:
		return "OK", "Request fulfilled, document follows"
	case StatusBadRequest:
		return "Bad Request", "Bad request syntax or unsupported method"
	case StatusForbidden:
		return "Forbidden", "Request forbidden -- authorization will not help"
	case StatusNotFound:
		return "Not Found", "Nothing matches the given URI"
	case StatusInternalServerError:
		return "Internal Server Error", "Server got itself in trouble"
	case StatusNotImplemented:
		return "Not Implemented", "Server does not support this operation"
	default:
		return "", ""
	}
}

// Formats the error response based on the error code and message
func GetErrorResponse(httpCode int, message string, explain string) string {
	return fmt.Sprintf(DEFAULT_ERROR_MESSAGE, httpCode, message, httpCode, explain)
}

// Get file permissions to verify if others have read access
func CheckReadAccess(fileInfo fs.FileInfo) bool {
	mode := fileInfo.Mode()
	return string(mode.String()[7]) == "r"
}
