#include<stdio.h>

char *strnCpy(char *d, char *s, int n);

int main() {
    char d[] = "Heybaby";
    char s[] = "ono";

    char *dStart = strnCpy(d, s, 4);

    printf("%s", dStart);
    return 0;
}


char *strnCpy(char *d, char *s, int n) {
    int i = 0;
    char *dStart = d;
    while (i++ < n && (*d++ = *s++));

    while (*d && i++< n)
        *d++ = '\0';
    
    return dStart;

}