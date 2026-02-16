package core

import "lemin/util"

type ant struct {
	id           int
	location     *Room
	finished     bool
	ignoredRooms RoomStates
}

func createAnts(farm *Farm) []*ant {
	ants := make([]*ant, farm.Ants)
	for i := 0; i < len(ants); i++ {
		ants[i] = &ant{
			id:           i,
			location:     farm.Start,
			ignoredRooms: util.MapWithValue(farm.Rooms, false),
			finished:     false,
		}
	}
	return ants
}

func (ant *ant) nextMove(farm *Farm, dst *Room, ignoredRooms RoomStates, ignoredTunnels TunnelStates) (*Room, *Tunnel) {
	src := ant.location

	path := dijkstra(farm, src, dst, mergedRoomStates(ant.ignoredRooms, ignoredRooms), ignoredTunnels)
	if len(path) == 0 {
		return nil, nil
	}

	next := path[0]
	tunnel := src.tunnelToRoom(next)
	return next, tunnel
}

func mergedRoomStates(a RoomStates, b RoomStates) RoomStates {
	c := make(map[*Room]bool, len(a))
	for k := range a {
		c[k] = a[k] || b[k]
	}
	return c
}
