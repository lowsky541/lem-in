package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	lemin "lemin/pkg"
	"log"
	"net/http"
	"os"
	"time"
)

func Fatal(e error) {
	fmt.Printf("Error: %s\n", e.Error())
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: lem-in FILEPATH")
		return
	}
	filename := args[0]
	environment := lemin.ParseEnvironment()

	start := time.Now()

	context := lemin.Context{}
	if err := context.Parse(filename); err != nil {
		Fatal(err)
		return
	}
	context.PrintBanner()

	fmt.Printf("Parsed input file in %v.\n\n", time.Since(start))

	if !environment.NoSanityCheck {
		if !lemin.IsSane(&context) {
			Fatal(lemin.ErrInsaneGraph)
			return
		}
	}

	var turns []lemin.Turn
	if !environment.NoPathFinding {
		start = time.Now()

		turns = lemin.Lemin(&context)
		fmt.Printf("\nFinished in %d turns.\n", len(turns))
		fmt.Printf("Time is %v.\n", time.Since(start))
	}

	if environment.EnableVisualizer {
		addr := fmt.Sprintf("127.0.0.1:%d", environment.Port)

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
