package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func hasCycle(head *ListNode) bool {
	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}

	return false
}

func main() {
	head := &ListNode{Val: 3}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 0}
	head.Next.Next.Next = &ListNode{Val: -4}
	head.Next.Next.Next.Next = head.Next

	if hasCycle(head) {
		println("Cycle detected")
	} else {
		println("No cycle")
	}

	head2 := &ListNode{Val: 1}
	head2.Next = &ListNode{Val: 2}
	if hasCycle(head2) {
		println("Cycle detected")
	} else {
		println("No cycle")
	}
}
