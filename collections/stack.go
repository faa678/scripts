package collections

import "errors"

type Stack struct {
	Slice []interface{}
	Size  int64
	Cap   int64
}

func NewStack(cap int64) *Stack {
	stack := new(Stack)
	stack.Cap = cap
	return stack
}

func (s *Stack) IsEmpty() bool {
	if len(s.Slice) == 0 {
		return true
	}
	return false
}

func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		return nil
	}
	return s.Slice[s.SizeOf()-1]
}

func (s *Stack) Offer(ele interface{}) error {
	if s.CapOf() > 0 && s.SizeOf() == s.CapOf() {
		return errors.New("stack full")
	}
	s.Slice = append(s.Slice, ele)
	s.Size++
	return nil
}

func (s *Stack) Poll() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("stack empty")
	}
	peek := s.Peek()
	s.Slice = s.Slice[:s.SizeOf()-1]
	s.Size--
	return peek, nil
}

func (s *Stack) SizeOf() int64 {
	return s.Size
}

func (s *Stack) CapOf() int64 {
	return s.Cap
}
