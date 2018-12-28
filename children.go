package btree

type children []*node

// 清掉指定位置，及其之后位置的node指针
func (s *children) truncate(index int) {
	*s = (*s)[:index]
}

func (s *children) insertAt(index int, n *node) {
	(*s) = append((*s), nil)
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = n
}

func (s *children) removeAt(index int) (n *node) {
	n = (*s)[index]
	if index+1 < len(*s) {
		copy((*s)[index:], (*s)[index+1:])
	}
	*s = (*s)[:len(*s)-1]
	return
}
