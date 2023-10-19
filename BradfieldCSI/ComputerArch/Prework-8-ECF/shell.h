#ifndef SHELL_H
#define SHELL_H

#include<string.h>
#include<errno.h>
#include<unistd.h>
#include<sys/types.h>
#include<sys/wait.h>
#include<signal.h>

#define MAX_INPUT 50
#define MAX_ARGS 4
#define BUILTIN_NO_EXEC 0
#define BUILTIN_EXEC 1
#define SHELL_SYMBOL "\U0001f525"
#define PROCESS_REAP_SUCCESS "child %d exited normally"
#define PROCESS_REAP_FAILURE "child %d exited abnormally"
#define FORK_ERR "fork error"
#define COMMAND_ERR "error executing command"
#define SHELL_PROFILE "./.shellprofile"
#define SLEEP_CMD "sleep"
#define EXIT_MSG "\n%sIf this is to end in fire, we should all burn together%s\n"
#define PATH "/bin"

struct Token { /* Call to ParseArgs always returns a token */
    char *command;
    char *args[MAX_ARGS];
    char *operator;
    char *s_next;
};

char *test_command = "pecho";
char *test_argv[] = {"my name is yuvi", "5"};

pid_t Fork(char *message);
int Execvp(char *command, char *args[], char *errmsg);
struct Token ParseArgs(char *s);
void ExecProgram(char *command, char *args[]);
int Builtin(char *command, char *args[]);
void Alias(char *args[]);
void AliasPrintAll(char *args[]);
void SigintHandler(int sig);
void SigtstpHandler(int sig);
void SigchldHandler(int sig);
char *GetOperator(char *s);
int ContinueExec(char *operator, int exitStatus);

#endif