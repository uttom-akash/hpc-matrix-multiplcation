package main

import (
	"fmt"
	"sync"
	"time"
)

type matrixType [][]int

// helper method to create matrix of dimension rows*rows with default value
func NewEmptyMatrix(rows int) matrixType {
	matrix := make(matrixType, rows)
	for i := range matrix {
		matrix[i] = make([]int, rows)
	}
	return matrix
}

// shared object
var InputMatrixA, InputMatrixB matrixType

// serial multiplication
func MultiplySerialy(matA matrixType, matB matrixType) matrixType {
	n := len(matA)
	matC := NewEmptyMatrix(n)

	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			for k := 0; k < n; k++ {
				matC[row][col] += matA[row][k] * matB[k][col]
			}
		}
	}

	return matC
}

// parallel matrix multiplication

// block metadata contains position of the block in original input matrix
type BlockMetadata struct {
	startRow int
	startCol int
}

func NewBlockMetadata(blocksize int, rowoffset int, coloffset int) *BlockMetadata {
	return &BlockMetadata{
		startRow: rowoffset,
		startCol: coloffset,
	}
}

// BlockMatrix contains the block metadata as cell in matrix
type BlockMatrix struct {
	Matrix    [][]*BlockMetadata
	BlockSize int
}

func NewEmptyBlockMatrix(rows int) [][]*BlockMetadata {
	chunkedMatrix := make([][]*BlockMetadata, rows)
	for i := range chunkedMatrix {
		chunkedMatrix[i] = make([]*BlockMetadata, rows)
	}
	return chunkedMatrix
}

// Create a block matrix containing blockmetadata as cell
func NewBlockMatrix(inputMatrix matrixType, blocksize int) *BlockMatrix {
	matrixDimention := len(inputMatrix)
	blocksInRow := matrixDimention / blocksize

	blockMetadataMatrix := &BlockMatrix{
		Matrix:    NewEmptyBlockMatrix(blocksInRow),
		BlockSize: blocksize,
	}

	// split the inputMatrix and update metadata in block matrix
	blockMetadataMatrix.UpdateBlockMetadataBySplitting(inputMatrix, matrixDimention)

	return blockMetadataMatrix
}

// split the input matrix into chunks and update the block metadata
func (blockMatrix *BlockMatrix) UpdateBlockMetadataBySplitting(inputMatrix matrixType, matrixDimention int) *BlockMatrix {

	for matRow := 0; matRow < matrixDimention; matRow += blockMatrix.BlockSize {
		for matCol := 0; matCol < matrixDimention; matCol += blockMatrix.BlockSize {

			//
			blockInRow, blockInCol := blockMatrix.calculateBlockPosition(matRow, matCol, blockMatrix.BlockSize)

			blockMatrix.Matrix[blockInRow][blockInCol] = NewBlockMetadata(blockMatrix.BlockSize, matRow, matCol)
		}
	}

	return blockMatrix
}

func (blockMatrix *BlockMatrix) calculateBlockPosition(matRow int, matCol int, blocksize int) (int, int) {
	// calculate the row number of the block from the row number of the cell in the input matrix
	blockInRow := matRow / blocksize
	blockInCol := matCol / blocksize
	return blockInRow, blockInCol
}

func (blockMatrixA *BlockMatrix) MultiplyParallely(blockMatrixB *BlockMatrix) matrixType {
	blocksInRows := len(blockMatrixA.Matrix)

	//result matrix
	matrixC := NewEmptyMatrix(blockMatrixA.BlockSize * blocksInRows)

	var wg sync.WaitGroup

	for blockI := 0; blockI < blocksInRows; blockI++ {
		for blockJ := 0; blockJ < blocksInRows; blockJ++ {
			for blockIjIter := 0; blockIjIter < blocksInRows; blockIjIter++ {
				wg.Add(1)

				// launch separate process and run in parallel
				go func(blockI, blockJ, blockIjIter int) {
					defer wg.Done()

					blockSize := blockMatrixA.BlockSize
					// extract blockA and blockB to multiply
					blockA := blockMatrixA.Matrix[blockI][blockIjIter]
					blockB := blockMatrixB.Matrix[blockIjIter][blockJ]

					// multiply block A and B
					for i := 0; i < blockSize; i++ {
						for j := 0; j < blockSize; j++ {
							for k := 0; k < blockSize; k++ {

								// calculate the original result matrix cell position from block position
								matCr := blockI*blockSize + i
								matCc := blockJ*blockSize + j
								matrixC[matCr][matCc] += InputMatrixA[blockA.startRow+i][blockA.startCol+k] * InputMatrixB[blockB.startRow+k][blockB.startCol+j]
							}
						}
					}
				}(blockI, blockJ, blockIjIter)
			}
		}
	}
	wg.Wait()

	return matrixC
}

func main() {

	// go run main.go

	InputMatrixB = matrixType{
		{3, 10, 12, 18},
		{12, 1, 4, 9},
		{9, 10, 12, 2},
		{3, 12, 4, 10}}

	InputMatrixA = matrixType{
		{5, 7, 9, 10},
		{2, 3, 3, 8},
		{8, 10, 2, 3},
		{3, 3, 4, 8}}

	serialResult := MultiplySerialy(InputMatrixA, InputMatrixB)

	blockMatrixA := NewBlockMatrix(InputMatrixA, 2)
	blockMatrixB := NewBlockMatrix(InputMatrixB, 2)
	parallelResult := blockMatrixA.MultiplyParallely(blockMatrixB)

	fmt.Println(serialResult)
	fmt.Println(parallelResult)

	// performance

	for blocksize := 1; blocksize <= 50; blocksize++ {
		dim := blocksize * blocksize

		fmt.Println("---------------------")
		fmt.Println("n =     ", dim)
		fmt.Println("block = ", blocksize)

		InputMatrixA = NewEmptyMatrix(dim)

		InputMatrixB = NewEmptyMatrix(dim)

		// start serial multiplication
		start := time.Now()

		MultiplySerialy(InputMatrixA, InputMatrixB)

		// end serial multiplication
		end := time.Now()
		serialDuration := end.Sub(start)
		fmt.Println("SerialMultiply duration:   ", serialDuration.Microseconds())

		// start parallel multiplication
		start = time.Now()

		blockMatrixA := NewBlockMatrix(InputMatrixA, blocksize)
		blockMatrixB := NewBlockMatrix(InputMatrixB, blocksize)
		blockMatrixA.MultiplyParallely(blockMatrixB)

		// end parallel multiplication
		end = time.Now()
		parallelDuration := end.Sub(start)
		fmt.Println("parallelMultiply duration: ", parallelDuration.Microseconds())
	}
}
