/* Write a program that reads a C program and prints in alphabetical order 
each group of variable names that are identical in the first 6 characters 
but different somewhere thereafter. Don't count words within strings and comments. 
Make 6 a parameter that can be set from the command line. */


/* 

Design the solution as a binary search tree structure of the form (if k = 6 i.e first 6 chars)

                      abcdefk-> abcdefl
                    /        \
            abcded -> NULL     abcdg -> NULL

struct node {
    char *word;
    struct node *next;
    struct node *left;
    strcut node *right;
}

- Each node of the binary search tree is also the head of a linked list, the elems of the linked list are all words containing the same first 'k' letters.
- Except for the 'head' of the linked list, all other nodes of the linked list have no 'left and right' elems
- When a new node is to be added to a linked list, the it takes over as the head of the linked list and takes the left and right pointers of the previous head.
 
 Functions:
 main() - read input using readLine, get input back (in a global array of char pointers), trigger traverse function for each word
 readLine() - read input, ignore comments and put it into a line
 extractVars() - extract variables from a given line and put them into an array of char pointers
 traverse() - traverse BST and put word into correct position
 salloc,nalloc
 
 Global variables:
 
 struct node {...};
 struct node head = NULL;
 
 TODOS:
 1. readline doesn't properly handle comments (fails for certain multiline comments cases)
 2. extractVars handles only basic single line declarations of the form char a,b,c; (no proper error handling + doesn't even handle spaces between multiple variables)
 3. traverse() doesn't group words within a group alphabetically i.e identifies correct word group but always attaches new word to the head of the group (pretty easy to fix)
 */

#include<stdio.h>
#include<malloc.h>
#include<string.h>
#include "util.h"

#define MAXVARIABLES 20 /* Max number of variables in a line */
#define MAXLINELENGTH 90 /* Maximum length of a line i.e until semicolon occurs */

char *datatypes[] = {
    "int",
    "float",
    "char",
    NULL
};

struct node {
    char *word;
    struct node *next;
    struct node *left;
    struct node *right;
};

int readLine(char line[], int lim); /* Reads line and puts it into char array */
void extractVars(char line[], char *lineVars[]);
struct node *traverse(struct node *head, char *p, int k); /* Words of length greater than k expected */
char *salloc(char *sp, char *ep); /* Allocates space for string from sp to ep + adds an additional \0 at the end */
struct node *nalloc(); /* Allocate space for a new node */
void dfs(struct node *head);

int main() {
    char *lineVars[MAXVARIABLES];
    char line[MAXLINELENGTH];
    int varCount;
    struct node *head = NULL;

    while ((readLine(line, MAXLINELENGTH)) >= 0) { /* First read line */
        printf("%s e\n", line);
        extractVars(line, lineVars); /* Extract variables from that line */
        for (char **p = lineVars; *p != NULL; *p++) {
            printf("Traverse for %s\n", *p);
            head = traverse(head, *p, 3);
            // printf("Start dfs\n");
            // dfs(head);
        }
        dfs(head);
    }

    printf("Input over!");

    return 0;
}

void dfs(struct node *head) {
    /* Base case */
    if (head == NULL)
        return;


    /* Recursive case */
    for (struct node *curNode = head; curNode != NULL; curNode = curNode -> next) { /* Prints all the matching k-letter words of current group */
        printf("%s->", curNode->word);
    }
    printf("\n");
    dfs(head->left);
    dfs(head->right);
}

/* Attaches nodes with their alphabetical 'k-length' group (but not in alphabetical order within the group) */
struct node *traverse(struct node *head, char *p, int k) {
    int comparison; /* To compare current nodes word with current word */
    struct node *tNode = head; /* traversal node */
    struct node *prevTNode = head; /* previous traversal node - required when adding a new node */
    struct node *newNode = nalloc(); /* Allocate new node and set its fields */
    newNode->word = p;
    newNode->left = NULL;
    newNode->right = NULL;
    newNode ->next = NULL;

