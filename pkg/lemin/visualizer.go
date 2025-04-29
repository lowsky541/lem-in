package lemin

import (
	"encoding/json"
	"net/http"
)

func ServeVisualizer(assets http.FileSystem, farm *Farm, addr string, turns []Turn) error {
	srv := http.Server{Addr: addr}

	http.Handle("/", http.FileServer(assets))
	http.HandleFunc("/farm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(
			FarmResponse{
				Ants:    farm.Ants,
				Start:   farm.Start.Id,
				End:     farm.End.Id,
				Rooms:   farm.Rooms,
				Tunnels: farm.Tunnels,
				Turns:   turns,
			},
		)

		if err != nil {
			panic(err)
		}
	})

	return srv.ListenAndServe()
}
