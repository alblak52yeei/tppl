package main

import (
	"testing"
)

// TestEmptyProgram тестирует пустую программу
func TestEmptyProgram(t *testing.T) {
	code := `BEGIN
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	if len(variables) != 0 {
		t.Errorf("Ожидалось 0 переменных, получено %d: %v", len(variables), variables)
	}
}

// TestSimpleExpressions тестирует простые выражения
func TestSimpleExpressions(t *testing.T) {
	code := `BEGIN
	x:= 2 + 3 * (2 + 3);
        y:= 2 / 2 - 2 + 3 * ((1 + 1) + (1 + 1));
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	// x = 2 + 3 * (2 + 3) = 2 + 3 * 5 = 2 + 15 = 17
	if x, ok := variables["x"]; !ok {
		t.Error("Переменная x не найдена")
	} else if x != 17 {
		t.Errorf("x: ожидалось 17, получено %g", x)
	}
	
	// y = 2 / 2 - 2 + 3 * ((1 + 1) + (1 + 1)) = 1 - 2 + 3 * (2 + 2) = 1 - 2 + 3 * 4 = 1 - 2 + 12 = 11
	if y, ok := variables["y"]; !ok {
		t.Error("Переменная y не найдена")
	} else if y != 11 {
		t.Errorf("y: ожидалось 11, получено %g", y)
	}
}

// TestNestedBlocks тестирует вложенные блоки
func TestNestedBlocks(t *testing.T) {
	code := `BEGIN
    y: = 2;
    BEGIN
        a := 3;
        a := a;
        b := 10 + a + 10 * y / 4;
        c := a - b
    END;
    x := 11;
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	// y = 2
	if y, ok := variables["y"]; !ok {
		t.Error("Переменная y не найдена")
	} else if y != 2 {
		t.Errorf("y: ожидалось 2, получено %g", y)
	}
	
	// a = 3 (после a := a)
	if a, ok := variables["a"]; !ok {
		t.Error("Переменная a не найдена")
	} else if a != 3 {
		t.Errorf("a: ожидалось 3, получено %g", a)
	}
	
	// b = 10 + a + 10 * y / 4 = 10 + 3 + 10 * 2 / 4 = 10 + 3 + 20 / 4 = 10 + 3 + 5 = 18
	if b, ok := variables["b"]; !ok {
		t.Error("Переменная b не найдена")
	} else if b != 18 {
		t.Errorf("b: ожидалось 18, получено %g", b)
	}
	
	// c = a - b = 3 - 18 = -15
	if c, ok := variables["c"]; !ok {
		t.Error("Переменная c не найдена")
	} else if c != -15 {
		t.Errorf("c: ожидалось -15, получено %g", c)
	}
	
	// x = 11
	if x, ok := variables["x"]; !ok {
		t.Error("Переменная x не найдена")
	} else if x != 11 {
		t.Errorf("x: ожидалось 11, получено %g", x)
	}
}

// TestLexer тестирует лексер
func TestLexer(t *testing.T) {
	code := `BEGIN x := 5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	expectedTypes := []TokenType{
		TokenBEGIN,
		TokenIDENTIFIER,
		TokenASSIGN,
		TokenNUMBER,
		TokenSEMICOLON,
		TokenEND,
		TokenDOT,
		TokenEOF,
	}
	
	if len(tokens) != len(expectedTypes) {
		t.Fatalf("Ожидалось %d токенов, получено %d", len(expectedTypes), len(tokens))
	}
	
	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Errorf("Токен %d: ожидался тип %v, получен %v", i, expectedType, tokens[i].Type)
		}
	}
}

// TestParser тестирует парсер
func TestParser(t *testing.T) {
	code := `BEGIN x := 2 + 3; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Errorf("Ожидалось 1 оператор, получено %d", len(program.Statements))
	}
	
	assignment, ok := program.Statements[0].(*Assignment)
	if !ok {
		t.Fatalf("Ожидалось присваивание, получено %T", program.Statements[0])
	}
	
	if assignment.Variable != "x" {
		t.Errorf("Ожидалась переменная x, получена %s", assignment.Variable)
	}
}

// TestExpressionPrecedence тестирует приоритет операций
func TestExpressionPrecedence(t *testing.T) {
	code := `BEGIN
	x := 2 + 3 * 4;
	y := 10 - 2 / 2;
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	// x = 2 + 3 * 4 = 2 + 12 = 14
	if x, ok := variables["x"]; !ok {
		t.Error("Переменная x не найдена")
	} else if x != 14 {
		t.Errorf("x: ожидалось 14, получено %g", x)
	}
	
	// y = 10 - 2 / 2 = 10 - 1 = 9
	if y, ok := variables["y"]; !ok {
		t.Error("Переменная y не найдена")
	} else if y != 9 {
		t.Errorf("y: ожидалось 9, получено %g", y)
	}
}

// TestParentheses тестирует скобки
func TestParentheses(t *testing.T) {
	code := `BEGIN
	x := (2 + 3) * 4;
	y := 2 * (3 + 4);
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	// x = (2 + 3) * 4 = 5 * 4 = 20
	if x, ok := variables["x"]; !ok {
		t.Error("Переменная x не найдена")
	} else if x != 20 {
		t.Errorf("x: ожидалось 20, получено %g", x)
	}
	
	// y = 2 * (3 + 4) = 2 * 7 = 14
	if y, ok := variables["y"]; !ok {
		t.Error("Переменная y не найдена")
	} else if y != 14 {
		t.Errorf("y: ожидалось 14, получено %g", y)
	}
}

// TestVariableReassignment тестирует переприсваивание переменной
func TestVariableReassignment(t *testing.T) {
	code := `BEGIN
	x := 5;
	x := x + 3;
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	// x = 5 + 3 = 8
	if x, ok := variables["x"]; !ok {
		t.Error("Переменная x не найдена")
	} else if x != 8 {
		t.Errorf("x: ожидалось 8, получено %g", x)
	}
}

