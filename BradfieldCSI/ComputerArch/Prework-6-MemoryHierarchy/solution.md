- Note: All exercises performed on Intel Machine on Linux (Lubuntu)

## loop-order

### Observations:

### Og

- While both loops seem to have good temporal locality since they refernce the same array, option_1 has better spatial locality. It accesses adjacent memory addressess, and is like to have better cache utilisation. Hence option_1 will be faster.

- They seem to execute roughly the same number of instructions.

- Perf output:

- option_one: clearly faster:
```
             50.96 msec task-clock                       #    0.993 CPUs utilized             
                 0      context-switches                 #    0.000 /sec                      
                 0      cpu-migrations                   #    0.000 /sec                      
            15,677      page-faults                      #  307.604 K/sec                     
      15,72,34,077      cycles                           #    3.085 GHz                       
      25,64,28,497      instructions                     #    1.63  insn per cycle            
       3,22,05,924      branches                         #  631.925 M/sec                     
            40,184      branch-misses                    #    0.12% of all branches           

       0.051304759 seconds time elapsed

       0.027657000 seconds user
       0.023706000 seconds sys
```

- option_two: clearly slower:
```
            126.24 msec task-clock                       #    0.998 CPUs utilized             
                 2      context-switches                 #   15.843 /sec                      
                 1      cpu-migrations                   #    7.921 /sec                      
            15,676      page-faults                      #  124.175 K/sec                     
      39,00,95,570      cycles                           #    3.090 GHz                       
      24,98,88,589      instructions                     #    0.64  insn per cycle            
       3,08,85,931      branches                         #  244.658 M/sec                     
            53,716      branch-misses                    #    0.17% of all branches           

       0.126555446 seconds time elapsed

       0.093908000 seconds user
       0.032663000 seconds sys
```


### -O1

- Performance for both is suddenly similar:

- option_one:
```
 Performance counter stats for './a.out':

              0.45 msec task-clock                       #    0.627 CPUs utilized             
                 0      context-switches                 #    0.000 /sec                      
                 0      cpu-migrations                   #    0.000 /sec                      
                52      page-faults                      #  114.803 K/sec                     
         13,80,759      cycles                           #    3.048 GHz                       
         10,63,030      instructions                     #    0.77  insn per cycle            
          1,92,772      branches                         #  425.593 M/sec                     
             6,615      branch-misses                    #    3.43% of all branches           

       0.000721929 seconds time elapsed

       0.000788000 seconds user
       0.000000000 seconds sys
```

- option_two
```
Performance counter stats for './a.out':

              0.48 msec task-clock                       #    0.651 CPUs utilized             
                 0      context-switches                 #    0.000 /sec                      
                 0      cpu-migrations                   #    0.000 /sec                      
                53      page-faults                      #  111.190 K/sec                     
         13,89,172      cycles                           #    2.914 GHz                       
         10,74,025      instructions                     #    0.77  insn per cycle            
          1,94,314      branches                         #  407.655 M/sec                     
             6,586      branch-misses                    #    3.39% of all branches           

       0.000732270 seconds time elapsed

       0.000817000 seconds user
       0.000000000 seconds sys

```

- The assembly for -O1 is as follows, it simply runs the loop without any ops!
```
option_one:
        mov     edx, 4000
.L2:
        mov     eax, 4000
.L3:
        sub     eax, 1
        jne     .L3
        sub     edx, 1
        jne     .L2
        ret
option_two:
        mov     edx, 4000
.L7:
        mov     eax, 4000
.L8:
        sub     eax, 1
        jne     .L8
        sub     edx, 1
        jne     .L7
        ret
main:
        mov     eax, 0
        ret
```


### Slight change:


- If you change the functions to something like this, -O1 will again give the expected results


```
void option_three(int (*x)[4000]) {
  int i, j;
  // static int x[4000][4000];
  for (i = 0; i < 4000; i++) {
    for (j = 0; j < 4000; j++) {
      x[i][j] = i + j;
    }
  }
}

void option_four(int (*x)[4000]) {
  int i, j;
  // static int x[4000][4000];
  for (i = 0; i < 4000; i++) {
    for (j = 0; j < 4000; j++) {
      x[j][i] = i + j;
      // printf("%d", x[0]);
    }
  }
}

int main() {
  int (*x)[4000] = malloc(sizeof(int[4000][4000]));

  if (x == NULL) {
    fprintf(stderr, "Memory allocation failed.\n");
    return 1;
  }


  // option_one();
  // option_two();
  option_three(x);
  // option_four(x);
  return 0;
}

```


- option_three:

```

 Performance counter stats for './a.out':

             45.88 msec task-clock                       #    0.992 CPUs utilized             
                 1      context-switches                 #   21.798 /sec                      
                 0      cpu-migrations                   #    0.000 /sec                      
            15,679      page-faults                      #  341.775 K/sec                     
      14,18,20,248      cycles                           #    3.091 GHz                       
      16,72,55,144      instructions                     #    1.18  insn per cycle            
       3,04,40,198      branches                         #  663.544 M/sec                     
            38,371      branch-misses                    #    0.13% of all branches           

       0.046255030 seconds time elapsed

       0.007716000 seconds user
       0.038583000 seconds sys
```

