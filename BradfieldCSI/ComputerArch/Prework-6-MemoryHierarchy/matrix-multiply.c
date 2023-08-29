/*
Naive code for multiplying two matrices together.

There must be a better way!
*/

#include <stdio.h>
#include <stdlib.h>

#define TILE_SIZE 32
#define min(a, b) (((a) < (b)) ? (a) : (b))

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void transpose_matrix(double **A, int a_cols, int b_cols);

void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}


/* Single accoumulator
-Og
Naive: 0.865s
Fast: 0.596s
1.45x speedup

-O2
Naive: 0.221s
Fast: 0.215s
1.03x speedup
 */

/* void fast_matrix_multiply(double **c, double **a, double **b, int a_rows,
                          int a_cols, int b_cols) {
  
  double acc;
  int i, j, k;

  for (i = 0; i < a_rows; i++) {
    for (j = 0; j < b_cols; j++) {
      acc = 0;
      for (k = 0; k < a_cols; k++)
        acc += a[i][k] * b[k][j];
      c[i][j] = acc;
    }
  }

} */




/* Multiple accumulators i.e loop unrolling
Size 512
-Og and O2 roughly the same
Naive: 0.242s
Fast: 0.138s
1.76x speedup
However, matrix results did not match - probably rounding error when comparing

Size 1024
-O2 - not much
Naive: 2.523s
Fast: 2.245s
1.12x speedup

-Og
Naive: 10.489s
Fast: 6.960s
1.51x speedup

*/

/* void fast_matrix_multiply(double **c, double **a, double **b, int a_rows,
                          int a_cols, int b_cols) {
  
  double acc1, acc2, acc3, acc4;
  int i, j, k;
  int limit = a_cols - 3;

  for (i = 0; i < a_rows; i++) {
    for (j = 0; j < b_cols; j++) {
      acc1 = 0, acc2 = 0, acc3 = 0, acc4 = 0;
      for (k = 0; k < limit; k+=4) {
        acc1 += a[i][k] * b[k][j];
        acc2 += a[i][k+1] * b[k+1][j];
        acc3 += a[i][k+2] * b[k+2][j];
        acc4 += a[i][k+3] * b[k+3][j];
      }
      for (;k < a_cols; k++) {
        acc1 += a[i][k] * b[k][j];
      }
      c[i][j] = acc1 + acc2 + acc3 + acc4;
    }
  }
} */






/* Transpose and then multiply - there is also a version where you don't transpose (simply distribute weights as you go)
Size 512
-Og
Naive: 0.904s
Fast: 0.590s
1.53x speedup

-O2 minimal diff
Naive: 0.204s
Fast: 0.172s
1.19x speedup

Size 1024
-O2
Naive: 2.592s
Fast: 1.501s
1.73x speedup

Huge difference in terms of cache hits/misses shown in solution.md
 */
/* void fast_matrix_multiply(double **c, double **a, double **b, int a_rows,
                          int a_cols, int b_cols) {
  // TODO: write a faster implementation here!
  transpose_matrix(b, a_cols, b_cols); // considering square matrix anyway
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      c[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        c[i][j] += a[i][k] * b[j][k];
    }
  }
} */


/* Results in soln.md */
/* Cache blocking - separate into tiles 
-O2, Og 1024 size - 512 is still ~1.5
Naive: 2.580s
Fast: 1.109s
2.33x speedup
*/
/* void fast_matrix_multiply(double **C, double **A, double **B, int a_rows,
                          int a_cols, int b_cols) {
  int ti, tj, tk,          // indexes of the tile
      i, j, k,             // indexes within a tile
      i_end, j_end, k_end; // end when matrix dim not a multiple of tile_size

  for (ti = 0; ti < a_rows; ti += TILE_SIZE) {
    i_end = min(ti + TILE_SIZE, a_rows);

    for (tj = 0; tj < b_cols; tj += TILE_SIZE) {
      j_end = min(tj + TILE_SIZE, b_cols);

      for (tk = 0; tk < a_cols; tk += TILE_SIZE) {
        k_end = min(tk + TILE_SIZE, a_cols);

        // Compute this tile
        for (i = ti; i < i_end; i++)
          for (j = tj; j < j_end; j++)
            for (k = tk; k < k_end; k++)
              C[i][j] += A[i][k] * B[k][j];
      }
    }
  }
} */

/* Multi accumulators + Transpose 
Size 512

-Og
cc -Wall matrix-multiply.c benchmark.c && ./a.out 512
Naive: 0.890s
Fast: 0.356s
2.50x speedup

However, matrix results did not match! - Rounding check error

-O2 Same as Og 

Size 1024

-Og

Naive: 10.552s
Fast: 3.074s
3.43x speedup

-O2
Naive: 2.534s
Fast: 0.808s
3.13x speedup


*/
void fast_matrix_multiply(double **C, double **A, double **B, int a_rows,
                          int a_cols, int b_cols) {
  int ti, tj, tk,          // indexes of the tile
      i, j, k,             // indexes within a tile
      i_end, j_end, k_end; // end when matrix dim not a multiple of tile_size

  double acc1, acc2, acc3, acc4;

  for (ti = 0; ti < a_rows; ti += TILE_SIZE) {
    i_end = min(ti + TILE_SIZE, a_rows);

    for (tj = 0; tj < b_cols; tj += TILE_SIZE) {
      j_end = min(tj + TILE_SIZE, b_cols);

      for (tk = 0; tk < a_cols; tk += TILE_SIZE) {
        k_end = min(tk + TILE_SIZE, a_cols);

        // Compute this tile
        for (i = ti; i < i_end; i++)
          for (j = tj; j < j_end; j++) {
            acc1 = 0, acc2 = 0, acc3 = 0, acc4 = 0;
            for (k = tk; k < k_end - 3; k+=4) {
              acc1 += A[i][k] * B[k][j];
              acc2 += A[i][k+1] * B[k+1][j];
              acc3 += A[i][k+2] * B[k+2][j];
              acc4 += A[i][k+3] * B[k+3][j];
            }
            for (;k<k_end;k++) /* Picking up remaining since k_end might not be exact multiple of 4 */
              acc1 += A[i][k] * B[k][j];
            C[i][j] += acc1 + acc2 + acc3 + acc4;
          }
      }
    }
  }
}

    //   for (int i = 0; i < a_rows; i++) {
    //   for (int j = 0; j < b_cols; j++) {
    //     printf("%lf ", C[i][j]);
    //   }
    //   printf("\n");
    // }
    // printf("\n\n");

void transpose_matrix(double **A, int a_rows, int a_cols) {
  for (int i = 0; i < a_rows; i++)
    for (int j = i + 1; j < a_cols; j++) {
      double temp = A[i][j];
      A[i][j] = A[j][i];
      A[j][i] = temp;
    }
}



