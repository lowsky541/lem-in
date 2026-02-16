package core

import (
	"time"
)

type Result struct {
	// Both are in nanoseconds
	ParseTime       time.Duration
	PathfindingTime time.Duration
	TurnCount       int
	Farm            *Farm
	Turns           []Turn
}

func Run(farmDesc string) (*Result, error) {
	start := time.Now()
	farm, err := Parse(farmDesc)
	parseTime := time.Since(start)

	if err != nil {
		return &Result{
			ParseTime:       parseTime,
			PathfindingTime: 0,
			TurnCount:       0,
			Farm:            nil,
			Turns:           nil,
		}, err
	}

	start = time.Now()
	turns, err := Lemin(farm)
	pathfindingTime := time.Since(start)

	if err != nil {
		return &Result{
			ParseTime:       parseTime,
			PathfindingTime: pathfindingTime,
			TurnCount:       0,
			Farm:            farm,
			Turns:           nil,
		}, err
	}

	return &Result{
		ParseTime:       parseTime,
		PathfindingTime: pathfindingTime,
		TurnCount:       len(turns),
		Farm:            farm,
		Turns:           turns,
	}, nil
}