- option_four()
```

 Performance counter stats for './a.out':

            128.69 msec task-clock                       #    0.995 CPUs utilized             
                20      context-switches                 #  155.417 /sec                      
                 0      cpu-migrations                   #    0.000 /sec                      
            15,679      page-faults                      #  121.839 K/sec                     
      39,72,41,082      cycles                           #    3.087 GHz                       
      16,73,75,966      instructions                     #    0.42  insn per cycle            
       3,04,63,186      branches                         #  236.725 M/sec                     
            42,220      branch-misses                    #    0.14% of all branches           

       0.129349919 seconds time elapsed

       0.096762000 seconds user
       0.032254000 seconds sys


```


### Cachegrind

- valgrind --tool=cachegrind --cache-sim=yes ./a.out for both option_one and option_two

- Note: cg_annotate cachegrind.out.62677 --annotate gives an awesome large annotation!

- option_one()

```
==61640== Cachegrind, a cache and branch-prediction profiler
==61640== Copyright (C) 2002-2017, and GNU GPL'd, by Nicholas Nethercote et al.
==61640== Using Valgrind-3.21.0 and LibVEX; rerun with -h for copyright info
==61640== Command: ./a.out
==61640== 
--61640-- warning: L3 cache found, using its data for the LL simulation.
==61640== 
==61640== I refs:        240,173,912
==61640== I1  misses:          1,107
==61640== LLi misses:          1,091
==61640== I1  miss rate:        0.00%
==61640== LLi miss rate:        0.00%
==61640== 
==61640== D refs:        112,062,339  (96,046,094 rd   + 16,016,245 wr)
==61640== D1  misses:      1,002,173  (     1,584 rd   +  1,000,589 wr)
==61640== LLd misses:      1,001,906  (     1,357 rd   +  1,000,549 wr)
==61640== D1  miss rate:         0.9% (       0.0%     +        6.2%  )
==61640== LLd miss rate:         0.9% (       0.0%     +        6.2%  )
==61640== 
==61640== LL refs:         1,003,280  (     2,691 rd   +  1,000,589 wr)
==61640== LL misses:       1,002,997  (     2,448 rd   +  1,000,549 wr)
==61640== LL miss rate:          0.3% (       0.0%     +        6.2%  )
```



```
==61780== Cachegrind, a cache and branch-prediction profiler
==61780== Copyright (C) 2002-2017, and GNU GPL'd, by Nicholas Nethercote et al.
==61780== Using Valgrind-3.21.0 and LibVEX; rerun with -h for copyright info
==61780== Command: ./a.out
==61780== 
--61780-- warning: L3 cache found, using its data for the LL simulation.
==61780== 
==61780== I refs:        240,173,912
==61780== I1  misses:          1,107
==61780== LLi misses:          1,091
==61780== I1  miss rate:        0.00%
==61780== LLi miss rate:        0.00%
==61780== 
==61780== D refs:        112,062,339  (96,046,094 rd   + 16,016,245 wr)
==61780== D1  misses:     16,002,173  (     1,584 rd   + 16,000,589 wr)
==61780== LLd misses:      1,001,906  (     1,357 rd   +  1,000,549 wr)
==61780== D1  miss rate:        14.3% (       0.0%     +       99.9%  )
==61780== LLd miss rate:         0.9% (       0.0%     +        6.2%  )
==61780== 
==61780== LL refs:        16,003,280  (     2,691 rd   + 16,000,589 wr)
==61780== LL misses:       1,002,997  (     2,448 rd   +  1,000,549 wr)
==61780== LL miss rate:          0.3% (       0.0%     +        6.2%  )
```


### Does cache match our expectations?


- Let's take a look at our cache settings using 'lscpu | grep cache'
```
L1d cache:                       64 KiB (2 instances)
L1i cache:                       64 KiB (2 instances)
L2 cache:                        512 KiB (2 instances)
L3 cache:                        3 MiB (1 instance)
```

- Let's take a look at our cache line / cache block size:
```
getconf LEVEL1_ICACHE_LINESIZE
getconf LEVEL1_DCACHE_LINESIZE
getconf LEVEL2_CACHE_LINESIZE
getconf LEVEL3_CACHE_LINESIZE

All return 64 (bytes)
```

- With this information lets try to predict the number of cache misses.

- For option_one(), we know that we are accessing contiguous integer data to write:
- 4000 * 4000 = 16 * 10^6 writes
- Size of an int is 4 bytes and one line of cache = 64 bytes. So every line of cache will allow us access to 16 ints until a miss.
- Therefore number of cache misses = 16 million / 16 = 1 million write misses which is bang on the money!

