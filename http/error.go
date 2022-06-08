package http

type Error uint32

func (e Error) Error() string {
	switch e {
	case INVALID_HANDLE:
		return "Invalid handle"
	case MEMORY_ACCESS_ERROR:
		return "Memory access error"
	case BUFFER_TOO_SMALL:
		return "Buffer too small"
	case HEADER_NOT_FOUND:
		return "Header not found"
	case UTF8_ERROR:
		return "UTF-8 error"
	case DESTINATION_NOT_ALLOWED:
		return "Destination not allowed"
	case INVALID_METHOD:
		return "Invalid method"
	case INVALID_ENCODING:
		return "Invalid encoding"
	case INVALID_URL:
		return "Invalid URL"
	case REQUEST_ERROR:
		return "Request error"
	case RUNTIME_ERROR:
		return "Runtime error"
	case TOO_MANY_SESSIONS:
		return "Too many sessions"
	case INVALID_DRIVER:
		return "Invalid driver"
	case PERMISSION_DENY:
		return "Permission deny"
	default:
		return "Runtime error"
	}
}
