package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// OpCode представляет типы операций COW
type OpCode int

const (
	OpIncrement OpCode = iota // MoO
	OpDecrement               // MOo
	OpNext                    // moO
	OpPrev                    // mOo
	OpLoopStart               // moo
	OpLoopEnd                 // MOO
	OpOutput                  // OOM
	OpInput                   // oom
	OpExecute                 // mOO
	OpConditional             // Moo
	OpZero                    // OOO
)

// Instruction представляет инструкцию с кодом операции и позицией
type Instruction struct {
	Op  OpCode
	Pos int
}

// Interpreter представляет интерпретатор COW
type Interpreter struct {
	memory       []int
	pointer      int
	instructions []Instruction
	ip           int // instruction pointer
	reader       *bufio.Reader
	writer       io.Writer
}

// NewInterpreter создает новый интерпретатор
func NewInterpreter(reader io.Reader, writer io.Writer) *Interpreter {
	return &Interpreter{
		memory:       make([]int, 30000),
		pointer:      0,
		instructions: []Instruction{},
		ip:           0,
		reader:       bufio.NewReader(reader),
		writer:       writer,
	}
}

// Parse разбирает исходный код COW на инструкции
func (interp *Interpreter) Parse(code string) error {
	tokens := tokenize(code)
	instructions := []Instruction{}

	for i, token := range tokens {
		var op OpCode
		switch token {
		case "MoO":
			op = OpIncrement
		case "MOo":
			op = OpDecrement
		case "moO":
			op = OpNext
		case "mOo":
			op = OpPrev
		case "moo":
			op = OpLoopStart
		case "MOO":
			op = OpLoopEnd
		case "OOM":
			op = OpOutput
		case "oom":
			op = OpInput
		case "mOO":
			op = OpExecute
		case "Moo":
			op = OpConditional
		case "OOO":
			op = OpZero
		default:
			return fmt.Errorf("неизвестная инструкция: %s", token)
		}
		instructions = append(instructions, Instruction{Op: op, Pos: i})
	}

	interp.instructions = instructions
	return nil
}

// tokenize разбивает код на токены COW
func tokenize(code string) []string {
	var tokens []string
	words := strings.Fields(code)

	for _, word := range words {
		// Ищем все вхождения COW токенов в слове
		for len(word) >= 3 {
			found := false
			for _, pattern := range []string{"MoO", "MOo", "moO", "mOo", "moo", "MOO", "OOM", "oom", "mOO", "Moo", "OOO"} {
				if strings.HasPrefix(word, pattern) {
					tokens = append(tokens, pattern)
					word = word[len(pattern):]
					found = true
					break
				}
			}
			if !found {
				// Если нашли последовательность из 3+ символов M/o/O, но она не совпадает с паттерном
				// Возвращаем её как неизвестный токен
				if len(word) >= 3 && (strings.HasPrefix(word, "M") || strings.HasPrefix(word, "m") || strings.HasPrefix(word, "O") || strings.HasPrefix(word, "o")) {
					// Извлекаем первые 3 символа как потенциальный токен
					potentialToken := word[:3]
					tokens = append(tokens, potentialToken)
					word = word[3:]
				} else {
					word = word[1:]
				}
			}
		}
	}

	return tokens
}

// Run выполняет загруженные инструкции
func (interp *Interpreter) Run() error {
	loopStack := []int{}

	for interp.ip < len(interp.instructions) {
		inst := interp.instructions[interp.ip]

		switch inst.Op {
		case OpIncrement:
			interp.memory[interp.pointer]++

		case OpDecrement:
			interp.memory[interp.pointer]--

		case OpNext:
			interp.pointer++
			if interp.pointer >= len(interp.memory) {
				return fmt.Errorf("выход за границы памяти: pointer=%d", interp.pointer)
			}

		case OpPrev:
			interp.pointer--
			if interp.pointer < 0 {
				return fmt.Errorf("выход за границы памяти: pointer=%d", interp.pointer)
			}

		case OpLoopStart:
			if interp.memory[interp.pointer] == 0 {
				// Пропустить цикл, найти соответствующий MOO
				depth := 1
				for {
					interp.ip++
					if interp.ip >= len(interp.instructions) {
						return fmt.Errorf("не найден конец цикла для moo")
					}
					if interp.instructions[interp.ip].Op == OpLoopStart {
						depth++
					} else if interp.instructions[interp.ip].Op == OpLoopEnd {
						depth--
						if depth == 0 {
							break
						}
					}
				}
			} else {
				loopStack = append(loopStack, interp.ip)
			}

		case OpLoopEnd:
			if len(loopStack) == 0 {
				return fmt.Errorf("MOO без соответствующего moo")
			}
			loopStart := loopStack[len(loopStack)-1]
			if interp.memory[interp.pointer] != 0 {
				interp.ip = loopStart
			} else {
				loopStack = loopStack[:len(loopStack)-1]
			}

		case OpOutput:
			fmt.Fprintf(interp.writer, "%d", interp.memory[interp.pointer])

		case OpInput:
			var value int
			_, err := fmt.Fscanf(interp.reader, "%d", &value)
			if err != nil {
				return fmt.Errorf("ошибка ввода: %v", err)
			}
			interp.memory[interp.pointer] = value

		case OpExecute:
			targetIP := interp.memory[interp.pointer]
			if targetIP < 0 || targetIP >= len(interp.instructions) {
				return fmt.Errorf("неверный адрес инструкции: %d", targetIP)
			}
			// Сохраняем текущий IP и выполняем целевую инструкцию
			savedIP := interp.ip
			interp.ip = targetIP
			inst := interp.instructions[interp.ip]
			
			// Выполняем инструкцию (упрощенная версия без циклов)
			switch inst.Op {
			case OpIncrement:
				interp.memory[interp.pointer]++
			case OpDecrement:
				interp.memory[interp.pointer]--
			case OpNext:
				interp.pointer++
			case OpPrev:
				interp.pointer--
			case OpOutput:
				fmt.Fprintf(interp.writer, "%d", interp.memory[interp.pointer])
			case OpZero:
				interp.memory[interp.pointer] = 0
			}
			
			interp.ip = savedIP

		case OpConditional:
			if interp.memory[interp.pointer] == 0 {
				// Ввод
				var value int
				_, err := fmt.Fscanf(interp.reader, "%d", &value)
				if err != nil {
					return fmt.Errorf("ошибка ввода: %v", err)
				}
				interp.memory[interp.pointer] = value
			} else {
				// Вывод
				fmt.Fprintf(interp.writer, "%d", interp.memory[interp.pointer])
			}

		case OpZero:
			interp.memory[interp.pointer] = 0
		}

		interp.ip++
	}

	return nil
}

// GetMemory возвращает текущее состояние памяти (для тестов)
func (interp *Interpreter) GetMemory() []int {
	return interp.memory
}

// GetPointer возвращает текущую позицию указателя (для тестов)
func (interp *Interpreter) GetPointer() int {
	return interp.pointer
}

// SetMemory устанавливает значение в ячейку памяти (для тестов)
func (interp *Interpreter) SetMemory(index, value int) {
	if index >= 0 && index < len(interp.memory) {
		interp.memory[index] = value
	}
}

