#include<stdio.h>
#define MAXVAL 20
int stack[MAXVAL];
int sp = -1;

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
    printf("\n");
}

double pop(void) {
    printf("pop\n");
    for (int i = 0; i <= sp; i ++) {
        printf("%d ", stack[i]);
    }
    printf("\n");
    if (sp == -1) {
        printf("Stack is empty \n");
        return 0.0;
    }
    else
        return stack[sp--];
}
