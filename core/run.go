package core

import (
	"time"
)

type Result struct {
	// Nanoseconds taken for parsing
	ParseTime time.Duration
	// Nanoseconds taken for pathfinding
	PathfindingTime time.Duration
	TurnCount       int
	Farm            *Farm
	Turns           []Turn
}

// Run parses the provided farm description and executes the
// pathfinding algorithm.
//
// It returns a Result containing the parsed farm, computed turns,
// and timing information for parsing and pathfinding.
//
// If parsing fails, the returned Result contains only ParseTime.
// If pathfinding fails, the returned Result contains the parsed
// Farm and timing information up to that point.
func Run(farmDesc string) (*Result, error) {
	res := &Result{}

	start := time.Now()
	farm, err := Parse(farmDesc)
	res.ParseTime = time.Since(start)

	if err != nil {
		return res, err
	}
	res.Farm = farm

	start = time.Now()
	turns, err := Lemin(farm)
	res.PathfindingTime = time.Since(start)

	if err != nil {
		return res, err
	}
	res.Turns = turns
	res.TurnCount = len(turns)

	return res, nil
}
