package jsonparser

import (
	"bytes"
	"fmt"
	"strings"
)

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

type JSONEncoder struct {
	_isFirstKey []bool
	result      []string
}

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		_isFirstKey: []bool{true},
		result:      []string{},
	}
}

func (encode *JSONEncoder) SetString(key, value string) {
	encode.writeKey(key)
	encode.WriteString(value)
}

func (encode *JSONEncoder) SetBoolean(key string, value bool) {
	encode.writeKey(key)
	encode.WriteBoolean(value)
}

func (encode *JSONEncoder) SetNull(key string) {
	encode.writeKey(key)
	encode.write("null")
}

func (encode *JSONEncoder) SetInteger(key string, value int64) {
	encode.writeKey(key)
	encode.WriteInteger(value)
}

func (encode *JSONEncoder) SetFloat(key string, value float64) {
	encode.writeKey(key)
	encode.WriteFloat(value)
}

func (encode *JSONEncoder) PushArray(key string) {
	encode.writeKey(key)
	encode.write("[")
	encode._isFirstKey = append(encode._isFirstKey, true)
}

func (encode *JSONEncoder) PopArray() {
	encode.write("]")
	if len(encode._isFirstKey) > 1 {
		encode._isFirstKey = encode._isFirstKey[1:]
	}
}

func (encode *JSONEncoder) PushObject(key string) {
	encode.writeKey(key)
	encode.write("{")
	encode._isFirstKey = append(encode._isFirstKey, true)
}

func (encode *JSONEncoder) PopObject() {
	encode.write("}")
	if len(encode._isFirstKey) > 1 {
		encode._isFirstKey = encode._isFirstKey[1:]
	}
}

func (encode *JSONEncoder) isFirstKey() bool {
	return encode._isFirstKey[len(encode._isFirstKey)-1]
}

func (encode *JSONEncoder) ToString() string {
	return strings.Join(encode.result, "")
}

func (encode *JSONEncoder) WriteBoolean(value bool) {
	encode.write(fmt.Sprintf("%t", value))
}

func (encode *JSONEncoder) WriteInteger(value int64) {
	encode.write(fmt.Sprintf("%d", value))
}

func (encode *JSONEncoder) WriteFloat(value float64) {
	encode.write(fmt.Sprintf("%f", value))
}

func (encode *JSONEncoder) write(value string) {
	encode.result = append(encode.result, value)
}

func (encode *JSONEncoder) WriteString(str string) {
	encode.write("\"")
	var savedIndex = 0
	for i := 0; i < len(str); i++ {
		var char = str[i]
		var needsEscaping = char < 0x20 || char == '"' || char == '\\'
		if needsEscaping {
			encode.write(str[savedIndex:i])
			savedIndex = i + 1
			if char == '"' {
				encode.write("\\\"")
			} else if char == '\\' {
				encode.write("\\\\")
			} else if char == '\b' {
				encode.write("\\b")
			} else if char == '\n' {
				encode.write("\\n")
			} else if char == '\r' {
				encode.write("\\r")
			} else if char == '\t' {
				encode.write("\\t")
			} else {
				// TODO: Implement encoding for other contol characters
				// @ts-ignore integer does have toString
				panic(fmt.Sprintf("Unsupported control character code: %c", char))
			}
		}
	}
	encode.write(str[savedIndex:])
	encode.write("\"")
}

func (encode *JSONEncoder) writeKey(key string) {
	var isFirstKey = encode.isFirstKey()
	if !isFirstKey {
		encode.write(",")
	} else {
		encode._isFirstKey[len(encode._isFirstKey)-1] = false
	}
	if len(key) > 0 {
		encode.WriteString(key)
		encode.write(":")
	}
}
