package core

import "fmt"

type Room struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	IsStart   bool      `json:"isStart"`
	IsEnd     bool      `json:"isEnd"`
	Tunnels   []*Tunnel `json:"-"`
	tentative float64   `json:"-"`
}

type Tunnel struct {
	Id       int     `json:"id"`
	Distance float64 `json:"distance"`
	From     *Room   `json:"-"`
	To       *Room   `json:"-"`
}

type (
	RoomStates   = map[*Room]bool
	TunnelStates = map[*Tunnel]bool
)

func (r *Room) oppositeRoom(tunnel *Tunnel) *Room {
	if tunnel.From == r {
		return tunnel.To
	} else if tunnel.To == r {
		return tunnel.From
	}

	// This shouldn't happen
	panic(fmt.Sprintf("Tunnel %d isn't linked to room %d", tunnel.Id, r.Id))
}

func (r *Room) tunnelToRoom(next *Room) *Tunnel {
	for _, t := range r.Tunnels {
		if t.From == next || t.To == next {
			return t
		}
	}

	// This shouldn't happen either
	panic(fmt.Sprintf("There is no tunnel from room %d to %d", r.Id, next.Id))
}
