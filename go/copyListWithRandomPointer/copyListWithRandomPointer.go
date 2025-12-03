package main

type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

func copyRandomList(head *Node) *Node {
	if head == nil {
		return nil
	}

	current := head
	oldToNew := make(map[*Node]*Node)

	for current != nil {
		oldToNew[current] = &Node{Val: current.Val}
		current = current.Next
	}

	current = head
	for current != nil {
		oldToNew[current].Next = oldToNew[current.Next]
		oldToNew[current].Random = oldToNew[current.Random]
		current = current.Next
	}

	return oldToNew[head]
}

func main() {
	head := &Node{Val: 1}
	head.Next = &Node{Val: 2}
	head.Random = head.Next
	head.Next.Random = head
	copiedHead := copyRandomList(head)

	current := copiedHead
	for current != nil {
		print("Node Val: ", current.Val)
		if current.Random != nil {
			print(", Random Val: ", current.Random.Val)
		} else {
			print(", Random Val: nil")
		}
		print("\n")
		current = current.Next
	}
}
