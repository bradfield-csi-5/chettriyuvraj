#include<stdio.h>
#include<stdlib.h>
#include<ctype.h>
#include "shell.h"

int shell_pid;

int main() {
    char *s = NULL, *check = NULL;
    char input[MAX_INPUT];
    char *args[MAX_ARGS];
    int status, exit_status;
    struct Token token = {
        "\0",
        "\0"
    };
    sigset_t mask, prev_mask;
    shell_pid = getpid();
    

    signal(SIGINT, SigintHandler); /* Installing SIGINT handler  */
    signal(SIGTSTP, SigtstpHandler); /* Installing SIGTSTP handler */
    signal(SIGCHLD, SigchldHandler); /* Installing SIGTSTP handler */



    do {

        printf(SHELL_SYMBOL);

        if (*token.s_next != '\0' && ContinueExec(token.operator, exit_status)) { /* If operators were used, continue execution based on exit status */
            s = token.s_next;
        } else {
            s = fgets(input, MAX_INPUT, stdin);
        }

        if (s) {
            token = ParseArgs(s, args); /* parseArg will always return a token depending on operators used in expression */
            signal(SIGTTOU, SIG_IGN);
            if (!Builtin(args)) {
                pid_t pid = Fork(FORK_ERR);
                if (pid == 0) {
                    pid_t child_pid = getpid();
                    setpgid(child_pid, child_pid); /* We want to implement job control, so we will create separate pgid for each process */
                    if (*token.operator != '&' ) {
                        tcsetpgrp(STDIN_FILENO, child_pid);
                    }
                    signal(SIGTTOU, SIG_DFL);
                    Execvp(args, COMMAND_ERR);
                }
                

                
                
                
                
                if (*token.operator != '&' ) { /* Don't wait on background processes */
                    while (waitpid(pid, &status, WUNTRACED) > 0);  /* Reaping child process and tracking last exit_status*/
                    if (WIFEXITED(status)) {
                        exit_status = WEXITSTATUS(status);
                    }
                    tcsetpgrp(STDIN_FILENO, getpgid(shell_pid));
                    signal(SIGTTOU, SIG_DFL);
                }
            }
        }

        if (s && !strchr(s, '\n')) { /* If newline not found, print first MAX chars and flush remaining - why was this condition written (?)*/
            while ((s = fgets(input, MAX_INPUT, stdin)) && !strchr(s, '\n'));
            printf("\n");
        }
    
    } while(s != NULL);
    


    if (!feof(stdin)) { /* If EOF is not set i.e error */
        printf("\n Error while reading input");
        return 1;
    }

    printf("\n%sIf this is to end in fire, we should all burn together%s\n", SHELL_SYMBOL, SHELL_SYMBOL);
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
    int result_status;
    if ((result_status = execvp(args[0], args)) == -1) {
        printf("%s: %s", errmsg, strerror(errno));
    }
    return result_status;
}

/* Parse args from input string */
struct Token ParseArgs(char *s, char *args[]) {
    int arg_idx = 0;
    char *operator; /* To grab operators like && || */
    struct Token token = { 
        "\0",
        "\0"
    };

    for(;*s==' ';s++) ; /* Skip initial spaces */

    char *cur = s;
    while (*s != '\0' && *(operator = GetOperator(s)) == '\0') {

        while (*++cur != ' ' && *cur != '\n' && *cur != '\0'); /* Find delimiter */

        char *new_arg = (char *)malloc(sizeof(char) * (cur - s + 1)); /* Malloc for arg when delimiter found, add to to arg array */
        strncpy(new_arg, s, cur-s);
        new_arg[cur-s] = '\0';
        args[arg_idx++] = new_arg;

        for (;*cur==' ' || *cur=='\n'; cur++); /* Skip additional spaces */

        s = cur;
    }

    if (*operator != '\0') { /* If stopped parsing due to operator, return a token containing operator and start of next string */
        token = (struct Token) {
            s + strlen(operator), /* Might have additional spaces at the start, they will be handled by next call to ParseArgs */
            operator,
        };
    }

    args[arg_idx] = NULL;

    return token;
}

