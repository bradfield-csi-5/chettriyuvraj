/* Write the program expr , which evaluates a reverse Polish expression from the command line, where each operator or operand is a separate argument. For example,
expr 2 3 4 + * evaluates 2 x (3+4). */

/*
1. Reuse push and pop from Chapter 4 (stack.c) -> Using common header file util.h
2. Pass command line arguments to stack
    - If first digit is a number -> parse whole number (exit gracefully if any invalid char occurs)
    - Else evaluate as an operator -> ignore everything but the first char (exit gracefully if invalid operator)
3. When arguments over print stack top value as result

Notes:
1. For anything other than numbers (i.e operator expected) -> only the first char is checked regardless
2. Popping without any element in the stack leads to 0 being returned (So invalid reverse polish notation will lead to unexpected results)
3. Non-existent operators (except if they are a string with a valid first char) and numbers with invalid chars in between will both lead to termination
5. Operators like * and / lead to other args (eg entire directory name) being passed -> have to resolve
*/

#include<stdio.h>
#include<ctype.h>
#include<stdlib.h>
#include<math.h>
#include "util.h"

int main(int argc, char *argv[]) {
    double opLeft, opRight;
    char **p = argv;
    for(int i = 0; i<=argc; i++) {
        printf("%sl\n", *p++);
    }
    while (--argc > 0) { /* Include NULL character at the end as a signal to terminate */
        if (isdigit((*++argv)[0])) { /* If current val is a number */
            char *err;
            double curNum = strtod(*argv, &err);
            if (*err != '\0') {
                printf("Invalid number\n");
                return -1; /* Failure */
            }
            printf("Curnum %f\n", curNum);
            push(curNum);
        } else { /* Operator -> Only checking first char */
            printf("%c\n", (*argv)[0]);
            switch((*argv)[0]) {
                case '+':
                    push(pop() + pop());
                    break;
                case '*':
                    push(pop() * pop());
                    break;
                case '-':
                    opRight = pop();
                    opLeft = pop();
                    push(opLeft - opRight);
                    break;
                case '/':
                    opRight = pop();
                    opLeft = pop();
                    push(opLeft/opRight);
                    break;
                case '%':
                    opRight = pop();
                    opLeft = pop();
                    push(fmod(opLeft,opRight));
                    break;
                default:
                    printf("Invalid operator\n");
                    return -1;
            }
        }
    }

    printf("Result is %f\n", pop());
}


