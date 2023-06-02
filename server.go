package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Variables to be read from the command line
var (
	DOC_ROOT    string
	SERVER_PORT string
)

func init() {
	flag.StringVar(&DOC_ROOT, "document_root", "./", "Document root for all the files.")
	flag.StringVar(&SERVER_PORT, "port", "8080", "Port number where the server listens to.")
}

func main() {
	// Parse the command line arguments
	flag.Parse()

	// Setup the tcp address
	address := SERVER_HOST + ":" + SERVER_PORT
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic("Unable to start server at: " + address)
	}
	defer ln.Close()

	fmt.Printf("Serving HTTP on :: port %s (http://%s/) ...\n", SERVER_PORT, address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("Unable to connect to connections!")
		}

		// Parallel connection handler
		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {

	defer connection.Close()

	// Read the first line of the request
	requestData, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		sendErrorResponse(connection, requestData, StatusBadRequest)
		return
	}

	requestData = strings.TrimSuffix(requestData, "\r\n")

	// Parse the request to get the path and error status
	method, path, _, ok := ParseRequestLine(requestData)
	if !ok {
		sendErrorResponse(connection, requestData, StatusBadRequest)
		return
	}

	// Only GET requests are supported
	if method != "GET" {
		sendErrorResponse(connection, requestData, StatusNotImplemented)
		return
	}

	if path == "/" {
		path += "index.html"
	}

	filename := filepath.Clean(DOC_ROOT + path)

	// Check if file exists
	fileInfo, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		sendErrorResponse(connection, requestData, StatusNotFound)
		return
	}

	// Check if file can be read
	access := CheckReadAccess(fileInfo)
	if !access {
		sendErrorResponse(connection, requestData, StatusForbidden)
		return
	}

	// Open the file
	inputFile, err := os.Open(filename)
	if err != nil {
		sendErrorResponse(connection, requestData, StatusInternalServerError)
		return
	}
	defer inputFile.Close()

	// Send response headers
	sendBasicRespHeaders(connection, requestData, StatusOK)
	sendContentRespHeaders(connection, GetContentType(fileInfo.Name()), int(fileInfo.Size()))

	fileReader := bufio.NewReader(inputFile)
	clientConnWriter := bufio.NewWriter(connection)

	// Copy the file contents onto the connection
	io.Copy(clientConnWriter, fileReader)

}

// Send common HTTP response headers
func sendBasicRespHeaders(conn net.Conn, request string, httpCode int) {
	httpStatusMsg, _ := StatusMessage(httpCode)

	log.Printf("\"%s\" - %d", request, httpCode)

	conn.Write([]byte(DEFAULT_HTTP_PROTOCOL + " " + strconv.Itoa(httpCode) + " " + httpStatusMsg + "\r\n"))
	conn.Write([]byte("Connection: " + DEFAULT_CONNECTION_STATUS + "\r\n"))
	conn.Write([]byte("Date: " + time.Now().Format(TimeFormat) + "\r\n"))
}

// Send HTTP content headers
func sendContentRespHeaders(conn net.Conn, contentType string, contentLen int) {
	conn.Write([]byte("Content-Type: " + contentType + "\r\n"))
	conn.Write([]byte("Content-Length: " + strconv.Itoa(contentLen) + "\r\n\r\n"))
}

// Send error response for all error scenarios
func sendErrorResponse(conn net.Conn, request string, httpCode int) {
	httpStatusMsg, httpStatusExp := StatusMessage(httpCode)
	errResponse := GetErrorResponse(httpCode, httpStatusMsg, httpStatusExp)

	sendBasicRespHeaders(conn, request, httpCode)
	sendContentRespHeaders(conn, DEFAULT_ERROR_CONTENT_TYPE, len(errResponse))
	conn.Write([]byte(errResponse))
}
