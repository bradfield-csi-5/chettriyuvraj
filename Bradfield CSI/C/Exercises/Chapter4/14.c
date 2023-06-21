#include<stdio.h>

#define swap(t, x, y) ({t temp = x; x = y; y = temp; })

int main() {
    int x = 5;
    int y = 6;
    swap(int, x, y);
    printf("x val %d\n", x);
    printf("y val %d", y);
}