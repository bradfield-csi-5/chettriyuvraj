section .text
global binary_convert
		
binary_convert:
	xor rax, rax; empty rax register

compute:
	cmp byte [rdi], 0
	je done
	sal rax, 1; left shift al by 1 bit
	add al, byte [rdi]; move byte in rdi to al
	sub al, '0';; subtract by 0 to get numeric value
	inc rdi ; move rsi to next address 
	jmp compute
done:
	ret


	
