/*Write a program to determine ranges of char, short, int and long variables, both signed and unsigned 
using standard headers and by direct computation*/

#include<stdio.h>
#include<limits.h>
#include<float.h>

int main() {

    /* Using headers */
    printf("The max value of char is %d \n", CHAR_MAX);
    printf("The min value of char is %d \n", CHAR_MIN);
    printf("The max value of short is %d \n", SHRT_MAX);
    printf("The min value of short is %d \n", SHRT_MIN);
    printf("The max value of unsigned int is %u \n", UINT_MAX); /* unsigned */
    printf("The max value of long is %ld \n", LONG_MAX);
    printf("The min value of long is %ld \n", LONG_MIN);
    printf("The max value of float is %.12e \n", FLT_MAX); /* floats */
    printf("The min value of float is %.12e \n", FLT_MIN);

    /* Computing */
    char c = 1; /* char */
    while (c > 0) {
        c <<= 1;
        printf("%d \n", c);
    }
    printf("The max value of char is %d \n", -c - 1);
    printf("The min value of char is %d \n", c);
    printf("The max value of unsigned char is %d \n", -c * 2 - 1);

    int i = 1; /* int */
    while (i > 0) {
        i <<= 1;
        printf("%d \n", i);
    }

    printf("The max value of int is %d \n", -i - 1);
    printf("The min value of int is %d \n", i);
    printf("The max value of unsigned int is %lu \n", (unsigned long)(-i) * 2 - 1); // Converting to unsigned long to prevent overflow when making i -> -i



    return 0;
}