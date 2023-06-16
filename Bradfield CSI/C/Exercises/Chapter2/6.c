#include<stdio.h>

unsigned setbits(unsigned x, int p, int n, unsigned y);

int main() {
    /* n can be p + 1 at most */
    printf("The result of setbits is %u", setbits(0b101010101, 4, 2, 0b1101101));
    // printf("The result of setbits is %u", setbits(122, 5, 3, 147));

    return 0;
}

unsigned setbits(unsigned x, int p, int n, unsigned y) {
    unsigned yRightmostBits = y & ~(~0 << n); /* Grab rightmost n bits from y */
    unsigned xRightmostBits = (x & ~(~0 << (p + 1 - n))); /* Grab rightmost (p + 1 - n) bits from x */
    unsigned xLeftmostBits = ( x >> (p + 1)); /* Grab leftmost bits from x after removing starting p + 1 bits */

    unsigned res = xLeftmostBits;
    res <<= n; /* Move bits left by n to make space for n bits from y */
    res |= yRightmostBits; /* Append n rightmost bits of y */
    res <<= (p + 1 - n); /* Move bits left by p to make space for rightmost p + 1 - n bits from x */
    res |= xRightmostBits; /* Append p rightmost bits of x */

    return res;
}