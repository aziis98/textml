package utils

type Stack[T any] struct {
	stack []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{stack: []T{}}
}

func (s *Stack[T]) Top() T {
	return s.stack[len(s.stack)-1]
}

func (s *Stack[T]) Push(value T) {
	s.stack = append(s.stack, value)
}

func (s *Stack[T]) Pop() T {
	value := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return value
}
