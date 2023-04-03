package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GCodeParser struct {
	lastCmd           string
	previous_position Vector3d
	previous_feedrate float64
	lastParams        map[string]float64
}

// New GCode Parser
func NewGCodeParser() *GCodeParser {
	return &GCodeParser{}
}

type Command interface {
	Execute()
}

type G0 struct {
	X, Y, Z float64
}

func (g G0) Execute() {
	fmt.Printf("G0 Rapid: \t X=%8.4f, Y=%8.4f, Z=%8.4f\n", g.X, g.Y, g.Z)
}

type G1 struct {
	X, Y, Z, F float64
}

func (g G1) Execute() {
	fmt.Printf("G1 Linear: \t X=%8.4f, Y=%8.4f, Z=%8.4f, F=%8.4f\n", g.X, g.Y, g.Z, g.F)
}

type G2 struct {
	X, Y, Z, I, J, F float64
}

func (g G2) Execute() {
	fmt.Printf("G2: CW Arc: \t X=%8.4f, Y=%8.4f, Z=%8.4f, I=%8.4f, J=%8.4f, F=%8.4f\n", g.X, g.Y, g.Z, g.I, g.J, g.F)
}

type G3 struct {
	X, Y, Z, I, J, F float64
}

func (g G3) Execute() {
	fmt.Printf("G3: CCW Arc: \t X=%8.4f, Y=%8.4f, Z=%8.4f, I=%8.4f, J=%8.4f, F=%8.4f\n", g.X, g.Y, g.Z, g.I, g.J, g.F)
}

func (p *GCodeParser) parseCommand(line string) (Command, error) {
	tokens := strings.Fields(line)
	if len(tokens) < 1 {
		return nil, fmt.Errorf("empty line")
	}

	cmd := tokens[0]
	if !strings.HasPrefix(cmd, "G") {
		cmd = p.lastCmd
	} else {
		tokens = tokens[1:]
	}

	params := make(map[string]float64)
	for key, val := range p.lastParams {
		params[key] = val
	}

	for _, token := range tokens {
		if len(token) < 2 {
			return nil, fmt.Errorf("invalid token: %s", token)
		}

		val, err := strconv.ParseFloat(token[1:], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value: %s", token[1:])
		}

		params[string(token[0])] = val
	}

	p.lastCmd = cmd
	p.lastParams = params

	switch cmd {
	case "G0", "G1", "G2", "G3":
		switch cmd {
		case "G0":
			return G0{X: params["X"], Y: params["Y"], Z: params["Z"]}, nil
		case "G1":
			return G1{X: params["X"], Y: params["Y"], Z: params["Z"], F: params["F"]}, nil
		case "G2":
			return G2{X: params["X"], Y: params["Y"], Z: params["Z"], I: params["I"], J: params["J"], F: params["F"]}, nil
		case "G3":
			return G3{X: params["X"], Y: params["Y"], Z: params["Z"], I: params["I"], J: params["J"], F: params["F"]}, nil
		}
	}

	return nil, fmt.Errorf("unsupported command: %s", cmd)
}

func (g *GCodeParser) Parse(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var commands []Command

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == ';' || line[0] == '(' || line[0] == '%' {
			continue
		}

		command, err := g.parseCommand(line)
		if err != nil {
			fmt.Printf("Error parsing line: %v\n", err)
			continue
		}

		commands = append(commands, command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
		return
	}

	for _, command := range commands {
		command.Execute()
	}
}
