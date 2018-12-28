package btree

type BTree struct {
	degree uint
	root   *node
}

func (t *BTree) MaxCap() int {
	return int(2*t.degree - 1)
}

func (t *BTree) MinCap() int {
	return int(t.degree - 1)
}

func (t *BTree) ReplaceOrInsert(item Item) Item {

	//如果是空树，则创建根后插入
	if t.root == nil {
		t.root = t.newNode()
		t.root.items = append(t.root.items, item)
		return nil
	}

	// 从跟节点向下寻找合适的插入节点，对于新的item，最终的插入会落到叶子节点，下降过程中，每个经过的节点如果已满，聚会进行一次分裂动作，
	// 这样保证后续插入的时候，实际插入的节点肯定会有空闲空间供插入

	// 如果根结点已满，对根节点进行分裂处理
	if len(t.root.items) >= t.MaxCap() {
		midIndex := len(t.root.items) / 2
		upItem, newNode := t.root.split(midIndex)
		newRoot := t.newNode()
		newRoot.items = append(newRoot.items, upItem)
		newRoot.children = append(newRoot.children, t.root, newNode)
		t.root = newRoot
	}
	return t.root.insert(item, t.MaxCap())
}

func (t *BTree) newNode() *node {
	return &node{}
}

func (t *BTree) Delete(item Item) Item {
	return t.root.remove(item, t.MinCap())
}

func (t *BTree) Print() {
	t.root.print()
}

// func (t *BTree) PrintRoot() {
// 	fmt.Println("root::", t.root.items)
// 	for _, n := range t.root.children {
// 		fmt.Println("child 1:", n.items)
// 	}
// }

func (t *BTree) Depth() int {
	n := t.root
	if n == nil {
		return 0
	}
	d := 0
	for {
		d++
		if n.children == nil {
			return d
		} else {
			n = n.children[0]
		}
	}
}

func (t *BTree) Max() Item {

	return t.root.max()
}

func (t *BTree) Min() Item {

	return t.root.min()
}

func NewTree(degree uint) *BTree {
	return &BTree{
		degree: degree,
	}
}
