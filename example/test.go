package main

import (
	"fmt"
	"math/rand"
	"time"

	t "github.com/fsp1yjl/btree"
)

func main() {
	// t := b.New(2)

	// item := b.Int(3)
	// t.ReplaceOrInsert(item)
	// item = b.Int(5)
	// t.ReplaceOrInsert(item)
	// item = b.Int(7)
	// t.ReplaceOrInsert(item)
	// item = b.Int(9)
	// t.ReplaceOrInsert(item)
	// item = b.Int(11)
	// t.ReplaceOrInsert(item)

	// fmt.Println("len::", t.Len())

	test()
}

func init() {
	seed := time.Now().Unix()
	fmt.Println(seed)
	rand.Seed(seed)
}

// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []t.Item) {
	for _, v := range rand.Perm(n) {
		out = append(out, t.IntItem(v))
	}
	return
}

// rang returns an ordered list of Int items in the range [0, n).
func rang(n int) (out []t.Item) {
	for i := 0; i < n; i++ {
		out = append(out, t.IntItem(i))
	}
	return
}

func test() {

	tr := t.NewTree(3)
	fmt.Println("min:", tr.Min())

	const treeSize = 500
	// for i := 0; i < 10; i++ {
	if min := tr.Min(); min != nil {
		fmt.Println("empty min, got %+v", min)
	}
	if max := tr.Max(); max != nil {
		fmt.Println("empty max, got %+v", max)
	}

	for i := 0; i < 10; i++ {
		for _, item := range perm(treeSize) {
			fmt.Println("item:::insert:", item)
			if x := tr.ReplaceOrInsert(item); x != nil {
				fmt.Println("xxx:", x)
				fmt.Println("insert found item", item)
			}
		}
		for _, item := range perm(treeSize) {
			if x := tr.ReplaceOrInsert(item); x == nil {
				fmt.Println("insert didn't find item", item)
			}
		}

		for _, item := range perm(treeSize) {
			if x := tr.Delete(item); x == nil {
				fmt.Println("delete item failed", item)
			}
		}
	}

	fmt.Println("expect empty btree print")
	tr.Print()

}
