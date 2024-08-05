package veclib

import "math"

type Vec2Int32 struct {
	X int32
	Y int32
}

type Matrix[T any] struct {
	Rows int32
	Cols int32
	Data []T
}

func (m *Matrix[T]) Get(y int32, x int32) T {
	return m.Data[y*m.Cols+x]
}

func (m *Matrix[T]) Set(y int32, x int32, val T) {
	m.Data[y*m.Cols+x] = val
}

func NewMatrix[T any](rows int32, cols int32) *Matrix[T] {
	return &Matrix[T]{Rows: rows, Cols: cols, Data: make([]T, rows*cols)}
}

func FloorDiv(a float32, b float32) int32 {
	return int32(math.Floor(float64(a / b)))
}

func Mod(a int32, b int32) int32 {
	return (a%b + b) % b
}

func FloatModInt(a float32, b float32) int32 {
	return Mod(int32(a), int32(b))
}

func FloatModFloat(a float32, b float32) float32 {
	return float32(Mod(int32(a), int32(b)))
}
