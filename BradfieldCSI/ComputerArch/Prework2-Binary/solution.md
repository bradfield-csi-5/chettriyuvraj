## Exercise solutions:

### Hexadecimal

### 1.2: Say hello to hellohex

```
You have been given a file named `hellohex`. It is 17 bytes in size.

If you were to view a hexadecimal representation of the file, how many hexadecimal characters would you expect? Once you have answered this question, use `xxd -p hellohex` to confirm.
```

- 17 bytes is in theory 34 hex characters, counting it with 'xxd -p hellohex | wc -c' gives 35.

- The reason for this is that there appears to be (I did not know) an implicit \n at the end of the file. This is encoded as a nibble 0xA (or 0xD for carriage return).

- After ignoring this newline, the solution is 

```
echo -n $(xxd -p hellohex) | wc -c

Ans: 34
```


### Integers

### 2.3 

```
Given the following decimal values, determine their 8 bit two’s complement representations:

`127 -128 -1 1 -14`
```

- The intuition behind conversion is as follows:

```
The intuitive way of thinking about two’s complement conversions is that `0b10000...` is the most negative representable value, so when combined with other numbers provides a negative weight of 2 to the power of the number of available bits minus 1. For an 8 bit value, this provides a weight of -128. So if we wished to calculate -127, say, we could take the negative weight value and add 1 to it, giving us `0b10000001`.
```

- 127 = 01111111
- -128 = 10000000
- -1 = 11111111
- 1 = 00000001
- -14 = 11110010

### 2.5

```
It can be beneficial for our hardware to be able to detect overflow in two’s complement. To do so, we’d need a rule for determining—based solely on bit patterns—if overflow has occurred. Can you describe such a rule? Consider the following examples:

// In this case there is a carry out, but the result is correct
110000000 (carry row)
 11000000 (-64)
 01000000 (64)
--------------
 00000000 (0, but there was a carry out!)

// In this case the result is incorrect but there is no carry out
010000000 (carry row)
 01000000 (64)
 01000000 (64)
--------------
 10000000 (-128)

// In this case there is a carry out and the result is incorrect
100000000 (carry row)
 10000000 (-128)
 10000000 (-128)
----------------
 00000000 (0)

```

- The rule is that there is either a carry out bit or a carry in bit, but not both => XOR(carryIn, carryOut)

- Only a carry out bit: overflow into positive numbers. Both operands are negative as MSB is 1, so expected result is negative.

- Only carry in bit: overflow into negative numbers. Both operands are positive since MSB is 0, so expected result is positive.


### Byte Ordering

### 3.2, 3.3

- Exercises on TCP header, bitmap images and their different parts.

- Good exercises, refer to source file in metadata for solutions.

- Keep track of big/little endian for correct answer.


### Fleeting Notes/IEEE 754 Floating Point

### 4.1

```
For the largest fixed exponent, 11111110 == 254 - 127 = 127, what is the smallest (magnitude) incremental change that can be made to a number?
```

- It is a 1 in the fractional component, which would be 2^127 * (1 + 1/2^23)

```
For the smallest (most negative) fixed exponent, what is the smallest (magnitude) incremental change that can be made to a number?
```

- Most negative fixed exponent (denormalized value) would be -126 with all 0's in the fraction = 2 ^-126 * (0 + 0) = 0. Incrementing would give us 2^ -126 * (0 + 1/-23) = 2^-149. The difference is 2^-149.

```
What does this imply about the precision of IEEE Floating Point values?
```

- As the numbers become larger, precision decreases. This is a feature not a bug.

### 4.2

This is code to convert 32 bit value to next highest power of 2. How does this code snippet work? 
```
unsigned int const v; // Round this 32-bit value to the next highest power of 2
unsigned int r;       // Put the result here. (So v=3 -> r=4; v=8 -> r=8)

if (v > 1) 
{
  float f = (float)v;
  unsigned int const t = 1U << ((*(unsigned int *)&f >> 23) - 0x7f);
  r = t << (t < v);
}
else 
{
  r = 1;
}
```

Consider the example v = 3

- e: biased exponent value
- E: actual exponent value

- Binary representation of 3.0 = 11.0 = 1.1 * 2^1.

- Casting it to float converts it into IEEE 754 32 bit representation, so we can expect the bits 23 - 30 to represent 'Bias + E' = e i.e Bits 23-30 represent 127 + 1 = 128

- The value of 't' is as follows:
	- Innermost right bit shift by 23 digits brings its 'e' value to position of the least significant 8 bits (0 - 7)
	
	- Subtracting by 0x7f (0111 0111 in binary) removes the bias and gives us true value of exponent i.e converts 'e' to 'E'.
	
	- At this point, we shift the value '1U' i.e unsigned 1 left by 'E' bits to get a value 2^E (2^1 in case of 3).
	
- Finally, the value of r is computed such that if we started with a power of 2, we get back what we started with, otherwise we shift left by a further bit to get the next power of 2.

We could do the same with double precision floats and 64 bit ints by changing a few values.

### Character Encodings

### 5

- There is a possible additional cost to encoding a completely ASCII doc in UTF-8

- Extra space might be required as UTF-8 reserves bits for error checking + length count. 

- An ASCII document of all characters <= 7 bits would incur no extra cost, but each 8 bit ASCII value would require an additional byte to encode.

### 5.3 

- Simply printf ("\\a") will ring the bell.

## Additional Notes

- Big-Endian is an ordering in which bytes are stored in their 'intuitive' order. MSB is on the left.

- Little-Endian is an ordering in which bytes are stored in the reverse order. MSB is on the right.

- Note that the order of bytes is reversed, but the actual 4-bit nibbles in those bytes are ordered in the same way.

- his has implications, such as when interpreting a stream of bytes as hexadecimal, their Big and Little Endian interpretations are not exactly in reverse order.
