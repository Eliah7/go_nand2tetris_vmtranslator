package parser

import (
	"bufio"
	"os"
)

type Parser struct {
	inputStream *os.File
	scanner     *bufio.Scanner
}

type Command int

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

func MakeParser(file *os.File) Parser {
	return Parser{
		inputStream: file,
		scanner:     bufio.NewScanner(file),
	}
}

func (parser *Parser) commandType() Command {
	_ = parser.scanner.Text()
	// fmt.Println(line)
	// return line
	return C_ARITHMETIC // TODO: Get command from file
}

func (parser *Parser) HasMoreCommands() bool {
	return parser.scanner.Scan()
	// return true // return if it has more commands
}

func (parser *Parser) Advance() (Command, string, int) {
	cmd := parser.commandType()
	arg1 := parser.arg1()
	arg2 := parser.arg2()

	return cmd, arg1, arg2
}

func (parser *Parser) arg1() string {
	line := parser.scanner.Text()
	// fmt.Println(line)
	return line
}

func (parser *Parser) arg2() int {
	_ = parser.scanner.Text()
	// fmt.Println(line)
	// return line
	return 0
}
