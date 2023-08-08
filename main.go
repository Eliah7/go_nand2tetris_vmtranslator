package main

import (
	"fmt"
	"log"
	"os"
	p "vmtranlater/parser"
)

func main() {
	// fmt.Println("VMTranslater")
	if len(os.Args) == 1 {
		panic("Specify the Instructions source file")
	}

	instructionsFilePath := os.Args[1]
	f, err := os.Open(instructionsFilePath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	parser := p.MakeParser(f)

	filePath := fmt.Sprintf("%s%s", f.Name()[0:len(f.Name())-2], "asm")

	outputFile, err := os.Create(filePath)
	if err != nil {
		panic("Failed to create output file")
	}

	defer outputFile.Close()
	codeWriter := p.MakeCodeWriter(outputFile) // code_writer

	for parser.HasMoreCommands() {
		command, arg1, arg2 := parser.Advance()

		if command == p.C_ARITHMETIC {
			codeWriter.WriteArithmetic(arg1)
		} else if command == p.C_PUSH || command == p.C_POP {
			codeWriter.WritePushPop(command, arg1, arg2.(int))
		} else if command == p.C_LABEL {
			fmt.Println(command, arg1, arg2)
			// codeWriter.WriteLabel(arg1)
		}
	}

}
