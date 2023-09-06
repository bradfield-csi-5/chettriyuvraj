#include<stdio.h>
#include<stdlib.h>
#include<ctype.h>
#include<string.h>
#include<errno.h>
#include<unistd.h>
#include<sys/types.h>

#define MAXINPUT 50 /* Last \0 */
#define MAXARGS 4 /* Last NULL */
#define SHELLSYMBOL "\U0001f525"
#define FORKERRMSG "fork error"
#define COMMANDERRMSG "command does not exist"

char *testCommand = "pecho";
char *testArgv[] = {"my name is yuvi", "5"};

pid_t Fork(char *message);
int Execvp(char *command, char *argv[], char *errmsg);
void parseArgs(char *s, char *args[]);

int main() {
    char *s = NULL, *check = NULL;
    char input[MAXINPUT];
    char *args[MAXARGS];

    do {
        printf(SHELLSYMBOL);
        s = fgets(input, MAXINPUT, stdin);
        if (s) {
            // printf("%s", s);
            parseArgs(s, args);
            pid_t pid = Fork(FORKERRMSG);
            if (pid == 0) {
                Execvp(args[0], args, COMMANDERRMSG);
            }
        }
        if (s && !strchr(s, '\n')) { /* If newline not found, print first MAX chars and flush remaining */
            while ((s = fgets(input, MAXINPUT, stdin)) && !strchr(s, '\n'));
            printf("\n");
        }
    } while(s != NULL);
    


    if (!feof(stdin)) { /* If EOF is not set i.e error */
        printf("\n Error while reading input");
        return 1;
    }

    printf("\n%sIf this is to end in fire, we should all burn together %s", SHELLSYMBOL, SHELLSYMBOL);
}

/* Wrapper for fork() with error handling */
pid_t Fork(char *errmsg) {
    pid_t pid = fork();

    if (pid < 0) {
        printf("%s:%s", errmsg,strerror(errno));
    }
    return pid;
}

/* Wrapper for execvp(..) with error handling */
int Execvp(char *command, char *argv[], char *errmsg) {
    int resultStatus;
    if ((resultStatus = execvp(command, argv)) == -1) {
        printf("%s: %s", errmsg, strerror(errno));
    }
    return resultStatus;
}

/* Parse args from input string */
void parseArgs(char *s, char *args[]) {
    int argIdx = 0;

    for(;*s==' ';s++) ; /* Skip initial spaces */

    char *cur = s;
    while (*s != '\0') {

        while (*++cur != ' ' && *cur != '\n' && *cur != '\0'); /* Find delimiter */

        char *newArg = (char *)malloc(sizeof(char) * (cur - s + 1)); /* Malloc for arg when delimiter found, add to to arg array */
        strncpy(newArg, s, cur-s);
        newArg[cur-s] = '\0';
        args[argIdx++] = newArg;

        for (;*cur==' ' || *cur=='\n'; cur++); /* Skip additional spaces */

        s = cur;
    }
    args[argIdx] = NULL;
}

// cases ls a b c\0
// Idx = 0, cur = s = 'l', s != 0, cur skipped to idx 2, newArg[2+1], strnCpy(newArg, s, 2) i.e [l,s], newArg[2] = '\0', args[0++] = newarg;, cur = 3, s = 3
//  s = 3, cur skipped to 4, newArg[1+1], strncpy(newArg, s, 1), newArg[1] = '\0', args[1++] = 'a', args[2] =  NULL