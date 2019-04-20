package messageformat

import (
	"fmt"
	"testing"
)

func whitespaceResult(t *testing.T, start, end int, input []rune, expected_char rune, expected_pos int) int {
	char, pos := whitespace(start, end, input)
	if expected_pos != pos {
		t.Errorf("Expecting first non-whitespace found at %d but got %d", expected_pos, pos)
	} else if expected_char != char {
		t.Errorf("Expecting first non-whitespace was `%s` but got `%s`", string(expected_char), string(char))
	} else if testing.Verbose() {
		fmt.Printf("Successfully returns `%s`, %d\n", string(char), pos)
	}
	return pos
}

func TestIsWhitespace(t *testing.T) {
	for _, char := range []rune{' ', '\r', '\n', '\t'} {
		if true != isWhitespace(char) {
			t.Errorf("Do not returns true when receiving `%s`", string(char))
		}
	}

	if false != isWhitespace('a') {
		t.Errorf("Do not returns false when receiving `%s`", string('a'))
	}
}

func TestWhitespace(t *testing.T) {
	input := []rune(`  hello world`)
	start, end := 0, len(input)

	// should traverses the input, from "start" to "end"
	// until a non-whitespace char is encountered
	// and returns that char and its position
	pos := whitespaceResult(t, start, end, input, 'h', 2)

	// should returns the same char and position because the char
	// at "start" position is not a whitespace
	whitespaceResult(t, pos, end, input, 'h', 2)

	// should returns a \0 char when the end position is reached
	whitespaceResult(t, pos, pos, input, 0, 2)
}
