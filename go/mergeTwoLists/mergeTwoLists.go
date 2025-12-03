package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func isEmptyListNode(node *ListNode) bool {
	return node != nil && node.Val == 0 && node.Next == nil
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	if list1 != nil && *list1 == (ListNode{}) || list1 == nil {
		return list2
	}

	if list2 != nil && *list2 == (ListNode{}) || list2 == nil {
		return list1
	}

	head1 := list1
	head2 := list2
	if head1.Val > head2.Val {
		temp := list1
		head1 = list2
		head2 = list1
		list1 = list2
		list2 = temp
	}
	for head1.Next != nil {
		temp1 := head1.Next
		if head1.Val <= head2.Val && head1.Next.Val > head2.Val {
			head1.Next = head2

			temp2 := head2.Next
			head2.Next = temp1
			head2 = temp2
		}
		head1 = temp1
	}

	if head2 != nil {
		head1.Next = head2
	}

	return list1
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
		Val:  2,
		Next: nil,
		// Val: 1,
		// Next: &ListNode{
		// 	Val: 2,
		// 	Next: &ListNode{
		// 		Val: 4,
		// 	},
		// },
	}

	list2 := ListNode{
		Val:  1,
		Next: nil,
		// Val: 1,
		// Next: &ListNode{
		// 	Val: 3,
		// 	Next: &ListNode{
		// 		Val: 4,
		// 	},
		// },
	}

	// 	list1 =
	// [2]
	// list2 =
	// [1]

	PrintList(mergeTwoLists(&list, &list2))
}
