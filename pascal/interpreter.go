package main

import (
	"fmt"
)

// Interpreter представляет интерпретатор Pascal
type Interpreter struct {
	variables map[string]float64
}

// NewInterpreter создает новый интерпретатор
func NewInterpreter() *Interpreter {
	return &Interpreter{
		variables: make(map[string]float64),
	}
}

// Interpret выполняет программу
func (i *Interpreter) Interpret(program *Program) error {
	return i.executeStatements(program.Statements)
}

// executeStatements выполняет список операторов
func (i *Interpreter) executeStatements(statements []Statement) error {
	for _, stmt := range statements {
		err := i.executeStatement(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// executeStatement выполняет оператор
func (i *Interpreter) executeStatement(stmt Statement) error {
	switch s := stmt.(type) {
	case *Assignment:
		value, err := i.evaluateExpression(s.Value)
		if err != nil {
			return err
		}
		i.variables[s.Variable] = value
		return nil
	case *Block:
		return i.executeStatements(s.Statements)
	default:
		return fmt.Errorf("неизвестный тип оператора: %T", stmt)
	}
}

// evaluateExpression вычисляет значение выражения
func (i *Interpreter) evaluateExpression(expr Expression) (float64, error) {
	switch e := expr.(type) {
	case *Number:
		return e.Value, nil
	case *Identifier:
		value, ok := i.variables[e.Name]
		if !ok {
			// Переменная не инициализирована, считаем её равной 0
			return 0, nil
		}
		return value, nil
	case *BinaryOp:
		left, err := i.evaluateExpression(e.Left)
		if err != nil {
			return 0, err
		}
		
		right, err := i.evaluateExpression(e.Right)
		if err != nil {
			return 0, err
		}
		
		switch e.Operator {
		case TokenPLUS:
			return left + right, nil
		case TokenMINUS:
			return left - right, nil
		case TokenMULTIPLY:
			return left * right, nil
		case TokenDIVIDE:
			if right == 0 {
				return 0, fmt.Errorf("деление на ноль")
			}
			return left / right, nil
		default:
			return 0, fmt.Errorf("неизвестный оператор: %v", e.Operator)
		}
	default:
		return 0, fmt.Errorf("неизвестный тип выражения: %T", expr)
	}
}

// GetVariables возвращает словарь всех переменных
func (i *Interpreter) GetVariables() map[string]float64 {
	// Создаем копию, чтобы избежать изменений извне
	result := make(map[string]float64)
	for k, v := range i.variables {
		result[k] = v
	}
	return result
}



