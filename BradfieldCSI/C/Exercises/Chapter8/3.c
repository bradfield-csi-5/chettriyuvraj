/* Design and write _flushbuf, fflush and fclose */

/*
Plan:
_flushbuf
1. Define a structure (typedef) replicating 'FILE' with requisite fields (filemode, buffer etc)
2. Define any constants required (file modes et al)
3. Function checks all cases for file and depending on situation handles it
    - if file open for R/W with no errors
    - if no buffer allocated, allocate
    - write out entire buffer to file
    - set struct variables + add character provided to buffer as the first char
    - no check for unbuf writing 
*/

#include<unistd.h>
#include<stdlib.h>

#define BUFSIZE 1024
#define EOF (-1)

typedef struct _iobuf {
    int cnt; // chars left in buffer -> this indicates that 'cnt' number of characters can still be written to buffer 
    char *ptr; // next char pos
    char *base; // start of buffer
    int flag; // file mode
    int fd; // file descriptior
} FILE;

enum _flags {
    _WRITE = 2,
    _UNBUF = 4,
    _EOF = 8,
    _ERR = 16
};

int _flushbuf(int c, FILE *fp);

int main() {
    FILE *testFile = (FILE *) malloc (sizeof(FILE));

    char buf [] = {'a','b','c'}; // setting fields for testing
    testFile -> cnt = BUFSIZE - 3; // no of chars which can be written to buffer
    testFile -> base = buf;
    testFile -> ptr = buf + 2; // current char i.e latest written to buffer
    testFile -> fd = 1;
    testFile -> flag |= _WRITE; // set as open for write
    
    _flushbuf('y', testFile);
}

int _flushbuf(int c, FILE *fp) {

    int n; // to write buffer to file
    int charsToWrite; // number of chars to be emptied from buffer to file

    if ((fp -> flag & (_WRITE | _EOF | _ERR)) != _WRITE) // check if file is open for write with no errors
        return EOF;
    
    if (fp -> base == NULL) { // if no buffer allocated at all
        if ((fp -> base = (char *) malloc(BUFSIZE)) == NULL) { // unable to allocate for some reason
            fp -> flag |= _ERR; // some error or not all bytes written
            return EOF;
        }
        fp -> ptr = fp -> base;
    }
    

    charsToWrite = (fp -> ptr) - (fp -> base) + 1;
    if ( ((n = write(fp -> fd, fp -> base, charsToWrite)) < charsToWrite) || (n == -1) ) {  // write out entire buffer to file
        fp -> flag |= _ERR; // some error or not all bytes written
        return EOF;
    }

    // set all struct variables correctly and add the provided character as first char to buffer
    fp -> cnt = BUFSIZE; // indicates that 'cnt' chars can be stored in buffer
    *fp -> base = c;
    fp -> ptr = fp -> base ++;
    fp -> cnt--;


    return 0;
    
}

