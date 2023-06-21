/* 	Write the function strrindex(s,t) , which returns the position of the 
rightmost occurrence of t in s , or -1 if there is none. */

#include<stdio.h>
#define MAXLENGTH 200/* Max line length */
int getLine(char s[], int lim);
int strrindex(char s[], char t[]);

int main() {
    char s[MAXLENGTH];
    char t[] = "a";
    
    while (getLine(s, MAXLENGTH)) {
        /* printf("Enter pattern to search for\n"); Avoid inputting pattern using getLine -> it will always have a '\n' at the end which will not match*/ 
        /* getLine(t, MAXLENGTH); */
        printf("The rightmost occurrence of %s in %s is at %d", t, s, strrindex(s, t)); //index %d", strrindex(s, t))
    }

    return 0;
}

int getLine(char s[], int lim) {
    int i = 0;
    int c;
    while (i < lim - 1 && (c = getchar()) != EOF && c != '\n') {
        s[i++] = c;
    }
    if (c == '\n')
        s[i++] = c;
    s[i] = '\0';
    
    return i;
}

/* NOTE: striindex in the book is much more succinct - can be converted into this with 1 additional variable */
int strrindex(char s[], char t[]) {
    int i = 0, j = 0, matchIndex = -1, curMatchIndex = 0;

    /* Either of the two is empty */
    if (s[i] == '\0' || t[j] == '\0') {
        return -1;
    }

    while (s[i] != '\0') {
        if (s[i] == t[j]) { /* chars match */
            ++i;
            ++j;
            if (t[j] == '\0') { /* If complete match found */
                matchIndex = curMatchIndex;
                ++curMatchIndex;
                i = curMatchIndex; /* Start matching i again from next index */
                j = 0;
            }
        } else if (j != 0) { /* chars dont match and j not equal to 0 -> move j to start */
            j = 0;
        } else { /* chars dont match and j is equal to 0 -> increase i */
            ++curMatchIndex;
            i = curMatchIndex;
        }
    }

    return matchIndex;

}

