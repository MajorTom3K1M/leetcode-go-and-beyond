package main

type MinStack struct {
	stack    []int // Main stack to store all elements
	minStack []int // Auxiliary stack to keep track of minimums
}

func Constructor() MinStack {
	return MinStack{
		stack:    make([]int, 0), // Initialize empty main stack
		minStack: make([]int, 0), // Initialize empty minimum stack
	}
}

func (this *MinStack) Push(val int) {
	// Always push the new value to main stack
	this.stack = append(this.stack, val)

	// For minStack:
	// If it's empty OR new value is smaller than current minimum
	// Push the new value as the new minimum
	if len(this.minStack) == 0 || val <= this.minStack[len(this.minStack)-1] {
		this.minStack = append(this.minStack, val)
	}
}

func (this *MinStack) Pop() {
	// If stack is empty, return (though problem guarantees this won't happen)
	if len(this.stack) == 0 {
		return
	}

	// If the value being popped is the current minimum,
	// we need to remove it from minStack too
	if this.stack[len(this.stack)-1] == this.minStack[len(this.minStack)-1] {
		this.minStack = this.minStack[:len(this.minStack)-1]
	}

	// Remove the top element from main stack
	this.stack = this.stack[:len(this.stack)-1]
}

func (this *MinStack) Top() int {
	// Return the last element in the stack
	// Problem guarantees stack won't be empty when called
	return this.stack[len(this.stack)-1]
}

func (this *MinStack) GetMin() int {
	// Return the top of minStack which maintains the current minimum
	// Problem guarantees stack won't be empty when called
	return this.minStack[len(this.minStack)-1]
}

/**
 * Your MinStack object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Push(val);
 * obj.Pop();
 * param_3 := obj.Top();
 * param_4 := obj.GetMin();
 */
