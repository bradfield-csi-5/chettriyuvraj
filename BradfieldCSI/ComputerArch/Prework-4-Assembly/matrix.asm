section .text
global index
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex
	; assuming integer matrix, we can access i,jth element as xA + L(C*i + j) = xA + 4(C*i + j)
	

	imul rcx, rdx					; C*i
	lea rcx, [rdi+rcx*4]				; xA + 4*C*i
	lea rcx, [rcx+r8*4]				; xA + 4*C*i + 4j
	mov rax, [rcx]

	ret
