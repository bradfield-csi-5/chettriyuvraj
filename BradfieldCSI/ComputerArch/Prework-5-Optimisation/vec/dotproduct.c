#include <time.h>
#include <stdio.h>

#include "vec.h"

// #define TEST_LOOPS 1000000
#define TEST_LOOPS 10000

data_t dotproduct(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;

   for (long i = 0; i < vec_length(u); i++) { // we can assume both vectors are same length
      get_vec_element(u, i, &u_val);
      get_vec_element(v, i, &v_val);
      sum += u_val * v_val;
   }   
   return sum;
}


/* Remove vec_length invariant from loop - 
   4.46s to run 1000 tests (4460332.00ns per test)
   5.69s to run 1000 baseline tests (5688042.00ns per test)
   Speedup: 1.23s to run 1000 tests (1227710.00ns per test)
 */
/* data_t dotproductfast(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;
   long v_len = vec_length(u);

   for (long i = 0; i < v_len; i++) { // we can assume both vectors are same length
      u_val = get_vec_element(u, i, &u_val);
      v_val = get_vec_element(v, i, &v_val);
      sum += u_val * v_val;
   }   
   return sum;
} */


/* Replace get_vec_element function call with inline 
2.15s to run 1000 tests (2148923.00ns per test)
5.75s to run 1000 baseline tests (5754735.00ns per test)
Speedup: 3.61s to run 1000 baseline tests (3605812.00ns per test)
 */

/* data_t dotproductfast(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;
   long v_len = vec_length(u);


   for (long i = 0; i < v_len; i++) { // we can assume both vectors are same length
      u_val = GET_VEC_ELEMENT(u, i, &u_val);
      v_val = GET_VEC_ELEMENT(v, i, &v_val);
      sum += u_val * v_val;
   }   
   return sum;
} */

/* Eliminate memory redundancy of accessing v -> data each time
1.59s to run 1000 tests (1591183.00ns per test)
5.73s to run 1000 baseline tests (5731384.00ns per test)
Speedup: 4.14s to run 1000 tests (4140201.00ns per test)
 */
/* data_t dotproductfast(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;
   long v_len = vec_length(u); // Remove invariants from loop
   data_t * v_start = get_vec_start(v); // Remove invariant + this is also a memory redundancy
   data_t * u_start = get_vec_start(u);

   for (long i = 0; i < v_len; i++) { // we can assume both vectors are same length
      u_val = u_start[i];
      v_val = v_start[i];
      sum += u_val * v_val;
   }   
   return sum;
} */

/* Loop unrolling 2x1
0.97s to run 1000 tests (971237.00ns per test)
5.69s to run 1000 baseline tests (5685803.00ns per test)
Speedup: 4.71s to run 1000 tests (4714566.00ns per test)
 */
/* data_t dotproductfast(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;
   long v_len = vec_length(u), i; // Remove invariants from loop
   data_t * v_start = get_vec_start(v); // Remove invariant + this is also a memory redundancy
   data_t * u_start = get_vec_start(u);

   for (i = 0; i < v_len-1; i+=2) { // we can assume both vectors are same length
      sum += u_start[i] * v_start[i] + u_start[i + 1] * v_start[i + 1];
   }

   for (; i < v_len; i++) { // n - k + 1
      sum += u_start[i] * v_start[i];
   }
   return sum;
} */

/* Loop unrolling 3x1 - Faster than 2x1
0.90s to run 1000 tests (904007.00ns per test)
5.68s to run 1000 baseline tests (5678659.00ns per test)
Speedup: 4.77s to run 1000 tests (4774652.00ns per test)
 */
/* data_t dotproductfast(vec_ptr u, vec_ptr v) {
   data_t sum = 0, u_val, v_val;
   long v_len = vec_length(u), i; // Remove invariants from loop
   data_t * v_start = get_vec_start(v); // Remove invariant + this is also a memory redundancy
   data_t * u_start = get_vec_start(u);

   for (i = 0; i < v_len - 2; i+=3) { // we can assume both vectors are same length
      sum += u_start[i] * v_start[i] + u_start[i + 1] * v_start[i + 1] + u_start[i+2] * v_start[i+2];
   }

   for (; i < v_len; i++) { // n - k + 1
      sum += u_start[i] * v_start[i];
   }
   return sum;
} */

