package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	current := dummy

	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			current.Next = list1
			current = list1
			list1 = list1.Next
		} else {
			current.Next = list2
			current = list2
			list2 = list2.Next
		}
	}

	if list1 != nil {
		current.Next = list1
	} else if list2 != nil {
		current.Next = list2
	}

	return dummy.Next
}

func main() {
	list1 := &ListNode{Val: 1}
	list1.Next = &ListNode{Val: 3}
	list1.Next.Next = &ListNode{Val: 5}

	list2 := &ListNode{Val: 2}
	list2.Next = &ListNode{Val: 4}
	list2.Next.Next = &ListNode{Val: 6}

	mergedList := mergeTwoLists(list1, list2)

	current := mergedList
	for current != nil {
		print(current.Val)
		if current.Next != nil {
			print(" -> ")
		}
		current = current.Next
	}
}
