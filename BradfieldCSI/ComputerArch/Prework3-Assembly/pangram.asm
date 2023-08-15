section .text
global pangram
pangram:
	push rbx
	xor rbx, rbx; using first 25 bits as a hashmap
	xor rax, rax; result stored in rax at the end + used to hold each bit of sentence
	
	mov rsi, rdi; move string address to rsi
	
		
	jmp check;

setlowerupper:
	cmp byte [rsi], 90; check if byte within decimal 90
	jg lowercase

uppercase:; try to map byte to decimal 0 - 25 by subtracting ascii 'A'
	lodsb
	sub al, 0x41
	jmp compute


lowercase:; try to map byte to decimal 0-25 by subtracting ascii 'a'
	lodsb
	sub al, 0x61

compute:; check if byte falls between decimal 0-25, if yes create a bitmask
	cmp al, 0
	jl check	

	cmp al, 25
	jg check

	mov r8, 1
	jmp bitmask

check:
	cmp byte [rsi], 0x00
	jne setlowerupper

result:
	xor rax, rax; result stored here
	mov r8, 1 << 26
	dec r8 ; this is hashmap for a pangram i.e 2^26 - 1 = 1111...
	cmp r8, rbx ; comparing our map with that of a pangram to determine if we have a pangram
	jne done
	inc rax; input is a pangram

done:
	pop rbx
	ret

bitmask:
	mov cl, al
	sal r8, cl; our bitmask starts as 1, if al is the 'ith' index, shift left i times to get mask

mask:
	or rbx, r8; apply bitmask to our hash table
	jmp check


