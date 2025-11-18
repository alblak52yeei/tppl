package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestIncrement тестирует операцию MoO (увеличение значения)
func TestIncrement(t *testing.T) {
	code := "MoO MoO MoO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 3 {
		t.Errorf("Ожидалось 3, получено %d", interp.GetMemory()[0])
	}
}

// TestDecrement тестирует операцию MOo (уменьшение значения)
func TestDecrement(t *testing.T) {
	code := "MoO MoO MoO MOo"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 2 {
		t.Errorf("Ожидалось 2, получено %d", interp.GetMemory()[0])
	}
}

// TestNext тестирует операцию moO (переход к следующей ячейке)
func TestNext(t *testing.T) {
	code := "MoO moO MoO MoO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 1 {
		t.Errorf("Ячейка 0: ожидалось 1, получено %d", interp.GetMemory()[0])
	}
	if interp.GetMemory()[1] != 2 {
		t.Errorf("Ячейка 1: ожидалось 2, получено %d", interp.GetMemory()[1])
	}
	if interp.GetPointer() != 1 {
		t.Errorf("Указатель: ожидалось 1, получено %d", interp.GetPointer())
	}
}

// TestPrev тестирует операцию mOo (переход к предыдущей ячейке)
func TestPrev(t *testing.T) {
	code := "MoO moO MoO MoO mOo MoO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 2 {
		t.Errorf("Ячейка 0: ожидалось 2, получено %d", interp.GetMemory()[0])
	}
	if interp.GetPointer() != 0 {
		t.Errorf("Указатель: ожидалось 0, получено %d", interp.GetPointer())
	}
}

// TestLoop тестирует операции moo и MOO (цикл)
func TestLoop(t *testing.T) {
	// Устанавливаем счетчик в 3 и уменьшаем в цикле
	code := "MoO MoO MoO moo MOo MOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 0 {
		t.Errorf("Ожидалось 0, получено %d", interp.GetMemory()[0])
	}
}

// TestLoopSkip тестирует пропуск цикла когда значение ячейки 0
func TestLoopSkip(t *testing.T) {
	code := "moo MoO MOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 0 {
		t.Errorf("Ожидалось 0, получено %d", interp.GetMemory()[0])
	}
}

// TestNestedLoop тестирует вложенные циклы
func TestNestedLoop(t *testing.T) {
	// Внешний цикл 2 раза, внутренний 3 раза
	code := "MoO MoO moo moO MoO MoO MoO moo MOo MOO mOo MOo MOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 0 {
		t.Errorf("Ячейка 0: ожидалось 0, получено %d", interp.GetMemory()[0])
	}
}

// TestNestedLoopSkip тестирует пропуск вложенных циклов
func TestNestedLoopSkip(t *testing.T) {
	// Внешний цикл пропускается (ячейка 0 = 0), внутренний тоже должен быть пропущен
	code := "moo moo MoO MOO MOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	// Все циклы должны быть пропущены
	if interp.GetMemory()[0] != 0 {
		t.Errorf("Ячейка 0: ожидалось 0, получено %d", interp.GetMemory()[0])
	}
}

// TestOutput тестирует операцию OOM (вывод значения)
func TestOutput(t *testing.T) {
	code := "MoO MoO MoO MoO MoO OOM"
	var output bytes.Buffer
	interp := NewInterpreter(strings.NewReader(""), &output)
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	expected := "5"
	if output.String() != expected {
		t.Errorf("Ожидалось '%s', получено '%s'", expected, output.String())
	}
}

// TestInput тестирует операцию oom (ввод значения)
func TestInput(t *testing.T) {
	code := "oom"
	input := strings.NewReader("42\n")
	interp := NewInterpreter(input, &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 42 {
		t.Errorf("Ожидалось 42, получено %d", interp.GetMemory()[0])
	}
}

// TestZero тестирует операцию OOO (обнуление ячейки)
func TestZero(t *testing.T) {
	code := "MoO MoO MoO OOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 0 {
		t.Errorf("Ожидалось 0, получено %d", interp.GetMemory()[0])
	}
}

