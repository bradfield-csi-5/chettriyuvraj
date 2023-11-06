#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

#define MAXTHREADS 100
void *thread(void *arg);


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
}

void parallel_matrix_multiply(double **c, double **a, double **b, int a_rows, int a_cols, int b_cols) {
  pthread_t thread_id[MAXTHREADS];
  int thread_count = 5;
  for (int i = 0; i < thread_count; i++) {
    pthread_create(&thread_id[i], NULL,thread, NULL);
  }

}

void *thread (void *arg) {
  printf("\nThread id %ld", pthread_self());
  return NULL;
}


/* 
1. 

 */