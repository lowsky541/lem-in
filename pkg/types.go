package lemin

type ContextResponse struct {
	Ants    int       `json:"ants"`
	Start   uint      `json:"startId"`
	End     uint      `json:"endId"`
	Rooms   []*Room   `json:"rooms"`
	Tunnels []*Tunnel `json:"tunnels"`
	Turns   []Turn    `json:"turns"`
}

type Room struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	IsStart   bool      `json:"isStart"`
	IsEnd     bool      `json:"isEnd"`
	Tunnels   []*Tunnel `json:"-"`
	Tentative float64   `json:"-"`
}

type Tunnel struct {
	Id       uint    `json:"id"`
	Distance float64 `json:"distance"`
	FromId   uint    `json:"fromId"`
	ToId     uint    `json:"toId"`
	From     *Room   `json:"-"`
	To       *Room   `json:"-"`
}

type Move struct {
	Ant  int  `json:"ant"`
	From uint `json:"fromId"`
	To   uint `json:"toId"`
}

type Turn struct {
	Moves []Move `json:"moves"`
}
