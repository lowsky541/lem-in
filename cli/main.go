//go:build !js && !wasm

package main

import (
	"flag"
	"fmt"
	"lemin/core"
	"lemin/util"
	"os"
	"time"
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

	fmt.Println()
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

	start := time.Now()
	farm, err := core.ParseFromFilepath(filepath)
	if err != nil {
		fatal(err)
	}

	printFarm(farm)
	fmt.Printf("Parsed input file in %v.\n\n", time.Since(start))

	start = time.Now()
	turns, err := core.Lemin(farm)
	if err != nil {
		fatal(err)
	}

	printTurns(turns)
	fmt.Printf("\nFinished in %v and %d turns.\n", time.Since(start), len(turns))
}
