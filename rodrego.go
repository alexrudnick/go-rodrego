/*
Simple golang implementation of the RodRego register machine. See README.md for
more details.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
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

/* Replace any Mac newlines with Unix newlines. Because, y'know, there are Mac
 * newlines in the example programs. DOS newlines are handled by the golang
 * standard bufio.ScanLines function. */
func magicSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	innerline := regexp.MustCompile("\r([^\n])")
	endline := regexp.MustCompile("\r$")
	replaced := innerline.ReplaceAll(data, []byte("\n$1"))
	replaced = endline.ReplaceAll(replaced, []byte("\n"))
	advance, token, err = bufio.ScanLines(replaced, atEOF)
	return
}

/* Given a file name, load a rodrego program and find the initial line to
 * execute. */
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

	lineno := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lineno += 1
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		fields := strings.Fields(line)
		line_name := fields[0]
		inst_s := strings.ToLower(fields[1])

		if start == "" {
			start = line_name
		}

		var stmt Statement
		switch inst_s {
		case "end":
			stmt = Statement{END, 0, "", ""}
		case "deb":
			target, err := strconv.ParseInt(fields[2], 10, 64)
			branch := fields[3]
			elsebranch := fields[4]
			stmt = Statement{DEB, target, branch, elsebranch}
			if err != nil {
				fmt.Println("error on text line:", lineno)
				fmt.Println("Target registers must be valid integers.")
				os.Exit(1)
			}
		case "inc":
			target, err := strconv.ParseInt(fields[2], 10, 64)
			branch := fields[3]
			stmt = Statement{INC, target, branch, ""}
			if err != nil {
				fmt.Println("error on text line:", lineno)
				fmt.Println("Target registers must be valid integers.")
				os.Exit(1)
			}
		default:
			fmt.Println("error on text line:", lineno)
			fmt.Println("valid instructions are INC, DEB and END")
			os.Exit(1)
		}
		out[line_name] = stmt
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error on text line:", lineno)
		log.Fatal(err)
	}
	return
}

/* Read in the initial state of the registers from a file. */
func load_registers(registersfn string, registers *map[int64]int64) {
	f, err := os.Open(registersfn)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	scanner.Split(magicSplit)

	lineno := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lineno += 1
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Println("error on text line:", lineno)
			fmt.Println("required format for register file lines:")
			fmt.Println("<register number> <register value>")
			os.Exit(1)
		}
		register_name := fields[0]
		reg, err := strconv.ParseInt(register_name, 10, 64)
		if err != nil || reg < 0 {
			fmt.Println("error on text line:", lineno)
			fmt.Println("registers must be referenced with natural numbers.")
			os.Exit(1)
		}
		register_value := fields[1]
		val, err := strconv.ParseInt(register_value, 10, 64)
		if err != nil || val < 0 {
			fmt.Print("line ", lineno, ": ")
			fmt.Println("register values must be natural numbers.")
			os.Exit(1)
		}
		(*registers)[reg] = val
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

/* Method to print out Statement structs */
func (stmt Statement) Println() {
	switch stmt.inst {
	case END:
		fmt.Println("END")
	case INC:
		fmt.Printf("INC register %d and GOTO %s\n", stmt.target, stmt.branch)
	case DEB:
		fmt.Printf("DEB register %d and GOTO %s else %s\n",
			stmt.target, stmt.branch, stmt.elsebranch)
	}
}

/* Display the current values for all the set registers */
func printRegisters(registers *map[int64]int64) {
	if len(*registers) == 0 {
		fmt.Println("[ all registers empty ]")
	}
	for k, v := range *registers {
		fmt.Println("register", k, "=", v)
	}
}

/* Run a rodrego program, given a start state and a pointer to the registers.
* Will probably mutate the registers. */
func execute(program map[string]Statement, start string,
	registers *map[int64]int64, step bool) {
	current := start
	bio := bufio.NewReader(os.Stdin)
	for {
		stmt := program[current]
		fmt.Println("[[ now on line:", current, "]]")
		printRegisters(registers)
		fmt.Print("performing: ")
		stmt.Println()

		switch stmt.inst {
		case END:
			return
		case INC:
			(*registers)[stmt.target] += 1
			current = stmt.branch
		case DEB:
			if (*registers)[stmt.target] == 0 {
				current = stmt.elsebranch
			} else {
				(*registers)[stmt.target] -= 1
				current = stmt.branch
			}
		}

		if step {
			fmt.Println("ENTER to continue...")
			bio.ReadLine()
		}
	}
}

func main() {
	var infn string
	var valuesfn string
	var step bool
	flag.StringVar(&infn, "program", "",
		"filename of a rodrego program to execute (required)")
	flag.StringVar(&valuesfn, "values", "",
		"filename for a set of initial register values")
	flag.BoolVar(&step, "step", false,
		"if true, step through program one instruction at a time")
	flag.Parse()

	if infn == "" {
		fmt.Println("please specify a program to execute. See", os.Args[0],
			"-help")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	program, start := load_program(infn)

	registers := make(map[int64]int64)
	if valuesfn != "" {
		load_registers(valuesfn, &registers)
	}

	execute(program, start, &registers, step)
	fmt.Println("*** Final state of the world ***")
	printRegisters(&registers)
}
