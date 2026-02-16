//go:build !js && !wasm

package main

import (
	"flag"
	"fmt"
	"lemin/core"
	"lemin/util"
	"os"
)

func fatal(e error) {
	fmt.Printf("Error: %s\n", e.Error())
	os.Exit(1)
}

func printFarm(farm *core.Farm) {
	fmt.Println(farm.Ants)

	for _, r := range farm.Rooms {
		fmt.Printf("%s %d %d\n", r.Name, r.X, r.Y)
	}

	for _, t := range farm.Tunnels {
		fmt.Printf("%s-%s\n", t.From.Name, t.To.Name)
	}
}

func printTurns(turns []core.Turn) {
	for _, turn := range turns {
		for _, move := range turn {
			fmt.Printf("L%d-%s ", move.Ant, move.To.Name)
		}
		fmt.Println()
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s FILENAME\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	filepath := flag.Arg(0)
	if util.IsEmpty(filepath) {
		flag.Usage()
		return
	}

	fi, err := os.Stat(filepath)
	if err != nil || !fi.Mode().IsRegular() {
		fatal(fmt.Errorf("%s: not a regular file", filepath))
	}

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		fatal(fmt.Errorf("%s: couldn't be read", filepath))
	}

	result, err := core.Run(string(bytes))

	printFarm(result.Farm)

	fmt.Printf("¤ Parsing took %v.\n", result.ParseTime)

	printTurns(result.Turns)

	fmt.Printf("¤ Pathfinding took %v.\n", result.PathfindingTime)
	fmt.Printf("¤ All done in %v and %d turns.\n", result.ParseTime+result.PathfindingTime, result.TurnCount)
}
