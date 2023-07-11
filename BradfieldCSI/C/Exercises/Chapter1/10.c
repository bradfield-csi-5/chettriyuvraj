#include<stdio.h>
/*Program to copy input to output, replacing each tab by \t, backspace by \b and backslash by \\*/
main() {
    int curChar;
    while ((curChar = getchar()) != EOF) {
        if (curChar == '\t') {
            printf("\\t");
        } else if (curChar == '\b') {
            printf("\\b");
        } else if (curChar == '\\' ) {
            printf("\\\\");
        } else {
            putchar(curChar);
        }
    }
    
}