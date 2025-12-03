package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode = nil
	var current *ListNode = head
	for current != nil {
		temp := current.Next
		current.Next = prev
		prev = current
		current = temp

	}
	fmt.Println(prev)
	return prev
}

func PrintList(head *ListNode) {
	current := head
	for current != nil {
		if current.Next != nil {
			fmt.Printf("%d -> ", current.Val)
		} else {
			fmt.Printf("%d", current.Val)
		}
		current = current.Next
	}
	fmt.Println()
}

func main() {
	list := ListNode{
		Val: 1,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 3,
			},
		},
	}

	PrintList(reverseList(&list))
}