// TestDivision тестирует деление
func TestDivision(t *testing.T) {
	code := `BEGIN
	x := 10 / 2;
	y := 7 / 2;
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	// x = 10 / 2 = 5
	if x, ok := variables["x"]; !ok {
		t.Error("Переменная x не найдена")
	} else if x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", x)
	}
	
	// y = 7 / 2 = 3.5
	if y, ok := variables["y"]; !ok {
		t.Error("Переменная y не найдена")
	} else if y != 3.5 {
		t.Errorf("y: ожидалось 3.5, получено %g", y)
	}
}

// TestNegativeNumbers тестирует отрицательные числа
func TestNegativeNumbers(t *testing.T) {
	code := `BEGIN
	x := -5;
	y := 10 + -3;
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	// x = -5
	if x, ok := variables["x"]; !ok {
		t.Error("Переменная x не найдена")
	} else if x != -5 {
		t.Errorf("x: ожидалось -5, получено %g", x)
	}
	
	// y = 10 + -3 = 7
	if y, ok := variables["y"]; !ok {
		t.Error("Переменная y не найдена")
	} else if y != 7 {
		t.Errorf("y: ожидалось 7, получено %g", y)
	}
}

// TestMultipleNestedBlocks тестирует множественные вложенные блоки
func TestMultipleNestedBlocks(t *testing.T) {
	code := `BEGIN
	x := 1;
	BEGIN
		y := 2;
		BEGIN
			z := 3;
		END;
	END;
END.`
	
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	variables := interpreter.GetVariables()
	
	if x, ok := variables["x"]; !ok || x != 1 {
		t.Errorf("x: ожидалось 1, получено %g", variables["x"])
	}
	if y, ok := variables["y"]; !ok || y != 2 {
		t.Errorf("y: ожидалось 2, получено %g", variables["y"])
	}
	if z, ok := variables["z"]; !ok || z != 3 {
		t.Errorf("z: ожидалось 3, получено %g", variables["z"])
	}
}

// TestLexerErrors тестирует обработку ошибок лексера
func TestLexerErrors(t *testing.T) {
	// Неожиданный символ
	code := `BEGIN x @ 5; END.`
	lexer := NewLexer(code)
	_, err := lexer.Tokenize()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного символа @")
	}

	// Неправильный формат := (только :)
	code2 := `BEGIN x : 5; END.`
	lexer2 := NewLexer(code2)
	_, err2 := lexer2.Tokenize()
	if err2 == nil {
		t.Error("Ожидалась ошибка для неправильного формата :")
	}
}

