package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

var Arithmetic_Commands = []string{"add", "sub", "gt", "lt", "and", "or", "not"}

func MakeParser(file *os.File) Parser {
	return Parser{
		inputStream: file,
		scanner:     bufio.NewScanner(file),
	}
}

func (parser *Parser) HasMoreCommands() bool {
	hasMoreCommands := parser.scanner.Scan()
	line := parser.scanner.Text()

	comments_r, _ := regexp.Compile(`^\/\/\s*[\w\d\s",:=+/.([\])!#$%\^&\*_\-\{\}\|:;"><,.~]*$`) // remove all comments
	spaces_r, _ := regexp.Compile(`^\s*$`)

	if (spaces_r.MatchString(line) || comments_r.MatchString(line)) && hasMoreCommands { // line is a comment or space
		return parser.HasMoreCommands()
	} else {
		return hasMoreCommands
	}
}

func (parser *Parser) Advance() (Command, string, interface{}) {
	line := parser.scanner.Text()
	line = strings.Split(line, "//")[0]
	// fmt.Println(line)

	cmd := parser.commandType(line)
	arg1 := parser.arg1(line)
	arg2 := parser.arg2(line)
	return cmd, arg1, arg2
}

func (parser *Parser) commandType(line string) Command {
	cmps := strings.Split(line, "//")
	cmds := strings.Split(cmps[0], " ")

	if cmds[0] == "pop" {
		return C_POP
	} else if cmds[0] == "push" {
		return C_PUSH
	} else { // TODO: Add other commands as they are needed
		return C_ARITHMETIC
	}
}

func (parser *Parser) arg1(line string) string {
	cmps := strings.Split(line, "//")
	cmds := strings.Split(cmps[0], " ")
	if len(cmds) > 1 {
		return cmds[1]
	}
	return cmds[0]
}

func (parser *Parser) arg2(line string) interface{} {
	cmps := strings.Split(line, "//")
	cmds := strings.Split(cmps[0], " ")

	if len(cmds) > 2 {
		fmt.Println(cmps, " ", cmds[2])
		arg2, err := strconv.Atoi(cmds[2])
		if err != nil {
			panic("Could not convert arg2 to string")
		}
		return arg2
	}
	return nil
}
