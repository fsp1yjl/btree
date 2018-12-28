package btree

type IntItem int

// Less returns true if int(a) < int(b).
func (a IntItem) Less(b Item) bool {
	return a < b.(IntItem)
}
