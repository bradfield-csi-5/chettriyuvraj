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

fflush
1. Basically reuse _flushbuf and pass EOF as the character to clear the buffer

fclose(FILE *fp)
reverse the actions of fopen
1. check file slot arr (_iob) to see which one the file occupies, if yes reset it
2. flush buffers using fflush
3. set flag back to 0
3. call close() to release file descriptor
4. reset the file descriptor to -1
*/

#include<unistd.h>
#include<stdlib.h>
#include<fcntl.h>


#define BUFSIZE 1024
#define EOF (-1)
#define OPEN_MAX 6

typedef struct _iobuf {
    int cnt; // chars left in buffer -> this indicates that 'cnt' number of characters can still be written to buffer 
    char *ptr; // next char pos
    char *base; // start of buffer
    int flag; // file mode
    int fd; // file descriptior
} FILE;

enum _flags {
    _READ = 1,
    _WRITE = 2,
    _UNBUF = 4,
    _EOF = 8,
    _ERR = 16
};

FILE _iob[OPEN_MAX] = {
    {0, (char *) 0, (char *) 0, _READ, 0}, // stdin
    {0, (char *) 0, (char *) 0, _WRITE, 1}, // stdout
    {0, (char *) 0, (char *) 0, _WRITE, 2}, // stderr
    {3, (char *) 0, (char *) 0, _WRITE, 3} // can use open() system call and assign file descriptor returned to this entry - sample entry to test fclose()  
};

int _flushbuf(int c, FILE *fp);
int fflush(FILE *fp);
int fclose(FILE *fp);

int main() {
    FILE *testFile = (FILE *) malloc (sizeof(FILE));

    //test _flushbuf
    char buf [] = {'a','b','c'}; // setting fields for testing
    testFile -> cnt = BUFSIZE - 3; // no of chars which can be written to buffer
    testFile -> base = buf;
    testFile -> ptr = buf + 2; // current char i.e latest written to buffer
    testFile -> fd = 1;
    testFile -> flag |= _WRITE; // set as open for write
    _flushbuf('y', testFile);

    // test fclose
    fclose(_iob + 3);
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
    if ( (fp->ptr != fp->base) && ((n = write(fp -> fd, fp -> base, charsToWrite)) < charsToWrite) || (n == -1) ) {  // write out entire buffer to file
        fp -> flag |= _ERR; // some error or not all bytes written
        return EOF;
    }

    // set all struct variables correctly and add the provided character as first char to buffer
    fp -> cnt = BUFSIZE; // indicates that 'cnt' chars can be stored in buffer
    if (c != EOF) 
        *fp -> base = c;
    fp -> ptr = fp -> base ++;
    fp -> cnt--;


    return 0;
    
}

int fflush(FILE *fp) {
    return _flushbuf(EOF, fp);
}


int fclose(FILE *fp) {
    FILE *f;
    for(f = _iob; f < _iob + OPEN_MAX; f++) { // check if file occupies a slot (amongst files to be opened)
        if (f -> fd == fp -> fd)
            break;
    }

    if (f >= _iob + OPEN_MAX) // file doesn't seem to occupy a slot - no descriptor matches
        return EOF;
    
    if (fflush(fp) == EOF) // flush buffer
        return EOF;
    
    fp -> flag = 0; // set flags to 0

    if (close(fp -> fd) == -1) // release file descriptor -> since we don't have an actual fd assigned - at this point function will return EOF
        return EOF;

    fp -> fd = -1;

    return 0;


}