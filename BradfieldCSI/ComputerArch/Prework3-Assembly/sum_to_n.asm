section .text
global sum_to_n
 
sum_to_n: ; init function
	xor rax, rax; setting rax to 0
	mov rcx, rdi; setting rcx to n
	test rcx, rcx;
	jg compute; if n > 0 move to compute
        ret
compute:
	add rax, rcx; add n to rax
	sub rcx, 1; subtract 1 from n
	test rcx, rcx; 
	jg compute; if n > 0 move to compute
	ret
	



