import argparse
import sys
from collections import Counter

def analyze_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            lines = f.readlines()
        
        # Для подсчета символов и частотного словаря лучше прочитать файл целиком
        with open(filepath, 'r', encoding='utf-8') as f:
            full_content = f.read()

        # 1. Количество строк
        line_count = len(lines)

        # 2. Количество символов
        char_count = len(full_content)

        # 3. Количество пустых строк
        empty_line_count = 0
        for line in lines:
            if not line.strip():
                empty_line_count += 1
        
        # 4. Частотный словарь символов
        char_frequency = Counter(full_content)

        return {
            'line_count': line_count,
            'char_count': char_count,
            'empty_line_count': empty_line_count,
            'char_frequency': char_frequency
        }

    except FileNotFoundError:
        print(f"Ошибка: Файл не найден по пути '{filepath}'")
        return None
    except Exception as e:
        print(f"Произошла непредвиденная ошибка: {e}")
        return None

def display_results(stats, choices):
    print("\n--- Результаты анализа ---")

    if '1' in choices:
        print(f"1. Количество строк: {stats['line_count']}")
    
    if '2' in choices:
        print(f"2. Количество символов: {stats['char_count']}")

    if '3' in choices:
        print(f"3. Количество пустых строк: {stats['empty_line_count']}")

    if '4' in choices:
        print("4. Частотный словарь символов:")
        # Сортируем для более наглядного вывода
        sorted_freq = sorted(stats['char_frequency'].items())
        for char, count in sorted_freq:
            # Для специальных символов (как '\n') выводим их представление
            display_char = repr(char)[1:-1]
            print(f"  '{display_char}': {count}")

    print("------------------------\n")


def main():
    """
    Главная функция программы.
    """
    parser = argparse.ArgumentParser(
        description="Консольная утилита для анализа текстовых файлов."
    )
    parser.add_argument(
        "filename", 
        type=str, 
        help="Имя текстового файла для анализа."
    )
    args = parser.parse_args()

    # Анализируем файл
    statistics = analyze_file(args.filename)
    
    # Если анализ не удался, выходим
    if statistics is None:
        sys.exit(1)

    # Меню для пользователя
    menu_options = {
        '1': 'Количество строк',
        '2': 'Количество символов',
        '3': 'Количество пустых строк',
        '4': 'Частотный словарь символов'
    }

    while True:
        print("Выберите, какую информацию вы хотите видеть:")
        for key, value in menu_options.items():
            print(f"  {key}. {value}")
        
        print("\nВведите номера желаемых пунктов через запятую (например: 1,3,4).")
        print("Для выхода введите 'q' или 'exit'.")
        
        user_choice = input("Ваш выбор: ").strip().lower()

        if user_choice in ['q', 'exit']:
            print("Завершение работы.")
            break

        # Проверяем корректность ввода
        selected_options = [item.strip() for item in user_choice.split(',')]
        
        valid_choices = [opt for opt in selected_options if opt in menu_options]

        if not valid_choices:
            print("\nНекорректный ввод. Пожалуйста, выберите номера из списка.\n")
            continue

        # Выводим выбранные результаты
        display_results(statistics, valid_choices)

if __name__ == "__main__":
    main()