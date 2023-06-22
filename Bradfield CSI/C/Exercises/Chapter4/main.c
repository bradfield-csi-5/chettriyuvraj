/* Calculator - reverse polish notation

Functions:
getops()
getch()
ungetch()
push()
pop()
*/

#include<stdio.h>
#include<ctype.h>
#include "calc.h"

#define MAXVAL 20

int main() {
    int c;
    double val, op2;
    char s[MAXVAL];

    while((c = getops(s)) != EOF) {
        switch(c) {
            case NUMBER:
                push(atof(s));
                break;
            case '+':
                val = pop() + pop();
                push(val);
                break;
            case '*':
                val = pop() * pop();
                push(val);
                break;
            case '-':
                op2 = pop();
                val = pop() - op2;
                push(val);
                break;
            case '/':
                op2 = pop();
                val = pop()/op2;
                push(val);
                break;
            case '%':
                op2 = pop();
                val = (int)pop()%(int)op2;
                push(val);
                break;

            case '\n':
                printf("Result is %.8g", pop());
                break;

            default:
                printf("INVALID!!! \n");
                break;
        }
    }



}