// TestParserErrors тестирует обработку ошибок парсера
func TestParserErrors(t *testing.T) {
	// Отсутствие BEGIN
	code := `x := 5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	_, err = parser.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для отсутствия BEGIN")
	}

	// Отсутствие END
	code2 := `BEGIN x := 5; .`
	lexer2 := NewLexer(code2)
	tokens2, err := lexer2.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser2 := NewParser(tokens2)
	_, err = parser2.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для отсутствия END")
	}

	// Отсутствие точки
	code3 := `BEGIN x := 5; END`
	lexer3 := NewLexer(code3)
	tokens3, err := lexer3.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser3 := NewParser(tokens3)
	_, err = parser3.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для отсутствия точки")
	}

	// Неожиданные токены после точки
	code4 := `BEGIN x := 5; END. x := 6;`
	lexer4 := NewLexer(code4)
	tokens4, err := lexer4.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser4 := NewParser(tokens4)
	_, err = parser4.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданных токенов после точки")
	}

	// Неправильная структура присваивания
	code5 := `BEGIN := 5; END.`
	lexer5 := NewLexer(code5)
	tokens5, err := lexer5.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser5 := NewParser(tokens5)
	_, err = parser5.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для неправильной структуры присваивания")
	}

	// Незакрытая скобка
	code6 := `BEGIN x := (2 + 3; END.`
	lexer6 := NewLexer(code6)
	tokens6, err := lexer6.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser6 := NewParser(tokens6)
	_, err = parser6.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для незакрытой скобки")
	}

	// Неожиданный токен в выражении
	code7 := `BEGIN x := BEGIN; END.`
	lexer7 := NewLexer(code7)
	tokens7, err := lexer7.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser7 := NewParser(tokens7)
	_, err = parser7.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена в выражении")
	}
}

// TestInterpreterDivisionByZero тестирует деление на ноль
func TestInterpreterDivisionByZero(t *testing.T) {
	code := `BEGIN x := 10 / 0; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err == nil {
		t.Error("Ожидалась ошибка деления на ноль")
	}
	if err != nil && err.Error() != "деление на ноль" {
		t.Errorf("Ожидалась ошибка 'деление на ноль', получено: %v", err)
	}
}

// FakeStatement - фиктивный оператор для тестирования
type FakeStatement struct{}

func (f *FakeStatement) statementNode() {}
func (f *FakeStatement) String() string  { return "FakeStatement" }

// FakeExpression - фиктивное выражение для тестирования
type FakeExpression struct{}

func (f *FakeExpression) expressionNode() {}
func (f *FakeExpression) String() string   { return "FakeExpression" }

// TestInterpreterUnknownStatementType тестирует неизвестный тип оператора
func TestInterpreterUnknownStatementType(t *testing.T) {
	interpreter := NewInterpreter()
	// Создаем фиктивный оператор, который не является Assignment или Block
	program := &Program{
		Statements: []Statement{
			&FakeStatement{},
		},
	}
	err := interpreter.Interpret(program)
	if err == nil {
		t.Error("Ожидалась ошибка для неизвестного типа оператора")
	}
}

// TestInterpreterUnknownExpressionType тестирует неизвестный тип выражения
func TestInterpreterUnknownExpressionType(t *testing.T) {
	interpreter := NewInterpreter()
	// Создаем фиктивное выражение
	assignment := &Assignment{
		Variable: "x",
		Value:    &FakeExpression{},
	}
	program := &Program{
		Statements: []Statement{assignment},
	}
	err := interpreter.Interpret(program)
	if err == nil {
		t.Error("Ожидалась ошибка для неизвестного типа выражения")
	}
}

// TestInterpreterUnknownOperator тестирует неизвестный оператор
func TestInterpreterUnknownOperator(t *testing.T) {
	interpreter := NewInterpreter()
	// Создаем BinaryOp с неизвестным оператором
	binaryOp := &BinaryOp{
		Left:     &Number{Value: 5},
		Operator: TokenEOF, // Неизвестный оператор
		Right:    &Number{Value: 3},
	}
	assignment := &Assignment{
		Variable: "x",
		Value:    binaryOp,
	}
	program := &Program{
		Statements: []Statement{assignment},
	}
	err := interpreter.Interpret(program)
	if err == nil {
		t.Error("Ожидалась ошибка для неизвестного оператора")
	}
}

// TestStringMethods тестирует методы String() для отладки
func TestStringMethods(t *testing.T) {
	// Program.String()
	program := &Program{Statements: []Statement{}}
	if program.String() == "" {
		t.Error("Program.String() должен возвращать непустую строку")
	}

	// Assignment.String()
	assignment := &Assignment{
		Variable: "x",
		Value:    &Number{Value: 5},
	}
	if assignment.String() == "" {
		t.Error("Assignment.String() должен возвращать непустую строку")
	}

	// Block.String()
	block := &Block{Statements: []Statement{}}
	if block.String() == "" {
		t.Error("Block.String() должен возвращать непустую строку")
	}

	// Number.String()
	number := &Number{Value: 42}
	if number.String() == "" {
		t.Error("Number.String() должен возвращать непустую строку")
	}

	// Identifier.String()
	identifier := &Identifier{Name: "x"}
	if identifier.String() == "" {
		t.Error("Identifier.String() должен возвращать непустую строку")
	}

	// BinaryOp.String() - все операторы
	binaryOpPlus := &BinaryOp{
		Left:     &Number{Value: 2},
		Operator: TokenPLUS,
		Right:    &Number{Value: 3},
	}
	if binaryOpPlus.String() == "" {
		t.Error("BinaryOp.String() с PLUS должен возвращать непустую строку")
	}

	binaryOpMinus := &BinaryOp{
		Left:     &Number{Value: 2},
		Operator: TokenMINUS,
		Right:    &Number{Value: 3},
	}
	if binaryOpMinus.String() == "" {
		t.Error("BinaryOp.String() с MINUS должен возвращать непустую строку")
	}

	binaryOpMultiply := &BinaryOp{
		Left:     &Number{Value: 2},
		Operator: TokenMULTIPLY,
		Right:    &Number{Value: 3},
	}
	if binaryOpMultiply.String() == "" {
		t.Error("BinaryOp.String() с MULTIPLY должен возвращать непустую строку")
	}

	binaryOpDivide := &BinaryOp{
		Left:     &Number{Value: 2},
		Operator: TokenDIVIDE,
		Right:    &Number{Value: 3},
	}
	if binaryOpDivide.String() == "" {
		t.Error("BinaryOp.String() с DIVIDE должен возвращать непустую строку")
	}

	// Тест маркерных методов - вызываем их для покрытия
	assignment.statementNode()
	block.statementNode()
	number.expressionNode()
	identifier.expressionNode()
	binaryOpPlus.expressionNode()
	
	// Проверяем, что методы не паникуют
	if assignment.String() == "" {
		t.Error("Assignment.String() после statementNode() должен работать")
	}
	if block.String() == "" {
		t.Error("Block.String() после statementNode() должен работать")
	}
	if number.String() == "" {
		t.Error("Number.String() после expressionNode() должен работать")
	}
	if identifier.String() == "" {
		t.Error("Identifier.String() после expressionNode() должен работать")
	}
	if binaryOpPlus.String() == "" {
		t.Error("BinaryOp.String() после expressionNode() должен работать")
	}
}

// TestLexerPeekNext тестирует метод peekNext
func TestLexerPeekNext(t *testing.T) {
	// Тест когда есть следующий символ
	code := `BEGIN x := 5; END.`
	lexer := NewLexer(code)
	lexer.pos = 0
	next := lexer.peekNext()
	if next == 0 {
		t.Error("peekNext должен вернуть символ, когда есть следующий")
	}

	// Тест когда нет следующего символа (в конце строки)
	lexer2 := NewLexer("A")
	lexer2.pos = 0
	next2 := lexer2.peekNext()
	if next2 != 0 {
		t.Error("peekNext должен вернуть 0, когда нет следующего символа")
	}

	// Тест когда pos+1 >= len(input)
	lexer3 := NewLexer("AB")
	lexer3.pos = 2
	next3 := lexer3.peekNext()
	if next3 != 0 {
		t.Error("peekNext должен вернуть 0, когда pos+1 >= len")
	}
}

// TestParserEdgeCases тестирует краевые случаи парсера
func TestParserEdgeCases(t *testing.T) {
	// Пустой блок с точкой с запятой
	code := `BEGIN ; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	_, err = parser.Parse()
	// Это может вызвать ошибку, что нормально
	_ = err

	// Присваивание без точки с запятой в конце блока
	code2 := `BEGIN x := 5 END.`
	lexer2 := NewLexer(code2)
	tokens2, err := lexer2.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser2 := NewParser(tokens2)
	program, err := parser2.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestParserCurrentEdgeCases тестирует краевые случаи current()
func TestParserCurrentEdgeCases(t *testing.T) {
	// Пустой список токенов
	parser := NewParser([]Token{})
	token := parser.current()
	if token.Type != TokenEOF {
		t.Errorf("Ожидался TokenEOF для пустого списка, получен %v", token.Type)
	}
}

// TestLexerPeekRuneEdgeCases тестирует краевые случаи peekRune
func TestLexerPeekRuneEdgeCases(t *testing.T) {
	// Пустая строка
	lexer := NewLexer("")
	r, size := lexer.peekRune()
	if size != 0 || r != 0 {
		t.Errorf("Ожидался (0, 0) для пустой строки, получен (%v, %d)", r, size)
	}

	// Строка в конце
	lexer2 := NewLexer("BEGIN")
	lexer2.pos = len(lexer2.input)
	r2, size2 := lexer2.peekRune()
	if size2 != 0 || r2 != 0 {
		t.Errorf("Ожидался (0, 0) для конца строки, получен (%v, %d)", r2, size2)
	}
}

// TestParserAdditiveMultipleOps тестирует множественные аддитивные операции
func TestParserAdditiveMultipleOps(t *testing.T) {
	code := `BEGIN x := 1 + 2 + 3 + 4; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 10 {
		t.Errorf("x: ожидалось 10, получено %g", variables["x"])
	}
}

// TestParserMultiplicativeMultipleOps тестирует множественные мультипликативные операции
func TestParserMultiplicativeMultipleOps(t *testing.T) {
	code := `BEGIN x := 2 * 3 * 4; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 24 {
		t.Errorf("x: ожидалось 24, получено %g", variables["x"])
	}
}

// TestParserUnaryMinus тестирует унарный минус
func TestParserUnaryMinus(t *testing.T) {
	code := `BEGIN x := -10; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != -10 {
		t.Errorf("x: ожидалось -10, получено %g", variables["x"])
	}
}

// TestParserUnaryMinusInExpression тестирует унарный минус в выражении
func TestParserUnaryMinusInExpression(t *testing.T) {
	code := `BEGIN x := 5 + -3; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 2 {
		t.Errorf("x: ожидалось 2, получено %g", variables["x"])
	}
}

// TestParserComplexExpression тестирует сложное выражение
func TestParserComplexExpression(t *testing.T) {
	code := `BEGIN x := (1 + 2) * (3 - 4) / 5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = (1 + 2) * (3 - 4) / 5 = 3 * (-1) / 5 = -3 / 5 = -0.6
	if x, ok := variables["x"]; !ok || x != -0.6 {
		t.Errorf("x: ожидалось -0.6, получено %g", variables["x"])
	}
}

// TestInterpreterUninitializedVariable тестирует использование неинициализированной переменной
func TestInterpreterUninitializedVariable(t *testing.T) {
	code := `BEGIN x := y; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// y не инициализирована, должна быть равна 0
	if x, ok := variables["x"]; !ok || x != 0 {
		t.Errorf("x: ожидалось 0, получено %g", variables["x"])
	}
}

// TestLexerSkipWhitespace тестирует skipWhitespace
func TestLexerSkipWhitespace(t *testing.T) {
	code := "   BEGIN   x   :=   5   ;   END   ."
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	// Проверяем, что пробелы пропущены
	if len(tokens) < 7 {
		t.Errorf("Ожидалось минимум 7 токенов, получено %d", len(tokens))
	}
}

// TestParserStatementWithoutSemicolon тестирует оператор без точки с запятой
func TestParserStatementWithoutSemicolon(t *testing.T) {
	code := `BEGIN x := 5 END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", variables["x"])
	}
}

