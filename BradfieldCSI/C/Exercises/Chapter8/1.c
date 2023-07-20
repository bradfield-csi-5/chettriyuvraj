/*Rewrite the program cat from Chapter 7 using read, write, open and close instead of their standard library equivalents.
Perform experiment to determine the relative speeds of the two versions*/

/*Scope: 
- If no filename, then keep copying standard output
- Print error to stderr if error in opening file
- Release file descriptor one by one as soon as copying done
*/

/*
Plan:
- Protoype for function to copy input to output
- If no args exist copy stdin to stdout
- Else create cat until arguments exist, open files using unix funcs and close using them (release fdescriptior + handle errors)
*/

#include<stdio.h>
#include<stdlib.h>
#include<fcntl.h> // for open file modes
#include<unistd.h> // for open system call

#define BUFSIZE 1024

int fcopy(int fd1, int fd2);

int main(int argc, char *argv[]) {
    int fd;
    int curFileIndex = 0; // currently opened file index in argv

    if (argc == 1) { // copy stdin to stdout
        fcopy(0,1);
    } else {
        while (--argc > 0) { // open files one after the other + handle errors if req
            
            fd = open(*++argv, O_RDONLY, 0);

            if (fd == -1) { // if file does not exist or some other error
                fprintf(stderr, "Unable to open filename %s\n", argv[curFileIndex]);
                continue;
            }

            if (fcopy(1, fd) < 0) { // cat file to stdout if no errors - else state error to stderr
                fprintf(stderr, "Unable to concatenate %s to standard output \n", argv[curFileIndex]);
            }

            if (fd > 0) //release file descriptior
                close(fd);

            curFileIndex++; // -> apparently 0th argument in my implementation is not the filename -> argv[0] is the first arg itself - so increasing it at the end
        }
    }


    return 0;
}

/*Copy fd2 to fd1  */
int fcopy(int fd1, int fd2) {
    int n;
    char buf[BUFSIZE];

    // Both files opened successfully
    while ((n = read(fd2, buf, BUFSIZE)) > 0) { // As long as input being read
        if (write(fd1, buf, n) != n) // write error
            return -1;
    }

    if (n < 0) { // read error
        return -1;
    }

    return 0;
}