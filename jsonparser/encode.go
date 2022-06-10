package jsonparser

import "bytes"

func Encode(s string) string {
	buf := bytes.Buffer{}
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch b {
		case '"':
			buf.WriteByte('\\')
		case '\t':
			buf.WriteString("\t")
			continue
		case '\n':
			buf.WriteString("\n")
			continue
		case '\r':
			buf.WriteString("\r")
			continue
		}
		buf.WriteByte(b)
	}
	return buf.String()
}
