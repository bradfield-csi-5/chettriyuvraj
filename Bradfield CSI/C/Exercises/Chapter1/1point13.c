/*Program to print vertical histogram of words in input*/

#define IN 1
#define OUT 0
#define MAX_WORD_COUNT 100 

#include<stdio.h>

main() {
    int c;
    int maxWordLen = 0; // length of longest word
    int state = OUT; // to check boundaries of words
    int curWordLen = 0;
    int curIndex = 0;
    int wordLength[MAX_WORD_COUNT]; // store length of each word(duplicate words will have separate bars)

    // taking input
    while ((c = getchar()) != '\n') {
        if (state == OUT) { // currently iterating over emptiness
            if (c != '\n' && c != '\t' && c != ' ') { // new word encountered 
                ++curWordLen; // set curWord length to 1
                state = IN;
            }
        } else { // iterating over a word currently
            if (c != '\n' && c != '\t' && c != ' ') { // word continues
                ++curWordLen;
            } else { // word has ended
                if (maxWordLen < curWordLen) { // check if max word length to be changed
                    maxWordLen = curWordLen;
                }
                wordLength[curIndex] = curWordLen;
                ++curIndex;
                curWordLen = 0;
                state = OUT;
            }
        }
    }
    
    //if last word exists
    if (state == IN) {
        if (maxWordLen < curWordLen) { // check if max word length to be changed
            maxWordLen = curWordLen;
        }
        wordLength[curIndex] = curWordLen;
        ++curIndex;
    }


    // decoration
    printf("\n\n\n\n");
    for(int i = 0; i < curIndex; i++)
    {
        if (i == curIndex/2) {
            printf("HISTOGRAM----------");
        } else {
            printf("%*s", maxWordLen + 10, "----------");
        }
    }
    printf("\n\n");

    // printing length numbers
    for (int i = 0; i < curIndex; i++) { // iterate over array
        printf("%*d", maxWordLen + 10 , wordLength[i]);
    }

    printf("\n"); //Moving on to the histogram
    for (int curLevel = maxWordLen; curLevel > 0; curLevel--) {
        for (int i = 0; i < curIndex; i++) { // iterate over array each time
            int checkLevel = wordLength[i]; // compare level with current level
            if (checkLevel >= curLevel) {
                printf("%*s", maxWordLen + 10 , ".");
            } else {
                printf("%*s", maxWordLen + 10, "");
            }
        }
        printf("\n");
    }

    // decoration
    for(int i = 0; i < curIndex; i++) 
        printf("%*s", maxWordLen + 10, "----------");

}