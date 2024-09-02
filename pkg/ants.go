package lemin

type Ant struct {
	Id           int
	Location     *Room
	Finished     bool
	IgnoredRooms map[*Room]bool
}

func newIgnoredRooms(context *Context) map[*Room]bool {
	out := map[*Room]bool{}
	for _, r := range context.Rooms {
		out[r] = false
	}
	return out
}

func createAnts(context *Context) []*Ant {
	ants := make([]*Ant, context.Ants)
	for i := 0; i < len(ants); i++ {
		ants[i] = &Ant{
			Id:           i,
			Location:     context.Start,
			Finished:     false,
			IgnoredRooms: newIgnoredRooms(context),
		}
	}
	return ants
}

func mergeIgnored(a map[*Room]bool, b map[*Room]bool) map[*Room]bool {
	c := make(map[*Room]bool, len(a))
	for k := range a {
		c[k] = a[k] || b[k]
	}
	return c
}

func (ant *Ant) NextMove(ctx *Context, dst *Room, ignoredRooms map[*Room]bool, ignoredTunnels map[*Tunnel]bool) (*Room, *Tunnel) {
	merged := mergeIgnored(ant.IgnoredRooms, ignoredRooms)

	src := ant.Location
	path := Dijkstra(ctx, src, dst, merged, ignoredTunnels)
	if len(path) == 0 {
		return nil, nil
	}

	next := path[0]
	tunnel := GetTunnel(ctx, src, next)

	return next, tunnel
}
