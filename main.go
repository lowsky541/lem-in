package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	lemin "lemin/pkg"
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
	if isEmpty(assetsDir) {
		fsys, err := fs.Sub(assets, "visualizer/dist")
		if err != nil {
			panic(err)
		}
		return http.FS(fsys)
	} else {
		return http.Dir(assetsDir)
	}
}

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func main() {
	var fileInfo os.FileInfo
	var err error
	var bind, assetsDir string
	var visualize, port uint

	flag.UintVar(&visualize, "visualize", 0, "Enable visualizer HTTP server; 0 = disable, 1 = enable; 2 = enable without pathfinding")
	flag.StringVar(&bind, "bind", "127.0.0.1", "Specify bind address for visualizer")
	flag.UintVar(&port, "port", 3000, "Specify port for visualizer")
	flag.StringVar(&assetsDir, "assets-dir", "", "Serve this directory instead of embedded assets; enables visualizer; useful for debugging")
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

	fileInfo, err = os.Stat(filename)
	if err != nil {
		fatal(err)
	}

	if !fileInfo.Mode().IsRegular() {
		fatal(fmt.Errorf("%s: not a regular file", filename))
	}

	if !isEmpty(assetsDir) {
		// Doesn't make sense to set the assets dir without
		// enabling the visualizer
		visualize = 1

		fileInfo, err = os.Stat(assetsDir)
		if err != nil {
			fatal(err)
		}

		if !fileInfo.IsDir() {
			fatal(fmt.Errorf("%s: not a directory", assetsDir))
		}
	}

	start := time.Now()
	context := lemin.Context{}
	if err := context.Parse(filename); err != nil {
		fatal(err)
	}
	context.PrintBanner()

	fmt.Printf("Parsed input file in %v.\n", time.Since(start))
	if visualize < 2 {
		fmt.Println()
	}

	// Check if an ant can go from the start to the
	// end node by running Dijkstra.
	if visualize < 2 && !lemin.IsSane(&context) {
		fatal(lemin.ErrInsaneGraph)
	}

	turns := []lemin.Turn{}
	if visualize < 2 {
		start = time.Now()

		turns = lemin.Lemin(&context)
		fmt.Printf("\nFinished in %v and %d turns.\n", time.Since(start), len(turns))
	}

	if visualize > 0 {
		fsys := assetsFilesystem(assetsDir)
		addr := fmt.Sprintf("%s:%d", bind, port)

		fmt.Printf("\nAll done, will now serve at http://%s...\n", addr)
		lemin.ServeVisualizer(fsys, &context, addr, turns)
	}
}
