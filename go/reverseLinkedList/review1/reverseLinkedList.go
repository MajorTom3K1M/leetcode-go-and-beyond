package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	previous := (*ListNode)(nil)
	current := head

	for current != nil {
		nextTemp := current.Next
		current.Next = previous
		previous = current
		current = nextTemp
	}

	return previous
}

func main() {
	head := &ListNode{Val: 1}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 3}
	head.Next.Next.Next = &ListNode{Val: 4}
	head.Next.Next.Next.Next = &ListNode{Val: 5}

	newHead := reverseList(head)

	current := newHead
	for current != nil {
		print(current.Val)
		if current.Next != nil {
			print(" -> ")
		}
		current = current.Next
	}
}
