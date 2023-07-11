/* Write the function itob(n,s,b) that converts the integer n 
into a base b character representation in the string s . 
In particular, itob(n,s,16) formats n as a hexadecimal integer in s

Note: Negative numbers return their two's complement hex representations
*/


#include<stdio.h>

void itob(unsigned n, char s[], int b);
char hexMapping []= {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'};

int main() {
    char s[9]; /* Assuming maximum hex repr is FFFFFFFF + 1 additional for \0 */
    itob(12234, s, 16);
    printf("%s", s);
    
}

void itob(unsigned n, char s[], int b) {
    int curNum = 0; /* We iterate 4 bits at a time, then assign the requisite hex character for it*/
    int sIndex = 0;
    while (n != 0) {
        curNum |= ((~(~0 << 4)) & n); /* Grab first 4 digit of n - from the right */
        printf("curNum %x \n", curNum);
        n >>= 4; /* Move n to the right by 4 bits */
        s[sIndex++] = hexMapping[curNum];
        curNum = 0;
    }

    /* Reverse hex string */
    for (int i = 0, j = sIndex - 1; i < j; i++, j--)
    {
        char temp = s[i];
        s[i] = s[j];
        s[j] = temp;
    }
    s[sIndex] = '\0';
}