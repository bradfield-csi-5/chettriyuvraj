/* 	Write a program entab that replaces strings of blanks with the minimum number of tabs and blanks to achieve the same spacing. Use the same stops as for detab . When either a tab or a single blank would suffice to reach a tab stop, which should be given preference? */

#include<stdio.h>
#define TABSTOP 8
#define COLUMNS 100

int main() {
    int c, blankCount = 0, blankStartCol, col = 0;
    while((c = getchar()) != EOF && c != '\n') {
        col++;
        if (c == ' ') {
            if (blankCount == 0) {
                blankStartCol = col - 1;
            }
            blankCount++;
        } else { /* Any other character */
            if (blankCount > 0) { /* Replace blanks with minimum spaces and tabs first */
                int distanceToNextTabStop = TABSTOP - (blankStartCol%TABSTOP); /* Distance from 'blankEndCol' to next tabstop - excluding current printed char */
                if (blankCount < distanceToNextTabStop) {/* Print all spaces */
                    for (int i = 0; i < blankCount; i++)
                        printf(" ");
                } else if (blankCount == distanceToNextTabStop){ /* Print one tab */
                    printf("\t");
                } else { /* If blankCount > distanceToNextTabStop  */
                    /* Print one tab - then print as many tabs as possible, then finally print remaining spaces */
                    int remainingDistance = blankCount - distanceToNextTabStop;
                    int tabs = 1 + (remainingDistance/TABSTOP);
                    int spaces = remainingDistance%TABSTOP;
                    for(int i = 0; i<tabs; i++)
                        printf("\t");
                    for(int j = 0; j<spaces; j++)
                        printf(" ");
                }
                blankCount = 0;
            }
            putchar(c);
        }
    }

    return 0;
}