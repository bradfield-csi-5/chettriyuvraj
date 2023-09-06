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
#define SHELLPROFILE "./.shellprofile"

char *testCommand = "pecho";
char *testArgv[] = {"my name is yuvi", "5"};

pid_t Fork(char *message);
int Execvp(char *args[], char *errmsg);
void parseArgs(char *s, char *args[]);
int builtin(char *args[]);
void alias(char *args[]);
void aliasPrintAll(char *args[]);

int main() {
    char *s = NULL, *check = NULL;
    char input[MAXINPUT];
    char *args[MAXARGS];

    do {
        printf(SHELLSYMBOL);
        s = fgets(input, MAXINPUT, stdin);
        if (s) {
            parseArgs(s, args);
            if (!builtin(args)) {
                pid_t pid = Fork(FORKERRMSG);
                if (pid == 0) {
                    Execvp(args, COMMANDERRMSG);
                }
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

    printf("\n%sIf this is to end in fire, we should all burn together%s\n", SHELLSYMBOL, SHELLSYMBOL);
    return 0;
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
int Execvp(char *args[], char *errmsg) {
    int resultStatus;
    if ((resultStatus = execvp(args[0], args)) == -1) {
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

/* Execute and return 1 if builtin, 0 otherwise */
int builtin(char *args[]) {
    int returnStatus = 0;

    if (strcmp("alias", args[0]) == 0) {
        alias(args);
        returnStatus = 1;
    } else if (strcmp("exit", args[0]) == 0) {
        exit(0);
    }

    return returnStatus;
}

/*   Alias -p and alias without arguments handled for now  */

void alias(char *args[]) {
    FILE *f;
    int curChar;
    int areAllPrinted = 0; /* Multiple -p flags will lead to all alias-es being printed once */
    int isMatch; /* Used for overwriting aliases value */
    char *valueOffsetArg;
    char *matchStr;
    fpos_t *valueOffset;

    if (*(args+1) == NULL) { /* No arguments provided to alias - print all aliases */
        aliasPrintAll(args);
        return;
    }

    // if ((f = fopen(SHELLPROFILE, "a+")) == NULL) {
    //     printf("alias error: %s", strerror(errno));
    //     return;
    // }
    


    while (*++args != NULL)  { /* Handle each arg independently  */
        char *curStr = *args;

        if(strcmp(curStr, "-p") == 0 && !areAllPrinted) { /* If -p encountered and all alias-es not already printed once */
            areAllPrinted = 1;
            aliasPrintAll(args);
            continue;
        }
        else { // only for now
            continue;
        }
        
        /* handling arguments other than -p --------INCOMPLETE --------------------------------------------------------------------------- */
        if((valueOffsetArg = strrchr(curStr, '=')) != NULL) { /* If assignment statement found, search in file and replace */
            f = fopen(SHELLPROFILE, "r+");
            curChar = '\0';

            while (curChar != EOF) {
                isMatch = 1;
                matchStr = curStr; /* Compare char by char */
                
                while((curChar = fgetc(f)) != '\n' && curChar != EOF) {
                    if (curChar == '=') {
                        fgetpos(f, valueOffset); /* Record start of string position; TODO: handle error */
                    }
                    if (*matchStr++ != curChar){ /* If mismatch at any point, co */
                        isMatch = 0;
                        break;
                    }
                }

                if (isMatch) { /* Start writing from valueOffsetIdx+1 in curStr to valueOffset in file */
                    fsetpos(f, valueOffset);
                    for (char *s = valueOffsetArg + 1; *s != '\0'; s++) {
                        fputc(*s, f); /* TODO: handle err */
                    }
                    fseek(f, 0, SEEK_END); /* TODO handle err */
                } else { /* Iterate till EOF/EOL */
                    while((curChar = fgetc(f)) != '\n' && curChar != EOF);
                }
            }

        } else { /* Search line wise */

        }

        /* handling arguments other than -p --------INCOMPLETE --------------------------------------------------------------------------- */

    }

}


void aliasPrintAll(char *args[]) {
    FILE *f;
    int curChar;
    if ((f = fopen(SHELLPROFILE, "r")) == NULL) {
        printf("alias error: %s", strerror(errno));
        return;
    }
    while ((curChar = fgetc(f)) != EOF){ /* Make sure \n exists in file itself */
        fputc(curChar, stdout);
    }
    return;
}