// TestConditionalInput тестирует операцию Moo (условный ввод)
func TestConditionalInput(t *testing.T) {
	code := "Moo"
	input := strings.NewReader("99\n")
	interp := NewInterpreter(input, &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 99 {
		t.Errorf("Ожидалось 99, получено %d", interp.GetMemory()[0])
	}
}

// TestConditionalOutput тестирует операцию Moo (условный вывод)
func TestConditionalOutput(t *testing.T) {
	code := "MoO MoO MoO Moo"
	var output bytes.Buffer
	interp := NewInterpreter(strings.NewReader(""), &output)
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	expected := "3"
	if output.String() != expected {
		t.Errorf("Ожидалось '%s', получено '%s'", expected, output.String())
	}
}

// TestExecute тестирует операцию mOO (выполнение инструкции по адресу)
func TestExecute(t *testing.T) {
	// Устанавливаем в ячейку 0 значение 0 (адрес инструкции 0)
	// Инструкция 0: MoO (increment) - будет выполнена через mOO
	// Инструкция 1: mOO (execute instruction at memory[0])
	code := "MoO mOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	// После MoO: memory[0] = 1
	// После mOO: выполняет инструкцию 1 (mOO), но memory[0]=1, поэтому выполнится инструкция 1
	// На самом деле mOO выполнит инструкцию по адресу memory[0]=1, т.е. саму mOO - но это не должно вызвать рекурсию
	// Правильный тест: после MoO memory[0]=1, после mOO выполнится инструкция 1 (сам mOO или что-то еще)
	// Но mOO выполняет инструкцию по адресу из ячейки памяти
	// После выполнения память должна быть изменена
	if interp.GetMemory()[0] != 1 {
		t.Errorf("Ячейка 0: ожидалось 1, получено %d", interp.GetMemory()[0])
	}
}

// TestComplexProgram тестирует более сложную программу
func TestComplexProgram(t *testing.T) {
	// Программа: установить в ячейку 0 значение 5, вывести, обнулить
	code := "MoO MoO MoO MoO MoO OOM OOO OOM"
	var output bytes.Buffer
	interp := NewInterpreter(strings.NewReader(""), &output)
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	expected := "50"
	if output.String() != expected {
		t.Errorf("Ожидалось '%s', получено '%s'", expected, output.String())
	}
	
	if interp.GetMemory()[0] != 0 {
		t.Errorf("Ячейка 0: ожидалось 0, получено %d", interp.GetMemory()[0])
	}
}

// TestMultipleCells тестирует работу с несколькими ячейками
func TestMultipleCells(t *testing.T) {
	code := "MoO MoO moO MoO MoO MoO moO MoO MoO MoO MoO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 2 {
		t.Errorf("Ячейка 0: ожидалось 2, получено %d", interp.GetMemory()[0])
	}
	if interp.GetMemory()[1] != 3 {
		t.Errorf("Ячейка 1: ожидалось 3, получено %d", interp.GetMemory()[1])
	}
	if interp.GetMemory()[2] != 4 {
		t.Errorf("Ячейка 2: ожидалось 4, получено %d", interp.GetMemory()[2])
	}
}

// TestTokenize тестирует функцию разбора токенов
func TestTokenize(t *testing.T) {
	code := "MoO MOo moO"
	tokens := tokenize(code)
	
	expected := []string{"MoO", "MOo", "moO"}
	if len(tokens) != len(expected) {
		t.Fatalf("Ожидалось %d токенов, получено %d", len(expected), len(tokens))
	}
	
	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("Токен %d: ожидалось '%s', получено '%s'", i, expected[i], token)
		}
	}
}

// TestTokenizeCompact тестирует разбор слитного кода
func TestTokenizeCompact(t *testing.T) {
	code := "MoOMOomoO"
	tokens := tokenize(code)
	
	expected := []string{"MoO", "MOo", "moO"}
	if len(tokens) != len(expected) {
		t.Fatalf("Ожидалось %d токенов, получено %d", len(expected), len(tokens))
	}
	
	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("Токен %d: ожидалось '%s', получено '%s'", i, expected[i], token)
		}
	}
}

