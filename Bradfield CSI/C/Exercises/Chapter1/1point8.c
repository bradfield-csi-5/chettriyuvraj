#include<stdio.h>

main() {
    printf("<--- Program to count blanks, tabs and newlines --->\n");
    printf("Please input your sentence:\n");
    int blankCount = 0, tabCount = 0, newlineCount = 0, charVal;
    while((charVal = getchar()) != EOF) {
        switch(charVal) {
            case '\n':
                ++newlineCount;
                break;
            case '\t':
                ++tabCount;
                break;
            case ' ':
                ++blankCount;
                break;
            default:
                continue;
        }
    }
    printf("The count of blanks, tabs and newlines is %d, %d, and %d respectively", blankCount, tabCount, newlineCount);
}