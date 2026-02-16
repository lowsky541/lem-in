package core

import "lemin/util"

type Farm struct {
	Ants    int       `json:"ants"`
	Start   *Room     `json:"start"`
	End     *Room     `json:"end"`
	Rooms   []*Room   `json:"rooms"`
	Tunnels []*Tunnel `json:"tunnels"`
}

type Turn = []Move

type Move struct {
	Ant    int
	From   *Room
	To     *Room
	Tunnel *Tunnel
}

func IsSane(farm *Farm) bool {
	path := dijkstra(farm, farm.Start, farm.End, nil, nil)
	return len(path) != 0
}

func Lemin(farm *Farm) ([]Turn, error) {
	if !IsSane(farm) {
		return nil, ErrInsaneFarm
	}

	turns := []Turn{}
	ants := createAnts(farm)

	currentTurn := 1

	arrived := 0
	total := farm.Ants

	lastHasWaited := false

	// Difference between Used and Ignored:
	// When a room or a tunnel is Used the path-finding
	// can still go through it.
	// If it is Ignored, the path-finder will execute
	// like the room or the tunnel does not exists.
	// We can't directly Ignore rooms or tunnels because,
	// for example, if the tunnel is the only way to the end,
	// then the path-finder will fail even if the ant
	// can go one room forward.

	usedRoomsMap := util.MapWithValue(farm.Rooms, false)

	for arrived != total {

		// Loop until all the ants arrive to the end room.
		// This loop is considered as a single turn.

		// The informations about the current turn.
		turn := Turn{}

		// The ignored tunnels for this turn.
		ignoredTunnels := map[*Tunnel]bool{}

		// The ignored rooms for AN (single) ant.
		// It is reset for the next ant.
		ignoredRooms := map[*Room]bool{}

		// Keep track of which tunnel has been used.
		usedTunnelsMap := util.MapWithValue(farm.Tunnels, false)

		for currentAnt := 0; currentAnt != total; {
			// Skip ants that are already at end
			ant := ants[currentAnt]
			if ant.finished {
				currentAnt++
				continue
			}

			if currentAnt == total-1 && total-arrived <= 2 && ant.location == farm.Start && !lastHasWaited {
				// Handle an exception,
				// if the ant is the last one, make it wait for the shortest path to free
				// I know, that's not really optimized, barely the same code two lines under...
				lastHasWaited = true

				next, tunnel := ant.nextMove(farm, farm.End, nil, nil)
				if next == nil || usedTunnelsMap[tunnel] || usedRoomsMap[next] {
					break
				}
			}

			// Declaration of variables mainly for readability purpose
			last := ant.location
			next, tunnel := ant.nextMove(farm, farm.End, ignoredRooms, ignoredTunnels)

			if next == nil {
				// No more paths, restart path-finding
				currentAnt++
				continue
			} else if usedTunnelsMap[tunnel] {
				// The tunnel has already been used for this turn,
				// restart path-finding
				ignoredTunnels[tunnel] = true
				continue
			} else if usedRoomsMap[next] {
				ignoredRooms[next] = true
				continue
			}

			// Add the move to the current turn
			turn = append(turn, Move{
				Ant:    ant.id + 1,
				From:   last,
				To:     next,
				Tunnel: tunnel,
			})

			// Release the room usage
			usedRoomsMap[last] = false

			// Move the ant
			ant.location = next

			// Forbid the ant to go back
			ant.ignoredRooms[last] = true

			// Mark this tunnel and room as used
			usedTunnelsMap[tunnel] = true
			if !next.IsStart && !next.IsEnd {
				usedRoomsMap[next] = true
			}

			if next.IsEnd {
				// Mark the ant as arrived: skip it for the future turns
				arrived++
				ant.finished = true
			}

			// We are going to work on the next ant,
			// all prohibited rooms must be reset
			ignoredRooms = map[*Room]bool{}

			// This ant has moved, work on the next ant
			// NOTE: We are currently in a for-loop, there is no need for
			// a `continue`
			currentAnt++
		}

		// Mark the turn as done and append it to the list of turns
		currentTurn++
		turns = append(turns, turn)
	}

	return turns, nil
}
