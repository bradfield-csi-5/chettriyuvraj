default rel

section .text
global volume

volume:

	movss xmm2, [pi]		
	movss xmm3, [const_3]		
	mulss xmm0, xmm0		; xmm0 = r^2
	mulss xmm0, xmm1		; xmm0 = r^2 * h
	mulss xmm0, xmm2		; xmm0 = r^2 * h * pi
	divss xmm0, xmm3		; xmm0 = (r^2 * h * pi)/3

	ret

section .data
pi: dd 3.14159
const_3: dd 3.0
