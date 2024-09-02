package lemin

import (
	"fmt"
)

const PATH_REPULSION = 999

func NewUsedTunnelsMap(ctx *Context) map[*Tunnel]bool {
	out := map[*Tunnel]bool{}
	for _, t := range ctx.Tunnels {
		out[t] = false
	}
	return out
}

func NewUsedRoomsMap(ctx *Context) map[*Room]bool {
	out := map[*Room]bool{}
	for _, r := range ctx.Rooms {
		out[r] = false
	}
	return out
}

func IsSane(ctx *Context) bool {
	path := Dijkstra(ctx, ctx.Start, ctx.End, nil, nil)
	return len(path) != 0
}

func Lemin(ctx *Context) []Turn {
	turns := []Turn{}
	ants := createAnts(ctx)

	currentTurn := 1

	arrived := 0
	total := ctx.Ants

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

	usedRoomsMap := NewUsedRoomsMap(ctx)

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
		usedTunnelsMap := NewUsedTunnelsMap(ctx)

		for currentAnt := 0; currentAnt != total; {
			// Skip ants that are already at end
			ant := ants[currentAnt]
			if ant.Finished {
				currentAnt++
				continue
			}

			if currentAnt == total-1 && total-arrived <= 2 && ant.Location == ctx.Start && !lastHasWaited {
				// Handle an exception,
				// if the ant is the last one, make it wait for the shortest path to free
				// I know, that's not really optimized, barely the same code two lines under...
				lastHasWaited = true

				next, tunnel := ant.NextMove(ctx, ctx.End, nil, nil)
				if next == nil || usedTunnelsMap[tunnel] || usedRoomsMap[next] {
					break
				}
			}

			// Declaration of variables mainly for readability purpose
			last := ant.Location
			next, tunnel := ant.NextMove(ctx, ctx.End, ignoredRooms, ignoredTunnels)

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
			turn.Moves = append(turn.Moves, Move{
				Ant:  ant.Id + 1,
				From: last.Id,
				To:   next.Id,
			})
			fmt.Printf("L%d-%s ", ant.Id+1, next.Name)

			// Release the room usage
			usedRoomsMap[last] = false

			// Move the ant
			ant.Location = next

			// Forbid the ant to go back
			ant.IgnoredRooms[last] = true

			// Mark this tunnel and room as used
			usedTunnelsMap[tunnel] = true
			if !next.IsStart && !next.IsEnd {
				usedRoomsMap[next] = true
			}

			if next.IsEnd {
				// Mark the ant as arrived: skip it for the future turns
				arrived++
				ant.Finished = true
			}

			// We are going to work on the next ant,
			// all prohibited rooms must be reset
			ignoredRooms = map[*Room]bool{}

			// This ant has moved, work on the next ant
			// NOTE: We are currently in a for-loop, there is no need for
			// a `continue`
			currentAnt++
		}

		fmt.Println()

		// Mark the turn as done and append it to the list of turns
		currentTurn++
		turns = append(turns, turn)
	}

	return turns
}
