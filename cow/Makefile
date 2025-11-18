.PHONY: build test clean coverage run-examples

# Сборка интерпретатора
build:
	go build -o cow

# Запуск всех тестов
test:
	go test -v

# Проверка покрытия
coverage:
	go test -coverprofile=coverage.out
	@echo "Общее покрытие (включая main.go):"
	@go tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "Покрытие без main.go (только interpreter.go):"
	@go tool cover -func=coverage.out | grep "^cow/interpreter" | awk '{gsub(/%/, "", $$NF); sum+=$$NF; count++} END {if (count > 0) printf "%.1f%%\n", sum/count; else print "0%"}'

# HTML отчет о покрытии
coverage-html:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Очистка
clean:
	rm -f cow coverage.out

# Запуск примеров
run-examples: build
	@echo "=== Пример 1: Вывод числа 10 ==="
	./cow examples/hello.cow
	@echo "\n=== Пример 2: Цикл (5,4,3,2,1) ==="
	./cow examples/loop.cow
	@echo "\n=== Пример 3: Ввод и вывод (введите число) ==="
	@echo "42" | ./cow examples/input_output.cow

# Форматирование кода
fmt:
	go fmt ./...

# Проверка линтером
lint:
	go vet ./...

# Полная проверка перед коммитом
check: fmt lint test
	@echo "Все проверки пройдены успешно!"

