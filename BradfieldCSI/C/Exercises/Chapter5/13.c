/* Write the program tail, which prints the last n lines of its input.
By default, n is 10, say, but it can be changed by an optional argument, 
so that tail -n prints the last n lines. 
The program should behave rationally no matter how unreasonable the input or the value of n. 
Write the program so it makes the best use of available storage; lines should be stored as in the sorting program of Section 5.6, 
not in a two-dimensional array of fixed size. */


/* Note: Assuming flags assigned at the start always + any additional flags ignored */



#include <stdio.h>
#include <stdlib.h>
#define MAXNUMBEROFLINES 20
#define MAXLENGTH 20
#define DEFAULTN 5 /* Print last 5 lines by default */

char  *getLine(int lim);

int main(int argc, char *argv[]) {
    char *lineptr, *err, *lines[MAXNUMBEROFLINES], *linesCount;
    int n = DEFAULTN, lineCount = 0, lineIndex;

    if (--argc > 0) { /* Checking for flags - if any command line args exist at all */
        while ((*++argv)[0] == '-') { /* While flags exist + skipping default first filename argv */
            n = strtol(*argv + 1, &err, 10);
            if (*err != '\0'){
                printf("Invalid number of lines to print \n");
                return -1;
            }
            break; /* Once first flag is encountered all others ignored */
        }
    }

    
    while ((lineptr = getLine(MAXLENGTH)) != NULL) { /* Taking input lines */
        lines[lineIndex++] = lineptr;
        lineIndex %= n;
        lineCount++;
    }

    if (lineCount > lineIndex) /* If lineCount exceeds lineIndex, first print initial (n - lineIndex) lines */
        for (int i = lineIndex; i < n; i++)
            printf("%s", *(lines + i));
    
    for (int i = 0; i < lineIndex; i++) /* Print first 'lineIndex' lines */
        printf("%s", *(lines + i));

    return 0;
}

/* Allocates and returns a pointer to a new input string each time */
char *getLine(int lim) {
    char *s = (char *)malloc(lim * sizeof(char)); /* Allocating block for string so that it persists in main */
    int i = 0;
    int c;
    while (i < lim - 1 && (c = getchar()) != EOF && c != '\n') {
        s[i++] = c;
    }
    if (c == '\n')
        s[i++] = c;
    s[i] = '\0';
    
    return (i > 0) ? s : NULL;
}


