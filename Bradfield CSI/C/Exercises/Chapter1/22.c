/* Write a program to fold long input lines into two or more shorter lines ...*/

/* 
Testcases:
1:
                                tabs                    and                     just                    p
lenty           of                              tabs    hereisatestcasewithoutanyspacesandwithoutanybullshit



2:
lets    test    for     tabs


3.
notabsnotabsnotabsnotabsnotabsnotabsnotabsnotabsnotabsnotabsnotabsnotabsnotabs

*/


#define MAXLINELENGTH 20 /* Assuming MAXLINELENGTH always > TABSTOP */
#define TABSTOP 8
#define IN 1
#define OUT 0

#include<stdio.h>

int main() {
     int charCount = 0, curChar, prevChar = '\0', state = OUT;
     int remainingSpaces; /* Used to compute distance to next tabstop */

    for(int i = 0; i < MAXLINELENGTH; i++) /* Print columns for reference */
        printf("%d", i%10);
    printf("\n");

    while((curChar = getchar()) != EOF) {
        if (curChar == '\n') { /* Early exit into next iteration - move to nextline */
            prevChar = '\n';
            charCount = 0;
            printf("\n");
            continue;
        }

        if (charCount == 0 ) { /* Either very first char - or state = IN/OUT and moving into next line */
            if (prevChar != '\n' && prevChar != '\t' && prevChar != ' ' && prevChar != '\n' && prevChar != '\0' && curChar != '\n' && curChar != '\t' && curChar != ' ') { /* state = IN and moving into next line */
                printf("-\n-");
                ++charCount;
            } else if (prevChar == ' ' || curChar == ' ') { /* Space separated input - other cases such as tabs, newlines and continued input handled separately */
                printf("\n");
            }
        }

        /* Handling state */
        if (curChar == ' ' || curChar == '\t') { /* In the midst of emptiness */
            state = OUT;
        } else { /* In the midst of character input - any other character */
            state = IN;
        }

        
        if (curChar == '\t') { /* Handling tabs separately */
            remainingSpaces = TABSTOP - (charCount % TABSTOP); /* Distance to next tabstop */
            charCount += remainingSpaces;
            if (charCount >= MAXLINELENGTH) { /* If tab spaced exceeds current line */
                printf("\n"); /* Go to nextline then insert tab */
                charCount = TABSTOP; 
            }
            putchar(curChar); /* Print tab */
        } else { /* All other characters */
            putchar(curChar);
            ++charCount;
            charCount %= MAXLINELENGTH;
        }

        prevChar = curChar;

    }
    

}