- Further evidence is that if we change the loop to access 2000 * 2000 elements, the Drefs Writes ~ 4 million and cache write misses ~ 250000 which is again very very close.



- Another tiny example, we create a function option_five():
```
void option_five() {
  int i, j;
  static int x[4000][8];
  for (i = 0; i < 8; i++) {
    for (j = 0; j < 4000; j++) {
      x[j][i] = i + j;
    }
  }
}

Cachegrind results:
==70890== I refs:        625,968
==70890== I1  misses:      1,107
==70890== LLi misses:      1,089
==70890== I1  miss rate:    0.18%
==70890== LLi miss rate:    0.17%
==70890== 
==70890== D refs:        270,371  (226,118 rd   + 44,253 wr)
==70890== D1  misses:     18,173  (  1,584 rd   + 16,589 wr)
==70890== LLd misses:      3,835  (  1,296 rd   +  2,539 wr)
==70890== D1  miss rate:     6.7% (    0.7%     +   37.5%  )
==70890== LLd miss rate:     1.4% (    0.6%     +    5.7%  )
==70890== 
==70890== LL refs:        19,280  (  2,691 rd   + 16,589 wr)
==70890== LL misses:       4,924  (  2,385 rd   +  2,539 wr)
==70890== LL miss rate:      0.5% (    0.3%     +    5.7%  )
```

- We have 32000 ints, since one line of cache holds 16 ints(64 bytes), we can load 2 rows in cache at a time. But remember, we are accessing column wise, so every alternte element would be a cache miss ~ 16000 cache write missess which corresponds to the data.

- LLd cache misses do not correspond to the data, why? Is it because sometime prior the data has already been fetched in memory??



## Matrix Multiplication
- A normal run on our hardware yields the following results:
```
Naive: 31.306s
Fast: 31.158s
1.00x speedup
```

### Cache miss expectation

- For 512 * 512 matrix multiplcation of size double, cache misses estimate would be as follows.
- Consider the line:
```
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
```
- Considering D1 misses.
- We are reading C row wise, Read misses for C[i][j] would be one every 8 element. Since there are total 512 * 512 = 262144 elements, 1 element = 8 bytes. 1 line of cache = 64 bytes, we load 8 elements at a time. Cache read misses = 262144/8 ~ 32768.
- On the next line, we are reading A[i][k] and B[k][j],
A[i][k] is being read row wise, total times it is being read => for each elem C[i][j] (512 * 512), we will have 512 reads in a row. 512 ^ 3 = 134,217,728. Let's take a ballpark and estimate that there is a cache miss once every 8th element, although it is more (since reading B[i][j] pushes stuff out of the cache as well). = 134,217,728 / 8 = 16777216 misses
For each solution element C[i][j](512 * 512), we will have 512 reads in a column k, so total reads = 512 * 512 * 512 = 134,217,728, we consider all misses (since moving column wise). So total misses = 134,217,728 + 16777216 = 150994944. 

```
After commenting out fast_multiply from benchmark.c

cc -g -Wall matrix-multiply.c benchmark.c

cg_annotate cachegrind.out.74487 --annotate

This gives us:

D1mw           DLmw
32,896 (11.2%) 32,896 (11.2%)        C[i][j] = 0; // agrees with our result

D1mr                   DLmr
167,278,077 (100.0%)   65,922 (95.5%)  C[i][j] = A[i][k] * B[k][j] // 1st level read agrees with our result, last level reads will ofcourse be lesser




```

- The more complex the program, the harder it becomes to estimate this (?), since there are now multiple accessess.



### Cache Improvement

- To lower cache misses, let's transpose RHS matrix so we can simply multiply corresponding elements of matrices to get results.
- so now let's consider the ballmark that we will have cache misses only once every 8 elements for both lhs and rhs.
- Since total accesses = 512 * 512 * 512 (for each C[i][j], we access 512 column elements B[j][k]), let's divide this by 8
- D1Mr = 134217728 / 8 = 16777216 (maybe slightly more) -> data agrees -> check results below
- Lets say for transposing, we have 512 * 512 / 2 elements to transpose, if we have all cache misses = 512 * 256 ~ 131 k misses, however actual figure will be lesser (data shows 78297 + 16491 = 94788) since there are hits ofc

```
On timing:

Naive: 0.987s
Fast: 0.588s
1.68x speedup
```


- On cachegrind

```
Transpose matrix cache misses f
D1mr             DLmr
16,491  (0.1%)  16,130 (23.3%)          double temp = A[i][j];
78,297  (0.5%)  16,831 (24.3%)          A[i][j] = A[j][i];
0                                       A[j][i] = temp



Fast multiply cache misses

D1mr                   DLmr
16,908,993 (99.4%)     32,961 (47.7%)       c[i][j] += a[i][k] * b[j][k];

D1mw                   DLmr
32,897 (11.2%)        32,897 (11.2%)        c[i][j] = 0





- Tiling is another solution (did not check on cachegrind) to reduce cache misses further

```