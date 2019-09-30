package main

/*
Калькулятор

Программа похожим образом работает
как интерпритатор python3, только без объявлений и
инициализаций переменных
Весь код переписан с С++ и изменен (не из википедии)

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

type Parser struct {
	chNum  int
	char   rune
	line   []rune
	stack  Stack
	tokens []string
	pos    int
	curLex Lex
	poliz  []Lex
	minus  bool
}

const (
	PLUS = iota
	MINUS
	MUL
	DIV
	NUM
	ASSIGN
	LEX
	LPAREN
	RPAREN
	LEXEND
)

var lexs map[string]int = map[string]int{
	"+": PLUS,
	"-": MINUS,
	"*": MUL,
	"/": DIV,
	"=": ASSIGN,
	"(": LPAREN,
	")": RPAREN,
}

func (p *Parser) nextToken() {
	p.curLex = p.getLex()
}

func (p *Parser) F() {
	//fmt.Printf("F %s\n", toprint[p.curLex.name])
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
	//fmt.Printf("T %s \n", toprint[p.curLex.name])
	p.F()
	for p.curLex.name == MUL || p.curLex.name == DIV {
		p.stack.Push(Lex{name: p.curLex.name})
		p.nextToken()
		p.F()
		p.checkOp()
	}
}

func (p *Parser) E() {
	//fmt.Printf("E %s\n", toprint[p.curLex.name])
	p.T()
	for p.curLex.name == PLUS || p.curLex.name == MINUS {
		p.stack.Push(Lex{name: p.curLex.name})
		p.nextToken()
		p.T()
		p.checkOp()
	}
}

func (p *Parser) gc() {
	if p.chNum < len(p.line) {
		p.char = p.line[p.chNum]
		p.chNum++
	} else {
		p.char = '\n'
	}
}

func (p *Parser) getLex() (lex Lex) {
	st := 0
	var buf []rune
	for {
		//fmt.Println("ST ", st, "CHAR", p.char)
		switch st {
		case 0:
			if p.char == ' ' || p.char == '\t' {
				p.gc()
				continue
			} else if p.char == '\n' {
				return Lex{name: LEXEND}
			} else if unicode.IsDigit(p.char) {
				st = 1
			} else {
				st = 2
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
				panic("getLex")
			}
		case 2:
			if _, ok := lexs[string(p.char)]; ok {
				if p.chNum == 1 || p.char == '(' {
					p.minus = true
				}
				if p.char == '-' && p.minus {
					st = 1
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

/*
var toprint map[int]string = map[int]string{
	PLUS:   "+",
	MINUS:  "-",
	MUL:    "*",
	DIV:    "/",
	NUM:    "num",
	ASSIGN: "=",
	LEX:    "lex",
	LPAREN: "(",
	RPAREN: ")",
	LEXEND: "lexend",
}*/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf(">>> ")
	for scanner.Scan() {
		parser := Parser{line: []rune(scanner.Text())}
		lexs := parser.parse()
		fmt.Println(executePoliz(&lexs))
		fmt.Printf(">>> ")
	}
}
