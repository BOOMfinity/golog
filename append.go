package golog

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

var (
	debugCode = []byte("\\e[3m\u001b[35m")
	infoCode  = []byte("\u001b[36m")
	warnCode  = []byte("\u001b[1m\u001b[33m")
	errorCode = []byte("\u001b[1m\u001b[31m")
	resetCode = []byte("\u001B[0m")
)

func appendType(b []byte, v interface{}) []byte {
	switch x := v.(type) {
	case int:
		return strconv.AppendInt(b, int64(x), 10)
	case int8:
		return strconv.AppendInt(b, int64(x), 10)
	case int16:
		return strconv.AppendInt(b, int64(x), 10)
	case int32:
		return strconv.AppendInt(b, int64(x), 10)
	case int64:
		return strconv.AppendInt(b, x, 10)
	case uint:
		return strconv.AppendInt(b, int64(x), 10)
	case uint8:
		return strconv.AppendInt(b, int64(x), 10)
	case uint16:
		return strconv.AppendInt(b, int64(x), 10)
	case uint32:
		return strconv.AppendInt(b, int64(x), 10)
	case uint64:
		return strconv.AppendInt(b, int64(x), 10)
	case float32:
		return strconv.AppendFloat(b, float64(x), 'f', -1, 32)
	case float64:
		return strconv.AppendFloat(b, x, 'f', -1, 64)
	case string:
		return append(b, []byte(v.(string))...)
	case error:
		return append(b, []byte(v.(error).Error())...)
	case []byte:
		return append(b, v.([]byte)...)
	case bool:
		if v.(bool) {
			return append(b, []byte("true")...)
		}
		return append(b, []byte("false")...)
	default:
		return appendFormat(b, "%v", v)
	}
}

func appendFormat(dest []byte, format string, values ...interface{}) []byte {
	x := bufferPool.Get().(*bytes.Buffer)
	x.Reset()
	defer putBuffer(x)
	_, _ = fmt.Fprintf(x, format, values...)
	return append(dest, x.Bytes()...)
}

// AppendTime appends a time.Time to a slice of bytes
func appendTime(dst []byte, t time.Time, format string) []byte {
	if format == "" {
		return appendType(dst, t.UnixNano()/1000000)
	}
	return t.AppendFormat(dst, format)
}

func appendLevel(dst []byte, level Level) []byte {
	switch level {
	case LevelDebug:
		return append(dst, []byte("DEBUG")...)
	case LevelInfo:
		return append(dst, []byte("INFO")...)
	case LevelWarn:
		return append(dst, []byte("WARN")...)
	case LevelError:
		return append(dst, []byte("ERROR")...)
	case LevelFatal:
		return append(dst, []byte("FATAL")...)
	default:
		return append(dst, []byte("UNKNOWN")...)
	}
}

func appendColors(dst []byte, level Level) []byte {
	switch level {
	case LevelDebug:
		return append(dst, debugCode...)
	case LevelInfo:
		return append(dst, infoCode...)
	case LevelWarn:
		return append(dst, warnCode...)
	case LevelError, LevelFatal:
		return append(dst, errorCode...)
	default:
		return dst
	}
}

func appendReset(dst []byte) []byte {
	return append(dst, resetCode...)
}

func unsafeBytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}
