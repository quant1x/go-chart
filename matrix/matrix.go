package matrix

import (
	"bytes"
	"errors"
	"math"
)

const (
	// DefaultEpsilon represents the minimum precision for matrix math operations.
	DefaultEpsilon = 0.000001
)

var (
	// ErrDimensionMismatch is a typical error.
	ErrDimensionMismatch = errors.New("dimension mismatch")
)

// New returns a new matrix.
func New(rows, cols int, values ...float64) *Matrix {
	if len(values) == 0 {
		return &Matrix{
			stride:   cols,
			epsilon:  DefaultEpsilon,
			elements: make([]float64, rows*cols),
		}
	}
	elems := make([]float64, rows*cols)
	copy(elems, values)
	return &Matrix{
		stride:   cols,
		epsilon:  DefaultEpsilon,
		elements: elems,
	}
}

// Identity returns the identity matrix of a given order.
func Identity(order int) *Matrix {
	m := New(order, order)
	for i := 0; i < order; i++ {
		m.Set(i, i, 1)
	}
	return m
}

// Zero returns a matrix of a given size zeroed.
func Zero(rows, cols int) *Matrix {
	return New(rows, cols)
}

// Ones returns an matrix of ones.
func Ones(rows, cols int) *Matrix {
	ones := make([]float64, rows*cols)
	for i := 0; i < (rows * cols); i++ {
		ones[i] = 1
	}

	return &Matrix{
		stride:   cols,
		epsilon:  DefaultEpsilon,
		elements: ones,
	}
}

// Eye returns the eye matrix.
func Eye(n int) *Matrix {
	m := Zero(n, n)
	for i := 0; i < len(m.elements); i += n + 1 {
		m.elements[i] = 1
	}
	return m
}

// NewFromArrays creates a matrix from a jagged array set.
func NewFromArrays(a [][]float64) *Matrix {
	rows := len(a)
	if rows == 0 {
		return nil
	}
	cols := len(a[0])
	m := New(rows, cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			m.Set(row, col, a[row][col])
		}
	}
	return m
}

// Matrix represents a 2d dense array of floats.
type Matrix struct {
	epsilon  float64
	elements []float64
	stride   int
}

// String returns a string representation of the matrix.
func (m *Matrix) String() string {
	buffer := bytes.NewBuffer(nil)
	rows, cols := m.Size()

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			buffer.WriteString(f64s(m.Get(row, col)))
			buffer.WriteRune(' ')
		}
		buffer.WriteRune('\n')
	}
	return buffer.String()
}

// Epsilon returns the maximum precision for math operations.
func (m *Matrix) Epsilon() float64 {
	return m.epsilon
}

// WithEpsilon sets the epsilon on the matrix and returns a reference to the matrix.
func (m *Matrix) WithEpsilon(epsilon float64) *Matrix {
	m.epsilon = epsilon
	return m
}

// Each applies the action to each element of the matrix in
// rows => cols order.
func (m *Matrix) Each(action func(row, cow int, value float64)) {
	rows, cols := m.Size()
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			action(row, col, m.Get(row, col))
		}
	}
}

// Round rounds all the values in a matrix to it epsilon,
// returning a reference to the original
func (m *Matrix) Round() *Matrix {
	rows, cols := m.Size()
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			m.Set(row, col, roundToEpsilon(m.Get(row, col), m.epsilon))
		}
	}
	return m
}

// Arrays returns the matrix as a two dimensional jagged array.
func (m *Matrix) Arrays() [][]float64 {
	rows, cols := m.Size()
	a := make([][]float64, rows)

	for row := 0; row < rows; row++ {
		a[row] = make([]float64, cols)

		for col := 0; col < cols; col++ {
			a[row][col] = m.Get(row, col)
		}
	}
	return a
}

// Size returns the dimensions of the matrix.
func (m *Matrix) Size() (rows, cols int) {
	rows = len(m.elements) / m.stride
	cols = m.stride
	return
}

// IsSquare returns if the row count is equal to the column count.
func (m *Matrix) IsSquare() bool {
	return m.stride == (len(m.elements) / m.stride)
}

// IsSymmetric returns if the matrix is symmetric about its diagonal.
func (m *Matrix) IsSymmetric() bool {
	rows, cols := m.Size()

	if rows != cols {
		return false
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < i; j++ {
			if m.Get(i, j) != m.Get(j, i) {
				return false
			}
		}
	}
	return true
}

// Get returns the element at the given row, col.
func (m *Matrix) Get(row, col int) float64 {
	index := (m.stride * row) + col
	return m.elements[index]
}

// Set sets a value.
func (m *Matrix) Set(row, col int, val float64) {
	index := (m.stride * row) + col
	m.elements[index] = val
}

// Col returns a column of the matrix as a vector.
func (m *Matrix) Col(col int) Vector {
	rows, _ := m.Size()
	values := make([]float64, rows)
	for row := 0; row < rows; row++ {
		values[row] = m.Get(row, col)
	}
	return Vector(values)
}