char *GetOperator(char *s) {
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

/* Execute and return 1 if Builtin, 0 otherwise */
int Builtin(char *args[]) {
    int return_status = 0;

    if (strcmp("alias", args[0]) == 0) {
        Alias(args);
        return_status = 1;
    } else if (strcmp("exit", args[0]) == 0) {
        printf("\n%sIf this is to end in fire, we should all burn together%s\n", SHELL_SYMBOL, SHELL_SYMBOL);
        exit(0);
    }

    return return_status;
}

/*     */

/**
 * Builtin alias command implementation
 * 
 * Notes:
 * - only alias -p and alias without arguments handled for now
 * - multiple -p flags will lead to all alias-es being printed once
 * - no arguments provided to alias = print all aliases
 * 
 * 
 **/

void Alias(char *args[]) {
    FILE *f;
    fpos_t *value_offset;
    int cur_char, are_all_printed = 0, is_match; 
    char *value_offset_arg, *match_str;
    

    if (*(args+1) == NULL) {
        AliasPrintAll(args);
        return;
    }

    // if ((f = fopen(SHELL_PROFILE, "a+")) == NULL) {
    //     printf("alias error: %s", strerror(errno));
    //     return;
    // }
    


    while (*++args != NULL)  { /* Handle each arg independently  */
        char *cur_str = *args;

        if(strcmp(cur_str, "-p") == 0 && !are_all_printed) { /* If -p encountered and all alias-es not already printed once */
            are_all_printed = 1;
            AliasPrintAll(args);
            continue;
        }
        else { // only for now
            continue;
        }
        
        /* For handling arguments other than -p */
        /* ------------------------------------------------------INCOMPLETE --------------------------------------------------------------------------- */
        if((value_offset_arg = strrchr(cur_str, '=')) != NULL) { /* If assignment statement found, search in file and replace */
            f = fopen(SHELL_PROFILE, "r+");
            cur_char = '\0';

            while (cur_char != EOF) {
                is_match = 1;
                match_str = cur_str; /* Compare char by char */
                
                while((cur_char = fgetc(f)) != '\n' && cur_char != EOF) {
                    if (cur_char == '=') {
                        fgetpos(f, value_offset); /* Record start of string position; TODO: handle error */
                    }
                    if (*match_str++ != cur_char){ /* If mis_match at any point, co */
                        is_match = 0;
                        break;
                    }
                }

                if (is_match) { /* Start writing from value_offsetIdx+1 in cur_str to value_offset in file */
                    fsetpos(f, value_offset);
                    for (char *s = value_offset_arg + 1; *s != '\0'; s++) {
                        fputc(*s, f); /* TODO: handle err */
                    }
                    fseek(f, 0, SEEK_END); /* TODO handle err */
                } else { /* Iterate till EOF/EOL */
                    while((cur_char = fgetc(f)) != '\n' && cur_char != EOF);
                }
            }

        } else { /* Search line wise */

        }

        /* ------------------------------------------------------INCOMPLETE --------------------------------------------------------------------------- */

    }

}


void AliasPrintAll(char *args[]) {
    FILE *f;
    int cur_char;
    if ((f = fopen(SHELL_PROFILE, "r")) == NULL) {
        printf("alias error: %s", strerror(errno));
        return;
    }
    while ((cur_char = fgetc(f)) != EOF){ /* Make sure \n exists in file itself */
        fputc(cur_char, stdout);
    }
    return;
}

void SigintHandler(int sig) {
    const char msg[] = "SIGINT caught\n";
    write(STDOUT_FILENO, msg, sizeof(msg)-1);
    return;
}

void SigtstpHandler(int sig) {
    const char msg[] = "SIGTSTP caught\n";
    write(STDOUT_FILENO, msg, sizeof(msg)-1);
    tcsetpgrp(STDIN_FILENO, getpgid(shell_pid));
    exit(0);
}

void SigchldHandler(int sig) {
    const char msg[] = "SIGCHLDcaught\n";
    write(STDOUT_FILENO, msg, sizeof(msg)-1);
    printf("%d shellPgid", getpgid(shell_pid));
    tcsetpgrp(STDIN_FILENO, getpgid(shell_pid));
    // exit(0);
    return;
}


/* Determine whether to continue execution or not depending on exit_status and operator */
int ContinueExec(char *operator, int exit_status) {
    if (strcmp(operator, "&&") == 0) {
        return exit_status == 0;
    } else if (strcmp(operator, "||") == 0) {
        return exit_status != 0;
    } else if (strcmp(operator, "&") == 0) {
        return 1;
    }
    return 0; /* Unknown operator */
}



/**
 * Naming convention
 *
 * THIS_IS_MY_CONVENTION for macros, enum members
 * ThisIsMyConvention for file name, object name (class, struct, enum, union...), function name, method name, typedef
 * this_is_my_convention global and local variables,
 * parameters, struct and union elements
 * thisismyconvention [optional] very local and temporary variables (such like a for() loop index)
 * 
 * 
 **/