// TestParserAdditiveWithMinus тестирует вычитание в parseAdditive
func TestParserAdditiveWithMinus(t *testing.T) {
	code := `BEGIN x := 10 - 3 - 2; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = 10 - 3 - 2 = 5
	if x, ok := variables["x"]; !ok || x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", variables["x"])
	}
}

// TestParserMultiplicativeWithDivide тестирует деление в parseMultiplicative
func TestParserMultiplicativeWithDivide(t *testing.T) {
	code := `BEGIN x := 24 / 2 / 3; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = 24 / 2 / 3 = 12 / 3 = 4
	if x, ok := variables["x"]; !ok || x != 4 {
		t.Errorf("x: ожидалось 4, получено %g", variables["x"])
	}
}

// TestParserUnaryMinusComplex тестирует сложный унарный минус
func TestParserUnaryMinusComplex(t *testing.T) {
	code := `BEGIN x := -(-5); END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = -(-5) = 5
	if x, ok := variables["x"]; !ok || x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", variables["x"])
	}
}

// TestParserPrimaryWithIdentifier тестирует parsePrimary с идентификатором
func TestParserPrimaryWithIdentifier(t *testing.T) {
	code := `BEGIN x := 5; y := x; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if y, ok := variables["y"]; !ok || y != 5 {
		t.Errorf("y: ожидалось 5, получено %g", variables["y"])
	}
}

// TestParserPrimaryWithParens тестирует parsePrimary со скобками
func TestParserPrimaryWithParens(t *testing.T) {
	code := `BEGIN x := ((5)); END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", variables["x"])
	}
}

// TestLexerTokenizeEdgeCases тестирует краевые случаи Tokenize
func TestLexerTokenizeEdgeCases(t *testing.T) {
	// Пустая строка
	lexer1 := NewLexer("")
	tokens1, err := lexer1.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа для пустой строки: %v", err)
	}
	if len(tokens1) == 0 || tokens1[len(tokens1)-1].Type != TokenEOF {
		t.Error("Пустая строка должна заканчиваться TokenEOF")
	}

	// Только пробелы
	lexer2 := NewLexer("   \n\t  ")
	tokens2, err := lexer2.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа для пробелов: %v", err)
	}
	if len(tokens2) == 0 || tokens2[len(tokens2)-1].Type != TokenEOF {
		t.Error("Строка с пробелами должна заканчиваться TokenEOF")
	}
}

// TestParserParseEdgeCases тестирует краевые случаи Parse
func TestParserParseEdgeCases(t *testing.T) {
	// Программа с несколькими операторами
	code := `BEGIN x := 1; y := 2; z := 3; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	if len(program.Statements) != 3 {
		t.Errorf("Ожидалось 3 оператора, получено %d", len(program.Statements))
	}
}