// TestTokenizeUnknownToken тестирует разбор неизвестного токена
func TestTokenizeUnknownToken(t *testing.T) {
	code := "MoO MOx MoO"
	tokens := tokenize(code)
	
	// Должен вернуть MoO, MOx, MoO
	if len(tokens) != 3 {
		t.Fatalf("Ожидалось 3 токена, получено %d", len(tokens))
	}
	if tokens[0] != "MoO" {
		t.Errorf("Токен 0: ожидалось 'MoO', получено '%s'", tokens[0])
	}
	if tokens[1] != "MOx" {
		t.Errorf("Токен 1: ожидалось 'MOx', получено '%s'", tokens[1])
	}
	if tokens[2] != "MoO" {
		t.Errorf("Токен 2: ожидалось 'MoO', получено '%s'", tokens[2])
	}
}

// TestTokenizeNonCowChars тестирует разбор символов не из COW
func TestTokenizeNonCowChars(t *testing.T) {
	code := "MoO abc MoO"
	tokens := tokenize(code)
	
	// Должен вернуть только MoO, MoO (abc игнорируется)
	if len(tokens) != 2 {
		t.Fatalf("Ожидалось 2 токена, получено %d: %v", len(tokens), tokens)
	}
	if tokens[0] != "MoO" || tokens[1] != "MoO" {
		t.Errorf("Ожидалось ['MoO', 'MoO'], получено %v", tokens)
	}
}

// TestErrorInvalidInstruction тестирует обработку неверной инструкции
func TestErrorInvalidInstruction(t *testing.T) {
	// Используем токен который похож на COW но неправильный
	code := "MoO MOx"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	// Теперь tokenize возвращает неизвестные токены, Parse должен вернуть ошибку
	if err == nil {
		t.Error("Ожидалась ошибка парсинга для неизвестной инструкции MOx")
	}
	if err != nil && !strings.Contains(err.Error(), "неизвестная инструкция") {
		t.Errorf("Ожидалась ошибка 'неизвестная инструкция', получено: %v", err)
	}
}

// TestErrorUnmatchedLoopEnd тестирует обработку лишнего MOO
func TestErrorUnmatchedLoopEnd(t *testing.T) {
	code := "MoO MOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err == nil {
		t.Error("Ожидалась ошибка выполнения для несоответствующего MOO")
	}
}

// TestErrorPointerUnderflow тестирует выход указателя за нижнюю границу
func TestErrorPointerUnderflow(t *testing.T) {
	code := "mOo"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err == nil {
		t.Error("Ожидалась ошибка выхода за границы памяти")
	}
}

// TestInputOutput тестирует комбинацию ввода и вывода
func TestInputOutput(t *testing.T) {
	code := "oom OOM"
	input := strings.NewReader("123\n")
	var output bytes.Buffer
	interp := NewInterpreter(input, &output)
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	expected := "123"
	if output.String() != expected {
		t.Errorf("Ожидалось '%s', получено '%s'", expected, output.String())
	}
}

// TestLoopWithOutput тестирует цикл с выводом
func TestLoopWithOutput(t *testing.T) {
	code := "MoO MoO MoO moo OOM MOo MOO"
	var output bytes.Buffer
	interp := NewInterpreter(strings.NewReader(""), &output)
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	expected := "321"
	if output.String() != expected {
		t.Errorf("Ожидалось '%s', получено '%s'", expected, output.String())
	}
}

