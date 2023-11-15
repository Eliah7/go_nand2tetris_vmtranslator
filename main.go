package main

import (
	"fmt"
	"log"
	"os"
	"path"
	p "vmtranlater/parser"
)

func main() {
	// fmt.Println("VMTranslater")
	if len(os.Args) == 1 {
		panic("Specify the Instructions source file")
	}

	directoryPath := os.Args[1]
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	var vmFiles []string
	for _, e := range entries {
		extension := path.Ext(e.Name())
		if extension == ".vm" {
			vmFiles = append(vmFiles, fmt.Sprintf("%s%s", directoryPath, e.Name()))
		}

	}

	// fmt.Println(vmFiles)

	// convert everything below here to a function
	// if the input path is a directory then iterate across all files and
	// execute the function
	for _, file_ := range vmFiles {
		translateFile(file_)
	}

}

func translateFile(file string) {
	f, err := os.Open(file)

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
			// fmt.Println(command, arg1, arg2)
			codeWriter.WriteLabel(arg1)
		}
	}
}
