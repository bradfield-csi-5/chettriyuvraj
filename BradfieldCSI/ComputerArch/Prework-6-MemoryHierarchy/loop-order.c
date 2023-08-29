  /*

Two different ways to loop over an array of arrays.

Spotted at:
http://stackoverflow.com/questions/9936132/why-does-the-order-of-the-loops-affect-performance-when-iterating-over-a-2d-arra

*/

#include<stdio.h>
#include<malloc.h>


void option_one() {
  int i, j;
  static int x[4000][4000];
  for (i = 0; i < 4000; i++) {
    for (j = 0; j < 4000; j++) {
      x[i][j] = i + j;
    }
  }
}

void option_two() {
  int i, j;
  static int x[4000][4000];
  for (i = 0; i < 4000; i++) {
    for (j = 0; j < 4000; j++) {
      x[j][i] = i + j;
    }
  }
}

void option_three(int (*x)[4000]) {
  int i, j;
  for (i = 0; i < 4000; i++) {
    for (j = 0; j < 4000; j++) {
      x[i][j] = i + j;
    }
  }
}

void option_four(int (*x)[4000]) {
  int i, j;
  for (i = 0; i < 4000; i++) {
    for (j = 0; j < 4000; j++) {
      x[j][i] = i + j;
      // printf("%d", x[0]);
    }
  }
}

void option_five() {
  int i, j;
  static int x[4000][8];
  for (i = 0; i < 8; i++) {
    for (j = 0; j < 4000; j++) {
      x[j][i] = i + j;
    }
  }
}

/// 32000 ints to access, every alternate row one cache miss => 16000 cache write misses?

int main() {
  // int (*x)[4000] = malloc(sizeof(int[4000][4000]));

  // if (x == NULL) {
  //   fprintf(stderr, "Memory allocation failed.\n");
  //   return 1;
  // }


  // option_one();
  // option_two();
  // // option_three(x);
  // // option_four(x);
  option_five();
  return 0;
}

// option_one - negligible misses from instructions
// 64 -> size of one line 16000 data accesssess in sequence,
// so every 64 lines 1 miss => 16000/64 ~ 16000/60 = 1600 / 6 = 800/3 = 266.66
// 250 data misses ??



