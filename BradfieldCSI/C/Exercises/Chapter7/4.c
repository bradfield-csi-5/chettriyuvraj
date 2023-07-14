/* Write a private version of scanf analogous to minprintf from the previous section */
/* Scope: only int assignments -> no precision/format specifiers et al*/


#include<stdio.h>
#include<stdarg.h>
#include<ctype.h>
#include<stdlib.h>

int minScanf(char *fmt, ...); 
int assignInt(int *i);
int assignDouble(double *d); // not implmenented
int assignStr(char *s); // not implmenented

int main() {
    int a, b, c;

    int assigned = minScanf("%d %d %d", &a, &b, &c);

    printf("%d %d %d", a, b, c);

    return 0;
}

int minScanf(char *fmt, ...) {
    va_list ap; // to traverse through the variable args list - note this will contain only pointers
    char *t; // to traverse format string
    

    int *i; // 3 vars to store different types of pointers - char, double not implemented
    char *s;
    double *d;

    va_start(ap, fmt);// init ap to point to the format string

    int assignSuccess; // indicates if assignment to argument successful i.e type specifier matches value provided in input
    int assignCount = 0;


    for(t = fmt; *t; t++) {

        if (*t != '%') { // continue if a type specifier is not encountered
            putchar(*t);
            continue; 
         }
        
        switch(*++t) { //if type specifier encountered, find what type it is
            case 'd':
                i = va_arg(ap, int*); // obtain the argument value i.e pointer and step ap to next argument
                assignSuccess = assignInt(i);
                break;
        }

        assignCount += assignSuccess;
    }

    va_end(ap);

    return assignCount;
}


/*
Assigns int value to pointer argument provided by user in scanf
1. 1 if success
1. 0 if no success i.e invalid argument
*/

int assignInt(int *i) {
    int c;
    char intStr[10]; // max length of 32 bit integer string i.e 2^32 order
    char *p = intStr;

    while ((c = getchar()) != EOF && isspace(c)); // removing all whitespace chars
    ungetc(c, stdin); // final char pushed back into stream else it would be missed

    while ((c = getchar()) != EOF && !isspace(c)) {
        if(isnumber(c)){ // valid numeric input
            *(p++) = c;
        } else { // invalid input - not number
            while ((c = getchar()) != EOF && !(isspace(c))); // discard entire input until space or EOF - no need to ungetc here as c will either be a space or EOF
            
            return 0; // unsuccessful assignment
        }
    }

    if (p == intStr) // no numbers found
        return 0;
    
    *p = '\0';
    *i = atoi(intStr); // successfully put integer into array intStr - convert it to number and put it in the argument pointer provided
    return 1;
}

