package btree

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	seed := time.Now().Unix()
	fmt.Println(seed)
	rand.Seed(seed)
}

// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []Item) {
	for _, v := range rand.Perm(n) {
		out = append(out, IntItem(v))
	}
	return
}

// rang returns an ordered list of Int items in the range [0, n).
func rang(n int) (out []Item) {
	for i := 0; i < n; i++ {
		out = append(out, IntItem(i))
	}
	return
}

func TestTree(t *testing.T) {
	tr := NewTree(3)
	fmt.Println("degree:", tr.degree)
	if tr.degree != 3 {
		t.Fatalf("error new tree")
	}
}

// func TestInsert(t *testing.T) {
// 	tr := NewTree(3)

// 	item1 := IntItem(1)
// 	fmt.Println("start")
// 	tr.ReplaceOrInsert(item1)
// 	if len(tr.root.items) != 1 {
// 		t.Fatalf("error insert")
// 	}
// }

func TestBTree(t *testing.T) {
	tr := NewTree(32)
	const treeSize = 1000
	for i := 0; i < 10; i++ {
		if min := tr.Min(); min != nil {
			t.Fatalf("empty min, got %+v", min)
		}
		if max := tr.Max(); max != nil {
			t.Fatalf("empty max, got %+v", max)
		}
		for _, item := range perm(treeSize) {
			if x := tr.ReplaceOrInsert(item); x != nil {
				t.Fatal("insert found item", item)
			}
		}
		for _, item := range perm(treeSize) {
			if x := tr.ReplaceOrInsert(item); x == nil {
				t.Fatal("insert didn't find item", item)
			}
		}
		if min, want := tr.Min(), Item(IntItem(0)); min != want {
			t.Fatalf("min: want %+v, got %+v", want, min)
		}
		if max, want := tr.Max(), Item(IntItem(treeSize-1)); max != want {
			t.Fatalf("max: want %+v, got %+v", want, max)
		}

		// got := all(tr)
		// want := rang(treeSize)
		// if !reflect.DeepEqual(got, want) {
		// 	t.Fatalf("mismatch:\n got: %v\nwant: %v", got, want)
		// }

		// gotrev := allrev(tr)
		// wantrev := rangrev(treeSize)
		// if !reflect.DeepEqual(gotrev, wantrev) {
		// 	t.Fatalf("mismatch:\n got: %v\nwant: %v", got, want)
		// }

		// for _, item := range perm(treeSize) {
		// 	if x := tr.Delete(item); x == nil {
		// 		t.Fatalf("didn't find %v", item)
		// 	}
		// }
		// if got = all(tr); len(got) > 0 {
		// 	t.Fatalf("some left!: %v", got)
		// }
	}
}
