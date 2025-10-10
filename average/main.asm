; main.asm (версия с DEFAULT REL для macOS)

DEFAULT REL ; Эта директива решает проблему адресации для всей программы

section .data
    x dd 5, 3, 2, 6, 1, 7, 4
    y dd 0, 10, 1, 9, 2, 8, 5
    ARRAY_LEN equ ($ - y) / 4

    msg_result db "Arithmetic mean of differences: "
    msg_len equ $ - msg_result
    newline db 10

section .bss
    result_str resb 20

section .text
    global _main

_main:
    xor r9, r9          ; r9 - общая сумма разниц
    xor rcx, rcx        ; rcx - счетчик цикла

    lea rdi, [x]        ; Загружаем адрес x (rel будет добавлено автоматически)
    lea rsi, [y]        ; Загружаем адрес y

loop_start:
    cmp rcx, ARRAY_LEN
    jge loop_end

    mov eax, [rdi + rcx * 4]
    sub eax, [rsi + rcx * 4]
    
    movsx r8, eax
    add r9, r8

    inc rcx
    jmp loop_start

loop_end:
    mov rax, r9
    mov rbx, ARRAY_LEN
    cqo
    idiv rbx
    mov rdi, rax

    ; --- Вывод результата ---
    push rdi                   ; Сохраняем результат деления
    mov rax, 0x2000004         ; syscall write
    mov rdi, 1                 ; stdout
    lea rsi, [msg_result]      ; Загружаем адрес сообщения
    mov rdx, msg_len
    syscall
    
    pop rdi                    ; Восстанавливаем результат деления
    lea rsi, [result_str]      ; Загружаем адрес буфера
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

    ; --- Завершение программы ---
    mov rax, 0x2000001         ; syscall exit
    xor rdi, rdi
    syscall

; Функция int_to_string не меняется
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