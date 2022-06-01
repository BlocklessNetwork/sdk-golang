package http

type Error int

func (e Error) Error() string {
	switch e {
	case EOF:
		return "end of file"
	case TOO_BIG:
		return "too big"
	case ACCES:
		return "access errror"
	case ADDRINUSE:
		return "address in use"
	case ADDRNOTAVAIL:
		return "address not avail"
	case BAD_PARAMS:
		return "bad parameters"
	case AGAIN:
		return "again"
	case BADF:
		return "bad file descriptor"
	case BAD_CONNECT:
		return "bad connection"
	case BAD_DRIVER:
		return "bad driver"
	case BAD_OPEN:
		return "open fail"
	default:
		return "unknow error"
	}
}
