// Write a function expand(s1,s2) that expands shorthand notations 
// like a-z in the string s1 into the equivalent complete list abc…xyz in s2. 
// Allow for letters of either case and digits, and be prepared to handle cases 
// like a-b-c and a-z0-9 and -a-z. Arrange that a leading or trailing -is taken literally.


/* 
1. Whenever we see a -,hold off printing it until further verification
2. -If prevChar was '-', check if prevPrevChar is 'matching' with curChar i.e valid range for explansion 
    1. If yes expand <dash><char>
    2. If not print - and char (invalid range so we print values as-is )
    
    -Else print curChar
 */

#include<stdio.h>
#include<ctype.h>

int isRangeValid(char left, char right);

int main() {
    int curChar = '\0', prevChar = '\0', prevPrevChar = '\0';

    while((curChar = getchar()) != EOF) {
        if (prevChar == '-' && isalnum(curChar)) { /* If prevChar is a hyphen and curChar is a valid letter or digit */
            if (isRangeValid(prevPrevChar, curChar)) { /* If valid range for expansion - expand*/
            
                for(char i = prevPrevChar + 1; i <= curChar; i++) /* Expand from prevPrevChar + 1  */
                    printf("%c",i);
    
            } else { /* Simply print hyphen and curChar */
                printf("-%c", curChar);
            }
        } else if (curChar != '-') { /* Hold off printing hyphen until next char is checked for expansion - print all other chars*/
            printf("%c", curChar);
        }
        prevPrevChar = prevChar;
        prevChar = curChar;
    }



}

/* Checks if range specified for expansion valid: a-9 or z-a would be an example of an invalid range*/
int isRangeValid(char left, char right) {
    if (left < right && (islower(left) && islower(right)) || (isupper(left) && isupper(right)) || (isdigit(left) && isdigit(right)))
        return 1;

    return 0;
}
