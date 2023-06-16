/* Write a function htoi(s) to convert a string of hexadecimal characters into equivalent intger value */

/* Test case:

AAf41 - returns 700225

*/

#include<stdio.h>
#include<string.h>

 int htoi(char s []);

 int main() {
    int curChar;
    printf("%d", htoi("AAF41"));
    return 0;
 }

 int htoi(char s []) {
    int val = 0;
    for (int i = 0; i < strlen(s); i++) {
        int curChar = s[i];
        int charVal = curChar >= 'A' && curChar <= 'F' ? curChar - 'A' + 10 : (curChar >= 'a' && curChar <= 'f' ? curChar - 'a' + 10 : curChar - '0');
        val = (val * 16) + (charVal);
    }
    return val;
 }

