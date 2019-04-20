package messageformat

import (
	"fmt"
	"strconv"
)

type Data map[string]interface{}

// ValueStr retrieves a value from the given map and tries to return a string representation.
//
// It will returns an error if the value's type is not <nil/string/int/float64>.
func (data Data) ValueStr(key string) (string, error) {
	if v, ok := data[key]; ok {
		switch v.(type) {
		default:
			return "", fmt.Errorf("ValueStr: Unsupported type: %T", v)

		case nil:
			return "", nil

		case string:
			return v.(string), nil

		case int:
			return fmt.Sprintf("%d", v.(int)), nil

		case float64:
			return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil
		}
	}
	return "", nil
}
