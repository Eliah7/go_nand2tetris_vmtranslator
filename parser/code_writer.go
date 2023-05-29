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
	boolCount  int
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
		boolCount:  0,
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
	// op, ok := arithmeticOpMap[cmd] // TODO: Implement logical operators and other arithmetic operators
	// if !ok {
	// 	panic("Operation is not add or sub")
	// }
	var assemblyInstructions string
	// fmt.Println(cmd)
	if cmd == "add" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
			vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nD=M" + "+" + "D\n" +
			"@SP\nA=M\nM=D\n" + vmToAssembly["SP++"]
	} else if cmd == "sub" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
			vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nD=D" + "-" + "M\n" +
			"@SP\nA=M\nM=D\n" + vmToAssembly["SP++"]
	} else if cmd == "or" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
			vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nD=D" + "|" + "M\n" +
			"@SP\nA=M\nM=D\n" + vmToAssembly["SP++"]
	} else if cmd == "and" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
			vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nD=D" + "&" + "M\n" +
			"@SP\nA=M\nM=D\n" + vmToAssembly["SP++"]
	} else if cmd == "eq" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
			vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\n" +
			fmt.Sprintf("D=D-M\n@IS_TRUE_%d\nD;JEQ\n", codeWriter.boolCount) +
			fmt.Sprintf("@IS_FALSE_%d\nD;JNE\n", codeWriter.boolCount) +
			fmt.Sprintf("(IS_TRUE_%[1]d)\n@SP\nA=M\nM=-1\n@END_%[1]d\n0;JMP\n(IS_FALSE_%[1]d)\n@SP\nA=M\nM=0\n@END_%[1]d\n0;JMP\n(END_%[1]d)\n", codeWriter.boolCount) + vmToAssembly["SP++"]
	} else if cmd == "gt" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
			vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\n" +
			fmt.Sprintf("D=D-M\n@IS_TRUE_%d\nD;JGT\n", codeWriter.boolCount) +
			fmt.Sprintf("@IS_FALSE_%d\nD;JLE\n", codeWriter.boolCount) +
			fmt.Sprintf("(IS_TRUE_%[1]d)\n@SP\nA=M\nM=-1\n@END_%[1]d\n0;JMP\n(IS_FALSE_%[1]d)\n@SP\nA=M\nM=0\n@END_%[1]d\n0;JMP\n(END_%[1]d)\n", codeWriter.boolCount) + vmToAssembly["SP++"]
	} else if cmd == "lt" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\nM=D\n" +
			vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R13\n" +
			fmt.Sprintf("D=D-M\n@IS_TRUE_%[1]d\nD;JLT\n", codeWriter.boolCount) +
			fmt.Sprintf("@IS_FALSE_%[1]d\nD;JGE\n", codeWriter.boolCount) +
			fmt.Sprintf("(IS_TRUE_%[1]d)\n@SP\nA=M\nM=-1\n@END_%[1]d\n0;JMP\n(IS_FALSE_%[1]d)\n@SP\nA=M\nM=0\n@END_%[1]d\n0;JMP\n(END_%[1]d)\n", codeWriter.boolCount) + vmToAssembly["SP++"]
	} else if cmd == "not" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "M=!M\n" + vmToAssembly["SP++"]
	} else if cmd == "neg" {
		assemblyInstructions = vmToAssembly["SP--"] + vmToAssembly["*SP"] + "M=-M\n" + vmToAssembly["SP++"]
	}

	codeWriter.outputFile.WriteString("// " + cmd + "\n")
	codeWriter.outputFile.WriteString(assemblyInstructions)
	codeWriter.boolCount += 1
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
			pushI := fmt.Sprintf("@%d\nD=A\n@%s\nA=D+M\nD=M\n", index, memorySegmentBaseAddress[segment]) + "@SP\nA=M\nM=D\n" + vmToAssembly["SP++"]

			codeWriter.outputFile.WriteString(fmt.Sprintf("// push %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(pushI)
		}
	case C_POP:
		{
			popI := vmToAssembly["SP--"] +
				fmt.Sprintf("@%d\nD=A\n@%s\nA=D+M\nD=A\n@R14\nM=D\n", index, memorySegmentBaseAddress[segment]) + // @%d\nD=A\n@%s\nA=D+M\nD=A\n@R15\nM=D\n + "@R15\nA=M\nM=D\n"
				vmToAssembly["*SP"] + "@R14\nA=M\nM=D\n"
			codeWriter.outputFile.WriteString(fmt.Sprintf("// pop %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(popI)
		}
	}
}

func (codeWriter *CodeWriter) writePushPopConstant(cmd Command, segment string, index int) {
	switch cmd {
	case C_PUSH:
		{
			pushConst := fmt.Sprintf("@%d\nD=A\n@SP\nA=M\nM=D\n", index) + vmToAssembly["SP++"]

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
			pushI := fmt.Sprintf("@%d\nD=A\n@5\nA=D+A\nD=M\n@SP\nA=M\nM=D\n", index) + vmToAssembly["SP++"]

			codeWriter.outputFile.WriteString(fmt.Sprintf("// push %s %d\n", segment, index))
			codeWriter.outputFile.WriteString(pushI)
		}
	case C_POP:
		{
			popI := fmt.Sprintf("@%d\nD=A\n@5\nD=D+A\n@R15\nM=D\n", index) + vmToAssembly["SP--"] + vmToAssembly["*SP"] + "@R15\nA=M\nM=D\n"

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
