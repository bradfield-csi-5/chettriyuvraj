section .text
global pangram
pangram:
	push rbx
	xor rbx, rbx; using first 25 bits as a hashmap
	xor rax, rax; result stored in rax at the end + some computations
	xor rcx, rcx; 2 jobs, check result and bit masking

	mov rsi, rdi; move string address to rsi
	
	jmp check;

uppercheck: ;; check for uppercase characters
	cmp byte [rsi], 65;
	jl not_an_alphabet

	cmp byte [rsi], 90;
	jg lowercheck

	lodsb; stores [rsi] in al and increments rsi
        sub al, 0x41

	xor rcx, rcx
	inc rcx; put 1 in rcx - this will be our bitmask to the pangram 'hash table'
	
	jmp bitmask

lowercheck: ;; check for lowercase characters
	cmp byte [rsi], 97;
	jl not_an_alphabet

	cmp byte [rsi], 122;
	jg not_an_alphabet

	lodsb; store character eg 'a' in al from [rsi] and inc rsi
	sub al, 0x61; between 97-122; map to 0-25	

	xor rcx, rcx
	inc rcx; put 1 in rcx - this will be our bitmask to the pangram 'hash table'
	
	jmp bitmask

not_an_alphabet: ; increment rsi since not a char
	lodsb
	

check:
	cmp byte [rsi], 0x00; end of string reached
	jne uppercheck

result:
	xor rax, rax; result stored here
	xor rcx, rcx
	add rcx, 1 << 26
	sub rcx, 1; this is hashmap for a pangram i.e 2^26 - 1 = 1111111...
	cmp rcx, rbx
	jne done; not a pangram
	inc rax; is a pangram
	jmp done

done:
	pop rbx
	ret

bitmask:
	cmp al, 0
	je mask; if a is already 0, we have our mask ready
	sal rcx, 1; our bitmask starts as 1, if al is the 'ith' index, shift left i times to get mask
	dec al
	jmp bitmask

mask:
	or rbx, rcx; apply bitmask to our hash table
	jmp check
	
