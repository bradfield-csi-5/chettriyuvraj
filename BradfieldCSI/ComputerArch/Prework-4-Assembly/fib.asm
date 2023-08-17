section .text
global fib
fib:
	mov eax, edi				; store n
	cmp eax, $1 				; if n <= 1, return
	jle done 				
	
	mov ecx, edi				; preserve n in rcx
	push rcx 				; push n to the stack
	dec rdi 				; n -> n - 1 for recursive call
	call fib				; recurse

	mov rdx, rax				; store fibo (n-1) in rdx
	pop rcx					; pop and grab value of n
	lea rdi, [rcx - 2]			; n -> n - 2 for recursive call
	push rdx				; push fibo (n-1) to stack
	call fib				; recurse

	pop rdx					; pop fibo(n-1) from stack
	lea rax, [rax + rdx]			; fibo (n) = fibo (n) + fibo (n-1)

done:
	ret
