package math

type MathUtils struct {
}

func New() *MathUtils {
	return &MathUtils{}
}

func (m *MathUtils) Add(x int, y int) int {
	return x + y
}

func (m *MathUtils) Substract(x int, y int) int {
	return x - y
}

func (m *MathUtils) Multiply(x int, y int) int {
	return x * y
}

func (m *MathUtils) Divide(x int, y int) int {
	if y == 0 {
		panic("can't divide by zero")
	}
	return x / y
}