/* 
Loop unrolling 3x3 - In the same ballpark as 2x1
0.99s to run 1000 tests (989719.00ns per test)
5.75s to run 1000 baseline tests (5750566.00ns per test)
Speedup: 4.76s to run 1000 tests (4760847.00ns per test)
 */
/* data_t dotproductfast(vec_ptr u, vec_ptr v) {
   data_t sum = 0, sum_1 = 0, sum_2 = 0, u_val, v_val;
   long v_len = vec_length(u), i; // Remove invariants from loop
   data_t * v_start = get_vec_start(v); // Remove invariant + this is also a memory redundancy
   data_t * u_start = get_vec_start(u);

   for (i = 0; i < v_len - 2; i+=3) { // we can assume both vectors are same length
      sum += u_start[i] * v_start[i];
      sum_1 += u_start[i + 1] * v_start[i + 1];
      sum_2 +=  u_start[i + 2] * v_start[i + 2];
   }

   for (; i < v_len; i++) { // n - k + 1
      sum += u_start[i] * v_start[i];
   }
   return sum + sum_1 + sum_2;
} */

/* Loop unrolling 4x4 */
/* Is infact, slower
1.32s to run 1000 tests (1317981.00ns per test)
5.75s to run 1000 baseline tests (5751620.00ns per test)
Speedup: 4.43s to run 1000 tests (4433639.00ns per test)
 */
data_t dotproductfast(vec_ptr u, vec_ptr v) {
   data_t sum = 0, sum_1 = 0, sum_2 = 0, sum_3 = 0, u_val, v_val;
   long v_len = vec_length(u), i; // Remove invariants from loop
   data_t * v_start = get_vec_start(v); // Remove invariant + this is also a memory redundancy
   data_t * u_start = get_vec_start(u);

   for (i = 0; i < v_len - 3; i+=3) { // we can assume both vectors are same length
      sum += u_start[i] * v_start[i];
      sum_1 += u_start[i + 1] * v_start[i + 1];
      sum_2 +=  u_start[i + 2] * v_start[i + 2];
      sum_3 +=  u_start[i + 3] * v_start[i + 3];
   }

   for (; i < v_len; i++) { // n - k + 1
      sum += u_start[i] * v_start[i];
   }
   return sum + sum_1 + sum_2 + sum_3;
}



int main() {
  clock_t test_start, test_end, baseline_start, baseline_end;
  double clock_elapsed, time_elapsed;
  int i;


   // init long vector
  long n = 1000000;
  vec_ptr u = new_vec(n);
  vec_ptr v = new_vec(n);

  for (long i = 0; i < n; i++) {
    set_vec_element(u, i, i + 1);
    set_vec_element(v, i, i + 1);
  }

  baseline_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) { // baseline will be the slowest here
    dotproduct(u, v);
  }
  baseline_end = clock();

  test_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) { // with each optimization test will get faster
    dotproductfast(u, v);
  }
  test_end = clock();

  clock_elapsed = test_end - test_start;
  time_elapsed = clock_elapsed/CLOCKS_PER_SEC;

  printf("%.2fs to run %d tests (%.2fns per test)\n", time_elapsed, TEST_LOOPS, 
         time_elapsed * 1e9 / TEST_LOOPS);

  clock_elapsed = baseline_end - baseline_start;
  time_elapsed = clock_elapsed/CLOCKS_PER_SEC;

  printf("%.2fs to run %d baseline tests (%.2fns per test)\n", time_elapsed, TEST_LOOPS,
         time_elapsed * 1e9 / TEST_LOOPS);


  clock_elapsed = baseline_end - baseline_start - (test_end - test_start);
  time_elapsed = clock_elapsed/CLOCKS_PER_SEC;

  printf("Speedup: %.2fs to run %d tests (%.2fns per test)\n", time_elapsed, TEST_LOOPS, // how much speedup with each test
   time_elapsed * 1e9 / TEST_LOOPS);


  free_vec(u);
  free_vec(v);
  
  return 0;
}
