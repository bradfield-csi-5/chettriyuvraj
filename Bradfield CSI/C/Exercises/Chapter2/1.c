/*Write a program to determine ranges of char, short, int and long variables, both signed and unsigned 
using standard headers and by direct computation*/

#include<stdio.h>
#include<limits.h>
#include<float.h>

int main() {

    // Using headers - all signed
    printf("The max value of char is %d \n", CHAR_MAX);
    printf("The min value of char is %d \n", CHAR_MIN);
    printf("The max value of short is %hi \n", SHRT_MAX);
    printf("The min value of short is %hi \n", SHRT_MIN);
    printf("The max value of unsigned int is %u \n", UINT_MAX); // unsigned
    printf("The max value of long is %ld \n", LONG_MAX);
    printf("The min value of long is %ld \n", LONG_MIN);

    //floats
    printf("The max value of float is %.12e \n", FLT_MAX);
    printf("The min value of float is %.12e \n", FLT_MIN);

    unsigned int i = 1;
    unsigned int prevI = 0;
    while (prevI < i) {
        prevI = i;
        i *= 2;
        printf("%d \n", i);
    }




    return 0;
}