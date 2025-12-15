package main

import (
	"os"
	"testing"
)

// TestRunInterpreterFileNotFound тестирует runInterpreter с несуществующим файлом
func TestRunInterpreterFileNotFound(t *testing.T) {
	err := runInterpreter("nonexistent.pas")
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего файла")
	}
}

// TestRunInterpreterLexerError тестирует runInterpreter с ошибкой лексера
func TestRunInterpreterLexerError(t *testing.T) {
	// Создаем временный файл с невалидным кодом
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("BEGIN x @ 5; END.") // Невалидный символ @
	tmpfile.Close()

	err = runInterpreter(tmpfile.Name())
	if err == nil {
		t.Error("Ожидалась ошибка лексера")
	}
}

// TestRunInterpreterParserError тестирует runInterpreter с ошибкой парсера
func TestRunInterpreterParserError(t *testing.T) {
	// Создаем временный файл с невалидным синтаксисом
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("x := 5; END.") // Нет BEGIN
	tmpfile.Close()

	err = runInterpreter(tmpfile.Name())
	if err == nil {
		t.Error("Ожидалась ошибка парсера")
	}
}

// TestRunInterpreterInterpreterError тестирует runInterpreter с ошибкой интерпретатора
func TestRunInterpreterInterpreterError(t *testing.T) {
	// Создаем временный файл с делением на ноль
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("BEGIN x := 10 / 0; END.")
	tmpfile.Close()

	err = runInterpreter(tmpfile.Name())
	if err == nil {
		t.Error("Ожидалась ошибка интерпретатора")
	}
}

// TestRunInterpreterSuccess тестирует успешное выполнение runInterpreter
func TestRunInterpreterSuccess(t *testing.T) {
	// Создаем временный файл с валидным кодом
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("BEGIN x := 5; END.")
	tmpfile.Close()

	err = runInterpreter(tmpfile.Name())
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestRunInterpreterSuccessEmpty тестирует успешное выполнение runInterpreter с пустой программой
func TestRunInterpreterSuccessEmpty(t *testing.T) {
	// Создаем временный файл с пустой программой
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("BEGIN END.")
	tmpfile.Close()

	err = runInterpreter(tmpfile.Name())
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestRunInterpreterSuccessWithVariables тестирует успешное выполнение runInterpreter с переменными
func TestRunInterpreterSuccessWithVariables(t *testing.T) {
	// Создаем временный файл с переменными
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("BEGIN x := 2 + 3; y := 5 * 2; END.")
	tmpfile.Close()

	err = runInterpreter(tmpfile.Name())
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestRunInterpreterSuccessFloatOutput тестирует успешное выполнение runInterpreter с дробными числами
func TestRunInterpreterSuccessFloatOutput(t *testing.T) {
	// Создаем временный файл с дробными числами
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("BEGIN x := 7 / 2; END.")
	tmpfile.Close()

	err = runInterpreter(tmpfile.Name())
	if err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
}

// TestMainWithExitCodeNoArgs тестирует mainWithExitCode без аргументов
func TestMainWithExitCodeNoArgs(t *testing.T) {
	// Сохраняем оригинальные аргументы
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"pascal"}
	exitCode := mainWithExitCode()
	if exitCode != 1 {
		t.Errorf("Ожидался код выхода 1, получен %d", exitCode)
	}
}

// TestMainWithExitCodeError тестирует mainWithExitCode с ошибкой
func TestMainWithExitCodeError(t *testing.T) {
	// Сохраняем оригинальные аргументы
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"pascal", "nonexistent.pas"}
	exitCode := mainWithExitCode()
	if exitCode != 1 {
		t.Errorf("Ожидался код выхода 1, получен %d", exitCode)
	}
}

// TestMainWithExitCodeSuccess тестирует успешное выполнение mainWithExitCode
func TestMainWithExitCodeSuccess(t *testing.T) {
	// Сохраняем оригинальные аргументы
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Создаем временный файл
	tmpfile, err := os.CreateTemp("", "test_*.pas")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("BEGIN x := 5; END.")
	tmpfile.Close()

	os.Args = []string{"pascal", tmpfile.Name()}
	exitCode := mainWithExitCode()
	if exitCode != 0 {
		t.Errorf("Ожидался код выхода 0, получен %d", exitCode)
	}
}


