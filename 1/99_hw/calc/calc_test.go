package main

import (
	"fmt"
	"testing"
)

func TestExecutePoliz(t *testing.T) {
	var testingVars map[string]int = map[string]int{
		"0":                        0,
		"-123":                     -123,
		"1 + 1":                    2,
		"31 / (3 - (4 + (5 * 6)))": -1,
		"-2 * (-12-(-8)-(-12))":    -16,
		"1 / 0":                    0,
		"1/(-1+1)":                 0,
		"1 / (1-1)":                0,
		"12 / 1":                   12,
		"-2 * ((-12-(-8)-(-12)) - 31 / (3 - (4 + (5 * 6))))": -18,
	}
	for k, _ := range testingVars {
		fmt.Printf("RUN \"%s\"\n", k)
		parser := Parser{line: []rune(k)}
		lexs := parser.parse()
		result := executePoliz(&lexs)
		if testingVars[k] != result {
			t.Error("expected", testingVars[k], "have", result)
		}
		fmt.Printf(" %d\n-----\n OK\n-----", result)
	}
	var testingBadVars map[string]int = map[string]int{
		"adfasdfasdf":                     0,
		"*15a54adf8a7-96+*--*":            0,
		"1 ++ 1":                          0,
		"31 /121fdas (3 - (4 + (5 * 6)))": 0,
		"-2 * (-12-f(-8)-(-12))":          0,
		"1 / 0qw":                         0,
		"1/(-qwe1+1)":                     0,
		"1e / (1-1)":                      0,
		"qw12 / 1":                        0,
		"(132fv - (4f + (5 * 6))))":       0,
		"(1 + 1":                          0,
		"--111111":                        0,
	}
	fmt.Printf("Bad runs : \n")
	for k, _ := range testingBadVars {
		fmt.Printf("RUN \"%s\"\n", k)
		parser := Parser{line: []rune(k)}
		lexs := parser.parse()
		result := executePoliz(&lexs)
		if testingVars[k] != result {
			t.Error("expected", testingBadVars[k], "have", result)
		}
		fmt.Printf("OK\n-----")
	}
}
