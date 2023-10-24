#include<stdio.h>
#include<stdlib.h>
#include<ctype.h>
#include "shell.h"

/* Global variable to keep track of forked children */
pid_t child_pid = 0;

int main() {
    char input[MAX_INPUT], *s = NULL;
    /* Keep track of exit_status for continuing/discontinuing execution of commands with operators eg && */
    int exit_status;
    struct Token token = {NULL, NULL, NULL, NULL};
    
    /* Installing signal handlers */
    signal(SIGINT, SigintHandler);

    while(1) {

        /* Determine if previous statements continuation exists eg. "ls && ps" or to wait for new command */
        if (token.s_next != NULL && ContinueExec(token.operator, exit_status)) {
            s = token.s_next;
        } else {
            printf(SHELL_SYMBOL);
            s = fgets(input, MAX_INPUT, stdin);
        }

        /* Parse token from string*/
        token = ParseArgs(s);

        /* Check and execute depending on Builtin or command */
        if (Builtin(token.command,token.args) == BUILTIN_EXEC) {
            continue;
        } else {
            exit_status = ExecProgram(token.command, token.args);
        }

    }
}

/**
 * Forks a subprocess, executes program and returns exit_status
 **/
int ExecProgram(char *command, char *args[]) {
    pid_t pid;
    int status, exit_status;

    /* Fork child, create new process group for it and then execute */
    if ((pid = Fork(FORK_ERR)) == 0) {
        setpgid(0,0);
        Execvp(command, args, COMMAND_ERR);
    }

    /* Set global variable child_pid */
    child_pid = pid;

    /* Reap child process - is while loop required here (?)*/
    while ((pid = wait(&status)) > 0) {
        if WIFEXITED(status) {
            printf(PROCESS_REAP_SUCCESS, pid);
            /* Exit status of normally terminated process */
            exit_status = WEXITSTATUS(status);
        } else {
            printf(PROCESS_REAP_FAILURE, pid);
            /* Not grabbing exact exit status of abnormally terminated process currently - set it to arbitrary value */
            exit_status = 2;
        }
    }
    return exit_status;
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
int Execvp(char *command, char *args[], char *errmsg) {

    /* Allocate new array to store command + args */
    char **command_and_args = (char**)malloc(sizeof(char**) * MAX_ARGS);
    command_and_args[0] = command;
    for (int i = 1; args[i-1] != NULL; i++) {
        command_and_args[i] = args[i-1];
    }

    /* Execute command */
    int result_status;
    if ((result_status = execvp(command, command_and_args)) == -1) {
        printf("%s: %s", errmsg, strerror(errno));
    }
    return result_status;
}

/**
 * Parse args from input string:
 * 
 * - Parses string till an operator is encountered ie given a string "ls -a && ps", parses "ls", parses "-a" grabs "&&", and notes the position of "ps"
 * - Returns this info as a token, token contains: 
 *      a) command ie "ls"
 *      b) args ie ["-a"] in an array
 *      c) s_next ie pointer to next command "ps"
 *      d) *operator ie "&&"
 **/
struct Token ParseArgs(char *s) {

    struct Token token = {NULL, NULL, NULL, NULL};
    int arg_idx = 0;

    if (s == NULL) {
        return token;
    }

    /* Skip initial spaces */
    for(;*s==' ';s++);

    /* Use a separate variable to traverse so we can calculate length of the command/arg we are traversing*/
    char *cur_s = s;

    /* Iterate until string ends or an operator is found */
    while (*s != '\0' && *(token.operator = GetOperator(s)) == '\0') {

        /* Find delimiter */
        for(;*cur_s != ' ' && *cur_s != '\n' && *cur_s != '\0'; cur_s++);

        /* Allocate space for command/arg */
        char *new_arg = (char *)malloc(sizeof(char) * (cur_s - s + 1));
        strncpy(new_arg, s, cur_s-s);
        new_arg[cur_s-s] = '\0';

        /* Determine whether we have parsed a command or an arg */
        if (token.command == NULL) {
            token.command = new_arg;
        } else {
            token.args[arg_idx++] = new_arg;
        }

        /* Skip additional spaces */
        for (;*cur_s==' ' || *cur_s=='\n'; cur_s++);

        /* Bring s to par with cur_s's position */
        s = cur_s;
    }

    /* If stopped parsing due to operator, additional tokens may exist, provide index to start of next command */
    if (token.operator != NULL) { 
        token.s_next = s + strlen(token.operator);
    }

    token.args[arg_idx] = NULL;

    return token;
}


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

/****** HELPERS ******/

/**
 * Helper to execute Builtins
 *
 * - Implemented only the alias builtin for now
 * - TODO: implement more builtins
 **/

int Builtin(char *command, char *args[]) {
    int return_status = BUILTIN_NO_EXEC;

    if (command == NULL) {
        return return_status;
    } 
    
    if (strcmp("alias", command) == 0) {
        Alias(args);
        return_status = BUILTIN_EXEC;
    } else if (strcmp("exit", command) == 0) {
        printf(EXIT_MSG, SHELL_SYMBOL, SHELL_SYMBOL);
        exit(0);
    }

    return return_status;
}

/**
 * Recognizes only valid operators, parses the valid operator and returns it
 * Currently valid operators are:
 * - &&
 * - || 
 * -  &
 * -  |
 **/

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


/**
 * Determine if execution is to be continued,
 * depending on if a known operator is found
 **/
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

/****** SIGNAL HANDLERS ******/

void SigintHandler(int sig) {
    const char msg[] = "SIGINT caught - terminating child\n";
    /* If child_pid is set, terminate it */
    if (child_pid > 0) {
        kill(child_pid, SIGTERM);
        child_pid = 0;
    }
    write(STDOUT_FILENO, msg, sizeof(msg)-1);
    return;
}

void SigintChildHandler(int sig) {
    printf("Sitgint child\n");
    exit(0);
}


/****** NOTES ******/

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




/**
 * Rules:
 * 
 * - Commands, flags, operators must be separated by a space delimiter ie "ls&&ps" is invalid
 * - Valid operators: &&, ||, &, |
 **/





// int main() {
//     char *s = NULL, *check = NULL;
//     char input[MAX_INPUT];
//     char *args[MAX_ARGS];
//     int status, exit_status;
//     struct Token token = {NULL, NULL, NULL, NULL};
//     sigset_t mask, prev_mask;
//     shell_pid = getpid();
    

//     // signal(SIGINT, SigintHandler); /* Installing SIGINT handler  */
//     // signal(SIGTSTP, SigtstpHandler); /* Installing SIGTSTP handler */
//     // signal(SIGCHLD, SigchldHandler); /* Installing SIGTSTP handler */


//     while(1) {
        

//          /* Determine if previous command's execution to be continued or to wait for new command */
//         if (token.s_next != NULL && ContinueExec(token.operator, exit_status)) {
//             s = token.s_next;
//         } else {
//             printf(SHELL_SYMBOL);
//             s = fgets(input, MAX_INPUT, stdin);
//         }

//         // if (s) {
//             // token = ParseArgs(s, args); /* parseArg will always return a token depending on operators used in expression */
//             // signal(SIGTTOU, SIG_IGN);
//             // pid_t pid = Fork(FORK_ERR);
//             // if (pid == 0) {
//             //     pid_t child_pid = getpid();
//             //     setpgid(child_pid, child_pid); /* We want to implement job control, so we will create separate pgid for each process */
//             //     if (*token.operator != '&' ) {
//             //         tcsetpgrp(STDIN_FILENO, child_pid);
//             //     }
//             //     signal(SIGTTOU, SIG_DFL);
//             //     Execvp(args, COMMAND_ERR);
//             // }
            
//             // if (*token.operator != '&' ) { /* Don't wait on background processes */
//             //     while (waitpid(pid, &status, WUNTRACED) > 0);  /* Reaping child process and tracking last exit_status*/
//             //     if (WIFEXITED(status)) {
//             //         exit_status = WEXITSTATUS(status);
//             //     }
//             //     tcsetpgrp(STDIN_FILENO, getpgid(shell_pid));
//             //     signal(SIGTTOU, SIG_DFL);
//             // }
//         // }

//         token = ParseArgs(s);

//         /* Check and execute if command is a builtin */
//         if (Builtin(token.command,token.args) == BUILTIN_EXEC) {
//             continue;
//         } else {
//             ExecToken(token.command, token.args);
//         }

//         // if (s && !strchr(s, '\n')) { /* If newline not found, print first MAX chars and flush remaining - why was this condition written (?)*/
//         //     while ((s = fgets(input, MAX_INPUT, stdin)) && !strchr(s, '\n'));
//         //     printf("\n");
//         // }

//     }
// }


/* 
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
} */