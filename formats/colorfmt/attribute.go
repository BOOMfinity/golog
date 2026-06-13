package colorfmt

import (
	"encoding/json"
	"strconv"
)

func appendAttribute(buff []byte, val any) []byte {
	switch v := val.(type) {
	case []byte:
		return append(buff, v...)
	case string:
		return append(buff, v...)
	case bool:
		return strconv.AppendBool(buff, v)
	case int:
		return strconv.AppendInt(buff, int64(v), 10)
	case int8:
		return strconv.AppendInt(buff, int64(v), 10)
	case int16:
		return strconv.AppendInt(buff, int64(v), 10)
	case int32:
		return strconv.AppendInt(buff, int64(v), 10)
	case int64:
		return strconv.AppendInt(buff, v, 10)
	case uint:
		return strconv.AppendUint(buff, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(buff, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(buff, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(buff, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(buff, v, 10)
	case float32:
		return strconv.AppendFloat(buff, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(buff, v, 'f', -1, 64)
	default:
		d, err := json.Marshal(v)
		if err != nil {
			d = json.RawMessage("null")
		}
		return append(buff, d...)
	}
}
