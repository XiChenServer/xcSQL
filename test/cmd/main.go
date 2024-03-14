package main

import (
	"fmt"
)

type SkipListNode struct {
	Value int
	Next  *SkipListNode
}

type SkipList struct {
	Head *SkipListNode
}

func NewSkipList() *SkipList {
	return &SkipList{Head: nil}
}

// 插入节点到跳表中
func (sl *SkipList) Insert(value int) {
	newNode := &SkipListNode{Value: value, Next: nil}

	if sl.Head == nil {
		sl.Head = newNode
		return
	}

	if value < sl.Head.Value {
		newNode.Next = sl.Head
		sl.Head = newNode
		return
	}

	current := sl.Head
	for current.Next != nil && current.Next.Value < value {
		current = current.Next
	}

	newNode.Next = current.Next
	current.Next = newNode
}

// 合并跳表
func MergeSkipLists(skipLists []*SkipList) *SkipList {
	mergedList := NewSkipList()

	for _, skipList := range skipLists {
		current := skipList.Head
		for current != nil {
			mergedList.Insert(current.Value)
			current = current.Next
		}
	}

	return mergedList
}

// 将合并后的跳表分割成多个跳表，每个跳表最多包含 maxSize 个元素
func SplitSkipList(mergedList *SkipList, maxSize int) []*SkipList {
	var result []*SkipList
	current := mergedList.Head

	for current != nil {
		skipList := NewSkipList()
		for i := 0; i < maxSize && current != nil; i++ {
			skipList.Insert(current.Value)
			current = current.Next
		}
		result = append(result, skipList)
	}

	return result
}

// 打印跳表
func (sl *SkipList) Print() {
	current := sl.Head
	for current != nil {
		fmt.Printf("%d ", current.Value)
		current = current.Next
	}
	fmt.Println()
}

func main() {
	// 初始化三个跳表
	skipList1 := NewSkipList()
	skipList2 := NewSkipList()
	skipList3 := NewSkipList()

	// 插入元素到每个跳表
	for _, v := range []int{1, 5, 7, 9} {
		skipList1.Insert(v)
	}
	for _, v := range []int{10, 13, 15, 17} {
		skipList2.Insert(v)
	}
	for _, v := range []int{20, 26, 28, 39} {
		skipList3.Insert(v)
	}

	// 合并跳表
	mergedList := MergeSkipLists([]*SkipList{skipList1, skipList2, skipList3})

	//// 将合并后的跳表分割成每个跳表最多包含 4 个元素
	//splitLists := SplitSkipList(mergedList, 4)
	//
	//// 打印分割后的跳表
	//fmt.Println("分割后的跳表：")
	//for i, skipList := range splitLists {
	//	fmt.Printf("跳表%d：", i+1)
	//	skipList.Print()
	//}

	// 插入新的跳表
	newSkipList := NewSkipList()
	for _, v := range []int{4, 8, 14, 23} {
		newSkipList.Insert(v)
	}

	// 合并新跳表和原先的跳表
	mergedList = MergeSkipLists(append([]*SkipList{skipList1, skipList2, skipList3}, newSkipList))
	// 将合并后的跳表分割成每个跳表最多包含 4 个元素
	splitLists := SplitSkipList(mergedList, 4)

	// 打印分割后的跳表
	fmt.Println("插入新跳表后的分割结果：")
	for i, skipList := range splitLists {
		fmt.Printf("跳表%d：", i+1)
		skipList.Print()
	}
}
