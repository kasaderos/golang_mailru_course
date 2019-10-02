package main

/*
Калькулятор

Программа похожим образом работает
как интерпритатор python3, только без объявлений и
инициализаций переменных

В случае неправильного ввода выводит 0
Описание принимающих выражений и соответствующих вызовов

E -> T {[ + | - ] T}
T -> F {[ * | / ] F}
F -> N | (E)
N -> R | NR
R -> 0 | 1 | .. | 9
*/

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

const (
	PLUS = iota
	MINUS
	MUL
	DIV
	NUM
	LPAREN
	RPAREN
	LEXEND
)

var lexs map[string]int = map[string]int{
	"+": PLUS,
	"-": MINUS,
	"*": MUL,
	"/": DIV,
	"(": LPAREN,
	")": RPAREN,
}

type Lex struct {
	name int
	val  int
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

type Parser struct {
	charPos int
	char    rune
	line    []rune
	stack   Stack
	curLex  Lex
	poliz   []Lex
	minus   bool
}

func (p *Parser) nextToken() {
	p.curLex = p.getLex()
}

func (p *Parser) F() {
	if p.curLex.name == LPAREN {
		p.nextToken()
		p.E()
		if p.curLex.name != RPAREN {
			panic("Error: rparen not closed")
		}
		p.nextToken()
	} else {
		p.stack.Push(p.curLex)
		p.poliz = append(p.poliz, p.curLex)
		p.nextToken()
	}
}

func (p *Parser) T() {
	p.F()
	for p.curLex.name == MUL || p.curLex.name == DIV {
		p.stack.Push(Lex{name: p.curLex.name})
		p.nextToken()
		p.F()
		p.checkOp()
	}
}

func (p *Parser) E() {
	p.T()
	for p.curLex.name == PLUS || p.curLex.name == MINUS {
		p.stack.Push(Lex{name: p.curLex.name})
		p.nextToken()
		p.T()
		p.checkOp()
	}
}

func (p *Parser) gc() {
	if p.charPos < len(p.line) {
		p.char = p.line[p.charPos]
		p.charPos++
	} else {
		p.char = '\n'
	}
}

func (p *Parser) getLex() (lex Lex) {
	state := 0
	var buf []rune
	for {
		switch state {
		case 0:
			if p.char == '\n' {
				return Lex{name: LEXEND}
			} else if p.char == ' ' || p.char == '\t' {
				p.gc()
				break
			} else if unicode.IsDigit(p.char) {
				state = 1
			} else {
				state = 2
			}
		case 1:
			p.minus = false
			if unicode.IsDigit(p.char) {
				buf = append(buf, p.char)
				p.gc()
			} else {
				if v, err := strconv.Atoi(string(buf)); err == nil {
					return Lex{name: NUM, val: v}
				}
				panic("getLex: lex analysis failed")
			}
		case 2:
			if _, ok := lexs[string(p.char)]; ok {
				if p.charPos == 1 || p.char == '(' {
					p.minus = true
				}
				if p.char == '-' && p.minus {
					state = 1
					buf = append(buf, '-')
					p.gc()
					break
				}
				lex = Lex{name: lexs[string(p.char)]}
				p.gc()
				return lex
			}
			panic("getLex: lex analysis failed")
		}
	}
}

func (p *Parser) parse() []Lex {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from parse: ", r)
		}
	}()
	p.gc()
	p.stack = Stack{
		data: make([]*Lex, 0, 16),
		size: 0,
	}
	p.nextToken()
	if p.curLex.name == LEXEND {
		return nil
	}
	p.E()
	if p.curLex.name != LEXEND {
		panic("Error LEXEND")
	}
	return append(p.poliz, Lex{name: LEXEND})
}

func executePoliz(poliz *[]Lex) (res int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from executePoliz: ", r)
		}
	}()
	stack := Stack{
		data: make([]*Lex, 0, 16),
		size: 0,
	}
	for _, lex := range *poliz {
		switch lex.name {
		case LEXEND:
			return stack.Pop().val
		case PLUS:
			stack.Push(Lex{name: NUM, val: stack.Pop().val + stack.Pop().val})
		case MINUS:
			v1 := stack.Pop().val
			stack.Push(Lex{name: NUM, val: stack.Pop().val - v1})
		case MUL:
			stack.Push(Lex{name: NUM, val: stack.Pop().val * stack.Pop().val})
		case DIV:
			v1 := stack.Pop().val
			stack.Push(Lex{name: NUM, val: stack.Pop().val / v1})
		case NUM:
			stack.Push(lex)
		default:
			fmt.Println("Error : unknown expression")
		}
	}
	return
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
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf(">>> ")
	for scanner.Scan() {
		parser := Parser{line: []rune(scanner.Text())}
		lexs := parser.parse()
		fmt.Println(executePoliz(&lexs))
		fmt.Printf(">>> ")
	}
	fmt.Printf("\n")
}
