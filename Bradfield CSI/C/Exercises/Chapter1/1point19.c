/*Write a function s that reverses a character string s - use it to reverse one line at a time*/
/*
Plan:
1. function reverse(s, strLen) -> utilize length of string computed while taking input then use reverse algo to reverse it -> assume last char is \0
2. function getline() to take input and handle edge cases \n + terminate input using \0
3. main() - declare maxlength, take input, pass input to reverse and print the soln
*/


#include<stdio.h>

void reverse(char s[], int strLen);
int getLine(char s[], int limit);

#define MAXLENGTH 10

int main() {
    char input[MAXLENGTH];
    int len;

    while ((len = getLine(input, MAXLENGTH)) > 0) {
        // strings that exceed limit - leftover part is sent to next iteration of input and reverse
        reverse(input, len);
        printf("%s", input);
        printf("\n");
    }

    return 0;
}

int getLine(char s[], int limit) {
    int i, c;
    for (i = 0; i < limit - 1 && (c = getchar()) != EOF && c != '\n'; ++i){
        s[i] = c;
    }
    if (c == '\n') {
        s[i] = c;
        ++i;
    }
    s[i] = '\0';
    return i; // i is the length of the string, i + 1th char contains \0
}

void reverse(char s[], int strLen) {
    for (int l = 0, r = strLen - 1; l < r; ++l,--r) {
        char temp = s[l];
        s[l] = s[r];
        s[r] = temp;
    }
}

