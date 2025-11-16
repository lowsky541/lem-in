package lemin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type (
	FarmRes struct {
		Ants    int                `json:"ants"`
		Start   string             `json:"start"`
		End     string             `json:"end"`
		Rooms   map[string]*Room   `json:"rooms"`
		Tunnels map[string]*Tunnel `json:"tunnels"`
		Turns   []TurnRes          `json:"turns"`
	}
	TurnRes = []MoveRes
	MoveRes struct {
		AntId  int    `json:"antId"`
		From   string `json:"from"`
		To     string `json:"to"`
		Tunnel string `json:"tunnel"`
	}
)

func intoMapRes[T any](keyfixer func(e T) string, arr []T) map[string]T {
	newMap := make(map[string]T, len(arr))
	for _, elem := range arr {
		key := keyfixer(elem)
		if _, ok := newMap[key]; ok {
			panic("duplicates within intoMap")
		}
		newMap[key] = elem
	}
	return newMap
}

func intoTurnRes(turns []Turn) []TurnRes {
	newTurns := make([]TurnRes, len(turns))
	for turnIndex, turn := range turns {
		newTurn := make(TurnRes, len(turn))
		for moveIndex, move := range turn {
			newTurn[moveIndex] = MoveRes{
				AntId:  move.Ant,
				From:   fmt.Sprintf("/rooms/%d", move.From.Id),
				To:     fmt.Sprintf("/rooms/%d", move.From.Id),
				Tunnel: fmt.Sprintf("/tunnels/%d", move.Tunnel.Id),
			}
		}
		newTurns[turnIndex] = newTurn
	}
	return newTurns
}

func ServeVisualizer(assets http.FileSystem, farm *Farm, addr string, turns []Turn) error {
	srv := http.Server{Addr: addr}

	http.Handle("/", http.FileServer(assets))
	http.HandleFunc("/farm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(
			FarmRes{
				Ants:    farm.Ants,
				Start:   "/rooms/" + strconv.Itoa(farm.Start.Id),
				End:     "/rooms/" + strconv.Itoa(farm.End.Id),
				Rooms:   intoMapRes(func(e *Room) string { return fmt.Sprintf("/rooms/%d", e.Id) }, farm.Rooms),
				Tunnels: intoMapRes(func(e *Tunnel) string { return fmt.Sprintf("/tunnels/%d", e.Id) }, farm.Tunnels),
				Turns:   intoTurnRes(turns),
			},
		)

		if err != nil {
			panic(err)
		}
	})

	return srv.ListenAndServe()
}
