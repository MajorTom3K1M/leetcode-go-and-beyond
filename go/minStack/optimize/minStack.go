package main

type element struct {
	val int
	min int
}

type MinStack struct {
	data []*element
	min  int
}

func Constructor() MinStack {
	return MinStack{
		data: make([]*element, 0),
	}
}

func (this *MinStack) Push(val int) {
	if len(this.data) == 0 {
		this.data = append(this.data, &element{
			val: val,
			min: val,
		})
		return
	}

	if val < this.data[len(this.data)-1].min {
		this.data = append(this.data, &element{
			val: val,
			min: val,
		})
	} else {
		this.data = append(this.data, &element{
			val: val,
			min: this.data[len(this.data)-1].min,
		})
	}
}

func (this *MinStack) Pop() {
	if len(this.data) == 0 {
		return
	}

	this.data[len(this.data)-1] = nil
	this.data = this.data[:len(this.data)-1]
}

func (this *MinStack) Top() int {
	if len(this.data) == 0 {
		return 0
	}
	return (this.data[len(this.data)-1]).val
}

func (this *MinStack) GetMin() int {
	return this.data[len(this.data)-1].min
}

/**
 * Your MinStack object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Push(val);
 * obj.Pop();
 * param_3 := obj.Top();
 * param_4 := obj.GetMin();
 */
