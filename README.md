#### HPC_Matrix_Multiplcation

I used Golang to implement the functionality because it takes less time to launch a go routine equivalent to a thread.

##### Sequencial Matrix Multiplication

The `MultiplySerialy` method will multiply two matrices `A` and `B` serially. You will find it in the source code. This is the classical implementation of matrix multiplication. 

`Time complexity: O(N^3)`
`Space complexity: O(N^2) (no extra space other than the result matrix)`

##### Parallel Matrix Multiplication

-   I split the input matrix into smaller blocks. I tried to keep the block size and number of blocks across a row as close as possible. 

-   I do not copy the matrix cells into blocks instead create bock metadata for mapping the position with the input matrix.

-   For each block calculation of the result matrix, I launched the n/block_size number of goroutines (threads). 

-   In the picture below, it is shown how the first block of result matrix C is calculated from the A and B matrix.

`Time complexity: O(blocksInRows^3 + blockSize^3) (here blocksInRows = n/block_size)`
`Space complexity: O(2*blocksInRows^2) (metadata)`


##### Performance Comparison: (duration in microseconds)

Here `n = number of rows in the matrix`

Serial Multiplication takes less time for a small input matrix for n < 49 with blocksize = sqrt(n)
--------------------
```

n =      1
Block size =  1
Serial Multiply duration:    0
Parallel Multiply duration:  7

```



Parallel Multiplication takes less time for a larger input matrix (for n >= 49):
---------------------
```

n =      49
Block size  =  7
Serial Multiply duration:    202
Parallel Multiply duration:  222

```

---------------------
```

n =      64
block size =  8
Serial Multiply duration:    455
Parallel Multiply duration:  306

```

---------------------
```

n =      2500
block size =  50
Serial Multiply duration:    155627318
Parallel Multiply duration:    20585884
```


