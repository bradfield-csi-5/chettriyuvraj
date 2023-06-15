#include<stdio.h>

main() {
    /* Program to convert Celsius to Fahrenheit - with a header */


    int lower, upper, step;
    float celsius, fahr; //formula c *(9/5) + 32

    lower = 0;
    upper = 300;
    step = 13;
    celsius = lower;

    printf("%30s", "Celsius to Fahrenheit Table\n\n");

    while (celsius <= upper) {
        fahr = celsius * (9.0/5.0) + 32;
        printf("%10.0f %6.2f\n", celsius, fahr);
        celsius += step;
    }
    
    
}