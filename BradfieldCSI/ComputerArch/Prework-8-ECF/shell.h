#ifndef SHELL_H
#define SHELL_H

#include<string.h>
#include<errno.h>
#include<unistd.h>
#include<sys/types.h>
#include<sys/wait.h>
#include<signal.h>

#define MAX_INPUT 50 /* Last \0 */
#define MAX_ARGS 4 /* Last NULL */
#define SHELL_SYMBOL "\U0001f525"
#define FORK_ERR "fork error"
#define COMMAND_ERR "command does not exist"
#define SHELL_PROFILE "./.shellprofile"
#define SLEEP_CMD "sleep"

struct Token { /* Call to ParseArgs always returns a token */
    char *command;
    char *args[MAX_ARGS];
    char *operator;
    char *s_next;
};

char *test_command = "pecho";
char *test_argv[] = {"my name is yuvi", "5"};

pid_t Fork(char *message);
int Execvp(char *args[], char *errmsg);
struct Token ParseArgs(char *s, char *args[]);
int Builtin(char *args[]);
void Alias(char *args[]);
void AliasPrintAll(char *args[]);
void SigintHandler(int sig);
void SigtstpHandler(int sig);
void SigchldHandler(int sig);
char *GetOperator(char *s);
int ContinueExec(char *operator, int exitStatus);

#endif