// Row returns a row of the matrix as a vector.
func (m *Matrix) Row(row int) Vector {
	_, cols := m.Size()
	values := make([]float64, cols)
	for col := 0; col < cols; col++ {
		values[col] = m.Get(row, col)
	}
	return Vector(values)
}

// Copy returns a duplicate of a given matrix.
func (m *Matrix) Copy() *Matrix {
	m2 := &Matrix{stride: m.stride, epsilon: m.epsilon, elements: make([]float64, len(m.elements))}
	copy(m2.elements, m.elements)
	return m2
}

// DiagonalVector returns a vector from the diagonal of a matrix.
func (m *Matrix) DiagonalVector() Vector {
	rows, cols := m.Size()
	rank := minInt(rows, cols)
	values := make([]float64, rank)

	for index := 0; index < rank; index++ {
		values[index] = m.Get(index, index)
	}
	return Vector(values)
}

// Diagonal returns a matrix from the diagonal of a matrix.
func (m *Matrix) Diagonal() *Matrix {
	rows, cols := m.Size()
	rank := minInt(rows, cols)
	m2 := New(rank, rank)

	for index := 0; index < rank; index++ {
		m2.Set(index, index, m.Get(index, index))
	}
	return m2
}

// Equals returns if a matrix equals another matrix.
func (m *Matrix) Equals(other *Matrix) bool {
	if other == nil && m != nil {
		return false
	} else if other == nil {
		return true
	}

	if m.stride != other.stride {
		return false
	}

	msize := len(m.elements)
	m2size := len(other.elements)

	if msize != m2size {
		return false
	}

	for i := 0; i < msize; i++ {
		if m.elements[i] != other.elements[i] {
			return false
		}
	}
	return true
}

// L returns the matrix with zeros below the diagonal.
func (m *Matrix) L() *Matrix {
	rows, cols := m.Size()
	m2 := New(rows, cols)
	for row := 0; row < rows; row++ {
		for col := row; col < cols; col++ {
			m2.Set(row, col, m.Get(row, col))
		}
	}
	return m2
}

// U returns the matrix with zeros above the diagonal.
// Does not include the diagonal.
func (m *Matrix) U() *Matrix {
	rows, cols := m.Size()
	m2 := New(rows, cols)
	for row := 0; row < rows; row++ {
		for col := 0; col < row && col < cols; col++ {
			m2.Set(row, col, m.Get(row, col))
		}
	}
	return m2
}

// math operations

// Multiply multiplies two matrices.
func (m *Matrix) Multiply(m2 *Matrix) (m3 *Matrix, err error) {
	if m.stride*m2.stride != len(m2.elements) {
		return nil, ErrDimensionMismatch
	}

	m3 = &Matrix{epsilon: m.epsilon, stride: m2.stride, elements: make([]float64, (len(m.elements)/m.stride)*m2.stride)}
	for m1c0, m3x := 0, 0; m1c0 < len(m.elements); m1c0 += m.stride {
		for m2r0 := 0; m2r0 < m2.stride; m2r0++ {
			for m1x, m2x := m1c0, m2r0; m2x < len(m2.elements); m2x += m2.stride {
				m3.elements[m3x] += m.elements[m1x] * m2.elements[m2x]
				m1x++
			}
			m3x++
		}
	}
	return
}

// Decompositions

// QR performs the qr decomposition.
func (m *Matrix) QR() (q, r *Matrix) {
	defer func() {
		q = q.Round()
		r = r.Round()
	}()

	rows, cols := m.Size()
	qr := m.Copy()
	q = New(rows, cols)
	r = New(rows, cols)

	var i, j, k int
	var norm, s float64

	for k = 0; k < cols; k++ {
		norm = 0
		for i = k; i < rows; i++ {
			norm = math.Hypot(norm, qr.Get(i, k))
		}

		if norm != 0 {
			if qr.Get(k, k) < 0 {
				norm = -norm
			}

			for i = k; i < rows; i++ {
				qr.Set(i, k, qr.Get(i, k)/norm)
			}
			qr.Set(k, k, qr.Get(k, k)+1.0)

			for j = k + 1; j < cols; j++ {
				s = 0
				for i = k; i < rows; i++ {
					s += qr.Get(i, k) * qr.Get(i, j)
				}
				s = -s / qr.Get(k, k)
				for i = k; i < rows; i++ {
					qr.Set(i, j, qr.Get(i, j)+s*qr.Get(i, k))

					if i < j {
						r.Set(i, j, qr.Get(i, j))
					}
				}

			}
		}

		r.Set(k, k, -norm)

	}

	//Q Matrix:
	i, j, k = 0, 0, 0

	for k = cols - 1; k >= 0; k-- {
		q.Set(k, k, 1.0)
		for j = k; j < cols; j++ {
			if qr.Get(k, k) != 0 {
				s = 0
				for i = k; i < rows; i++ {
					s += qr.Get(i, k) * q.Get(i, j)
				}
				s = -s / qr.Get(k, k)
				for i = k; i < rows; i++ {
					q.Set(i, j, q.Get(i, j)+s*qr.Get(i, k))
				}
			}
		}
	}

	return
}
