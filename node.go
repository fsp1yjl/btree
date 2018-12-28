package btree

import "fmt"

type node struct {
	items    Items
	children children
}

func (n *node) insert(item Item, maxCap int) Item {
	index, found := n.items.find(item)

	if found {
		return item
	}

	if len(n.children) == 0 {
		n.items.insertAt(index, item)
		return nil
	}

	// 不是叶子节点，看看i处的子Node是否需要分裂
	if n.maybeSplitChild(index, maxCap) {
		// 分裂了，导致当前node的变化，需要重新定位i
		inTree := n.items[index] // 获取新升级的item
		switch {
		case item.Less(inTree):
			// 要插入的item比分裂产生的item小，i没改变
		case inTree.Less(item):
			index++ // 要插入的item比分裂产生的item大，i++
		default:
			// 分裂升level的item和插入的item一致，替换
			out := n.items[index]
			n.items[index] = item
			return out
		}
	}

	return n.children[index].insert(item, maxCap)
}

func (n *node) maybeSplitChild(childIndex int, maxCap int) bool {
	l := len(n.children[childIndex].items)

	if l >= maxCap {
		upItem, newNode := n.children[childIndex].split(l / 2)

		n.items.insertAt(childIndex, upItem)
		n.children.insertAt(childIndex+1, newNode)

		return true
	} else {
		return false
	}

}

func (n *node) split(mid int) (Item, *node) {

	if len(n.items)-1 < mid {
		panic("error index")
	}
	upItem := n.items[mid]
	newNode := &node{}

	newNode.items = append(newNode.items, n.items[mid+1:]...)
	if n.children != nil {
		newNode.children = append(newNode.children, n.children[mid+1:]...)
	}

	n.items = n.items[:mid]
	if n.children != nil {
		n.children = n.children[:mid+1]
	}
	return upItem, newNode
}

func (n *node) remove(t Item, minCap int) Item {

	i, found := n.items.find(t)

	// 如果node节点的子树i中items小于等于最小键个数，则先对其进行一次调整处理, 保证到叶子节点时其可以直接删除

	// 如果子树i需要调整，则调整节点后，重新进行remove操作
	if len(n.children) != 0 && len(n.children[i].items) <= minCap {
		return n.growChildAndRemove(i, t, minCap)
	}

	//
	if found {
		if len(n.children) == 0 {
			// 如果查找到元素，且元素在叶子节点，则直接删除并返回
			// 因为在进入叶子节点之前，已经通过growChildAndRemove做了扩容处理，保证到达叶子是，一定可以直接删除而不破坏结构
			return n.items.removeAt(i)
		} else {
			// 从n.childrend[i]为根的子树依次向下，最终删除一个最大值，替换n.items[i-1]原来的值

			// 如果查找到要删除的元素不在叶子节点上， 则删除之，并从children[i]为根的子树上提取最大元素放入要删除的位置
			out := n.items[i]
			n.items[i] = n.children[i].removeMax(minCap)
			return out
		}
	} else {
		if len(n.children) == 0 {
			return nil
		} else {
			child := n.children[i]
			return child.remove(t, minCap)
		}

	}
}

func (parent *node) growChildAndRemove(index int, t Item, minCap int) Item {

	// 从左兄弟子树偷一个尾部item放入parent item[i-1],然后源parent item[i-1]下移插入到parent.children[i]的头部
	if index > 0 && len(parent.children[index-1].items) > minCap {

		stoleNode := parent.children[index-1]
		stoleNodeIndex := len(parent.children[index-1].items) - 1
		stoleItem := stoleNode.items.removeAt(stoleNodeIndex)
		parent.children[index].items.insertAt(0, parent.items[index-1])
		parent.items[index-1] = stoleItem
		if len(stoleNode.children) != 0 {
			lastIndex := len(stoleNode.children) - 1
			nodeLastChildren := stoleNode.children.removeAt(lastIndex)
			parent.children[index].children.insertAt(0, nodeLastChildren)
		}

		return parent.remove(t, minCap)

	}

	// 从右子树偷
	if index < len(parent.children)-1 && len(parent.children[index+1].items) > minCap {

		stoleNode := parent.children[index+1]
		parent.children[index].items = append(parent.children[index].items, parent.items[index])
		stoleItem := stoleNode.items.removeAt(0)
		parent.items[index] = stoleItem

		if len(stoleNode.children) != 0 {
			nodeFirstChildren := stoleNode.children.removeAt(0)
			parent.children[index].children = append(parent.children[index].children, nodeFirstChildren)
		}

		return parent.remove(t, minCap)

	}

	//无可偷兄弟节点，需要做一次兄弟合并
	if index > 0 {
		// 非最左节点， 则合并左兄弟
		n := &node{}
		n.items = append(n.items, parent.children[index-1].items...)
		n.items = append(n.items, parent.items[index-1])
		n.items = append(n.items, parent.children[index].items...)

		if len(parent.children[index-1].children) != 0 {
			n.children = append(n.children, parent.children[index-1].children...)
			n.children = append(n.children, parent.children[index].children...)
		}

		// 处理根节点被删空的情况
		if len(parent.items) == 1 {
			*parent = *n
		} else {
			parent.items.removeAt(index - 1)
			parent.children.removeAt(index - 1)
			parent.children[index-1] = n //合并后index前移一位
		}
		return parent.remove(t, minCap)
	}

	if index+1 < len(parent.children) {
		// 最左边不可偷节点，合并右侧兄弟节点
		n := &node{}
		n.items = append(n.items, parent.children[index].items...)
		n.items = append(n.items, parent.items[index])
		n.items = append(n.items, parent.children[index+1].items...)

		if len(parent.children[index].children) != 0 {
			n.children = append(n.children, parent.children[index].children...)
			n.children = append(n.children, parent.children[index+1].children...)
		}

		if len(parent.items) == 1 {
			*parent = *n
		} else {
			parent.items.removeAt(index)
			parent.children.removeAt(index + 1)
			parent.children[index] = n
		}

		return parent.remove(t, minCap)
	}

	return nil

}

