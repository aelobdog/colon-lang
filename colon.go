package main

import (
	"colon/colinterp"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	version   = "1.0.0Alpha"
	author    = "Ashwin Godbole"
	email     = "dev.godboleashwin@gmail.com"
	linkToSrc = "https://github.com/ashvin-godbole/colon-lang"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}
	code, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file : " + os.Args[1])
		return
	}
	colinterp.Interpret(string(code))
}

func usage() {
	fmt.Println()
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("             Colon Programming Language Interpreter")
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("Interpreter version : %s\n", version)
	fmt.Printf("Author              : %s\n", author)
	fmt.Printf("Email               : %s\n", email)
	fmt.Printf("Link to source code : %s\n", linkToSrc)
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("Usage:")
	fmt.Println("       colon <filename>.col")
	fmt.Println("------------------------------------------------------------------")
}
