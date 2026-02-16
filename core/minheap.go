package core

type vertexHeap []*Room

func (h vertexHeap) Len() int           { return len(h) }
func (h vertexHeap) Less(i, j int) bool { return h[i].tentative < h[j].tentative }
func (h vertexHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *vertexHeap) Push(x any) {
	*h = append(*h, x.(*Room))
}

func (h *vertexHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
