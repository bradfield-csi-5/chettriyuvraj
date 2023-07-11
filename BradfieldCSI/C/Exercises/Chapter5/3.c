#include<stdio.h>

void strCat(char *s, char *p);

int main() {
    char s[] = "Hey bro";
    char p[] = "";

    strCat(s,p);

    printf("%s", s);
    return 0;

}

void strCat(char *s, char *p) {
    while (*s != '\0')
        s++;

    while (*s++ = *p++);
}