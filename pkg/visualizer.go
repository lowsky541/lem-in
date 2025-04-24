package lemin

import (
	"encoding/json"
	"net/http"
)

type ContextResponse struct {
	Ants    int       `json:"ants"`
	Start   uint      `json:"start"`
	End     uint      `json:"end"`
	Rooms   []*Room   `json:"rooms"`
	Tunnels []*Tunnel `json:"tunnels"`
	Turns   []Turn    `json:"turns"`
}

func ServeVisualizer(assets http.FileSystem, parser *Context, addr string, turns []Turn) error {
	srv := http.Server{Addr: addr}

	http.Handle("/", http.FileServer(assets))
	http.HandleFunc("/context", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(
			ContextResponse{
				Ants:    parser.Ants,
				Start:   parser.Start.Id,
				End:     parser.End.Id,
				Rooms:   parser.Rooms,
				Tunnels: parser.Tunnels,
				Turns:   turns,
			},
		)

		if err != nil {
			panic(err)
		}
	})

	return srv.ListenAndServe()
}
