package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{0, head}
	adead := dummy
	behind := dummy

	for i := 0; i <= n; i++ {
		adead = adead.Next
	}

	for adead != nil {
		adead = adead.Next
		behind = behind.Next
	}

	behind.Next = behind.Next.Next

	return dummy.Next
}

func main() {
	head := &ListNode{Val: 1}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 3}
	head.Next.Next.Next = &ListNode{Val: 4}
	head.Next.Next.Next.Next = &ListNode{Val: 5}

	n := 2

	newHead := removeNthFromEnd(head, n)
	current := newHead
	for current != nil {
		print(current.Val)
		if current.Next != nil {
			print(" -> ")
		}
		current = current.Next
	}
}
