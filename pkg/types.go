package lemin

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
	FromId   uint    `json:"from"`
	ToId     uint    `json:"to"`
	From     *Room   `json:"-"`
	To       *Room   `json:"-"`
}

type Move struct {
	Ant  int  `json:"ant"`
	From uint `json:"from"`
	To   uint `json:"to"`
}

type Turn struct {
	Moves []Move `json:"moves"`
}
