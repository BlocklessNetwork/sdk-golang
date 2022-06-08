package http

const (
	SUCCESS Error = iota
	INVALID_HANDLE
	MEMORY_ACCESS_ERROR
	BUFFER_TOO_SMALL
	HEADER_NOT_FOUND
	UTF8_ERROR
	DESTINATION_NOT_ALLOWED
	INVALID_METHOD
	INVALID_ENCODING
	INVALID_URL
	REQUEST_ERROR
	RUNTIME_ERROR
	TOO_MANY_SESSIONS
	INVALID_DRIVER
	PERMISSION_DENY
)
