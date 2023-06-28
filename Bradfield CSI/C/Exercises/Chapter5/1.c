#include<stdio.h>
#include<ctype.h>
#include "util.h"

int getint(int *p);
#define ARRSIZE 5
#define NOTANUMBER 'a'

int main() {
    int arr[ARRSIZE];
    for (int i = 0; i < ARRSIZE && getint(&arr[i]) != EOF; i++ );

    for (int i = 0; i < ARRSIZE; i++)
        printf("%d \n", arr[i]);
    return 0;
}

int getint(int *p) {
    int c, sign;

    while (isspace(c=getch())); /* Remove whitespaces */

    if (!isdigit(c) && c != EOF && c != '+' && c != '-') { /* Not a number */
        ungetch(c);
        return 0;
    }

    sign = c == '-' ? -1 : 1;
    if (c == '+' || c == '-')  /* Move to next elem if sign */
    {
        c = getch();
        if (!isdigit(c)) { /* If sign not followed by a digit - not a number, push it back to input */
            if (c != EOF){
                ungetch(NOTANUMBER); /* ungetch NOT A NUMBER to terminate operation similar to what we did earlier for number - using NOTANUMBER to account for blank space after sign */
            }
                
            return 0;
        }
    }


    for (*p = 0; isdigit(c); c = getch())
        *p = 10 * *p + (c - '0');
    
    *p *= sign;

    if (c != EOF)
        ungetch(c);
    
    return c;

}