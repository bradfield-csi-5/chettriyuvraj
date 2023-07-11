/* Write a program to remove all comments - don't forget to handle quoted strings and character constants */
/* Assumption - all programs syntactically valid */

/* Test case: (Input taken line by line)

Bottom of this page: https://clc-wiki.net/wiki/K%26R2_solutions:Chapter_1:Exercise_23#Uncomment03_.28Gregory_Pietsch.29
Note: Did not account for wrapping a comment using '\' symbol

*/



#define OUT 0   /* Define valid states */
#define INQUOTES 1
#define INSINGLECOMMENTS 2 
#define INMULTICOMMENTS 3

#include<stdio.h>
int main() {
    int curChar, prevChar = '\0';
    int state = OUT;

    while ((curChar = getchar()) != EOF) {
        if (curChar == '\n') { /* Print newline regardless of state */
            putchar(curChar);
        } else if (state == OUT) { /* If not inside any quotes/comments - check for any possible change in states */
            if (curChar == '/' && prevChar == '/') {
                state = INSINGLECOMMENTS;
            } else if (curChar == '*' && prevChar == '/')  {
                state = INMULTICOMMENTS;
            } else if (curChar == '\"' || curChar == '\'') { /* If quotes start - print them + change state */
                state = INQUOTES;
                putchar(curChar);
            } else if (curChar != '\\' && curChar != '/') { /* If no change in state, simply print char */
                putchar(curChar);
            }
        } else  { /* Currently inside comments or quotes */
            if ((state == INSINGLECOMMENTS && curChar == '\n') || (state == INMULTICOMMENTS && curChar == '/' && prevChar == '*')) { /* If comment ends */
                state = OUT;
            } else if (state == INQUOTES)   { /* If quotes end */
                if (curChar == '\'' || curChar == '\"') {
                    state = OUT;
                }
                putchar(curChar);
            }
        }
        prevChar = curChar; /* record current char as the previous char for next iteration */
    }
}
