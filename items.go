package btree

import (
	"sort"
)

type Items []Item

type Item interface {
	Less(than Item) bool
}

// 如果找到，则返回对应的item索引和ture, 否则返回第一个大于item的索引和false
func (items Items) find(t Item) (index int, found bool) {

	i := sort.Search(len(items), func(i int) bool {
		return t.Less(items[i])
	})
	if i > 0 && !items[i-1].Less(t) {
		return i - 1, true
	}
	return i, false
}

// 清掉指定位置，及其之后位置的Item
func (s *Items) truncate(index int) {
	// var toClear Items
	*s = (*s)[:index]
}

// items指定位置插入一个Item
func (s *Items) insertAt(index int, item Item) {
	*s = append(*s, nil) // 扩展一个位置出来
	if index < len(*s) { // 可能index==len，需要插入的位置就在最后一个
		copy((*s)[index+1:], (*s)[index:]) // 位置后面的往后挪一下
	}
	(*s)[index] = item // 覆盖指定位置的数据
}

// children指定位置删除一个node指针
// func (s *children) removeAt(index int) *node {
// 	n := (*s)[index]
// 	copy((*s)[index:], (*s)[index+1:])
// 	(*s)[len(*s)-1] = nil
// 	*s = (*s)[:len(*s)-1]
// 	return n
// }

func (s *Items) removeAt(index int) (out Item) {
	out = (*s)[index]
	if index+1 < len(*s) {
		copy((*s)[index:], (*s)[index+1:])
	}

	*s = (*s)[:len(*s)-1]
	return
}
