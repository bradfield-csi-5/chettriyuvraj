#include<stdio.h>

unsigned setbits(unsigned x, int p, int n, unsigned y);

int main() {
    printf("The result of setbits is %u", setbits(0b101010101, 2, 4, 0b1101101));
    return 0;
}

unsigned setbits(unsigned x, int p, int n, unsigned y) {
    unsigned yRightmostBits = y & ~(~0 << n); /* Grab rightmost n bits from y */
    unsigned xRightmostBits = (x & ~(~0 << p)); /* Grab rightmost p bits from x */
    unsigned xLeftmostBits = ( x >> (p + n)); /* Grab leftmost bits from x after removing starting p + n bits */

    unsigned res = xLeftmostBits;
    res <<= n; /* Move bits left by n to make space for n bits from y */
    res |= yRightmostBits; /* Append n rightmost bits of y */
    res <<= p; /* Move bits left by p to make space for rightmost p bits from x */
    res |= xRightmostBits; /* Append p rightmost bits of x */

    return res;
}