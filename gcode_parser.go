package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GCodeParser struct {
	lastCommand       string
	previous_position Vector3d
	previous_feedrate float64
	lastParams        map[string]float64
	allowedParamsMap  map[string]map[string]bool
}

type GCodeCommand struct {
	description string
	command     string
	params      map[string]float64
}

// New GCode Parser
func newGCodeParser() *GCodeParser {
	var allowedParamsMap = map[string]map[string]bool{
		"G0": {"X": true, "Y": true, "Z": true, "A": true, "B": true, "C": true, "U": true, "V": true, "W": true, "F": true},
		"G1": {"X": true, "Y": true, "Z": true, "A": true, "B": true, "C": true, "U": true, "V": true, "W": true, "F": true},
		"G2": {"X": true, "Y": true, "Z": true, "I": true, "J": true, "K": true, "R": true, "P": true, "F": true},
		"G3": {"X": true, "Y": true, "Z": true, "I": true, "J": true, "K": true, "R": true, "P": true, "F": true},
		/*"G4":  {"P": true, "S": true},
		"G20": {},
		"G21": {},
		"G28": {"X": true, "Y": true, "Z": true, "A": true, "B": true, "C": true, "U": true, "V": true, "W": true},
		"G30": {"X": true, "Y": true, "Z": true, "A": true, "B": true, "C": true, "U": true, "V": true, "W": true},
		"G90": {},
		"G91": {},
		"G92": {"X": true, "Y": true, "Z": true, "A": true, "B": true, "C": true, "U": true, "V": true, "W": true, "E": true},
		"G93": {},
		"G94": {},
		"M0":  {"P": true},
		"M1":  {"P": true},
		"M2":  {"P": true},
		"M30": {"P": true},
		"M3":  {"S": true},
		"M4":  {"S": true},
		"M5":  {"S": true},
		"M6":  {"T": true},
		"M7":  {"P": true},
		"M8":  {"P": true},
		"M9":  {"P": true},
		"M48": {},
		"M49": {},*/
	}

	var GCodeParser = GCodeParser{

		allowedParamsMap: allowedParamsMap,
	}

	return &GCodeParser
}

func (p *GCodeParser) parseCommand(line string) *GCodeCommand {

	if len(line) == 0 || line[0] == ';' || line[0] == '(' || line[0] == '%' {
		return &GCodeCommand{line, "Comment", nil}
	}

	tokens := strings.Fields(line)
	if len(tokens) < 1 {
		fmt.Println("Empty line")
		return nil
	}

	command := tokens[0]
	if !strings.HasPrefix(command, "G") {
		command = p.lastCommand
	} else {
		tokens = tokens[1:]
	}

	params := make(map[string]float64)
	for key, val := range p.lastParams {
		params[key] = val
	}

	for _, token := range tokens {
		if len(token) < 2 {
			fmt.Println("invalid token: %s", token)
			return nil
		}

		val, err := strconv.ParseFloat(token[1:], 64)
		if err != nil {
			fmt.Println("invalid value: %s", token[1:])
			return nil
		}

		params[string(token[0])] = val
	}

	p.lastCommand = command
	p.lastParams = params

	validParams := p.validateAndFilterParams(command, params)

	return &GCodeCommand{line, command, validParams}
}

func (p *GCodeParser) validateAndFilterParams(command string, params map[string]float64) map[string]float64 {
	validParams := make(map[string]float64)

	for key, val := range params {
		if allowed, ok := p.allowedParamsMap[command][key]; ok && allowed {
			validParams[key] = val
		}
	}

	return validParams
}

func (g *GCodeParser) fromString(gcodeStringList []string) []GCodeCommand {

	var gcodeLines []GCodeCommand

	for _, line := range gcodeStringList {
		gcodeLine := g.parseCommand(line)

		if gcodeLine != nil {
			gcodeLines = append(gcodeLines, *gcodeLine)
		}
	}

	return gcodeLines

}

func (g *GCodeParser) fromFile(filename string) []GCodeCommand {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Array of GCodeCommand

	var gcodeStringList []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		gcodeStringList = append(gcodeStringList, line)
	}

	gCodeList := g.fromString(gcodeStringList)

	return gCodeList

}
