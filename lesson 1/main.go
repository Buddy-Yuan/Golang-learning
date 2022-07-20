package main

import "fmt"

type LinkedListNode struct {
	Val  int
	Next *LinkedListNode
}

func mergeTwoLists(l1 *LinkedListNode, l2 *LinkedListNode) *LinkedListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}
	if l1.Val < l2.Val {
		l1.Next = mergeTwoLists(l1.Next, l2)
		return l1
	}
	l2.Next = mergeTwoLists(l1, l2.Next)
	return l2
}

func (l *LinkedListNode) Traverse() {
	for l != nil {
		fmt.Println(l.Val)
		l = l.Next
	}
}

func main() {
	l1 := &LinkedListNode{
		Val: 1, Next: &LinkedListNode{
			Val: 2, Next: &LinkedListNode{
				Val: 8, Next: nil,
			},
		},
	}
	l2 := &LinkedListNode{
		Val: 1, Next: &LinkedListNode{
			Val: 4, Next: &LinkedListNode{
				Val: 5, Next: &LinkedListNode{
					Val: 6, Next: nil,
				},
			},
		},
	}
	l3 := mergeTwoLists(l1, l2)
	l3.Traverse()
}
