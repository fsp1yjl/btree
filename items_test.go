package btree

import "testing"

func TestItemRemove(t *testing.T) {
	var items Items
	items = append(items, IntItem(1))
	items = append(items, IntItem(2))
	items = append(items, IntItem(3))
	items = append(items, IntItem(4))

	// remove fisrt
	items.removeAt(0)
	if len(items) != 3 {
		t.Fatalf("error 111")
	}

	if items[0].Less(IntItem(2)) || IntItem(2).Less(items[0]) {
		t.Fatalf("error 2222")
	}

	// remove last
	items.removeAt(len(items) - 1)
	if len(items) != 2 {
		t.Fatalf("error 111")
	}
	if items[len(items)-1].Less(IntItem(3)) || IntItem(3).Less(items[len(items)-1]) {
		t.Fatalf("error 3333")
	}
}

func TestItemFind(t *testing.T) {
	var items Items
	items = append(items, IntItem(2))
	items = append(items, IntItem(4))
	items = append(items, IntItem(6))
	items = append(items, IntItem(9))

	leftOverFlow := IntItem(1)
	var i int
	var found bool
	i, found = items.find(leftOverFlow)
	if i != 0 || found {
		t.Fatalf("left overflow test failed")
	}

	firstItem := IntItem(2)
	i, found = items.find(firstItem)
	if i != 0 || !found {
		t.Fatalf(" find first item  test failed")
	}

	midItemNotFound := IntItem(3)
	i, found = items.find(midItemNotFound)
	if i != 1 || found {
		t.Fatalf(" mid item not found  test failed")
	}

	midItemFound := IntItem(6)
	i, found = items.find(midItemFound)
	if i != 2 || !found {
		t.Fatalf(" mid item  found  test failed")
	}

	lastItemFound := IntItem(9)
	i, found = items.find(lastItemFound)
	if i != 3 || !found {
		t.Fatalf(" last item  found  test failed")
	}

	rightOverFlow := IntItem(20)
	i, found = items.find(rightOverFlow)
	if i != 4 || found {
		t.Fatalf(" last item  not found  test failed")
	}
}

func TestTruncate(t *testing.T) {
	var items Items
	items = append(items, IntItem(2))
	items = append(items, IntItem(4))
	items = append(items, IntItem(6))
	items = append(items, IntItem(9))

	var temp Items
	temp = append(temp, items...)

	items.truncate(3)
	if len(items) != 3 {
		t.Fatalf("truncate test 1-1 failed")
	}
	for i, _ := range items {
		if !(items[i].(IntItem) == temp[i].(IntItem)) {
			t.Fatalf("truncate test 1-2 failed")
		}
	}
}
