/* Modify the programs entab and detab (written as exercises in Chapter 1) to accept a list of tab stops as arguments. Use the default tab settings if there are no arguments. */



/*
DETAB!
1. Store the 'tabstop columns' as 'tabstops' by computing the difference between 'prevStop' and 'curStop' -> Eg. 5,8,11 gets stored as [5,3,3,8] where the last value is the default tab stop
2. Now start maintaining the charCount -> each time storing it as charCount % currentTabStop -> changing tab stops as and when we reach them -> keeping the last tab stop as the default once all others have been exhausted
PS. Print a long string of zeroes at the start for reference
*/

#define TABSTOP 8 /* Default tabstop */
#define COLUMNS 100

#include<stdio.h>
#include<stdlib.h>
int main(int argc, char *argv[]) {
    int curChar, curTabStop, prevTabStop = 0, charCount = 0, tabStopIndex = 0, tabStops[argc]; /* One extra index to put default tabstop */
    int argCount = argc; /* For later reference */
    char *err;

    /*Column bar for reference*/
    for(int i = 0; i < COLUMNS; i++)
        printf("0");
    printf("\n");

    /* Put non default tab stops in array */
    while (--argc > 0)
    { 	
        /* Convert to digit */
        curTabStop = (int)strtol(*++argv, &err, 10);
        printf("%d \n", curTabStop);
        if (*err != '\0') { /* error check */
            printf("Invalid number \n");
            return -1;
        }
        tabStops[tabStopIndex] = curTabStop - prevTabStop;
        prevTabStop = curTabStop;
        ++tabStopIndex;
    }
    tabStops[tabStopIndex] = TABSTOP; /* Consider tabstops of 5,8,11 and default 8 -> we have arr = [5,2,3,8] */
    tabStopIndex  = 0;

    // for (int i = 0; i < argCount; i++)
    //     printf("%d ", tabStops[i]);
    // printf("\n");


    /* Handle default tabstop */
    curTabStop = tabStops[tabStopIndex++];
    while((curChar = getchar()) != EOF) {
        if (curChar == '\n') { /* Move to the next line */
            charCount = 0;
            tabStopIndex = 0;
            curTabStop = tabStops[tabStopIndex];
        } else if (curChar == '\t') { /* Add remaning spaces to get to tab stop */
            int remainingSpaces = curTabStop - charCount;
            for(int i = 0; i < remainingSpaces; i++)
                printf(" ");
            charCount = 0;
            curTabStop = tabStops[tabStopIndex++];
        } else {  /* Any other character */
            putchar(curChar);
            ++charCount;
        }

        charCount %= curTabStop; /* Since we are counting the distance to the next tab stop */
        if (charCount == 0 && tabStopIndex < argCount - 1) { /* Move to the next tab stop if we are not on the last i.e default */
            charCount = 0;
            curTabStop = tabStops[tabStopIndex++];
        }
    }


    return 0;
}
