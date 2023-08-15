section .text
global pangram
pangram:
	push rbx
	xor rbx, rbx; using first 25 bits as a hashmap
	xor rax, rax; result stored in rax at the end + used to hold each bit of sentence

check:
     cmp byte[rdi], 0x00
     je result	


forcelower: ; if alphabets, force to lowercase and constrict them to decimal 0-25
	movzx eax, byte [rdi]
	or eax, 32;
        sub eax, 'a'

compute:; create a bitmask and then use it to update hashmap
	mov r8, 1
	mov cl, al
	sal r8, cl
	or rbx, r8
	inc rdi
        jmp check; circle back to start of the loop

result:; invoked when loop condition fails
	xor rax, rax; result stored here
	and rbx, 0x3ffffff; keep only bits 0-25 of hashmap
	cmp rbx, 0x3ffffff ; comparing our map with that of a pangram to determine if we have a pangram
	jne done
	inc rax; input is a pangram

done:
	pop rbx
	ret


