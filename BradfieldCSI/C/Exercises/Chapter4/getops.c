#include<stdio.h>
#include<ctype.h>
#include "calc.h"


int getops(char s[]) {
    int i = 0, c;
    while ((s[0] = c = getch()) == ' ' || c == '\t');
    s[1] = '\0';

    i = 0;
    if (!isdigit(c) && c != '.' && c != '-') /* Neither digit nor decimal nor negative sign which is handled separate */
    {
        return c;
    }

    if (c == '-') {
        c = getch();
        ungetch(c); /* We have to ungetch the char regardless of whether it is a digit or not - for further processing for next getch() */
        if (!isdigit(c)) /* Is not a negative number, but an operand */
            return '-';
    }

    
    if (isdigit(c)) /* Get integer part */
        while (isdigit(s[++i] = c = getch()));
    
    if (c == '.') /* Get fractional part */
        while (isdigit(s[++i] = c = getch()));
    
    s[i] = '\0';

    if (c != EOF)
        ungetch(c);
    return NUMBER;
}


double atof(char s[]) {
    printf("atof %s\n", s);
    double sum = 0, pow = 1.0;
    int i, sign;
    
    for (i = 0; isspace(s[i]); i++); /* Skip whitespace */

    sign = s[i] == '-'? -1 : 1;
    if (sign == -1)
        i++;


    for (; isdigit(s[i]); i++) /* Add integer part */
        sum = sum * 10 + (s[i] - '0');


    if (s[i++] == '.')
        for (; isdigit(s[i]); i++) /* Add fractional part */
        {
            sum = sum * 10 + (s[i] - '0');
            pow *= 10.0;
        }
    
    return sign * sum/pow;



}
