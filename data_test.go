package messageformat

import (
	"fmt"
	"testing"
)

func valueStrResult(t *testing.T, data Data, key, expected string) {
	result, err := data.ValueStr(key)

	if err != nil {
		t.Errorf("Expecting `%s` but got an error `%s`", expected, err.Error())
	} else if expected != result {
		t.Errorf("Expecting `%s` but got `%s`", expected, result)
	} else if testing.Verbose() {
		fmt.Printf("Successfully returns the expected value: `%s`\n", expected)
	}
}

func valueStrError(t *testing.T, data Data, key string) {
	result, err := data.ValueStr("B")

	if nil == err {
		t.Errorf("Expecting an error but got `%s`", result)
	} else if testing.Verbose() {
		fmt.Printf("Successfully returns an error `%s`\n", err.Error())
	}
}

func TestValueStr(t *testing.T) {
	data := Data{
		"S": "I am a string",
		"I": 42,
		"F": 0.305,
		"B": true,
		"N": nil,
	}

	// should returns an empty string when the key does not exists
	valueStrResult(t, data, "NAME", "")

	// should returns an empty string when the value is nil
	valueStrResult(t, data, "N", "")

	// should returns an error when the value's type is not supported
	valueStrError(t, data, "B")

	// should otherwise returns a string representation (string, int, float)
	valueStrResult(t, data, "S", "I am a string")
	valueStrResult(t, data, "I", "42")
	valueStrResult(t, data, "F", "0.305")
}