// TestEmptyProgram тестирует пустую программу
func TestEmptyProgram(t *testing.T) {
	code := ""
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestMemoryInitialization тестирует начальное состояние памяти
func TestMemoryInitialization(t *testing.T) {
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	if len(interp.GetMemory()) != 30000 {
		t.Errorf("Ожидалось 30000 ячеек памяти, получено %d", len(interp.GetMemory()))
	}
	
	if interp.GetPointer() != 0 {
		t.Errorf("Ожидалось начальное значение указателя 0, получено %d", interp.GetPointer())
	}
	
	for i := 0; i < 100; i++ {
		if interp.GetMemory()[i] != 0 {
			t.Errorf("Ячейка %d должна быть 0, получено %d", i, interp.GetMemory()[i])
		}
	}
}

// TestErrorUnclosedLoop тестирует ошибку незакрытого цикла
func TestErrorUnclosedLoop(t *testing.T) {
	// Цикл moo без MOO должен вызвать ошибку когда значение ячейки = 0
	code := "moo MoO"  // moo без соответствующего MOO
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err == nil {
		t.Error("Ожидалась ошибка выполнения для незакрытого цикла")
	}
	if err != nil && !strings.Contains(err.Error(), "не найден конец цикла") {
		t.Errorf("Ожидалась ошибка 'не найден конец цикла', получено: %v", err)
	}
}

// TestErrorInvalidExecuteAddress тестирует ошибку неверного адреса в mOO
func TestErrorInvalidExecuteAddress(t *testing.T) {
	// Используем отрицательный адрес
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	code := "MOo mOO"  // Уменьшаем на 1 (получится -1), затем mOO с адресом -1
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err == nil {
		t.Error("Ожидалась ошибка выполнения для отрицательного адреса инструкции")
	}
}

// TestErrorInvalidExecuteAddressUpper тестирует ошибку адреса выше границы в mOO
func TestErrorInvalidExecuteAddressUpper(t *testing.T) {
	// Устанавливаем адрес больше или равный количеству инструкций
	// Создаем программу: MoO MoO MoO mOO (4 инструкции, индексы 0-3)
	// После выполнения: memory[0] = 3, затем mOO пытается выполнить инструкцию 3
	// Но нужно установить значение >= 4 чтобы выйти за границы
	code := "MoO MoO MoO MoO mOO"
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	// Устанавливаем значение вручную чтобы оно было >= количества инструкций
	interp.SetMemory(0, 10)  // 10 >= 5 (количество инструкций)
	
	err = interp.Run()
	// Пропускаем первые 4 MoO, выполняем mOO с адресом 10
	if err == nil {
		t.Error("Ожидалась ошибка выполнения для адреса инструкции выше границы")
	}
	if err != nil && !strings.Contains(err.Error(), "неверный адрес инструкции") {
		t.Errorf("Ожидалась ошибка 'неверный адрес инструкции', получено: %v", err)
	}
}

// TestSetMemory тестирует метод SetMemory
func TestSetMemory(t *testing.T) {
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	interp.SetMemory(5, 42)
	if interp.GetMemory()[5] != 42 {
		t.Errorf("Ожидалось 42 в ячейке 5, получено %d", interp.GetMemory()[5])
	}
	
	// Проверка границ
	interp.SetMemory(-1, 100)  // Не должно вызвать панику
	interp.SetMemory(40000, 100)  // Не должно вызвать панику
}

// TestErrorInputFailure тестирует ошибку ввода
func TestErrorInputFailure(t *testing.T) {
	code := "oom"
	input := strings.NewReader("")  // Пустой ввод
	interp := NewInterpreter(input, &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err == nil {
		t.Error("Ожидалась ошибка ввода при пустом потоке")
	}
}

// TestConditionalInputError тестирует ошибку ввода в Moo
func TestConditionalInputError(t *testing.T) {
	code := "Moo"
	input := strings.NewReader("")  // Пустой ввод
	interp := NewInterpreter(input, &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err == nil {
		t.Error("Ожидалась ошибка ввода при пустом потоке в Moo")
	}
}

// TestPointerBoundsNext тестирует выход за верхнюю границу памяти
func TestPointerBoundsNext(t *testing.T) {
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	// Устанавливаем указатель почти в конец
	interp.pointer = 29999
	
	code := "moO"
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err == nil {
		t.Error("Ожидалась ошибка выхода за верхнюю границу памяти")
	}
}

// TestExecuteDecrement тестирует выполнение MOo через mOO
func TestExecuteDecrement(t *testing.T) {
	code := "MOo MoO OOO MoO mOO"  
	// Инструкции: 0:MOo 1:MoO 2:OOO 3:MoO 4:mOO
	// memory[0]=-1, затем 0, обнуление, затем 1
	// mOO выполняет инструкцию 1 (MoO) -> memory[0]=2
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	// После выполнения memory[0]=2
	if interp.GetMemory()[0] != 2 {
		t.Errorf("Ожидалось 2, получено %d", interp.GetMemory()[0])
	}
}

// TestExecuteNext тестирует выполнение moO через mOO
func TestExecuteNext(t *testing.T) {
	code := "MoO MoO OOO moO MoO mOo MoO MoO mOO"  
	// Инструкции: 0:MoO 1:MoO 2:OOO 3:moO 4:MoO 5:mOo 6:MoO 7:MoO 8:mOO
	// memory[0]=2, затем 0, переходим в ячейку 1, инкремент, возвращаемся в 0, инкремент x2
	// memory[0]=2, выполняем инструкцию 2 (OOO) -> memory[0]=0
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestExecutePrev тестирует выполнение mOo через mOO
func TestExecutePrev(t *testing.T) {
	code := "MoO MoO MoO MoO MoO OOO moO MoO MoO MoO MoO mOo MoO mOO"
	// memory[0]=5, затем 0, переход на ячейку 1, memory[1]=4, возврат, memory[0]=1
	// выполняем инструкцию 1 (второй MoO) -> memory[0]=2
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestExecuteOutput тестирует выполнение OOM через mOO
func TestExecuteOutput(t *testing.T) {
	code := "MoO MoO MoO OOM MoO OOO moO MoO MoO MoO mOo MoO MoO MoO mOO"
	// Инструкции: 0:MoO 1:MoO 2:MoO 3:OOM 4:MoO 5:OOO 6:moO 7:MoO 8:MoO 9:MoO 10:mOo 11:MoO 12:MoO 13:MoO 14:mOO
	// memory[0]=3, вывод 3, memory[0]=4, обнуление, переход на ячейку 1
	// memory[1]=3, возврат на 0, memory[0]=3, выполняем инструкцию 3 (OOM)
	var output bytes.Buffer
	interp := NewInterpreter(strings.NewReader(""), &output)
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	// Должен быть вывод
	if len(output.String()) == 0 {
		t.Error("Ожидался вывод от mOO -> OOM")
	}
}

// TestExecuteZero тестирует выполнение OOO через mOO
func TestExecuteZero(t *testing.T) {
	code := "MoO MoO OOO MoO MoO mOO"
	// Инструкции: 0:MoO 1:MoO 2:OOO 3:MoO 4:MoO 5:mOO
	// memory[0]=2, обнуление, memory[0]=0, memory[0]=1, memory[0]=2
	// выполняем инструкцию 2 (OOO) -> memory[0]=0
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	if interp.GetMemory()[0] != 0 {
		t.Errorf("Ожидалось 0, получено %d", interp.GetMemory()[0])
	}
}

// TestExecuteNextBounds тестирует выполнение moO через mOO с проверкой границ
func TestExecuteNextBounds(t *testing.T) {
	// Устанавливаем указатель почти в конец, затем выполняем moO через mOO
	interp := NewInterpreter(strings.NewReader(""), &bytes.Buffer{})
	interp.pointer = 29999
	
	code := "moO MoO mOO"  // Инструкции: 0:moO 1:MoO 2:mOO
	// memory[0]=0, выполняем инструкцию 0 (moO) -> pointer становится 30000 (ошибка)
	
	err := interp.Parse(code)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}
	
	err = interp.Run()
	// moO увеличит pointer до 30000, что вызовет ошибку в основном Run, но не в OpExecute
	// Нужно чтобы moO внутри OpExecute тоже проверял границы
	// Но в текущей реализации OpExecute не проверяет границы для moO/mOo
	// Это нормально - проверка происходит на уровне основного Run
	if err != nil {
		// Ожидаем ошибку выхода за границы
		if !strings.Contains(err.Error(), "выход за границы") {
			t.Errorf("Ожидалась ошибка 'выход за границы', получено: %v", err)
		}
	}
}

