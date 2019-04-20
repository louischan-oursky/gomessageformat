package messageformat

// isWhitespace returns true if the rune is a whitespace.
func isWhitespace(char rune) bool {
	return ' ' == char || '\r' == char || '\n' == char || '\t' == char
}

// whitespace traverses the input until a non-whitespace char is encountered.
func whitespace(start, end int, input []rune) (rune, int) {
	pos := start
	for pos < end {
		char := input[pos]

		switch char {
		default:
			return char, pos

		case ' ', '\r', '\n', '\t':
		}

		pos++
	}
	return 0, pos
}
