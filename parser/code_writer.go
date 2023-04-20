package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// "SP": 0, "LCL": 1, "ARG": 2, "THIS": 3, "THAT": 4,
type CodeWriter struct {
	outputFile *os.File
	scanner    *bufio.Scanner
}

var arithmeticOpMap = map[string]string{
	"add": "+",
	"sub": "-",
}

var vmToAssembly = map[string]string{
	"SP--": "@SP\nM=M-1\n",
	"*SP":  "@SP\nA=M\nD=M\n",
	"SP++": "@SP\nM=M+1\n",
}

var memorySegmentBaseAddress = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
}

func MakeCodeWriter(output_file *os.File) CodeWriter {
	codeWriter := CodeWriter{
		outputFile: output_file,
		scanner:    bufio.NewScanner(output_file),
	}

	// initialize base addresses and stack pointer

	// initializeStack := "@256\nD=A\n@SP\nM=D\n"
	// _, err := codeWriter.outputFile.WriteString(initializeStack)
	// if err != nil {
	// 	panic("Failed to write to output file")
	// }
	return codeWriter
}

func (codeWriter *CodeWriter) WriteArithmetic(cmd string) {
	op, ok := arithmeticOpMap[cmd] // TODO: Implement logical operators
	if !ok {
		panic("Operation is not add or sub")
	}
	assemblyInstructions := vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
		vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nD=M" + op + "D\n" +
		"@SP\nA=M\nM=D\n" + vmToAssembly["SP++"]

	codeWriter.outputFile.WriteString("// " + cmd + "\n")
	codeWriter.outputFile.WriteString(assemblyInstructions)
}

func (codeWriter *CodeWriter) WritePushPop(cmd Command, segment string, index int) {
	if segment == "local" || segment == "argument" || segment == "this" || segment == "that" {
		codeWriter.writePushPopSegmentI(cmd, segment, index)
	} else if segment == "constant" {
		codeWriter.writePushPopConstant(cmd, segment, index)
	} else if segment == "static" {
		codeWriter.writePushPopStatic(cmd, segment, index)
	} else if segment == "temp" {
		codeWriter.writePushPopTemp(cmd, segment, index)
	} else if segment == "pointer" {
		codeWriter.writePushPopPointer(cmd, segment, index)
	}
}

func (codeWriter *CodeWriter) writePushPopSegmentI(cmd Command, segment string, index int) {
	switch cmd {
	case C_PUSH:
		{
			pushI := fmt.Sprintf("@%s\nA=M+%d\nD=A\n@R15\nM=D\n", memorySegmentBaseAddress[segment], index) +
				vmToAssembly["*SP"] + "@R15\nA=M\nM=D\n" + vmToAssembly["SP++"]

			codeWriter.outputFile.WriteString(fmt.Sprintf("// push %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(pushI)
		}
	case C_POP:
		{
			popI := fmt.Sprintf("@%s\nA=M+%d\nD=A\n@R15\nM=D\n", memorySegmentBaseAddress[segment], index) + vmToAssembly["SP--"] +
				vmToAssembly["*SP"] + "@R15\nA=M\nM=D\n"

			codeWriter.outputFile.WriteString(fmt.Sprintf("// pop %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(popI)
		}
	}
}

func (codeWriter *CodeWriter) writePushPopConstant(cmd Command, segment string, index int) {
	switch cmd {
	case C_PUSH:
		{
			pushConst := fmt.Sprintf("@SP\nA=M\nM=%d\n", index) + vmToAssembly["SP++"]

			codeWriter.outputFile.WriteString(fmt.Sprintf("// push %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(pushConst)
		}
	case C_POP:
		{
			panic("pop const not implemented")

		}
	}
}

func (codeWriter *CodeWriter) writePushPopStatic(cmd Command, segment string, index int) {
	// get file name
	fileInfo, err := codeWriter.outputFile.Stat()
	if err != nil {
		panic("failed to get file info")
	}

	fileName := strings.Split(fileInfo.Name(), ".")[0]

	switch cmd {
	case C_PUSH:
		{
			pushI := fmt.Sprintf("@%s.%d\nD=M\n@SP\nA=M\nM=D\n", fileName, index) + vmToAssembly["SP++"]

			codeWriter.outputFile.WriteString(fmt.Sprintf("// push %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(pushI)

		}
	case C_POP:
		{
			popI := vmToAssembly["SP--"] + vmToAssembly["*SP"] + fmt.Sprintf("@%s.%d\nM=D\n", fileName, index)

			codeWriter.outputFile.WriteString(fmt.Sprintf("// pop %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(popI)
		}
	}
}

func (codeWriter *CodeWriter) writePushPopTemp(cmd Command, segment string, index int) {
	switch cmd {
	case C_PUSH:
		{
			pushI := fmt.Sprintf("@5\nA=A+%d\nD=M\n@SP\nA=M\nM=D\n", index) + vmToAssembly["SP++"]

			codeWriter.outputFile.WriteString(fmt.Sprintf("// push %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(pushI)
		}
	case C_POP:
		{
			popI := vmToAssembly["SP--"] + vmToAssembly["*SP"] + fmt.Sprintf("@5\nA=A+%d\nM=D\n", index)

			codeWriter.outputFile.WriteString(fmt.Sprintf("// pop %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(popI)
		}
	}
}

func (codeWriter *CodeWriter) writePushPopPointer(cmd Command, segment string, index int) {
	op := ""
	if index == 0 {
		op = "THIS"
	} else {
		op = "THAT"
	}

	switch cmd {
	case C_PUSH:
		{
			pushI := "@" + op + "\nD=M\n@SP\nA=M\nM=D\n" + vmToAssembly["SP++"]

			codeWriter.outputFile.WriteString(fmt.Sprintf("// push %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(pushI)
		}
	case C_POP:
		{
			popI := vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@" + op + "\nM=D\n"

			codeWriter.outputFile.WriteString(fmt.Sprintf("// pop %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(popI)
		}
	}
}
