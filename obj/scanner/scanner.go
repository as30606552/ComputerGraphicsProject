package scanner

import (
	"io"
)

// One of the possible values that the Scanner.Next method returns.
type TokenType uint8

const (
	Word    TokenType = iota // Can consist of letters, numbers, and underscores. Cannot start with a number.
	Integer                  // Consists of digits. Can start with a minus.
	Float                    // Consists of digits with a dot between them. Can start with a minus.
	Slash                    // '/' character.
	Space                    // A sequence of spaces and/or tabs.
	EOL                      // '\n' character.
	EOF                      // Indicates that the end of the sequence of bytes being read has been reached.
	Unknown                  // Unknown type of token.
	Comment                  // Starts with the '#' character and ends with the character before the end of the line.
)

// Number of different token options.
const TokensCount = 9

// Converts the state of the finite state machine from which it moved to the initial state to the type of the read token.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Scanner.
var tokenTypeMap = [...]TokenType{Unknown, Comment, EOL, Space, Slash, Unknown, Unknown, Integer, Float, Word, Unknown}

// Converts a token type constant to its string representation.
var tokenTypeNamesMap = [...]string{"WORD", "INTEGER", "FLOAT", "SLASH", "SPACE", "EOL", "EOF", "UNKNOWN", "COMMENT"}

// Converts a token type constant to its string representation.
func (tokenType TokenType) String() string {
	return tokenTypeNamesMap[tokenType]
}

// Allows you to sequentially call the Next method to get tokens from a io.Reader that can occur in .obj files.
type Scanner interface {
	// Returns the next token read from the reader.
	// If all bytes are read from the reader before calling the method, the (EOF, "") is always returned.
	Next() (TokenType, string)
	// Skips all characters until the beginning of the next line.
	// LineString method can be called after to get the skipped line.
	SkipLine()
	// Returns the line fragment that was read by the Scanner.
	LineString() string
	// Returns the position of the character that was last processed by the Scanner
	// relative to the beginning of the sequence of bytes being read.
	Position() int
	// Returns the number of the line that was last processed by the Scanner.
	Line() int
	// Returns the position in the line that was last processed by the scanner.
	Column() int
	// Returns true if the Scanner will skip comments and will not return comment tokens.
	IsSkipComments() bool
	// You can use this method to enable or disable skipping comments.
	SkipComments(skipComments bool)
}

// One of the possible states of a finite state machine.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Scanner.
type stateType uint8

const (
	start      stateType = iota // Initial state.
	skipLine                    // Skipping all characters up to the '\n' character.
	foundEol                    // '\n' character found.
	foundSpace                  // Whitespace character found.
	foundSlash                  // '/' character found.
	foundMinus                  // '-' character was found at the beginning of the token, and a digit is expected.
	foundDot                    // A '.' character is found after an integer, a digit is expected.
	foundInt                    // '\n' character found.
	foundFloat                  // A sequence of characters satisfying the Float token is found, a digit is expected.
	foundWord                   // A sequence of characters satisfying the Word token is found.
	unknown                     // A sequence of characters that does not match the above types.
)

// One of the possible character types that can be contained in a sequence of bytes to be read.
type symbolType uint8

const (
	eol    symbolType = iota // '\n'
	space                    // ' ' or '\t'
	hash                     // '#'
	slash                    // '/'
	minus                    // '-'
	dot                      // '.'
	digit                    // '0' - '9'
	letter                   // 'a' - 'z' or 'A' - 'Z' or '_'
	other                    // Any other character.
)

// Calculates the character type.
func getSymbolType(symbol byte) symbolType {
	switch symbol {
	case '\n':
		return eol
	case ' ':
		return space
	case '\t':
		return space
	case '#':
		return hash
	case '/':
		return slash
	case '-':
		return minus
	case '.':
		return dot
	case '_':
		return letter
	}
	if '0' <= symbol && symbol <= '9' {
		return digit
	}
	if 'a' <= symbol && symbol <= 'z' || 'A' <= symbol && symbol <= 'Z' {
		return letter
	}
	return other
}

// The finite state machine table.
// See https://github.com/as30606552/ComputerGraphicsProject/wiki/Scanner.
var matrix = [9][11]stateType{
	{foundEol, start, start, start, start, start, start, start, start, start, start},
	{foundSpace, skipLine, start, foundSpace, start, start, start, start, start, start, start},
	{skipLine, skipLine, start, start, start, start, start, start, start, start, start},
	{foundSlash, skipLine, start, start, start, start, start, start, start, start, start},
	{foundMinus, skipLine, start, start, start, unknown, unknown, unknown, unknown, unknown, unknown},
	{unknown, skipLine, start, start, start, unknown, unknown, foundDot, unknown, unknown, unknown},
	{foundInt, skipLine, start, start, start, foundInt, foundFloat, foundInt, foundFloat, foundWord, unknown},
	{foundWord, skipLine, start, start, start, unknown, unknown, unknown, unknown, foundWord, unknown},
	{unknown, skipLine, start, start, start, unknown, unknown, unknown, unknown, unknown, unknown},
}

// The size of the buffer in which the scanner stores the read characters.
const bufsize uint8 = 255

