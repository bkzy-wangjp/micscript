package calc

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	FLOAT   = "FLOAT"

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	LPAREN = "("
	RPAREN = ")"
)

const (
	_ int = iota
	LOWEST
	SUM     // +, -
	PRODUCT // *, /
	PREFIX  // -X
	CALL    // (X)
)

var precedences = map[string]int{
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
	LPAREN:   CALL,
}

/*四则运算*/
func Calc(input string) (float64, error) {
	reg := regexp.MustCompile(`^(\s*\(*\s*\-?\+?\d+)[-0-9e\+\*\/\(\)\.\s]*((\d*\s*)|(\d*\s*\)\s*))$`)
	if reg.Match([]byte(input)) == false {
		return 0.0, fmt.Errorf("Expression error[表达式错误]:[%s]", input)
	}
	reg = regexp.MustCompile(`\(`)
	s1 := reg.FindAllString(input, -1)
	reg = regexp.MustCompile(`\)`)
	s2 := reg.FindAllString(input, -1)
	if len(s1) != len(s2) {
		return 0.0, fmt.Errorf("Number of '()' does not match[表达式中括号的数量不匹配]:[%s]", input)
	}
	lexer := NewLex(input)
	parser := NewParser(lexer)
	exp := parser.ParseExpression(LOWEST)
	return Eval(exp), nil
}

/*四则运算和比较运算*/
func CalcAndCompare(input string) (float64, error) {
	_mathSymbol := `((\>\=|\<\=|\=\=)|(\+|\-|\*|\/|\>|\<|\=|\(\s*|\s*\)){1})` //解析符号
	_number := `\d+\.*\d*e*`
	zz := fmt.Sprintf("%s|%s", _mathSymbol, _number)
	reg := regexp.MustCompile(zz)
	str := reg.FindAllString(input, -1)

	var leftstr, operater, rightstr string
	for _, v := range str {
		if v == ">" || v == "<" || v == ">=" || v == "<=" || v == "==" {
			operater = v
		} else {
			if len(operater) == 0 {
				leftstr = fmt.Sprintf("%s%s", leftstr, v)
			} else {
				rightstr = fmt.Sprintf("%s%s", rightstr, v)
			}
		}
	}
	var rt float64
	leftv, err := Calc(leftstr) //计算左边
	if err != nil {
		return 0.0, err
	}
	if len(rightstr) > 0 {
		rightv, err := Calc(rightstr) //计算右边
		if err != nil {
			return 0.0, err
		}
		switch operater {
		case ">":
			if leftv > rightv {
				rt = 1
			} else {
				rt = 0
			}
		case "<":
			if leftv < rightv {
				rt = 1
			} else {
				rt = 0
			}
		case ">=":
			if leftv >= rightv {
				rt = 1
			} else {
				rt = 0
			}
		case "<=":
			if leftv <= rightv {
				rt = 1
			} else {
				rt = 0
			}
		case "==":
			if leftv == rightv {
				rt = 1
			} else {
				rt = 0
			}
		default:
			rt = 0
		}
	} else {
		rt = leftv
	}
	return rt, nil
}

func Eval(exp Expression) float64 {
	switch node := exp.(type) {
	case *IntegerLiteralExpression:
		return node.Value
	case *PrefixExpression:
		rightV := Eval(node.Right)
		return evalPrefixExpression(node.Operator, rightV)
	case *InfixExpression:
		leftV := Eval(node.Left)
		rightV := Eval(node.Right)
		return evalInfixExpression(leftV, node.Operator, rightV)
	}

	return 0
}

func evalPrefixExpression(operator string, right float64) float64 {
	if operator != "-" {
		return 0
	}
	return -right
}

func evalInfixExpression(left float64, operator string, right float64) float64 {

	switch operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		if right != 0 {
			return left / right
		} else {
			return 0
		}
	default:
		return 0
	}
}

//词元
type Token struct {
	Type    string //对应我们上面的词元类型
	Literal string // 实际的string字符
}

func newToken(tokenType string, c byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(c),
	}
}

//词法分析器
type Lexer struct {
	input        string // 输入
	position     int    // 当前位置
	readPosition int    // 将要读取的位置
	ch           byte   //当前字符
	lastch       byte   //前一个字符(by wjp)
}

func NewLex(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

//词法分析器的核心函数,用于获取下一个词元。针对不同字符返回不同的Token结构值
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()
	if l.lastch == 'e' && (l.ch == '+' || l.ch == '-') {
		if isDigit(l.ch, l.lastch) {
			tok.Type = FLOAT
			tok.Literal = l.readNumber()
			return tok
		}
	} else {
		switch l.ch {
		case '(':
			tok = newToken(LPAREN, l.ch)
		case ')':
			tok = newToken(RPAREN, l.ch)
		case '+':
			tok = newToken(PLUS, l.ch)
		case '-':
			tok = newToken(MINUS, l.ch)
		case '/':
			tok = newToken(SLASH, l.ch)
		case '*':
			tok = newToken(ASTERISK, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = EOF
		default:
			if isDigit(l.ch, l.lastch) {
				tok.Type = FLOAT
				tok.Literal = l.readNumber()
				return tok
			} else {
				tok = newToken(ILLEGAL, l.ch)
			}
		}
	}

	l.readChar()
	return tok
}

//判断是否数字
func isDigit(ch, lastch byte) bool {
	return ('0' <= ch && ch <= '9') || '.' == ch || ch == 'e' || (lastch == 'e' && (ch == '+' || ch == '-'))
}

//读取数字
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch, l.lastch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

//不断读取字符，并且更新结构体的值
func (l *Lexer) readChar() {
	l.lastch = l.ch
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

//用于在读取时候直接跳过空白字符
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// ast

type Expression interface {
	String() string
}

type IntegerLiteralExpression struct {
	Token Token
	Value float64
}

func (il *IntegerLiteralExpression) String() string { return il.Token.Literal }

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// parser
type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l *Lexer

	curToken  Token
	peekToken Token

	prefixParseFns map[string]prefixParseFn
	infixParseFns  map[string]infixParseFn

	errors []string
}

func (p *Parser) registerPrefix(tokenType string, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType string, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[string]prefixParseFn)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[string]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	returnExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return returnExp
		}

		p.nextToken()
		returnExp = infix(returnExp)
	}

	return returnExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekError(t string) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instend",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t string) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekTokenIs(t string) bool {
	return p.peekToken.Type == t
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseFloatLiteral() Expression {

	lit := &IntegerLiteralExpression{Token: p.curToken}

	//value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() Expression {

	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.ParseExpression(PREFIX)
	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.ParseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseInfixExpression(left Expression) Expression {

	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	// // 通过降低优先级，来达到右结合
	//if expression.Operator == "+" {
	//	expression.Right = p.parseExpression(precedence - 1)
	//} else {
	//	expression.Right = p.parseExpression(precedence)
	//}
	expression.Right = p.ParseExpression(precedence)

	return expression
}

func (p *Parser) Errors() []string {
	return p.errors
}
