section .text
global pangram
pangram:
	push rbx; callee saved register
	xor rbx, rbx; using first 25 bits as a hashmap

forcelower: ; if alphabets, force to lowercase and constrict them to decimal 0-25
	movzx eax, byte [rdi]
	cmp eax, 0
	je result
	or eax, 32;
        sub eax, 'a'

compute:; create a bitmask and then use it to update hashmap
	bts rbx, rax
	inc rdi
        jmp forcelower; circle back to start of the loop

result:; invoked when loop condition fails
	xor rax, rax; result stored here
	and rbx, 0x3ffffff; keep only bits 0-25 of hashmap
	cmp rbx, 0x3ffffff ; comparing our map with that of a pangram to determine if we have a pangram
	sete al
done:
	pop rbx
	ret


