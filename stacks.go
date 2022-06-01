package main

import "fmt"

// stack represents stack that holds a slice
type Stack struct {
	items []string
}

// Push will add a value at the end
func (s *Stack) Push(i string) {
	s.items = append(s.items, i)
}

// Pop will remove a value at the end
// and RETURNS the removed value
func (s *Stack) Pop() string {
	l := len(s.items) - 1
	toRemove := s.items[l]
	s.items = s.items[:l]
	return toRemove
}

func main() {
	myStack := Stack{}
	fmt.Println(myStack)
	myStack.Push("L3-a")
	myStack.Push("L2-H")
	myStack.Push("L1-End")
	fmt.Println(myStack)
	myStack.Pop()
	fmt.Println(myStack.items)
	myStack.Push("L4-Start")
	for _, i := range myStack.items {
		fmt.Print(i + " ")
	}
}
