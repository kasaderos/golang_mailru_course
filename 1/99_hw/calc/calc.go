package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Lex struct {
	name   int
	val    int
	strval string
}

type Stack struct {
	data []*Lex
	size int
}

func (st *Stack) Push(lex Lex) {
	st.data = append(st.data, &Lex{name: lex.name, val: lex.val})
	st.size++
}

func (st *Stack) Pop() *Lex {
	st.size--
	res := st.data[st.size]
	st.data[st.size] = nil
	st.data = st.data[:st.size]
	return res
}

/*
E -> T {[ + | - ] T}
T -> F {[ * | / ] F}
F -> N | (E)
N -> R | NR
R -> 0 | 1 | .. | 9
*/

type Parser struct {
	stack  Stack
	tokens []string
	pos    int
	curLex Lex
	poliz  []Lex
}

const (
	PLUS = iota
	MINUS
	MUL
	DIV
	NUM
	ASSIGN
	LEX
)

var lexs map[string]int = map[string]int{
	"+": PLUS,
	"-": MINUS,
	"*": MUL,
	"/": DIV,
	"=": ASSIGN,
}

func (p *Parser) nextToken() {
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
	p.curLex = Lex{strval: p.tokens[p.pos]}
}

func (p *Parser) F() {
	fmt.Printf("F\n")
	var err error
	p.curLex.val, err = strconv.Atoi(p.curLex.strval)
	if err != nil {
		panic("Error: can't convert to int")
	}
	lex := Lex{name: NUM, val: p.curLex.val}
	p.stack.Push(lex)
	p.poliz = append(p.poliz, lex)
	p.nextToken()
}

func (p *Parser) T() {
	fmt.Printf("T\n")
	p.F()
	for p.curLex.strval == "*" || p.curLex.strval == "/" {
		p.stack.Push(Lex{name: lexs[p.curLex.strval]})
		p.nextToken()
		p.F()
		p.checkOp()
	}

}

func (p *Parser) E() {
	fmt.Printf("E\n")
	p.T()
	for p.curLex.strval == "+" || p.curLex.strval == "-" {
		p.stack.Push(Lex{name: lexs[p.curLex.strval]})
		p.nextToken()
		p.T()
		p.checkOp()
	}
}

func (p *Parser) parse(line string) []Lex {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered : ", r)
		}
	}()
	p.tokens = strings.Split(line, " ")
	p.pos = 0
	p.curLex = Lex{strval: p.tokens[p.pos]}
	p.stack = Stack{
		data: make([]*Lex, 0, 16),
		size: 0,
	}
	if lexs[p.curLex.strval] == ASSIGN {
		p.nextToken()
		p.E()
		p.poliz = append(p.poliz, Lex{name: ASSIGN})
	} else {
		panic("Error : where is the simbol '=' ?")
	}
	return p.poliz
}

func (p *Parser) checkOp() {
	p.stack.Pop()
	op := p.stack.Pop()
	p.stack.Pop()
	if op.name == PLUS || op.name == MINUS ||
		op.name == MUL || op.name == DIV {
		p.stack.Push(Lex{name: NUM})
	}
	p.poliz = append(p.poliz, Lex{name: op.name})
}
func main() {
	stack := Stack{
		data: make([]*Lex, 0, 16),
		size: 0,
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parser := Parser{}
		poliz := parser.parse(scanner.Text())
		//fmt.Println(poliz)
	POLIZ:
		for _, token := range poliz {
			switch token.name {
			case ASSIGN:
				fmt.Printf("Result = %d\n", stack.Pop().val)
			case PLUS:
				stack.Push(Lex{name: NUM, val: stack.Pop().val + stack.Pop().val})
			case MINUS:
				v1 := stack.Pop().val
				stack.Push(Lex{name: NUM, val: stack.Pop().val - v1})
			case MUL:
				stack.Push(Lex{name: NUM, val: stack.Pop().val * stack.Pop().val})
			case DIV:
				v1 := stack.Pop().val
				v2 := stack.Pop().val
				if v1 == 0 {
					fmt.Println("Error: division by zero")
					break POLIZ
				}
				stack.Push(Lex{name: NUM, val: v2 / v1})
			case NUM:
				stack.Push(token)
			default:
				break POLIZ
			}
		}
	}
}
