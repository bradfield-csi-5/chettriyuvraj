#include<stdio.h>
#include<stdlib.h>
#include<ctype.h>
#include<string.h>
#include<errno.h>
#include<unistd.h>
#include<sys/types.h>
#include<sys/wait.h>
#include<signal.h>

#define MAXINPUT 50 /* Last \0 */
#define MAXARGS 4 /* Last NULL */
#define SHELLSYMBOL "\U0001f525"
#define FORKERRMSG "fork error"
#define COMMANDERRMSG "command does not exist"
#define SHELLPROFILE "./.shellprofile"
#define SLEEPCOMMAND "sleep"

char *testCommand = "pecho";
char *testArgv[] = {"my name is yuvi", "5"};

pid_t Fork(char *message);
int Execvp(char *args[], char *errmsg);
struct token parseArgs(char *s, char *args[]);
int builtin(char *args[]);
void alias(char *args[]);
void aliasPrintAll(char *args[]);
void sigintHandler(int sig);
void sigtstpHandler(int sig);
void sigChldHandler(int sig);
char *getOperator(char *s);
int continueExec(char *operator, int exitStatus);

struct token { /* Call to parseArgs always returns a token */
    char *sNext;
    char *operator;
};

int shellPid;

int main() {
    char *s = NULL, *check = NULL;
    char input[MAXINPUT];
    char *args[MAXARGS];
    int status, exitStatus;
    struct token token = {
        "\0",
        "\0"
    };
    sigset_t mask, prevMask;
    shellPid = getpid();
    

    signal(SIGINT, sigintHandler); /* Installing SIGINT handler  */
    signal(SIGTSTP, sigtstpHandler); /* Installing SIGTSTP handler */



    do {

        printf(SHELLSYMBOL);

        if (*token.sNext != '\0' && continueExec(token.operator, exitStatus)) { /* If operators were used, continue execution based on exit status */
            s = token.sNext;
        } else {
            s = fgets(input, MAXINPUT, stdin);
        }

        if (s) {
            token = parseArgs(s, args); /* parseArg will always return a token depending on operators used in expression */
            signal(SIGTTOU, SIG_IGN);
            if (!builtin(args)) {
                pid_t pid = Fork(FORKERRMSG);
                if (pid == 0) {
                    pid_t childPid = getpid();
                    setpgid(childPid, childPid); /* We want to implement job control, so we will create separate pgid for each process */
                    if (*token.operator != '&' ) {
                        tcsetpgrp(STDIN_FILENO, childPid);
                    }
                    Execvp(args, COMMANDERRMSG);
                }

                
                
                
                
                if (*token.operator != '&' ) { /* Don't wait on background processes */
                    while (waitpid(pid, &status, 0) > 0);  /* Reaping child process and tracking last exitStatus*/
                    if (WIFEXITED(status)) {
                        exitStatus = WEXITSTATUS(status);
                    }
                    tcsetpgrp(STDIN_FILENO, getpgid(shellPid));
                }
            }
        }

        if (s && !strchr(s, '\n')) { /* If newline not found, print first MAX chars and flush remaining - why was this condition written (?)*/
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
struct token parseArgs(char *s, char *args[]) {
    int argIdx = 0;
    char *operator; /* To grab operators like && || */
    struct token token = { 
        "\0",
        "\0"
    };

    for(;*s==' ';s++) ; /* Skip initial spaces */

    char *cur = s;
    while (*s != '\0' && *(operator = getOperator(s)) == '\0') {

        while (*++cur != ' ' && *cur != '\n' && *cur != '\0'); /* Find delimiter */

        char *newArg = (char *)malloc(sizeof(char) * (cur - s + 1)); /* Malloc for arg when delimiter found, add to to arg array */
        strncpy(newArg, s, cur-s);
        newArg[cur-s] = '\0';
        args[argIdx++] = newArg;

        for (;*cur==' ' || *cur=='\n'; cur++); /* Skip additional spaces */

        s = cur;
    }

    if (*operator != '\0') { /* If stopped parsing due to operator, return a token containing operator and start of next string */
        token = (struct token) {
            s + strlen(operator), /* Might have additional spaces at the start, they will be handled by next call to parseArgs */
            operator,
        };
    }

    args[argIdx] = NULL;

    return token;
}

char *getOperator(char *s) {
    char *operator = "\0";
    if ((*s == '&' && *(s+1) == '&' && *(s+2) == ' ') || (*s == '|' && *(s+1) == '|' && *(s+2) == ' ')) {
        operator = (char *)malloc(3 * sizeof(char));
        strncpy(operator, s, 2);
        operator[2] = '\0';
    } else if ((*s == '|' && *(s+1) == ' ') || (*s == '&' && *(s+1) == ' ')) {
        operator = (char *)malloc(2 * sizeof(char));
        strncpy(operator, s, 1);
        operator[1] = '\0';
    }
    return operator;
}

/* Execute and return 1 if builtin, 0 otherwise */
int builtin(char *args[]) {
    int returnStatus = 0;

    if (strcmp("alias", args[0]) == 0) {
        alias(args);
        returnStatus = 1;
    } else if (strcmp("exit", args[0]) == 0) {
        printf("\n%sIf this is to end in fire, we should all burn together%s\n", SHELLSYMBOL, SHELLSYMBOL);
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
        
        /* For handling arguments other than -p */
        /* ------------------------------------------------------INCOMPLETE --------------------------------------------------------------------------- */
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

        /* ------------------------------------------------------INCOMPLETE --------------------------------------------------------------------------- */

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

void sigintHandler(int sig) {
    printf("\n Caught SIGINT!");
    return;
}

void sigtstpHandler(int sig) {
    printf("\n Caught SIGTSTP - setting controlling terminal to shell");
    tcsetpgrp(STDIN_FILENO, getpgid(shellPid));
    exit(0);
}


/* Determine whether to continue execution or not depending on exitStatus and operator */
int continueExec(char *operator, int exitStatus) {
    if (strcmp(operator, "&&") == 0) {
        return exitStatus == 0;
    } else if (strcmp(operator, "||") == 0) {
        return exitStatus != 0;
    } else if (strcmp(operator, "&") == 0) {
        return 1;
    }
    return 0; /* Unknown operator */
}