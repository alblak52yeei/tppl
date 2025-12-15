package main

import (
	"fmt"
)

// Node представляет узел AST
type Node interface {
	String() string
}

// Program представляет программу
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	return fmt.Sprintf("Program(%d statements)", len(p.Statements))
}

// Statement представляет оператор
type Statement interface {
	Node
	statementNode()
}

// Assignment представляет присваивание
type Assignment struct {
	Variable string
	Value    Expression
}

func (a *Assignment) statementNode() {
	_ = a // маркерный метод
}
func (a *Assignment) String() string {
	return fmt.Sprintf("Assignment(%s := %s)", a.Variable, a.Value)
}

// Block представляет блок BEGIN ... END
type Block struct {
	Statements []Statement
}

func (b *Block) statementNode() {
	_ = b // маркерный метод
}
func (b *Block) String() string {
	return fmt.Sprintf("Block(%d statements)", len(b.Statements))
}

// Expression представляет выражение
type Expression interface {
	Node
	expressionNode()
}

// Number представляет число
type Number struct {
	Value float64
}

func (n *Number) expressionNode() {
	_ = n // маркерный метод
}
func (n *Number) String() string {
	return fmt.Sprintf("Number(%g)", n.Value)
}

// Identifier представляет переменную
type Identifier struct {
	Name string
}

func (i *Identifier) expressionNode() {
	_ = i // маркерный метод
}
func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier(%s)", i.Name)
}

// BinaryOp представляет бинарную операцию
type BinaryOp struct {
	Left     Expression
	Operator TokenType
	Right    Expression
}

func (b *BinaryOp) expressionNode() {
	_ = b // маркерный метод
}
func (b *BinaryOp) String() string {
	op := ""
	switch b.Operator {
	case TokenPLUS:
		op = "+"
	case TokenMINUS:
		op = "-"
	case TokenMULTIPLY:
		op = "*"
	case TokenDIVIDE:
		op = "/"
	}
	return fmt.Sprintf("BinaryOp(%s %s %s)", b.Left, op, b.Right)
}

// Parser представляет парсер
type Parser struct {
	tokens []Token
	pos    int
}

// NewParser создает новый парсер
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Parse разбирает токены в AST
func (p *Parser) Parse() (*Program, error) {
	program := &Program{}
	
	// Ожидаем BEGIN
	if !p.match(TokenBEGIN) {
		return nil, fmt.Errorf("ожидался BEGIN на позиции %d", p.current().Pos)
	}
	
	// Парсим блок
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	
	program.Statements = block.Statements
	
	// Ожидаем END
	if !p.match(TokenEND) {
		return nil, fmt.Errorf("ожидался END на позиции %d", p.current().Pos)
	}
	
	// Ожидаем точку
	if !p.match(TokenDOT) {
		return nil, fmt.Errorf("ожидалась точка на позиции %d", p.current().Pos)
	}
	
	// Проверяем EOF
	if !p.match(TokenEOF) {
		return nil, fmt.Errorf("неожиданные токены после точки")
	}
	
	return program, nil
}

// parseBlock парсит блок BEGIN ... END
func (p *Parser) parseBlock() (*Block, error) {
	block := &Block{Statements: []Statement{}}
	
	for {
		// Проверяем, не конец ли блока
		if p.check(TokenEND) {
			break
		}
		
		// Парсим оператор
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		
		block.Statements = append(block.Statements, stmt)
		
		// Проверяем, не конец ли блока после точки с запятой
		if p.check(TokenEND) {
			break
		}
	}
	
	return block, nil
}

// parseStatement парсит оператор
func (p *Parser) parseStatement() (Statement, error) {
	// Проверяем, не вложенный ли блок
	if p.match(TokenBEGIN) {
		block, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		
		if !p.match(TokenEND) {
			return nil, fmt.Errorf("ожидался END на позиции %d", p.current().Pos)
		}
		
		// Точка с запятой после END (если не последний оператор)
		if p.check(TokenSEMICOLON) {
			p.advance()
		}
		
		return block, nil
	}
	
	// Парсим присваивание
	if p.check(TokenIDENTIFIER) {
		varName := p.current().Value
		p.advance()
		
		if !p.match(TokenASSIGN) {
			return nil, fmt.Errorf("ожидался := на позиции %d", p.current().Pos)
		}
		
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		
		// Точка с запятой после присваивания
		if p.check(TokenSEMICOLON) {
			p.advance()
		}
		
		return &Assignment{
			Variable: varName,
			Value:    expr,
		}, nil
	}
	
	return nil, fmt.Errorf("неожиданный токен на позиции %d: %v", p.current().Pos, p.current())
}

// parseExpression парсит выражение (с учетом приоритета операций)
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseAdditive()
}

// parseAdditive парсит аддитивные операции (+ и -)
func (p *Parser) parseAdditive() (Expression, error) {
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}
	
	for p.check(TokenPLUS) || p.check(TokenMINUS) {
		op := p.current().Type
		p.advance()
		
		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		
		left = &BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}
	
	return left, nil
}

// parseMultiplicative парсит мультипликативные операции (* и /)
func (p *Parser) parseMultiplicative() (Expression, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	
	for p.check(TokenMULTIPLY) || p.check(TokenDIVIDE) {
		op := p.current().Type
		p.advance()
		
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		
		left = &BinaryOp{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}
	
	return left, nil
}

// parseUnary парсит унарные выражения и первичные выражения
func (p *Parser) parseUnary() (Expression, error) {
	if p.check(TokenMINUS) {
		p.advance()
		expr, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &BinaryOp{
			Left:     &Number{Value: 0},
			Operator: TokenMINUS,
			Right:    expr,
		}, nil
	}
	
	return p.parsePrimary()
}

// parsePrimary парсит первичные выражения (числа, переменные, скобки)
func (p *Parser) parsePrimary() (Expression, error) {
	if p.check(TokenNUMBER) {
		value := 0.0
		fmt.Sscanf(p.current().Value, "%f", &value)
		p.advance()
		return &Number{Value: value}, nil
	}
	
	if p.check(TokenIDENTIFIER) {
		name := p.current().Value
		p.advance()
		return &Identifier{Name: name}, nil
	}
	
	if p.match(TokenLPAREN) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		
		if !p.match(TokenRPAREN) {
			return nil, fmt.Errorf("ожидалась закрывающая скобка на позиции %d", p.current().Pos)
		}
		
		return expr, nil
	}
	
	return nil, fmt.Errorf("неожиданный токен на позиции %d: %v", p.current().Pos, p.current())
}

func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) check(t TokenType) bool {
	return p.current().Type == t
}

func (p *Parser) match(t TokenType) bool {
	if p.check(t) {
		p.advance()
		return true
	}
	return false
}

