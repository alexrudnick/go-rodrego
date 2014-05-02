package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

type Instruction int

const (
	INC Instruction = iota
	DEB
)

type Statement struct {
	inst       Instruction
	target     int
	branch     string
	elsebranch string
}

func magicSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	newdata := bytes.Replace(data, []byte("\r"), []byte("\n"), -1)
	advance, token, err = bufio.ScanLines(newdata, atEOF)
	return
}

func load_program(infn string) (out map[string]Statement, start string) {
	out = make(map[string]Statement)
	start = ""

	f, err := os.Open(infn)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	scanner := bufio.NewScanner(reader)
	scanner.Split(magicSplit)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		fmt.Println(line)
		fields := strings.Fields(line)
		line_name := fields[0]
		inst := fields[1]
		target := fields[2]
		branch := fields[3]
		elsebranch := ""
		if strings.ToLower(inst) == "deb" {
			elsebranch = fields[3]
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("need more arguments")
		os.Exit(1)
		return
	}
	infn := os.Args[1]
	program, start := load_program(infn)
	fmt.Println(program)
	fmt.Println(start)
}
