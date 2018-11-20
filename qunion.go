package vision

// quickUnion is a union-find data structure implementation that uses some algorithmic techniques for time performance improvement
type quickUnion struct {
	ID   []int
	sz   []int
	Size int
}

// newQuickUnion returns a quickUnion of Size s
func newQuickUnion(s int) quickUnion {
	qu := quickUnion{
		ID:   make([]int, s),
		sz:   make([]int, s),
		Size: s,
	}
	for i := 0; i < qu.Size; i++ {
		qu.ID[i] = i
	}
	return qu
}

func (qu quickUnion) root(i int) int {
	for i != qu.ID[i] {
		qu.ID[i] = qu.ID[qu.ID[i]]
		i = qu.ID[i]
	}
	return i
}

// find returns true when p and q are connected by some path in the underlying QuickFind or false otherwise
func (qu quickUnion) find(p, q int) bool {
	return qu.root(p) == qu.root(q)
}

// unite makes p immediately connected to q in the underlying QuickFind
func (qu quickUnion) unite(p, q int) {
	i := qu.root(p)
	j := qu.root(q)
	if qu.sz[i] < qu.sz[j] {
		qu.ID[i] = j
		qu.sz[j] += qu.sz[i]
	} else {
		qu.ID[j] = i
		qu.sz[i] += qu.sz[j]
	}
}
