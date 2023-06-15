#include<stdio.h>

/*Testing value of getchar() != EOF */
main() {
    int c = 5;
    printf("Value of getchar != EOF %d", getchar() != EOF);
    printf("\n");
    printf("EOF %d", EOF);
    printf("\n");
}