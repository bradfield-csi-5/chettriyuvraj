/*Write a program detab that replaces tabs in the input with the proper number of blanks to space to the next tab stop. Assume a fixed set of tab stops, say every n columns. Should n be a variable or a symbolic parameter?*/

/*
0. Define a TABSTOP as any val eg. 8
1. Keep a count of the incoming characters (mod 8)
2. At any point, no of tabs to be added = (TABSTOP - charCount%8)

PS. Print a long string of zeroes at the start for reference
*/

#define TABSTOP 8
#define COLUMNS 100

#include<stdio.h>
int main() {
    int curChar, charCount = 0;

    /*Column bar for reference*/
    for(int i = 0; i < COLUMNS; i++)
        printf("%d", i%TABSTOP + 1);
    printf("\n");
    
    while((curChar = getchar()) != EOF) {
        charCount %= TABSTOP; /* Since we are counting the distance to the next tab stop */

        if (curChar == '\n') { /* Move to the next line */
            charCount = 0;
        } else if (curChar == '\t') { /* Add remaning spaces to get to tab stop */
            int remainingSpaces = TABSTOP - charCount;
            for(int i = 0; i < remainingSpaces; i++)
                printf(" ");
            charCount = 0;
        } else {  /* Any other character */
            ++charCount;
        }
    }


    return 0;
}
