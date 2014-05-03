package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instruction int

const (
	INC Instruction = iota
	DEB
	END
)

type Statement struct {
	inst       Instruction
	target     int64
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
		fields := strings.Fields(line)
		line_name := fields[0]
		inst_s := strings.ToLower(fields[1])

		if start == "" {
			start = line_name
		}

		var stmt Statement
		if inst_s == "end" {
			stmt = Statement{END, 0, "", ""}
		} else if inst_s == "deb" {
			target, err := strconv.ParseInt(fields[2], 10, 64)
			branch := fields[3]
			elsebranch := fields[3]
			stmt = Statement{DEB, target, branch, elsebranch}
			if err != nil {
				fmt.Println("Target registers must be valid integers.")
				os.Exit(1)
			}
		} else if inst_s == "inc" {
			target, err := strconv.ParseInt(fields[2], 10, 64)
			branch := fields[3]
			stmt = Statement{INC, target, branch, ""}
			if err != nil {
				fmt.Println("Target registers must be valid integers.")
				os.Exit(1)
			}
		}
		out[line_name] = stmt
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
