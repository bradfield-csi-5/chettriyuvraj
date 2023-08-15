section .text
global binary_convert
		
binary_convert:
	xor rax, rax; empty rax register
	mov rsi, rdi; move the address provided in rdi to rsi
	mov cl, 0x01
	cmp byte [rsi], 0;
	jne compute; binary string will always have ascii '0' and '1'
	ret

compute:
	sal rax, 1; left shift al by 1 bit
	add al, byte [rsi]; move byte in rsi to al
	sub al, '0';; subtract by 0 to get numeric value
	add rsi, 1; move rsi to next address 
	cmp byte [rsi], 0x00; check if null character reached
	jne compute
	ret