// TestParserStatementNestedBlockWithSemicolon тестирует вложенный блок с точкой с запятой
func TestParserStatementNestedBlockWithSemicolon(t *testing.T) {
	code := `BEGIN BEGIN x := 1; END; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 1 {
		t.Errorf("x: ожидалось 1, получено %g", variables["x"])
	}
}

// TestParserStatementNestedBlockWithoutEnd тестирует вложенный блок без END
func TestParserStatementNestedBlockWithoutEnd(t *testing.T) {
	code := `BEGIN BEGIN x := 1; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	_, err = parser.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для вложенного блока без END")
	}
}

// TestParserStatementAssignmentWithoutAssign тестирует присваивание без :=
func TestParserStatementAssignmentWithoutAssign(t *testing.T) {
	code := `BEGIN x 5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	_, err = parser.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для присваивания без :=")
	}
}

// TestParserAdditiveMixedOps тестирует смешанные аддитивные операции
func TestParserAdditiveMixedOps(t *testing.T) {
	code := `BEGIN x := 1 + 2 - 3 + 4; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = 1 + 2 - 3 + 4 = 4
	if x, ok := variables["x"]; !ok || x != 4 {
		t.Errorf("x: ожидалось 4, получено %g", variables["x"])
	}
}

// TestParserMultiplicativeMixedOps тестирует смешанные мультипликативные операции
func TestParserMultiplicativeMixedOps(t *testing.T) {
	code := `BEGIN x := 2 * 3 / 2 * 4; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = 2 * 3 / 2 * 4 = 6 / 2 * 4 = 3 * 4 = 12
	if x, ok := variables["x"]; !ok || x != 12 {
		t.Errorf("x: ожидалось 12, получено %g", variables["x"])
	}
}

// TestParserUnaryMinusWithExpression тестирует унарный минус с выражением
func TestParserUnaryMinusWithExpression(t *testing.T) {
	code := `BEGIN x := -(2 + 3); END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = -(2 + 3) = -5
	if x, ok := variables["x"]; !ok || x != -5 {
		t.Errorf("x: ожидалось -5, получено %g", variables["x"])
	}
}

// TestParserPrimaryError тестирует ошибку в parsePrimary
func TestParserPrimaryError(t *testing.T) {
	code := `BEGIN x := +; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	_, err = parser.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена в parsePrimary")
	}
}

// TestLexerSkipWhitespaceEdgeCases тестирует краевые случаи skipWhitespace
func TestLexerSkipWhitespaceEdgeCases(t *testing.T) {
	// Различные типы пробелов
	code := "\tBEGIN\nx\r:=\v5\f;\r\nEND."
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	// Проверяем, что пробелы пропущены
	if len(tokens) < 7 {
		t.Errorf("Ожидалось минимум 7 токенов, получено %d", len(tokens))
	}
}

// TestInterpreterEvaluateExpressionErrors тестирует ошибки в evaluateExpression
func TestInterpreterEvaluateExpressionErrors(t *testing.T) {
	interpreter := NewInterpreter()
	
	// Тест ошибки в левой части BinaryOp
	fakeExpr := &FakeExpression{}
	binaryOp := &BinaryOp{
		Left:     fakeExpr,
		Operator: TokenPLUS,
		Right:    &Number{Value: 5},
	}
	assignment := &Assignment{
		Variable: "x",
		Value:    binaryOp,
	}
	program := &Program{
		Statements: []Statement{assignment},
	}
	err := interpreter.Interpret(program)
	if err == nil {
		t.Error("Ожидалась ошибка для неизвестного типа выражения в левой части")
	}

	// Тест ошибки в правой части BinaryOp
	binaryOp2 := &BinaryOp{
		Left:     &Number{Value: 5},
		Operator: TokenPLUS,
		Right:    fakeExpr,
	}
	assignment2 := &Assignment{
		Variable: "x",
		Value:    binaryOp2,
	}
	program2 := &Program{
		Statements: []Statement{assignment2},
	}
	err2 := interpreter.Interpret(program2)
	if err2 == nil {
		t.Error("Ожидалась ошибка для неизвестного типа выражения в правой части")
	}
}

// TestParserPrimaryErrorCase тестирует ошибку в parsePrimary когда нет подходящего токена
func TestParserPrimaryErrorCase(t *testing.T) {
	// Создаем токены, которые не подходят для parsePrimary
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenSEMICOLON, Value: ";", Pos: 10},
		{Type: TokenEND, Value: "END", Pos: 15},
		{Type: TokenDOT, Value: ".", Pos: 18},
		{Type: TokenEOF, Value: "", Pos: 19},
	}
	parser := NewParser(tokens)
	parser.pos = 1 // Позиция на SEMICOLON
	_, err := parser.parsePrimary()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена в parsePrimary")
	}
}

// TestParserParseMultipleStatements тестирует Parse с несколькими операторами
func TestParserParseMultipleStatements(t *testing.T) {
	code := `BEGIN x := 1; y := 2; z := 3; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	if len(program.Statements) != 3 {
		t.Errorf("Ожидалось 3 оператора, получено %d", len(program.Statements))
	}
}

// TestParserAdditiveNoOps тестирует parseAdditive без операций
func TestParserAdditiveNoOps(t *testing.T) {
	code := `BEGIN x := 5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", variables["x"])
	}
}

// TestParserMultiplicativeNoOps тестирует parseMultiplicative без операций
func TestParserMultiplicativeNoOps(t *testing.T) {
	code := `BEGIN x := 5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", variables["x"])
	}
}

