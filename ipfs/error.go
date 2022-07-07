package ipfs

type Error uint32

func (e Error) Error() string {
	switch e {
	case INVALID_HANDLE:
		return "Invalid handle"
	case UTF8_ERROR:
		return "UTF-8 error"
	case INVALID_METHOD:
		return "Invalid method"
	case INVALID_PARAMETER:
		return "Invalid parameter"
	case INVALID_ENCODING:
		return "Invalid encoding"
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
	case BUFFER_TOO_SMALL:
		return "Buffer too small"
	default:
		return "Runtime error"
	}
}
