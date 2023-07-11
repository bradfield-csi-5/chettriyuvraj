#include<stdio.h>

#define BUFMAX 20
int bp = -1;
char buf[BUFMAX];

char getch(void) {
    return bp == -1 ? getchar() : buf[bp--];
}

void ungetch(char c) {
    if (bp == BUFMAX - 1)
        printf("Buffer max size exceeded \n");
    else
        buf[++bp] = c;
}
