package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	lemin "lemin/pkg"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Fatal(e error) {
	fmt.Printf("Error: %s\n", e.Error())
	os.Exit(1)
}

func main() {
	var visEnable, noPathFinding bool
	var visBind string
	var visPort int

	flag.BoolVar(&visEnable, "vis-enable", false, "Enable visualizer HTTP server")
	flag.StringVar(&visBind, "vis-bind", "127.0.0.1", "Specify bind address for visualizer")
	flag.IntVar(&visPort, "vis-port", 3000, "Specify port for visualizer")
	flag.BoolVar(&noPathFinding, "no-path-finding", false, "Do not run path finding")
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

	start := time.Now()
	context := lemin.Context{}
	if err := context.Parse(filename); err != nil {
		Fatal(err)
		return
	}
	context.PrintBanner()

	fmt.Printf("Parsed input file in %v.\n\n", time.Since(start))

	// Check if an ant can go from the start to the
	// end node by running Dijkstra.
	if !lemin.IsSane(&context) {
		Fatal(lemin.ErrInsaneGraph)
		return
	}

	turns := []lemin.Turn{}
	if !noPathFinding {
		start = time.Now()

		turns = lemin.Lemin(&context)
		fmt.Printf("\nFinished in %d turns.\n", len(turns))
		fmt.Printf("Time is %v.\n", time.Since(start))
	}

	if visEnable {
		addr := fmt.Sprintf("%s:%d", visBind, visPort)

		fmt.Printf("\nAll done, will now serve at http://%s...\n", addr)
		ServeVisualizer(&context, addr, turns)
	}
}

//go:embed visualizer/dist/assets
//go:embed visualizer/dist/index.html
var embeddedFS embed.FS

func ServeVisualizer(parser *lemin.Context, addr string, turns []lemin.Turn) {
	rootDir, err := fs.Sub(embeddedFS, "visualizer/dist")
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.FS(rootDir)))
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]any{
			`ants`:    parser.Ants,
			"start":   parser.Start.Id,
			"end":     parser.End.Id,
			"rooms":   parser.Rooms,
			"tunnels": parser.Tunnels,
			"turns":   turns,
		})

		if err != nil {
			panic(err)
		}
	})

	log.Fatalln(http.ListenAndServe(addr, nil))
}