func (parent *node) growChildAndRemoveMax(lastChildIndex int, minCap int) Item {

	// 从左兄弟子树偷一个尾部item放入parent item[i-1],然后源parent item[i-1]下移插入到parent.children[i]的头部
	if len(parent.children[lastChildIndex-1].items) > minCap {
		lastItemIndex := lastChildIndex - 1

		stoleNode := parent.children[lastChildIndex-1]
		stoleItemIndex := len(stoleNode.items) - 1
		stoleItem := stoleNode.items.removeAt(stoleItemIndex)
		downItem := parent.items[lastItemIndex]
		parent.children[lastChildIndex].items.insertAt(0, downItem)
		parent.items[lastItemIndex] = stoleItem
		if len(stoleNode.children) != 0 {
			lastIndex := len(stoleNode.children) - 1
			nodeLastChildren := stoleNode.children.removeAt(lastIndex)
			parent.children[lastChildIndex].children.insertAt(0, nodeLastChildren)
		}
	} else {
		// combing subling

		lastItemIndex := lastChildIndex - 1
		n := &node{}
		n.items = append(n.items, parent.children[lastChildIndex-1].items...)
		n.items = append(n.items, parent.items[lastItemIndex])
		n.items = append(n.items, parent.children[lastChildIndex].items...)

		if len(parent.children[lastChildIndex-1].children) != 0 {
			n.children = append(n.children, parent.children[lastChildIndex-1].children...)
			n.children = append(n.children, parent.children[lastChildIndex].children...)
		}

		if len(parent.items) == 1 {
			*parent = *n //注意这里要对对指针指向内容重新赋值
		} else {
			parent.items.removeAt(lastItemIndex)
			parent.children.removeAt(lastChildIndex)
			parent.children[lastChildIndex-1] = n
		}

	}

	return parent.removeMax(minCap)
}

func (n *node) removeMax(minCap int) Item {
	lastItemIndex := len(n.items) - 1
	if len(n.children) == 0 {
		return n.items.removeAt(lastItemIndex)
	} else if len(n.children[lastItemIndex+1].items) <= minCap {
		lastChildIndex := lastItemIndex + 1
		return n.growChildAndRemoveMax(lastChildIndex, minCap) //
	}

	lastChildIndex := lastItemIndex + 1
	lastChildNode := n.children[lastChildIndex]
	return lastChildNode.removeMax(minCap)
}

func (n *node) removeMin() {

}

func (n *node) print() {

	l := len(n.items)
	if len(n.children) != 0 {
		if len(n.items)+1 != len(n.children) {
			fmt.Println("tree faile-----------------")
		}
		for i, _ := range n.items {
			n.children[i].print()
			fmt.Println("nnn:", n.items[i])
		}

		n.children[l].print()
	} else {
		for i, _ := range n.items {
			fmt.Println("ddd:", n.items[i])
		}
	}

}

func (n *node) max() Item {
	if n == nil {
		return nil
	}
	for len(n.children) > 0 {
		n = n.children[len(n.children)-1]
	}
	if len(n.items) == 0 {
		return nil
	}
	return n.items[len(n.items)-1]
}

func (n *node) min() Item {
	if n == nil {
		return nil
	}
	for len(n.children) > 0 {
		n = n.children[0]
	}
	if len(n.items) == 0 {
		return nil
	}
	return n.items[0]
}
