package jsonfmt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func (h *handler) appendKey(buff []byte, key string) []byte {
	if len(buff) > 0 && buff[len(buff)-1] != '{' {
		buff = append(buff, ',')
	}
	buff = h.appendString(buff, key)
	return append(buff, ':')
}

func (h *handler) appendString(buff []byte, str string) []byte {
	buff = append(buff, '"')
	buff = append(buff, str...)
	buff = append(buff, '"')
	return buff
}

func (h *handler) appendEscapedString(buff []byte, str string) []byte {
	return strconv.AppendQuote(buff, str)
}

func (h *handler) appendTime(buff []byte, t time.Time, f string) []byte {
	buff = append(buff, '"')
	buff = t.AppendFormat(buff, f)
	return append(buff, '"')
}

func (h *handler) appendUint(buff []byte, v uint64) []byte {
	return strconv.AppendUint(buff, v, 10)
}

func (h *handler) appendBool(buff []byte, v bool) []byte {
	return strconv.AppendBool(buff, v)
}

func (h *handler) appendInt(buff []byte, v int64) []byte {
	return strconv.AppendInt(buff, v, 10)
}

func (h *handler) appendFloat(buff []byte, v float64) []byte {
	return strconv.AppendFloat(buff, v, 'f', -1, 64)
}

func (h *handler) appendAny(buff []byte, v any) []byte {
	switch val := v.(type) {
	case nil:
		return append(buff, "null"...)
	case []byte:
		return append(buff, val...)
	case string:
		if h.config.DisableAttributeEscaping {
			return h.appendString(buff, val)
		}
		return h.appendEscapedString(buff, val)
	case bool:
		return h.appendBool(buff, val)
	case uint:
		return h.appendUint(buff, uint64(val))
	case uint8:
		return h.appendUint(buff, uint64(val))
	case uint16:
		return h.appendUint(buff, uint64(val))
	case uint32:
		return h.appendUint(buff, uint64(val))
	case uint64:
		return h.appendUint(buff, val)
	case int:
		return h.appendInt(buff, int64(val))
	case int8:
		return h.appendInt(buff, int64(val))
	case int16:
		return h.appendInt(buff, int64(val))
	case int32:
		return h.appendInt(buff, int64(val))
	case int64:
		return h.appendInt(buff, val)
	case float32:
		return h.appendFloat(buff, float64(val))
	case float64:
		return h.appendFloat(buff, val)
	default:
		d, err := json.Marshal(v)
		if err != nil {
			println(fmt.Sprintf("[golog] cannot marshal attribute: %v", v))
			return append(buff, "null"...)
		}
		return append(buff, d...)
	}
}
