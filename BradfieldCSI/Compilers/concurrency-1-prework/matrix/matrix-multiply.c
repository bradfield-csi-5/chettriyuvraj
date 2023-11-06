#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

#define MAXTHREADS 32
void *thread(void *arg);
struct matrix_mult {
  double **A;
  double **B;
  double **C;
  int a_cols, b_cols, a_start_row, a_end_row;
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
  pthread_t thread_id[MAXTHREADS + 1];
  struct matrix_mult * m;
  int start_row=0, end_row=0, row_per_thread = a_rows/MAXTHREADS;
  int extra_thread = 0;

  /* Split rows up by thread */
  for (int i = 0; i < MAXTHREADS; i++) {
    end_row += row_per_thread;
    m = malloc(sizeof(struct matrix_mult));
    (*m).A = A;
    (*m).B = B;
    (*m).C = C;
    (*m).a_cols = a_cols;
    (*m).b_cols = b_cols;
    (*m).a_start_row = start_row;
    (*m).a_end_row = end_row;
    pthread_create(&thread_id[i], NULL, thread, m);
    start_row = end_row;
  }

  /* Leftover rows */
  if (end_row != a_rows) {
    extra_thread = 1;
    m = malloc(sizeof(struct matrix_mult));
    (*m).A = A;
    (*m).B = B;
    (*m).C = C;
    (*m).a_cols = a_cols;
    (*m).a_start_row = start_row;
    (*m).a_end_row = a_cols;
    pthread_create(&thread_id[MAXTHREADS], NULL, thread, m);
  }

  for (int i = 0; i < MAXTHREADS + extra_thread; i++) {
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

  for (int i = m.a_start_row; i < m.a_end_row; i++) {
    for (int j = 0; j < m.b_cols; j++) {
      m.C[i][j] = 0;
      for (int k = 0; k < m.a_cols; k++)
        m.C[i][j] += m.A[i][k] * m.B[k][j];
    }
  }
  return NULL;
}


/* 
- Single thread per computation:
n = 30 
Naive: 0.000s
Parallel: 3.774s
0.00x speedup between naive and parallel
- Results don't match, presumably rounding issue

- Dividing rows among threads: n = 500
;; Threads = 2
Naive: 0.579s
Parallel: 0.293s
1.98x speedup between naive and parallel

;; Threads = 4
Naive: 0.576s
Parallel: 0.146s
3.95x speedup between naive and parallel

;; Threads = 8
Naive: 0.574s
Parallel: 0.086s
6.71x speedup between naive and parallel

;; Threads = 16
Naive: 0.577s
Parallel: 0.078s
7.41x speedup between naive and parallel

Did not indrease twofolds after this point
 */