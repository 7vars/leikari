package query

import (
	"strconv"
	"strings"

	"github.com/7vars/leikari"
)

type tokenType int

const (
	illegal tokenType = iota
	eof
	lowest
	and
	or
	lparen
	rparen
	eq
	ne
	co
	sw
	ew
	gt
	ge
	lt
	le
	not
	pr
	ident
	int_type
	float_type
	bool_type
	string_type
)

var (
	keywords = map[string]tokenType {
		"EQ": eq, "eq": eq,
		"NE": ne, "ne": ne,
		"CO": co, "co": co, 
		"SW": sw, "sw": sw,
		"EW": ew, "ew": ew,
		"PR": pr, "pr": pr,
		"GT": gt, "gt": gt,
		"GE": ge, "ge": ge,
		"LT": lt, "lt": lt,
		"LE": le, "le": le,
		"AND": and, "and": and,
		"OR": or, "or": or,

		"true": bool_type, "false": bool_type,
	}
)

type token struct {
	ttype tokenType
	literal string
}

func newToken(ttype tokenType, literal string) token {
	return token{
		ttype: ttype,
		literal: literal,
	}
}

type lex struct {
	input string
	position int
	readPosition int
	ch byte
}

func newLex(input string) *lex {
	l := &lex{
		input: input,
	}
	l.readChar()
	return l
}

func (l *lex) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *lex) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *lex) nextToken() token {
	var tok token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(eq, string(ch) + string(l.ch))
		} else {
			tok = newToken(illegal, "")
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(ne, string(ch) + string(l.ch))
		} else {
			tok = newToken(illegal, "")
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(ge, string(ch) + string(l.ch))
		} else {
			tok = newToken(gt, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = newToken(le, string(ch) + string(l.ch))
		} else {
			tok = newToken(lt, string(l.ch))
		}
	case '(':
		tok = newToken(lparen, string(l.ch))
	case ')':
		tok = newToken(rparen, string(l.ch))
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = newToken(and, string(ch) + string(l.ch))
		} else {
			tok = newToken(illegal, "")
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = newToken(or, string(ch) + string(l.ch))
		} else {
			tok = newToken(illegal, "")
		}
	case '\'':
		alpha := l.readAlpha()
		tok = newToken(string_type, alpha)
		return tok
	case 0:
		tok = newToken(eof, "")
	default:
		if isLetter(l.ch) {
			ident := l.readAttribute()
			tok = newToken(lookupIdentifier(ident), ident)
			return tok
		} else if isDigit(l.ch) {
			num := l.readNumber()
			tok = newToken(numToken(num), num)
			return tok
		} else {
			tok = newToken(illegal, "")
		}
	}

	l.readChar()
	return tok
}

func (l *lex) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *lex) readAttribute() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *lex) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}
func (l *lex) readAlpha() string {
	position := l.position
	l.readChar()
	for l.ch != '\'' {
		l.readChar()
	}
	l.readChar()
	return l.input[position+1:l.position-1]
}

func lookupIdentifier(s string) tokenType {
	if tok, ok := keywords[s]; ok {
		return tok
	}
	return ident
}

