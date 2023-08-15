section .text
global sum_to_n
 
sum_to_n: ; init function
	xor rax, rax; setting rax to 0

compute:
	add rax, rdi; add n to rax
	dec rdi; decrement n by 1
	jg compute; no need to use test/cmp - dec modifies zero flag
	ret



