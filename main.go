package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	lemin "lemin/pkg/lemin"
	util "lemin/pkg/util"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed visualizer/dist
var assets embed.FS

func fatal(e error) {
	fmt.Printf("Error: %s\n", e.Error())
	os.Exit(1)
}

func assetsFilesystem(assetsDir string) http.FileSystem {
	if util.IsEmpty(assetsDir) {
		fsys, err := fs.Sub(assets, "visualizer/dist")
		if err != nil {
			panic(err)
		}
		return http.FS(fsys)
	} else {
		return http.Dir(assetsDir)
	}
}

func printFarm(farm *lemin.Farm) {
	fmt.Println(farm.Ants)

	for _, r := range farm.Rooms {
		fmt.Printf("%s %d %d\n", r.Name, r.X, r.Y)
	}

	for _, t := range farm.Tunnels {
		fmt.Printf("%s-%s\n", t.From.Name, t.To.Name)
	}

	fmt.Println()
}

func printTurns(turns []lemin.Turn) {
	for _, turn := range turns {
		for _, move := range turn {
			fmt.Printf("L%d-%s ", move.Ant, move.To.Name)
		}
		fmt.Println()
	}
}

func main() {
	var fi os.FileInfo
	var err error
	var bind, assetsDir string
	var visualize, port uint

	flag.UintVar(&visualize, "vis", 0, "Enable visualizer HTTP server; 0 = disable, 1 = enable; 2 = enable without pathfinding")
	flag.StringVar(&bind, "bind", "127.0.0.1", "Specify bind address for visualizer")
	flag.UintVar(&port, "port", 3000, "Specify port for visualizer")
	flag.StringVar(&assetsDir, "assets-dir", "", "Serve this directory instead of the embedded assets")
	flag.Usage = func() {
		fmt.Println("usage: lem-in [OPTION] FILENAME")
		flag.PrintDefaults()
	}
	flag.Parse()

	filename := flag.Arg(0)
	if strings.TrimSpace(filename) == "" {
		flag.Usage()
		return
	}

	fi, err = os.Stat(filename)
	if err != nil || !fi.Mode().IsRegular() {
		fatal(fmt.Errorf("%s: not a regular file", filename))
	}

	if !util.IsEmpty(assetsDir) {
		// Doesn't make sense to set the assets dir without
		// enabling the visualizer
		visualize = 1

		fi, err = os.Stat(assetsDir)
		if err != nil || !fi.IsDir() {
			fatal(fmt.Errorf("%s: not a directory", assetsDir))
		}
	}

	start := time.Now()
	farm, err := lemin.Parse(filename)
	if err != nil {
		fatal(err)
	}

	printFarm(farm)
	fmt.Printf("Parsed input file in %v.\n", time.Since(start))

	turns := []lemin.Turn{}
	if visualize < 2 {
		fmt.Println()

		start = time.Now()
		turns, err = lemin.Lemin(farm)
		if err != nil {
			fatal(err)
		}

		printTurns(turns)
		fmt.Printf("\nFinished in %v and %d turns.\n", time.Since(start), len(turns))
	}

	if visualize > 0 {
		fsys := assetsFilesystem(assetsDir)
		addr := fmt.Sprintf("%s:%d", bind, port)

		fmt.Printf("\nAll done, will now serve at http://%s...\n", addr)
		lemin.ServeVisualizer(fsys, farm, addr, turns)
	}
}
