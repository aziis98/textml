package utils

type Stack[T any] struct {
	stack []T
}

func NewStack[T any](initialValues ...T) *Stack[T] {
	return &Stack[T]{stack: append([]T{}, initialValues...)}
}

func (s *Stack[T]) Top() T {
	return s.stack[len(s.stack)-1]
}

func (s *Stack[T]) Peek(offset int) T {
	return s.stack[len(s.stack)-offset]
}

func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
}

func (s *Stack[T]) Pop() T {
	value := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return value
}
