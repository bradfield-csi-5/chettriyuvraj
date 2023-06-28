/* Simply allocates lines to a pointer array and prints them -> instead of implementing quicksort */


#include<stdio.h>
#include<string.h>

#define ALLOCBUFSIZE 20 /* Size of alloc buffer */
#define MAXLINES 10 /* Max number of lines */
int getLine(char line[], int lim);
int readLines( char *lineptr[], int maxlines, char *allocbuf);
int getLine(char line[], int lim);

int main() {
    int nlines;
    char allocbuf[ALLOCBUFSIZE];
    char *lineptr[MAXLINES];
    int i;

    if ((nlines = readLines(lineptr, MAXLINES, allocbuf))){ /* Prints even if lines exceeded */
        nlines = nlines > 0 ? nlines : MAXLINES;
        for (int i = 0; i < nlines; i++) {
            printf("%s \n", lineptr[i]);
        }
        
    }
    
}

#define MAXLEN 10 /* Max line length */

/* Returns -1 if too much input is presented, else number of lines - but takes in input within bounds */
int readLines( char *lineptr[], int maxlines, char *allocbuf) {
    char line[MAXLEN];
    int nlines = 0, lineLength;
    char *bufp = allocbuf, *p;


    while (nlines < maxlines && (lineLength = getLine(line, MAXLEN)) > 0 && allocbuf + ALLOCBUFSIZE >= bufp + lineLength + 1) { /* Also counts \0 character */
        line[lineLength - 1] = '\0';
        p = bufp;
        strcpy(p, line);
        lineptr[nlines++] = p;
        bufp += lineLength;
    }

    if (nlines >= MAXLINES) /* Number of lines exceeeded limit */
        return -1;
    
    return nlines;

}


int getLine(char line[], int lim) {
    int i = 0, c;

    while ( i < lim - 1 && (c = getchar()) != EOF && c != '\n') /* considering if EOF or \n might be at the end of getchar()  */
        line[i++] = c;

    if (c == '\n')
        line[i++] = c;
    
    line[i] = '\0';
    
    return i;
}