    if (head == NULL) { /* Edge case for very first node */
        head = newNode;
        return head;
    }
    while ((comparison = strncmp(p,tNode->word,k)) != 0) {
        // printf("word %s\n", p);
        // printf("tNode word %s\n",tNode->word);
        // printf("comparison %d\n", comparison);
        if(comparison < 0){ /* Move to left sub tree or terminate if nothing left */
            if (tNode->left == NULL) {
                tNode->left = newNode;
                return head;
            } else {
                prevTNode = tNode;
                tNode = tNode->left;
            }
        } else { /* greater than 0  - Move to right sub tree or terminate if nothing left*/
            if (tNode->right == NULL) {
                tNode->right = newNode;
                return head;
            } else {
                prevTNode = tNode;
                tNode = tNode->right;
            }
        }
    }

    /* matching node found, add new node as the head of current group, set all nodes which reference to old head to new head */
    newNode->left = tNode->left;
    newNode->right = tNode->right;
    newNode->next = tNode;

    tNode->left = NULL;
    tNode->right = NULL;

    if (prevTNode->left == tNode)
        prevTNode->left = newNode;
    else
        prevTNode->right = newNode;
    
    // printf("inside");
    // for (struct node *curNode = newNode; curNode != NULL; curNode = curNode -> next) {
    //     printf("%s->", curNode->word);
    // }
    // printf("\n");
    if(tNode == head)
        return newNode;
    return head;
}

/* 
Extracts vars into array of character pointers:
Note: Works only for valid declarations without spaces i.e char a,b,c; */

void extractVars(char line[], char *lineVars[]) {
    char *varStart;
    char **lineVarP = lineVars;
    for (char **p = datatypes; *p != NULL; p++){ /* Check each dataype if present  */
        // printf("checking for %s", *p);

        if ((varStart = strstr(line, *p)) != NULL) { /* Always assuming valid declaration */
            // printf("found for %s", varStart);
            for(line = varStart; *line != ' ' && *line != '\t'; line++); /* iterate till the next blank ' ' */

            for(;((*line) == ' ' || (*line) == '\t'); line++); /* iterate over all blanks to reach variables*/

            while (*line != '\0') {
                char *sp = line; /* starting pointer */
                for (;*line !=',' && *line != '\0'; line++); /* set bounds for start and end pointer of variable */
                
                // printf("First char of var %c and last char of var %c\n", *sp, *line);
                char *vp = salloc(sp, line); /* pointer to variable name */
                // printf("Allocated value %s\n",vp);
                *(lineVarP++) = vp; 
                if (*line == ',')
                    ++line;
            }
        }
    }
    *lineVarP = NULL; /* Set last pointer as null to indicate end */

}

/* Allocate space for a new word */
char *salloc(char *sp, char *ep) {
    char *p = (char *)malloc(ep - sp + 1);
    strncpy(p,sp, ep-sp);
    *(p + (ep-sp)) = '\0';
    return p;
}

/* Allocate space for a new node */
struct node *nalloc() {
    return (struct node *)malloc(sizeof(struct node));
}

/* Test cases for readLine: 

Improve this later on

fsdjlf fdskljfdls  fdsfjsdf;
    fdsfsd  fsd f;;;
dsfasfasfaffdsafasf      (exceed limit)


// comment
hi // comment
hello ; // comment

multiline comment
multiline comment followed by closing and then words
words followed by multiline comment



*/


int readLine(char line[], int lim) {
    int c, i;
    for (i = 0; i < lim - 1 && (c = getch()) != EOF && c != ';'; i++) {
        // printf("ch %c \n", c);
        // pr
        if (c == '/' && i > 0 && *(line -1) == '/') { /* If single line comment encountered - return */
            while ((c = getch()) != EOF || c != '\n');
            *(line - 1) = '\0';
            return i; 
        }

        if (c == '/' && i > 0 && *(line -1) == '*') { /* If multi-line comment encountered - try to parse until limit encountered or comment ends */
            *(line - 1) = '\0';
            char prevC;
            while ((c = getch()) != '/' && prevC != '*' && i++ < lim - 1)
                prevC = c;
            continue; /* Skip to next iteration */
        }

        *(line++) = c;
    }

    *line = '\0';

    return c == EOF ? -1 : (c == ';' ? i : 0); /* If not terminated by semicolon, line exceeds limit and is invalid */

}

