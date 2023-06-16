#include<stdio.h>

unsigned rot(unsigned x, int n);

int main() {
    printf("Integer 101 Rotated by 9, we have %u", rot(0b1100101, 9)); 
    return 0;
}

unsigned rot(unsigned x, int n) {
    int highestBitIndex = 0; /* Find index of the highest bit that is set */
    unsigned xCheck = x;
    while (xCheck >>= 1) {
        ++highestBitIndex;
    }
    int length = highestBitIndex + 1; /* Length of the bit representation */
    n %= length; /* a 7-length bit representation rotated by 9 =  same as 2 rotations */

    unsigned rightmostBits = (x & ~(~0 << n)); /* Grab rightmost n bits - will be lost when we shift right */
    x >>= n; /* Finally shift right by n bits */
    
    return x | (rightmostBits << (length - n)); /* Append rightmost bits which were lost to the left of the bit representation */
}