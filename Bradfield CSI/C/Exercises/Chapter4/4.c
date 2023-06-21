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

#define NUMBER '0'

#define MAXVAL 20
void push(double val);
double pop(void);
int stack[MAXVAL];
int sp = -1;

#define BUFMAX 20
int bp = -1;
char buf[BUFMAX];
char getch(void);
void ungetch(char c);

double atof(char s[]);
int getops(char s[]);


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

int getops(char s[]) {
    int i = 0, c;
    while ((s[0] = c = getch()) == ' ' || c == '\t');
    s[1] = '\0';

    i = 0;
    if (!isdigit(c) && c != '.' && c != '-') /* Neither digit nor decimal nor negative sign which is handled separate */
    {
        return c;
    }

    if (c == '-') {
        c = getch();
        ungetch(c); /* We have to ungetch the char regardless of whether it is a digit or not - for further processing for next getch() */
        if (!isdigit(c)) /* Is not a negative number, but an operand */
            return '-';
    }

    
    if (isdigit(c)) /* Get integer part */
        while (isdigit(s[++i] = c = getch()));
    
    if (c == '.') /* Get fractional part */
        while (isdigit(s[++i] = c = getch()));
    
    s[i] = '\0';

    if (c != EOF)
        ungetch(c);
    return NUMBER;
}

void push(double val) {

    printf("\n");
    if (sp == MAXVAL - 1)
        printf("Stack is full\n");
    else
        stack[++sp] = val;

    printf("pushed\n");
    for (int i = 0; i <= sp; i ++) {
        printf("%d ", stack[i]);
    }
}

double pop(void) {
    printf("pop\n");
    for (int i = 0; i <= sp; i ++) {
        printf("%d ", stack[i]);
    }
    printf("\n");
    if (sp == -1)
        printf("Stack is empty \n");
    else
        return stack[sp--];
}

char getch(void) {
    return bp == -1 ? getchar() : buf[bp--];
}

void ungetch(char c) {
    if (bp == BUFMAX - 1)
        printf("Buffer max size exceeded \n");
    else
        buf[++bp] = c;
}

double atof(char s[]) {
    printf("atof %s\n", s);
    double sum = 0, pow = 1.0;
    int i, sign;
    
    for (i = 0; isspace(s[i]); i++); /* Skip whitespace */

    sign = s[i] == '-'? -1 : 1;
    if (sign == -1)
        i++;


    for (; isdigit(s[i]); i++) /* Add integer part */
        sum = sum * 10 + (s[i] - '0');


    if (s[i++] == '.')
        for (; isdigit(s[i]); i++) /* Add fractional part */
        {
            sum = sum * 10 + (s[i] - '0');
            pow *= 10.0;
        }
    
    return sign * sum/pow;



}

