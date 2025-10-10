; main_file.asm - версия с чтением данных из файла

DEFAULT REL

section .data
    filename db "data.txt", 0
    filename_len equ $ - filename - 1
    
    msg_result db "Arithmetic mean of differences: "
    msg_len equ $ - msg_result
    newline db 10
    
    msg_error db "Error reading file", 10
    msg_error_len equ $ - msg_error

section .bss
    file_buffer resb 1024
    x_array resd 10
    y_array resd 10
    result_str resb 20

section .text
    global _main

_main:
    ; Открываем файл
    mov rax, 0x2000005         ; syscall open
    lea rdi, [filename]
    mov rsi, 0                 ; O_RDONLY
    syscall
    
    cmp rax, 0
    jl error_exit
    mov r8, rax                ; r8 - file descriptor
    
    ; Читаем файл
    mov rax, 0x2000003         ; syscall read
    mov rdi, r8                ; file descriptor
    lea rsi, [file_buffer]
    mov rdx, 1024
    syscall
    
    ; Закрываем файл
    mov rax, 0x2000006         ; syscall close
    mov rdi, r8
    syscall
    
    ; Парсим данные
    lea rdi, [file_buffer]
    lea rsi, [x_array]
    call parse_line
    
    ; Находим начало второй строки
    lea rdi, [file_buffer]
    call find_next_line
    lea rsi, [y_array]
    call parse_line
    
    ; Вычисляем среднее арифметическое
    call calculate_average
    
    ; Выводим результат
    call print_result
    
    ; Завершаем программу
    mov rax, 0x2000001         ; syscall exit
    xor rdi, rdi
    syscall

error_exit:
    mov rax, 0x2000004         ; syscall write
    mov rdi, 1                 ; stdout
    lea rsi, [msg_error]
    mov rdx, msg_error_len
    syscall
    
    mov rax, 0x2000001         ; syscall exit
    mov rdi, 1
    syscall

; Функция поиска следующей строки
find_next_line:
    mov al, [rdi]
    cmp al, 10                 ; символ новой строки
    je .found
    cmp al, 0                  ; конец строки
    je .found
    inc rdi
    jmp find_next_line
.found:
    inc rdi                    ; переходим к началу следующей строки
    ret

; Функция парсинга строки чисел
parse_line:
    xor rcx, rcx               ; rcx - счетчик чисел
    xor r8, r8                 ; r8 - текущее число
    
.parse_loop:
    mov al, [rdi]
    cmp al, 0                  ; конец файла
    je .done
    cmp al, 10                 ; конец строки
    je .done
    cmp al, ' '                ; пробел
    je .save_number
    cmp al, 9                  ; табуляция
    je .save_number
    
    ; Преобразуем символ в цифру
    sub al, '0'
    cmp al, 9
    ja .skip_char              ; пропускаем нецифровые символы
    
    ; Добавляем цифру к числу
    imul r8, 10
    add r8, rax
    jmp .next_char
    
.save_number:
    cmp r8, 0
    je .next_char              ; пропускаем пустые числа
    mov [rsi + rcx * 4], r8d   ; сохраняем число
    inc rcx
    xor r8, r8                 ; сбрасываем текущее число
    
.next_char:
    inc rdi
    jmp .parse_loop
    
.skip_char:
    inc rdi
    jmp .parse_loop
    
.done:
    cmp r8, 0
    je .return
    mov [rsi + rcx * 4], r8d   ; сохраняем последнее число
    inc rcx
    
.return:
    ret

; Функция вычисления среднего арифметического
calculate_average:
    xor r9, r9                 ; r9 - общая сумма разниц
    xor rcx, rcx               ; rcx - счетчик цикла
    
    lea rdi, [x_array]
    lea rsi, [y_array]
    
.loop_start:
    cmp rcx, 7                 ; количество элементов
    jge .loop_end
    
    mov eax, [rdi + rcx * 4]
    sub eax, [rsi + rcx * 4]
    
    movsx r8, eax
    add r9, r8
    
    inc rcx
    jmp .loop_start
    
.loop_end:
    mov rax, r9
    mov rbx, 7                 ; количество элементов
    cqo
    idiv rbx
    mov rdi, rax
    ret

; Функция вывода результата
print_result:
    push rdi                   ; Сохраняем результат
    mov rax, 0x2000004         ; syscall write
    mov rdi, 1                 ; stdout
    lea rsi, [msg_result]
    mov rdx, msg_len
    syscall
    
    pop rdi                    ; Восстанавливаем результат
    lea rsi, [result_str]
    call int_to_string
    
    mov rdx, rax
    mov rax, 0x2000004
    mov rdi, 1
    lea rsi, [result_str]
    syscall
    
    mov rax, 0x2000004
    mov rdi, 1
    lea rsi, [newline]
    mov rdx, 1
    syscall
    ret

; Функция преобразования числа в строку
int_to_string:
    mov rbx, 10
    xor rcx, rcx
    mov r8, rsi
    test rdi, rdi
    jns .positive
    neg rdi
    mov byte [rsi], '-'
    inc rsi
    inc r8
.positive:
    cmp rdi, 0
    jne .conversion_loop
    mov byte[rsi], '0'
    inc rsi
    mov rax, 1
    ret
.conversion_loop:
    mov rax, rdi
    xor rdx, rdx
    div rbx
    mov rdi, rax
    add rdx, '0'
    push rdx
    inc rcx
    test rdi, rdi
    jnz .conversion_loop
.write_to_buffer:
    pop rax
    mov [rsi], al
    inc rsi
    dec rcx
    jnz .write_to_buffer
    cmp byte [r8-1], '-'
    jne .calc_len
    mov rax, rsi
    sub rax, r8
    inc rax
    ret
.calc_len:
    mov rax, rsi
    sub rax, r8
    ret
