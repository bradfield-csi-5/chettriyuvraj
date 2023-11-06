#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

#define MAXTHREADS 100000
void *thread(void *arg);
struct matrix_mult {
  double **A;
  double **B;
  double **C;
  int i, j, k;
};

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols, int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }

  /* Print result */
  // for (int i = 0; i < a_rows; i++) {
  //   for (int j = 0; j < b_cols; j++)
  //     printf("%f ", C[i][j]);
  //   printf("\n");
  // }


}

void parallel_matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols, int b_cols) {
  pthread_t thread_id[MAXTHREADS];
  struct matrix_mult * m;
  int thread_count = 0;

  /* One thread for each computation */
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++) {
        m = malloc(sizeof(struct matrix_mult));
        // m =  {A, B, C, i, j, k};
        (*m).A = A;
        (*m).B = B;
        (*m).C = C;
        (*m).i = i;
        (*m).j = j;
        (*m).k = k;
        pthread_create(&thread_id[thread_count++], NULL, thread, m);
      }
    }
  }

  for (int i = 0; i<thread_count; i++) {
    pthread_join(thread_id[i], NULL);
  }

  /* Print result */
  // printf("\n\n Parallel Result\n");
  // for (int i = 0; i < a_rows; i++) {
  //   for (int j = 0; j < b_cols; j++) {
  //     printf("%f ", C[i][j]);
  //   }
  //   printf("\n");
  // }

}

void *thread(void *arg) {
  struct matrix_mult m = *((struct matrix_mult *) arg);
  m.C[m.i][m.j] += m.A[m.i][m.k] * m.B[m.k][m.j];
  return NULL;
}


/* 
- Single thread per computation:
n = 30 
Naive: 0.000s
Parallel: 3.774s
0.00x speedup between naive and parallel
- Results don't match, presumably rounding issue



 */