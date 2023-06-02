package main

// All HTTP status codes supported by this server
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusNotImplemented      = 501
)

// Default host where the server starts up
const SERVER_HOST string = "localhost"

// RFC1123 - DateTime format for http
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// Default error message template
const DEFAULT_ERROR_MESSAGE = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
        "http://www.w3.org/TR/html4/strict.dtd">
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
        <title>Error response</title>
    </head>
    <body>
        <h1>Error response</h1>
        <p>Error code: %d</p>
        <p>Message: %s</p>
        <p>Error code explanation: %d - %s.</p>
    </body>
</html>
`

// Default error message content type
const DEFAULT_ERROR_CONTENT_TYPE = "text/html;charset=utf-8"

// Default HTTP protocol in use
const DEFAULT_HTTP_PROTOCOL = "HTTP/1.0"

// Default HTTP connection status
const DEFAULT_CONNECTION_STATUS = "Close"