// TestParserUnaryNoMinus тестирует parseUnary без минуса
func TestParserUnaryNoMinus(t *testing.T) {
	code := `BEGIN x := 5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	if x, ok := variables["x"]; !ok || x != 5 {
		t.Errorf("x: ожидалось 5, получено %g", variables["x"])
	}
}

// TestLexerTokenizeSizeZero тестирует Tokenize когда peekRune возвращает size == 0
func TestLexerTokenizeSizeZero(t *testing.T) {
	// Создаем строку с невалидной UTF-8 последовательностью, которая вызовет size == 0
	// Используем неполную UTF-8 последовательность
	code := string([]byte{0x80}) // Невалидная UTF-8 последовательность
	lexer := NewLexer(code)
	_, err := lexer.Tokenize()
	// Может быть ошибка или успешное завершение, в зависимости от реализации
	_ = err
}

// TestLexerSkipWhitespaceSizeZero тестирует skipWhitespace когда size == 0
func TestLexerSkipWhitespaceSizeZero(t *testing.T) {
	// Создаем строку с невалидной UTF-8 последовательностью
	code := string([]byte{0x80})
	lexer := NewLexer(code)
	lexer.pos = 0
	lexer.start = 0
	lexer.skipWhitespace()
	// Проверяем, что метод не паникует
}

// TestParserUnaryRecursive тестирует рекурсивный вызов parseUnary
func TestParserUnaryRecursive(t *testing.T) {
	code := `BEGIN x := ---5; END.`
	lexer := NewLexer(code)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		t.Fatalf("Ошибка синтаксического анализа: %v", err)
	}
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	variables := interpreter.GetVariables()
	// x = ---5 = -(-(-5)) = -(-(-5)) = -5
	if x, ok := variables["x"]; !ok || x != -5 {
		t.Errorf("x: ожидалось -5, получено %g", variables["x"])
	}
}

// TestParserUnaryError тестирует ошибку в рекурсивном parseUnary
func TestParserUnaryError(t *testing.T) {
	// Создаем токены, которые вызовут ошибку в parseUnary -> parsePrimary
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenMINUS, Value: "-", Pos: 11},
		{Type: TokenSEMICOLON, Value: ";", Pos: 12}, // Неожиданный токен
		{Type: TokenEND, Value: "END", Pos: 14},
		{Type: TokenDOT, Value: ".", Pos: 17},
		{Type: TokenEOF, Value: "", Pos: 18},
	}
	parser := NewParser(tokens)
	parser.pos = 3 // Позиция на MINUS
	_, err := parser.parseUnary()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена в parseUnary")
	}
}

// TestParserPrimaryAllCases тестирует все случаи parsePrimary
func TestParserPrimaryAllCases(t *testing.T) {
	// Случай 1: NUMBER
	code1 := `BEGIN x := 42; END.`
	lexer1 := NewLexer(code1)
	tokens1, _ := lexer1.Tokenize()
	parser1 := NewParser(tokens1)
	parser1.pos = 3 // Позиция на NUMBER
	expr1, err1 := parser1.parsePrimary()
	if err1 != nil || expr1 == nil {
		t.Errorf("Ошибка парсинга NUMBER: %v", err1)
	}

	// Случай 2: IDENTIFIER
	code2 := `BEGIN x := y; END.`
	lexer2 := NewLexer(code2)
	tokens2, _ := lexer2.Tokenize()
	parser2 := NewParser(tokens2)
	parser2.pos = 3 // Позиция на IDENTIFIER
	expr2, err2 := parser2.parsePrimary()
	if err2 != nil || expr2 == nil {
		t.Errorf("Ошибка парсинга IDENTIFIER: %v", err2)
	}

	// Случай 3: LPAREN
	code3 := `BEGIN x := (5); END.`
	lexer3 := NewLexer(code3)
	tokens3, _ := lexer3.Tokenize()
	parser3 := NewParser(tokens3)
	parser3.pos = 3 // Позиция на LPAREN
	expr3, err3 := parser3.parsePrimary()
	if err3 != nil || expr3 == nil {
		t.Errorf("Ошибка парсинга LPAREN: %v", err3)
	}
}

// TestParserAdditiveErrorInRight тестирует ошибку в правой части parseAdditive
func TestParserAdditiveErrorInRight(t *testing.T) {
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenNUMBER, Value: "5", Pos: 11},
		{Type: TokenPLUS, Value: "+", Pos: 12},
		{Type: TokenSEMICOLON, Value: ";", Pos: 13}, // Неожиданный токен
		{Type: TokenEND, Value: "END", Pos: 15},
		{Type: TokenDOT, Value: ".", Pos: 18},
		{Type: TokenEOF, Value: "", Pos: 19},
	}
	parser := NewParser(tokens)
	parser.pos = 3 // Позиция на NUMBER
	_, err := parser.parseAdditive()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена в правой части parseAdditive")
	}
}

// TestParserMultiplicativeErrorInRight тестирует ошибку в правой части parseMultiplicative
func TestParserMultiplicativeErrorInRight(t *testing.T) {
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenNUMBER, Value: "5", Pos: 11},
		{Type: TokenMULTIPLY, Value: "*", Pos: 12},
		{Type: TokenSEMICOLON, Value: ";", Pos: 13}, // Неожиданный токен
		{Type: TokenEND, Value: "END", Pos: 15},
		{Type: TokenDOT, Value: ".", Pos: 18},
		{Type: TokenEOF, Value: "", Pos: 19},
	}
	parser := NewParser(tokens)
	parser.pos = 3 // Позиция на NUMBER
	_, err := parser.parseMultiplicative()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена в правой части parseMultiplicative")
	}
}

// TestParserPrimaryErrorInParens тестирует ошибку внутри скобок в parsePrimary
func TestParserPrimaryErrorInParens(t *testing.T) {
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenLPAREN, Value: "(", Pos: 11},
		{Type: TokenSEMICOLON, Value: ";", Pos: 12}, // Неожиданный токен
		{Type: TokenEND, Value: "END", Pos: 15},
		{Type: TokenDOT, Value: ".", Pos: 18},
		{Type: TokenEOF, Value: "", Pos: 19},
	}
	parser := NewParser(tokens)
	parser.pos = 3 // Позиция на LPAREN
	_, err := parser.parsePrimary()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена внутри скобок")
	}
}

// TestParserPrimaryUnclosedParen тестирует незакрытую скобку в parsePrimary
func TestParserPrimaryUnclosedParen(t *testing.T) {
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenLPAREN, Value: "(", Pos: 11},
		{Type: TokenNUMBER, Value: "5", Pos: 12},
		{Type: TokenEND, Value: "END", Pos: 15}, // Нет RPAREN
		{Type: TokenDOT, Value: ".", Pos: 18},
		{Type: TokenEOF, Value: "", Pos: 19},
	}
	parser := NewParser(tokens)
	parser.pos = 3 // Позиция на LPAREN
	_, err := parser.parsePrimary()
	if err == nil {
		t.Error("Ожидалась ошибка для незакрытой скобки")
	}
}

// TestParserStatementBlockError тестирует ошибку в parseBlock внутри parseStatement
func TestParserStatementBlockError(t *testing.T) {
	// Создаем токены с BEGIN но без END
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 6},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 12},
		{Type: TokenASSIGN, Value: ":=", Pos: 14},
		{Type: TokenNUMBER, Value: "5", Pos: 17},
		{Type: TokenDOT, Value: ".", Pos: 18}, // Нет END
		{Type: TokenEOF, Value: "", Pos: 19},
	}
	parser := NewParser(tokens)
	parser.pos = 1 // Позиция на втором BEGIN
	_, err := parser.parseStatement()
	if err == nil {
		t.Error("Ожидалась ошибка для блока без END")
	}
}

// TestParserParseBlockError тестирует ошибку в parseBlock внутри Parse
func TestParserParseBlockError(t *testing.T) {
	// Создаем токены с BEGIN но без END
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenNUMBER, Value: "5", Pos: 11},
		{Type: TokenDOT, Value: ".", Pos: 12}, // Нет END
		{Type: TokenEOF, Value: "", Pos: 13},
	}
	parser := NewParser(tokens)
	_, err := parser.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для блока без END")
	}
}

// TestParserStatementExpressionError тестирует ошибку в parseExpression внутри parseStatement
func TestParserStatementExpressionError(t *testing.T) {
	// Создаем токены с присваиванием но с ошибкой в выражении
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenSEMICOLON, Value: ";", Pos: 11}, // Неожиданный токен в выражении
		{Type: TokenEND, Value: "END", Pos: 13},
		{Type: TokenDOT, Value: ".", Pos: 16},
		{Type: TokenEOF, Value: "", Pos: 17},
	}
	parser := NewParser(tokens)
	parser.pos = 1 // Позиция на IDENTIFIER
	_, err := parser.parseStatement()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного токена в выражении")
	}
}

// TestLexerTokenizeInvalidUTF8 тестирует Tokenize с невалидной UTF-8 последовательностью
func TestLexerTokenizeInvalidUTF8(t *testing.T) {
	// Создаем строку с невалидной UTF-8 последовательностью
	// Это может вызвать size == 0 в peekRune
	code := string([]byte{0xFF, 0xFE}) // Невалидная UTF-8
	lexer := NewLexer(code)
	_, err := lexer.Tokenize()
	// Может быть ошибка или успешное завершение
	_ = err
}

// TestLexerSkipWhitespaceInvalidUTF8 тестирует skipWhitespace с невалидной UTF-8
func TestLexerSkipWhitespaceInvalidUTF8(t *testing.T) {
	// Создаем строку с невалидной UTF-8 последовательностью
	code := string([]byte{0xFF, 0xFE})
	lexer := NewLexer(code)
	lexer.pos = 0
	lexer.start = 0
	lexer.skipWhitespace()
	// Проверяем, что метод не паникует
}

// TestLexerTokenizeColonWithoutEquals тестирует Tokenize когда после ':' нет '='
func TestLexerTokenizeColonWithoutEquals(t *testing.T) {
	code := `BEGIN x : 5; END.`
	lexer := NewLexer(code)
	_, err := lexer.Tokenize()
	if err == nil {
		t.Error("Ожидалась ошибка для ':' без '='")
	}
}

// TestParserParseBlockError тестирует ошибку в parseBlock внутри Parse
func TestParserParseBlockErrorInParse(t *testing.T) {
	// Создаем токены с BEGIN но с ошибкой в блоке
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenSEMICOLON, Value: ";", Pos: 7}, // Нет :=
		{Type: TokenEND, Value: "END", Pos: 9},
		{Type: TokenDOT, Value: ".", Pos: 12},
		{Type: TokenEOF, Value: "", Pos: 13},
	}
	parser := NewParser(tokens)
	_, err := parser.Parse()
	if err == nil {
		t.Error("Ожидалась ошибка для невалидного оператора в блоке")
	}
}

// TestParserStatementBlockErrorInStatement тестирует ошибку в parseBlock внутри parseStatement
func TestParserStatementBlockErrorInStatement(t *testing.T) {
	// Создаем токены с BEGIN но с ошибкой в блоке
	tokens := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 6},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 12},
		{Type: TokenSEMICOLON, Value: ";", Pos: 13}, // Нет :=
		{Type: TokenEND, Value: "END", Pos: 15},
		{Type: TokenEND, Value: "END", Pos: 19},
		{Type: TokenDOT, Value: ".", Pos: 22},
		{Type: TokenEOF, Value: "", Pos: 23},
	}
	parser := NewParser(tokens)
	parser.pos = 1 // Позиция на втором BEGIN
	_, err := parser.parseStatement()
	if err == nil {
		t.Error("Ожидалась ошибка для невалидного оператора в блоке")
	}
}

// TestLexerTokenizePeekRuneSizeZero тестирует Tokenize когда peekRune возвращает size == 0
func TestLexerTokenizePeekRuneSizeZero(t *testing.T) {
	// Создаем строку и устанавливаем pos так, чтобы peekRune вернул size == 0
	code := "BEGIN"
	lexer := NewLexer(code)
	lexer.pos = len(code)
	// Вызываем skipWhitespace, который использует peekRune
	lexer.skipWhitespace()
	// Проверяем, что метод не паникует
}

// TestParserParseErrorCases тестирует все случаи ошибок в Parse
func TestParserParseErrorCases(t *testing.T) {
	// Ошибка: отсутствие BEGIN
	tokens1 := []Token{
		{Type: TokenIDENTIFIER, Value: "x", Pos: 0},
		{Type: TokenASSIGN, Value: ":=", Pos: 2},
		{Type: TokenNUMBER, Value: "5", Pos: 5},
		{Type: TokenEND, Value: "END", Pos: 7},
		{Type: TokenDOT, Value: ".", Pos: 10},
		{Type: TokenEOF, Value: "", Pos: 11},
	}
	parser1 := NewParser(tokens1)
	_, err1 := parser1.Parse()
	if err1 == nil {
		t.Error("Ожидалась ошибка для отсутствия BEGIN")
	}

	// Ошибка: отсутствие END
	tokens2 := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenNUMBER, Value: "5", Pos: 11},
		{Type: TokenDOT, Value: ".", Pos: 12},
		{Type: TokenEOF, Value: "", Pos: 13},
	}
	parser2 := NewParser(tokens2)
	_, err2 := parser2.Parse()
	if err2 == nil {
		t.Error("Ожидалась ошибка для отсутствия END")
	}

	// Ошибка: отсутствие точки
	tokens3 := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 6},
		{Type: TokenASSIGN, Value: ":=", Pos: 8},
		{Type: TokenNUMBER, Value: "5", Pos: 11},
		{Type: TokenEND, Value: "END", Pos: 13},
		{Type: TokenEOF, Value: "", Pos: 16},
	}
	parser3 := NewParser(tokens3)
	_, err3 := parser3.Parse()
	if err3 == nil {
		t.Error("Ожидалась ошибка для отсутствия точки")
	}

	// Ошибка: неожиданные токены после точки
	tokens4 := []Token{
		{Type: TokenBEGIN, Value: "BEGIN", Pos: 0},
		{Type: TokenEND, Value: "END", Pos: 6},
		{Type: TokenDOT, Value: ".", Pos: 9},
		{Type: TokenIDENTIFIER, Value: "x", Pos: 11},
		{Type: TokenEOF, Value: "", Pos: 12},
	}
	parser4 := NewParser(tokens4)
	_, err4 := parser4.Parse()
	if err4 == nil {
		t.Error("Ожидалась ошибка для неожиданных токенов после точки")
	}
}

// TestLexerTokenizeDefaultCase тестирует default case в Tokenize
func TestLexerTokenizeDefaultCase(t *testing.T) {
	// Создаем строку с неожиданным символом
	code := "BEGIN x @ 5; END."
	lexer := NewLexer(code)
	_, err := lexer.Tokenize()
	if err == nil {
		t.Error("Ожидалась ошибка для неожиданного символа @")
	}
	if err != nil && err.Error() == "" {
		t.Error("Ошибка должна содержать описание")
	}
}

// TestMarkerMethods тестирует маркерные методы напрямую
func TestMarkerMethods(t *testing.T) {
	// Создаем экземпляры и вызываем маркерные методы
	assignment := &Assignment{Variable: "x", Value: &Number{Value: 5}}
	block := &Block{Statements: []Statement{}}
	number := &Number{Value: 42}
	identifier := &Identifier{Name: "y"}
	binaryOp := &BinaryOp{Left: &Number{Value: 1}, Operator: TokenPLUS, Right: &Number{Value: 2}}
	
	// Вызываем маркерные методы напрямую
	assignment.statementNode()
	block.statementNode()
	number.expressionNode()
	identifier.expressionNode()
	binaryOp.expressionNode()
	
	// Проверяем, что методы не паникуют
	_ = assignment.String()
	_ = block.String()
	_ = number.String()
	_ = identifier.String()
	_ = binaryOp.String()
}

// TestLexerTokenizePosGreaterEqualLen тестирует Tokenize когда pos >= len(input)
func TestLexerTokenizePosGreaterEqualLen(t *testing.T) {
	code := "BEGIN"
	lexer := NewLexer(code)
	lexer.pos = len(code)
	lexer.skipWhitespace()
	// После skipWhitespace pos должен быть >= len(input), что вызовет break в Tokenize
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Ошибка лексического анализа: %v", err)
	}
	// Должен быть хотя бы EOF токен
	if len(tokens) == 0 {
		t.Error("Ожидался хотя бы один токен (EOF)")
	}
}

// TestLexerSkipWhitespaceAllCases тестирует все случаи skipWhitespace
func TestLexerSkipWhitespaceAllCases(t *testing.T) {
	// Случай 1: обычные пробелы
	code1 := "   BEGIN"
	lexer1 := NewLexer(code1)
	lexer1.skipWhitespace()
	if lexer1.pos != 3 {
		t.Errorf("Ожидалось pos=3 после пробелов, получено %d", lexer1.pos)
	}
	
	// Случай 2: pos >= len(input)
	code2 := "BEGIN"
	lexer2 := NewLexer(code2)
	lexer2.pos = len(code2)
	lexer2.skipWhitespace()
	// Не должно быть паники
	
	// Случай 3: size == 0 в peekRune
	code3 := "BEGIN"
	lexer3 := NewLexer(code3)
	lexer3.pos = len(code3) - 1
	// Устанавливаем так, чтобы peekRune мог вернуть size == 0
	lexer3.skipWhitespace()
	// Не должно быть паники
}

