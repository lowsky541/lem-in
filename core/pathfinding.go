package core

import (
	"container/heap"
	"math"
	"slices"
)

func dijkstra(farm *Farm, source *Room, destination *Room, ignoredRooms RoomStates, ignoredTunnels TunnelStates) []*Room {
	queue := vertexHeap{}
	previous := map[*Room]*Room{}
	visited := []*Room{}

	// Initialize tentatives
	for _, r := range farm.Rooms {
		if r == source {
			r.tentative = 0
		} else {
			r.tentative = math.Inf(0)
		}
	}

	// Initialize queue
	heap.Init(&queue)
	heap.Push(&queue, source)

	for queue.Len() != 0 {
		current := heap.Pop(&queue).(*Room)

		if slices.Contains(visited, current) {
			continue
		}

		if current == destination {
			break
		}

		for _, tunnel := range current.Tunnels {
			opposite := current.oppositeRoom(tunnel)

			if ignoredRooms[opposite] || ignoredTunnels[tunnel] {
				continue
			}

			alt := current.tentative + tunnel.Distance
			if alt < opposite.tentative {
				opposite.tentative = alt
				previous[opposite] = current
			}

			heap.Push(&queue, opposite)
		}

		visited = append(visited, current)
	}

	return extractPath(source, destination, previous)
}

func extractPath(source *Room, destination *Room, previous map[*Room]*Room) []*Room {
	path := []*Room{destination}

	for extract := destination; ; extract = previous[extract] {
		if extract == nil {
			return nil
		} else if extract == source {
			break
		}
		path = append(path, extract)
	}

	slices.Reverse(path)
	return path
}
