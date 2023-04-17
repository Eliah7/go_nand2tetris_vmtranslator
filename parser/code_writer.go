package parser

import (
	"bufio"
	"os"
)

type CodeWriter struct {
	outputFile *os.File
	scanner    *bufio.Scanner
}

func MakeCodeWriter(output_file *os.File) CodeWriter {
	return CodeWriter{
		outputFile: output_file,
		scanner:    bufio.NewScanner(output_file),
	}
}
