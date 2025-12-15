package main

import (
	"fmt"
	"os"
)

// runInterpreter выполняет интерпретацию Pascal программы из файла
func runInterpreter(filename string) error {
	code, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	// Лексический анализ
	lexer := NewLexer(string(code))
	tokens, err := lexer.Tokenize()
	if err != nil {
		return fmt.Errorf("ошибка лексического анализа: %v", err)
	}

	// Синтаксический анализ
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("ошибка синтаксического анализа: %v", err)
	}

	// Интерпретация
	interpreter := NewInterpreter()
	err = interpreter.Interpret(program)
	if err != nil {
		return fmt.Errorf("ошибка выполнения: %v", err)
	}

	// Вывод значений всех переменных
	variables := interpreter.GetVariables()
	if len(variables) == 0 {
		fmt.Println("{}")
	} else {
		fmt.Print("{")
		first := true
		for name, value := range variables {
			if !first {
				fmt.Print(", ")
			}
			// Выводим целое число, если оно целое
			if value == float64(int64(value)) {
				fmt.Printf("%s: %d", name, int64(value))
			} else {
				fmt.Printf("%s: %g", name, value)
			}
			first = false
		}
		fmt.Println("}")
	}
	return nil
}

func main() {
	os.Exit(mainWithExitCode())
}

// mainWithExitCode выполняет основную логику и возвращает код выхода
func mainWithExitCode() int {
	if len(os.Args) < 2 {
		fmt.Println("Использование: pascal <файл.pas>")
		return 1
	}

	err := runInterpreter(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}

