package main

type MinStack struct {
	stack    []int
	minStack []int
}

func Constructor() MinStack {
	return MinStack{
		stack:    make([]int, 0),
		minStack: make([]int, 0),
	}
}

func (minStack *MinStack) Push(val int) {
	minStack.stack = append(minStack.stack, val)

	if len(minStack.minStack) == 0 || val <= minStack.minStack[len(minStack.minStack)-1] {
		minStack.minStack = append(minStack.minStack, val)
	}
}

func (minStack *MinStack) Pop() {
	if len(minStack.stack) == 0 {
		return
	}

	if minStack.stack[len(minStack.stack)-1] == minStack.minStack[len(minStack.minStack)-1] {
		minStack.minStack = minStack.minStack[:len(minStack.minStack)-1]
	}

	minStack.stack = minStack.stack[:len(minStack.stack)-1]
}

func (minStack *MinStack) Top() int {
	return minStack.stack[len(minStack.stack)-1]
}

func (minStack *MinStack) GetMin() int {
	return minStack.minStack[len(minStack.minStack)-1]
}