func numToken(num string) tokenType {
	if strings.Contains(num, ".") {
		return float_type
	}
	return int_type
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

type prefixParseFunc func() (Node, error)
type infixParseFunc func(Node) (Condition, error)

type parser struct {
	l *lex
	curToken token
	peekToken token

	prefixParseFuncs map[tokenType]prefixParseFunc
	infixParseFuncs map[tokenType]infixParseFunc
}

func newParser(query string) *parser {
	p := &parser{
		l: newLex(query),
		prefixParseFuncs: make(map[tokenType]prefixParseFunc),
		infixParseFuncs: make(map[tokenType]infixParseFunc),
	}

	p.registerPrefixParseFunc(ident, p.parseIdentifier)
	p.registerPrefixParseFunc(int_type, p.parseInteger)
	p.registerPrefixParseFunc(float_type, p.parseFloat)
	p.registerPrefixParseFunc(bool_type, p.parseBool)
	p.registerPrefixParseFunc(string_type, p.parseString)
	p.registerPrefixParseFunc(lparen, p.parseGroup)


	p.registerInfixParseFunc(eq, p.parseInfix)
	p.registerInfixParseFunc(ne, p.parseInfix)
	p.registerInfixParseFunc(co, p.parseInfix)
	p.registerInfixParseFunc(sw, p.parseInfix)
	p.registerInfixParseFunc(ew, p.parseInfix)
	p.registerInfixParseFunc(gt, p.parseInfix)
	p.registerInfixParseFunc(ge, p.parseInfix)
	p.registerInfixParseFunc(lt, p.parseInfix)
	p.registerInfixParseFunc(le, p.parseInfix)
	p.registerInfixParseFunc(and, p.parseInfix)
	p.registerInfixParseFunc(or, p.parseInfix)

	p.nextToken()
	p.nextToken()

	return p
}

func Parse(query string) (Node, error) {
	parser := newParser(query)

	if parser.curToken.ttype != eof {
		return parser.parseNode(lowest)
	}

	return nil, leikari.Errorln("", "no valid query defined")
}

func (p *parser) parseNode(tt tokenType) (Node, error) {
	pf, ok := p.prefixParseFuncs[p.curToken.ttype]
	if !ok {
		return nil, leikari.Errorf("", "could not parse %s (%d)", p.curToken.literal, p.curToken.ttype)
	}

	left, err := pf()
	if err != nil {
		return nil, err
	}
	for tt < p.peekToken.ttype {
		ifi, ok := p.infixParseFuncs[p.peekToken.ttype]
		if !ok {
			return left, nil
		}

		p.nextToken()

		left, err = ifi(left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.nextToken()
}

func (p *parser) registerPrefixParseFunc(tok tokenType, f prefixParseFunc) {
	p.prefixParseFuncs[tok] = f
}

func (p *parser) registerInfixParseFunc(tok tokenType, f infixParseFunc) {
	p.infixParseFuncs[tok] = f
}

func (p *parser) parseInfix(left Node) (Condition, error) {
	if left == nil {
		return nil, leikari.Errorln("", "left node is nil")
	}

	tt := p.curToken.ttype

	p.nextToken()
	right, err := p.parseNode(tt) 
	if err != nil {
		return nil, err
	}

	switch tt {
	case eq,ne,co,sw,ew,gt,ge,lt,le:
		ident, ok := left.(Identifier)
		if !ok { 
			return nil, leikari.Errorf("", "identifier expected %T found: %v", left, left)
		}
		value, ok := right.(Value)
		if !ok {
			return nil, leikari.Errorf("", "value expected %T found: %v", right, right)
		}
		
		return newComparsion(ident, Operator(tt), value)
	case and:
		return And(left, right), nil
	case or:
		return Or(left, right), nil
	}

	return nil, leikari.Errorf("", "(%v) is not an operator", tt)
}

func (p *parser) parseIdentifier() (Node, error) {
	return newIdentifier(p.curToken.literal), nil
}

func (p *parser) parseInteger() (Node, error) {
	value, err := strconv.ParseInt(p.curToken.literal, 10, 64)
	if err != nil {
		return nil, err
	}
	return newIntValue(value), nil
}

func (p *parser) parseFloat() (Node, error) {
	value, err := strconv.ParseFloat(p.curToken.literal, 64)
	if err != nil {
		return nil, err
	}
	return newFloatValue(value), nil
}

func (p *parser) parseBool() (Node, error) {
	value, err := strconv.ParseBool(p.curToken.literal)
	if err != nil {
		return nil, err
	}
	return newBoolValue(value), nil
}

func (p *parser) parseString() (Node, error) {
	return newStringValue(p.curToken.literal), nil
}

func (p *parser) parseGroup() (Node, error) {
	p.nextToken()
	exp, err := p.parseNode(lowest)
	if err != nil {
		return nil, err
	}

	grp := Group(exp)

	if p.peekToken.ttype != rparen {
		return nil, leikari.Errorln("", "missing close paren")
	}

	p.nextToken()

	return grp, nil
}