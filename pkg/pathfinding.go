package lemin

import (
	"container/heap"
	"math"
)

func Opposite(current *Room, tunnel *Tunnel) *Room {
	if tunnel.From == current {
		return tunnel.To
	} else if tunnel.To == current {
		return tunnel.From
	} else {
		panic("A tunnel was not bound to the right room")
	}
}

func GetTunnel(ctx *Context, a *Room, b *Room) *Tunnel {
	for _, t := range ctx.Tunnels {
		if (t.From == a && t.To == b) || (t.From == b && t.To == a) {
			return t
		}
	}
	return nil
}

func Dijkstra(ctx *Context, source *Room, destination *Room, ignoredRooms map[*Room]bool, ignoredTunnels map[*Tunnel]bool) []*Room {
	queue := VertexHeap{}
	heap.Init(&queue)

	previous := map[*Room]*Room{}
	visited := []*Room{}

	// Initialize graph tentatives and queue
	for _, r := range ctx.Rooms {
		if r == source {
			r.Tentative = 0
		} else {
			r.Tentative = math.Inf(0)
		}
	}

	// Push the source node on the queue
	heap.Push(&queue, source)

	for queue.Len() != 0 {
		current := heap.Pop(&queue).(*Room)

		if Contains(current, visited) {
			continue
		}

		if current == destination {
			break
		}

		for _, t := range current.Tunnels {
			adjacent := Opposite(current, t)

			if ignoredRooms[adjacent] || ignoredTunnels[t] {
				continue
			}

			alt := current.Tentative + t.Distance
			if alt < adjacent.Tentative {
				adjacent.Tentative = alt
				previous[adjacent] = current
			}
			heap.Push(&queue, adjacent)
		}

		visited = append(visited, current)
	}

	// Extract the path
	extract := []*Room{destination}
	extractCurrent := destination
	for {
		extractCurrent = previous[extractCurrent]
		if extractCurrent == nil {
			return nil
		} else if extractCurrent == source {
			break
		}
		extract = append(extract, extractCurrent)
	}

	// Reverse the path extraction
	for i, j := 0, len(extract)-1; i < j; i, j = i+1, j-1 {
		extract[i], extract[j] = extract[j], extract[i]
	}

	return extract
}

///////////////////////////////////////////////////////////////////////////////////
//                                MIN-HEAP                                       //
///////////////////////////////////////////////////////////////////////////////////

type VertexHeap []*Room

func (h VertexHeap) Len() int           { return len(h) }
func (h VertexHeap) Less(i, j int) bool { return h[i].Tentative < h[j].Tentative }
func (h VertexHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *VertexHeap) Push(x interface{}) {
	*h = append(*h, x.(*Room))
}

func (h *VertexHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