// Implements the Scanner interface.
// Stores the scanner state and a buffer of read bytes.
type scanner struct {
	reader io.Reader // The io.Reader from which the tokens will be read.

	buffer  [bufsize]byte // Temporary storage for bytes extracted from the reader but not yet processed.
	bufpos  uint8         // The position of the currently processed byte in the buffer.
	buflast uint8         // The number of bytes contained in the buffer.

	lineStr      []byte // Current processed line string.
	switchLine   bool   // true if the scanner read the string to the end.
	lineNum      int    // The number of the currently processed line.
	posNum       int    // The position of the currently processed character relative to the beginning of the byte sequence.
	skipComments bool   // true if comments should be skipped.
}

// Creates a new Scanner that reads from the reader.
// Sets skipping comments by default.
func NewScanner(reader io.Reader) Scanner {
	var scanner = scanner{reader: reader, skipComments: true}
	// Initialization: allocating memory and filling the buffer.
	scanner.refreshBuffer()
	scanner.refreshLine()
	scanner.lineNum = 0
	return Scanner(&scanner)
}

// Reads new values to the buffer.
// The number of bytes read is stored in the buflast field.
// The current bufpos is reset to 0.
func (scanner *scanner) refreshBuffer() {
	var n, err = scanner.reader.Read(scanner.buffer[:])
	if err != nil && err != io.EOF {
		panic(err)
	}
	scanner.buflast = uint8(n)
	scanner.bufpos = 0
}

// Moving the scanner to the next line.
func (scanner *scanner) refreshLine() {
	scanner.lineStr = make([]byte, 0, 100)
	scanner.lineNum++
}

// Returns true if there is a next token.
func (scanner *scanner) has() bool {
	// The buffer is processed to the end.
	// It is necessary to read the new data to the buffer.
	if scanner.bufpos == scanner.buflast {
		// If the number of types in the buffer is less than the buffer size,
		// it means that the buffer was not fully filled the previous time when reading it.
		if scanner.buflast < bufsize {
			return false
		} else {
			scanner.refreshBuffer()
		}
	}
	return scanner.bufpos != scanner.buflast
}

// Returns the next character from the reader.
// Panics if it can't get the next character, because this method is only used if the next character is present.
func (scanner *scanner) peek() byte {
	if scanner.has() {
		return scanner.buffer[scanner.bufpos]
	}
	// Impossible situation.
	panic("cannot get the next byte")
}

// Moves to the next character.
// Calls the peek method without checking the existence of the next character,
// so it must only be called if the next character exists.
func (scanner *scanner) step() {
	if scanner.switchLine {
		scanner.refreshLine()
		scanner.switchLine = false
	}
	var symbol = scanner.peek()
	if symbol == '\n' {
		scanner.switchLine = true
	} else {
		scanner.lineStr = append(scanner.lineStr, symbol)
	}
	scanner.bufpos++
	scanner.posNum++
}

// Implementation of the Next method in the Scanner interface.
func (scanner *scanner) Next() (TokenType, string) {
	// If all bytes are read from the reader, the scanner always returns the (EOF, "").
	if !scanner.has() {
		return EOF, ""
	}
	var (
		state     stateType // Contains the current state of finite state machine.
		symbol    byte      // Contains the character currently being processed.
		tokenType TokenType
		buffer    = make([]byte, 0, 100) // Contains the characters that were read.
	)
	for scanner.has() {
		symbol = scanner.peek()
		// Skipping the '\r' character to handle line ends on Windows
		if symbol == '\r' {
			if scanner.has() {
				scanner.step()
				symbol = scanner.peek()
			} else {
				symbol = '\n'
			}
		}
		tokenType = tokenTypeMap[state]
		state = matrix[getSymbolType(symbol)][state] // The next state is contained in the matrix.
		// The transition to the start state means the end of the token.
		if state == start {
			// If the comments are omitted, the next token must be returned.
			if scanner.skipComments && tokenType == Comment {
				return scanner.Next()
			}
			return tokenType, string(buffer)
		}
		buffer = append(buffer, symbol)
		scanner.step()
	}
	// All bytes are read from the reader.
	return tokenTypeMap[state], string(buffer)
}

// Implementation of the SkipLine method in the Scanner interface.
func (scanner *scanner) SkipLine() {
	if scanner.switchLine {
		return
	}
	var symbol byte
	for scanner.has() {
		symbol = scanner.peek()
		scanner.step()
		if symbol == '\n' {
			return
		}
	}
}

// Implementation of the LineString method in the Scanner interface.
func (scanner *scanner) LineString() string {
	return string(scanner.lineStr)
}

// Implementation of the Position method in the Scanner interface.
func (scanner *scanner) Position() int {
	return scanner.posNum - 1
}

// Implementation of the Line method in the Scanner interface.
func (scanner *scanner) Line() int {
	return scanner.lineNum
}

// Implementation of the Column method in the Scanner interface.
func (scanner *scanner) Column() int {
	if scanner.switchLine || !scanner.has() {
		return len(scanner.lineStr)
	}
	return len(scanner.lineStr) - 1
}

// Implementation of the IsSkipComments method in the Scanner interface.
func (scanner *scanner) IsSkipComments() bool {
	return scanner.skipComments
}

// Implementation of the SetSkipComments method in the Scanner interface.
func (scanner *scanner) SkipComments(skipComments bool) {
	scanner.skipComments = skipComments
}
