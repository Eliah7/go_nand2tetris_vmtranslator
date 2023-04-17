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

	file_path := fmt.Sprintf("%s%s", f.Name()[0:len(f.Name())-2], "asm")

	output_file, err := os.Create(file_path)
	if err != nil {
		panic("Failed to create output file")
	}

	defer output_file.Close()
	// _ = p.MakeCodeWriter(output_file) // code_writer

	for parser.HasMoreCommands() {
		parser.Advance()
	}

}
