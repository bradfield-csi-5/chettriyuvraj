/*Program that prints one word per line*/
/*A word loosely defined as collection of characters not separated by space, tab or newlines*/

#define IN 1
#define OUT 0

#include<stdio.h>

main() {
    int state = OUT;
    int c;
    while ((c=getchar()) != EOF) {
        if (state == IN) { //in the midst of iterting through a word
            if (c != ' ' && c != '\t' && c != '\n') {
                putchar(c);
            } else { //word has ended, change state
                state = OUT;
                putchar('\n');
            }
        } else { //in the midst of iterating through emptiness
            if (c != ' ' && c != '\t' && c != '\n') {
                //print first char of word and change state
                putchar(c);
                state = IN;
            }
        }
    }
}

